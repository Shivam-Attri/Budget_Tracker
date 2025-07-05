// config/config.go
// Handles loading configuration from files and environment variables using Viper.
package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB          DatabaseConfig    `mapstructure:"database"`
	Server      ServerConfig      `mapstructure:"server"`
	RateLimiter RateLimiterConfig `mapstructure:"rate_limiter"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"`
}

type ServerConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	JWTSecret         string        `mapstructure:"jwt_secret"`
	EncryptionKey     string        `mapstructure:"encryption_key"`
	DefaultPageLimit  int           `mapstructure:"default_page_limit"`
	AccessTokenTTL    int           `mapstructure:"access_token_ttl_minutes"`
	RefreshTokenTTL   int           `mapstructure:"refresh_token_ttl_hours"`
	AllowedSortFields []string      `mapstructure:"allowed_sort_fields"`
	Timeout           TimeoutConfig `mapstructure:"timeout"`
}

type TimeoutConfig struct {
	Read             int `mapstructure:"read"`
	Write            int `mapstructure:"write"`
	Idle             int `mapstructure:"idle"`
	GracefulShutdown int `mapstructure:"graceful_shutdown"`
}

type RateLimiterConfig struct {
	Enabled bool    `mapstructure:"enabled"`
	RPS     float64 `mapstructure:"rps"`
	Burst   int     `mapstructure:"burst"`
}

func LoadConfig(path string) (config Config, err error) {
	// Set defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8000)
	viper.SetDefault("server.default_page_limit", 10)
	viper.SetDefault("server.access_token_ttl_minutes", 15)
	viper.SetDefault("server.refresh_token_ttl_hours", 168)
	viper.SetDefault("server.allowed_sort_fields", []string{"date", "amount"})
	viper.SetDefault("server.timeout.read", 15)
	viper.SetDefault("server.timeout.write", 15)
	viper.SetDefault("server.timeout.idle", 60)
	viper.SetDefault("server.timeout.graceful_shutdown", 5)
	viper.SetDefault("rate_limiter.enabled", true)
	viper.SetDefault("rate_limiter.rps", 2)
	viper.SetDefault("rate_limiter.burst", 4)

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("BUDGET")
	viper.BindEnv("database.url", "DB_URL")
	viper.BindEnv("server.jwt_secret", "JWT_SECRET")
	viper.BindEnv("server.encryption_key", "ENCRYPTION_KEY")

	err = viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); err != nil && !ok {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
