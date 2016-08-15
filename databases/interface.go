package databases

type Config struct {
	Address        string
	Bucket         string
	BucketPassword string
}

type Database interface {
	Init(config Config)

	Shutdown()

	Create(key string, value map[string]interface{}) error

	Read(key string) error

	Update(key string, value map[string]interface{}) error

	Delete(key string) error
}
