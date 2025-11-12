package settings

import "github.com/spf13/viper"

type HTTPServerSettings struct {
	Port              string `mapstructure:"port"`
	ReadTimeout       uint   `mapstructure:"read_timeout"`
	WriteTimeout      uint   `mapstructure:"write_timeout"`
	IdleTimeout       uint   `mapstructure:"idle_timeout"`
	ReadHeaderTimeout uint   `mapstructure:"read_header_timeout"`
}

func SetHTTPServerDefaults(v *viper.Viper, prefix string, defaultPort string) {
	v.SetDefault(prefix+".port", defaultPort)
	v.SetDefault(prefix+".read_timeout", 10)
	v.SetDefault(prefix+".write_timeout", 10)
	v.SetDefault(prefix+".idle_timeout", 300)
	v.SetDefault(prefix+".read_header_timeout", 5)
}
