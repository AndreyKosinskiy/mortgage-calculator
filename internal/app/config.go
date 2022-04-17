package app

type Config struct {
	Port       string
	DebugLevel int
	DbURL      string
}

func NewConfig() *Config {
	return &Config{":8080", 0, ""}
}
