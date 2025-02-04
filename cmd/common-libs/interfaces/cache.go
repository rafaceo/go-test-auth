package interfaces

import "time"

type CacheInterface interface {
	Set(key string, value interface{}, expiration time.Duration) (err error)
	Get(key string) (res string, err error)
	Exists(key string) (ok bool, err error)
	Delete(key string) (int64, error)
}
