package env

import "github.com/caarlos0/env/v10"

type config struct {
	MongoURL           string `env:"MONGO_URL" envDefault:"mongodb://localhost:27017"`
	MongodPath         string `env:"MONGOD_PATH" envDefault:"mongod"`
	UseMemoryMongo     bool   `env:"USE_MEMORY_MONGO" envDefault:"true"`
	MaxReturnArraySize int    `env:"MAX_RETURN_ARRAY_SIZE" envDefault:"100"`
	Port               int    `env:"PORT"  envDefault:"5000"`
}

func Load() (*config, error) {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
