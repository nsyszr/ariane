package config

type Config struct {
	NATSServerURL string `mapstructure:"NATS_SERVER_URL" yaml:"-"`

	BuildVersion string `yaml:"-"`
	BuildHash    string `yaml:"-"`
	BuildTime    string `yaml:"-"`
}
