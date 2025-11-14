package config

import (
	"strings"

	"github.com/ZaiiiRan/job_search_service/auth-service/internal/config/settings"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	GRPCServer            settings.GRPCServerSettings `mapstructure:"grpc_server"`
	HTTPGatewayServer     settings.HTTPServerSettings `mapstructure:"http_gateway_server"`
	UserServiceGRPCClient settings.GRPCClientSettings `mapstructure:"user_service_grpc_client"`
	DB                    settings.PostgresSettings   `mapstructure:"db"`
	Migrate               settings.MigrateSettings    `mapstructure:"migrate"`
	Redis                 settings.RedisSettings      `mapstructure:"redis"`
	Shutdown              settings.ShutdownSettings   `mapstructure:"shutdown"`
}

func LoadServerConfig() (*ServerConfig, error) {
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/auth-service")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setServerDefaults(v)

	var cfg ServerConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setServerDefaults(v *viper.Viper) {
	settings.SetGRPCServerDefaults(v, "grpc_server", ":50052")
	settings.SetHTTPServerDefaults(v, "http_gateway_server", ":8082")
	settings.SetGRPCClientDefaults(v, "user_service_grpc_client", "localhost:50051")
	settings.SetPostgresDefaults(v, "db")
	settings.SetMigrateDefaults(v, "migrate")
	settings.SetRedisDefaults(v, "redis")
	settings.SetShutdownDefaults(v, "shutdown")
}
