package migrations

type Config struct {
	RunOnStartup bool `mapstructure:"run_on_startup" default:"true"`
}
