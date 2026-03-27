package service

import "fmt"

const (
	UserCacheKeyPrefix = "user:"
	UserListCacheKey   = "user:list"
)

const (
	ProductCacheKeyPrefix = "product:"
	ProductListCacheKey = "user:list"
)

func GetUserCacheKey(userID uint) string {
	return fmt.Sprintf("%s%d", UserCacheKeyPrefix, userID)
}

func GetUserListCacheKey(limit, offset int) string {
	return fmt.Sprintf("%s:%d:%d", UserListCacheKey, limit, offset)
}

func GetProductCacheKey(productID uint) string {
	return fmt.Sprintf("%s%d",ProductCacheKeyPrefix,productID)
}

func GetProductListCacheKey(limit,offset int) string{
	return fmt.Sprintf("%s:%d:%d",ProductListCacheKey,limit,offset)
}
