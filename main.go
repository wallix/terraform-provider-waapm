package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/wallix/terraform-provider-waapm/waapm"
)

func main() {
        plugin.Serve(&plugin.ServeOpts{
                ProviderFunc: waapm.Provider,
        })
}
