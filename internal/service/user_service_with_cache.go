package service

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"huynhmanh.com/gin/internal/cache"
	"huynhmanh.com/gin/internal/model"
	"huynhmanh.com/gin/internal/repository"
)

const (
	UserCacheTTL     = 5 * time.Minute
	UserListCacheTTL = 1 * time.Minute
)

type CachedUserService struct {
	repo        repository.UserRepository
	cacheConn   cache.CacheClient
	logger      *zap.Logger
	baseService UserService
}

func NewCachedUserService(
	repo repository.UserRepository,
	cacheConn cache.CacheClient,
	logger *zap.Logger,
) UserService {
	baseService := NewUserService(repo)
	return &CachedUserService{
		repo:        repo,
		cacheConn:   cacheConn,
		logger:      logger,
		baseService: baseService,
	}
}

func (s *CachedUserService) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	return s.baseService.Register(ctx, name, email, password)
}

func (s *CachedUserService) Login(ctx context.Context, email, password string) (string, error) {
	return s.baseService.Login(ctx, email, password)
}

func (s *CachedUserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	key := GetUserCacheKey(id)

	cached, err := s.cacheConn.Get(ctx, key)
	if err == nil && cached != "" {
		s.logger.Info("cache_hit", zap.String("key", key))
		s.cacheConn.Increment(ctx, "cache:hits")

		var user model.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	s.logger.Info("cache_miss", zap.String("key", key))
	s.cacheConn.Increment(ctx, "cache:misses")

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	if data, err := json.Marshal(user); err == nil {
		_ = s.cacheConn.Set(ctx, key, string(data), UserCacheTTL)
	}
	return user, nil
}

func (s *CachedUserService) GetAllUsers(
	ctx context.Context,
	limit, offset int,
) ([]model.User, int64, error) {
	key := GetUserListCacheKey(limit, offset)

	cached, err := s.cacheConn.Get(ctx, key)
	if err == nil && cached != "" {
		s.logger.Info("cache_hit", zap.String("key", key))
		s.cacheConn.Increment(ctx, "cache:hits")

		var resp struct {
			Users []model.User
			Total int64
		}
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return resp.Users, resp.Total, nil
		}
	}

	users, total, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	for i := range users {
		users[i].Password = ""
	}

	resp := struct {
		Users []model.User
		Total int64
	}{
		Users: users,
		Total: total,
	}

	if data, err := json.Marshal(resp); err == nil {
		_ = s.cacheConn.Set(ctx, key, string(data), UserListCacheTTL)
	}

	return users, total, nil
}

func (s *CachedUserService) UpdateUser(
	ctx context.Context,
	id uint,
	name, email string,
) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.Name = name
	user.Email = email

	if _, err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	key := GetUserCacheKey(id)
	if err := s.cacheConn.Delete(ctx, key); err != nil {
		s.logger.Warn("cache_invalidation_failed", zap.String("key", key))
	}
	user.Password = ""
	return user, nil
}

func (s *CachedUserService) DeleteUser(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	key := GetUserCacheKey(id)
	_ = s.cacheConn.Delete(ctx, key)
	return nil
}
