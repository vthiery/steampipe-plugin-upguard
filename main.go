package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/vthiery/steampipe-plugin-upguard/upguard"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: upguard.Plugin})
}
