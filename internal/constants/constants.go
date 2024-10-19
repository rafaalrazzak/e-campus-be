package constants

import "time"

type RedisKeys struct {
	SessionKey   string
	UserCacheKey string
	ProductKey   string
}

type AppConstants struct {
	SessionExpiration time.Duration
}

var App = AppConstants{
	SessionExpiration: 24 * time.Hour,
}

var Redis = RedisKeys{
	SessionKey:   "ecampus:session::%d::%d",  // session:userId:sessionToken
	UserCacheKey: "ecampus:cache:user:%s",    // cache:user:userId
	ProductKey:   "ecampus:cache:product:%s", // cache:product:productId
}
