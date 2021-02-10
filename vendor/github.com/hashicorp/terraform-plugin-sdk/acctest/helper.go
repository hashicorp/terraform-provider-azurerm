package acctest

import (
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	tftest "github.com/hashicorp/terraform-plugin-test/v2"
)

var TestHelper *tftest.Helper

func UseBinaryDriver(name string, providerFunc plugin.ProviderFunc) {
	log.Println("[DEBUG] not using binary driver name, it's no longer needed")
	sourceDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if tftest.RunningAsPlugin() {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: providerFunc,
		})
		os.Exit(0)
	} else {
		TestHelper = tftest.AutoInitProviderHelper(sourceDir)
	}
}
