package cache

type ICache interface {
	GetALL() (map[string]interface{}, error)
	Get(key string) (interface{}, error)
	GetObject(key string, data interface{}) error
	Set(key string, value interface{}) error
	Delete(key string) error
	DeleteExpired() error
	GC()
}
