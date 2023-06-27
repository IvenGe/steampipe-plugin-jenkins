package main

import (
	"github.com/turbot/steampipe-plugin-jenkins/jenkins"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: jenkins.Plugin})
}
