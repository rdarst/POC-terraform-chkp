package main

import (
        "github.com/hashicorp/terraform/plugin"
      	"./chkp"
)

func main() {
        plugin.Serve(&plugin.ServeOpts{
                ProviderFunc: chkp.Provider,
                    })
}
