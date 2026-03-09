package upguard

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// upguardConfig stores the connection configuration for the plugin.
type upguardConfig struct {
	APIKey *string `hcl:"api_key"`
}

// ConfigInstance returns a new instance of upguardConfig (used by the SDK).
func ConfigInstance() interface{} {
	return &upguardConfig{}
}

// GetConfig retrieves and casts the connection config from the plugin query data.
func GetConfig(connection *plugin.Connection) upguardConfig {
	if connection == nil || connection.Config == nil {
		return upguardConfig{}
	}
	config, _ := connection.Config.(upguardConfig)
	return config
}
