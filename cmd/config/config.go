package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTPConfig   `mapstructure:"http"`
	LogLevel slog.Level   `mapstructure:"log_level" default:"INFO"`
	Ollama   OllamaConfig `mapstructure:"ollama"`
	Qdrant   QdrantConfig `mapstructure:"qdrant"`
}

type OllamaConfig struct {
	Host           string `mapstructure:"host" validate:"required,url"`
	EmbeddingModel string `mapstructure:"embedding_model" validate:"required"`
}

type QdrantConfig struct {
	Host string `mapstructure:"host" validate:"required"`
}

func ParseLogLevel(level string) (slog.Level, error) {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return slog.Level(0), fmt.Errorf("failed parsing log level from config %q: %w", level, err)
	}
	return l, nil
}

// decode hook for viper/mapstructure
func logLevelHook() mapstructure.DecodeHookFunc {
	return func(from, to reflect.Type, data any) (any, error) {
		if from.Kind() != reflect.String {
			return data, nil
		}
		if to != reflect.TypeOf(slog.Level(0)) {
			return data, nil
		}

		return ParseLogLevel(strings.ToLower(data.(string)))
	}
}

// initConfig reads in config file and ENV variables if set.
func Load(cfgFile string) (Config, error) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("mcp-rag-vector")
	}

	// Automatically register all config keys from struct tags (including defaults)
	registerConfigKeys()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg, func(c *mapstructure.DecoderConfig) {
		c.DecodeHook = logLevelHook()
	}); err != nil {
		log.Fatalf("unable to decode config into struct: %v", err)
	}

	// Validate using struct tags
	if err := validateConfig(&cfg); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	return cfg, nil
}

// registerConfigKeys uses reflection to automatically register all config keys
// from the Config struct's mapstructure tags, so viper knows to read them from env
func registerConfigKeys() {
	var cfg Config
	registerStructKeys(reflect.TypeOf(cfg), "")
}

func registerStructKeys(t reflect.Type, prefix string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get the mapstructure tag
		tag := field.Tag.Get("mapstructure")
		if tag == "" || tag == "-" {
			continue
		}

		// Build the full key path
		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}

		// If it's a struct, recurse
		if field.Type.Kind() == reflect.Struct {
			registerStructKeys(field.Type, key)
		} else {
			// Get default value from tag, or use empty string
			defaultVal := field.Tag.Get("default")
			viper.SetDefault(key, defaultVal)
		}
	}
}

func validateConfig(cfg *Config) error {
	validate := validator.New()

	// Register custom tag name for better error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("mapstructure")
		if name == "" {
			return fld.Name
		}
		return name
	})

	return validate.Struct(cfg)
}
