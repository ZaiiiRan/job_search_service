package settings

import "github.com/spf13/viper"

type PostgresSettings struct {
	ConnectionString          string `mapstructure:"connection_string"`
	MigrationConnectionString string `mapstructure:"migration_connection_string"`
}

func SetPostgresDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".connection_string", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	v.SetDefault(prefix+".migration_connection_string", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
}
