package service

import "fmt"

const (
	UserCacheKeyPrefix = "user:"
	UserListCacheKey   = "user:list"
)

func GetUserCacheKey(userID uint) string {
	return fmt.Sprintf("%s%d", UserCacheKeyPrefix, userID)
}

func GetUserListCacheKey(limit, offset int) string {
	return fmt.Sprintf("%s:%d:%d", UserListCacheKey, limit, offset)
}
