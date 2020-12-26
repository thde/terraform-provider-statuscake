package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/thde/terraform-provider-statuscake/statuscake"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: statuscake.Provider,
	})
}
