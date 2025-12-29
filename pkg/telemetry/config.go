package telemetry

type Config struct {
	Enabled      bool   `mapstructure:"enabled" default:"false"`
	ServiceName  string `mapstructure:"service_name" default:"project_template"`
	OTLPEndpoint string `mapstructure:"otlp_endpoint" default:"localhost:4317"`
	Insecure     bool   `mapstructure:"insecure" default:"true"`
}
