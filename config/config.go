package config

import "sync/atomic"

type Config struct {
	DSN                 string
	DriverName          string
	TableName           string
	StructName          string
	FilePath            string
	PackageName         string
	JsonTag             bool
	DateString          bool
	SelectKey           string
	CustomizeCodeBefore string
	CustomizeCodeIn     string
	CustomizeCodeAfter  string
	FirstToUpper        bool
	GenerateFunc        bool
}

var defaultConfig = Config{}

var (
	globalConf atomic.Value
)

// InitializeConfig initialize the global config handler.
// TODO: read config from file.
func InitializeConfig(enforceCmdArgs func(*Config)) {
	cfg := defaultConfig
	enforceCmdArgs(&cfg)
	StoreGlobalConfig(&cfg)
}

// GetGlobalConfig returns the global configuration for this server.
// It should store configuration from command line and configuration file.
// Other parts of the system can read the global configuration use this function.
func GetGlobalConfig() *Config {
	return globalConf.Load().(*Config)
}

// StoreGlobalConfig stores a new config to the globalConf. It mostly uses in the test to avoid some data races.
func StoreGlobalConfig(config *Config) {
	globalConf.Store(config)
}
