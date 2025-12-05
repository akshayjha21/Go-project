package Config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string
}
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}

// NOTE - MustLoad loads the configuration from a file.
// It first tries to get the path from the environment variable CONFIG_PATH.
// If not set, it looks for a command line flag -config.
// If the config path is missing or the file does not exist, it exits the program.
func MustLoad() *Config {
	var configPath string

	// Check if CONFIG_PATH environment variable is set
	configPath = os.Getenv("CONFIG_PATH")

	// If not set, try to read the -config flag
	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()
		configPath = *flags

		// Exit if config path is still empty
		if configPath == "" {
			log.Fatalf("Config path is not set")
		}
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", err.Error())
	}

	// Load the config file into the Config struct
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	// Return the loaded config
	return &cfg
}
