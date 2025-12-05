package config

type HTTPConfig struct {
	Port string `mapstructure:"port" default:"3000"`
}
