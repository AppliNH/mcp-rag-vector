package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTPConfig `mapstructure:"http"`
	LogLevel slog.Level `mapstructure:"log_level"`
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

		// Search config in home directory with name ".goa-boilerplate" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".goa-boilerplate")
	}

	// Set default values
	viper.SetDefault("log_level", "INFO")
	viper.SetDefault("http.port", "3000")

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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

	return cfg, nil
}
