package messagebus

// Config holds the configuration for the message bus.
type Config struct {
	// Backend specifies which message bus backend to use.
	// Supported values: "gochannel", "redis", "kafka"
	Backend string `default:"gochannel"`
}
