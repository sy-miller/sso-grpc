package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	CONFIG_PATH_ENV_KEY = "CONFIG_PATH"
	CONFIG_PATH_FLAG    = "config"

	// Log levels
	INFO  = "INFO"
	DEBUG = "DEBUG"
	ERROR = "ERROR"
	WARN  = "WARN"
)

type Config struct {
	Env         string     `json:"env"`
	LogLevel    string     `json:"logLevel"`
	StoragePath *string    `json:"storagePath"`
	TokenTTL    Duration   `json:"tokenTtl"`
	GRPC        GRPCConfig `json:"grpc"`
}

type GRPCConfig struct {
	Port    *int     `json:"port"`
	Timeout Duration `json:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	cfgFile, err := os.Open(path)
	if err != nil {
		panic("could not open config file: " + err.Error())
	}
	defer cfgFile.Close()

	cfgBytes, err := io.ReadAll(cfgFile)
	if err != nil {
		panic("could not read config file: " + err.Error())
	}

	// Iniitialise with default values
	cfg := Config{
		Env:      "local",
		LogLevel: "INFO",
		TokenTTL: Duration(time.Duration(1 * time.Hour)),
		GRPC: GRPCConfig{
			Timeout: Duration(time.Duration(5 * time.Second)),
		},
	}

	if err = json.Unmarshal(cfgBytes, &cfg); err != nil {
		panic("could not parse config data: " + err.Error())
	}

	// Validate required fields
	if cfg.StoragePath == nil {
		panic("storagePath is required")
	}
	if cfg.GRPC.Port == nil {
		panic("grpc.port is required")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, CONFIG_PATH_FLAG, "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv(CONFIG_PATH_ENV_KEY)
	}

	return res
}

// Duration is a custom type that embeds time.Duration
type Duration time.Duration

func (d *Duration) ToTimeDuration() time.Duration {
	return time.Duration(*d)
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return fmt.Errorf("invalid duration type %T, value: %s", v, b)
	}
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}
