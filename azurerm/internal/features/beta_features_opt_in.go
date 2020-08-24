package features

import (
	"os"
)

// VMSSExtensionsBeta returns whether or not the beta for VMSS Extensions for Linux and Windows VMSS resources is
// enabled.
//
// Set to any non-empty value to enable
func VMSSExtensionsBeta() bool {
	value := os.Getenv("ARM_PROVIDER_VMSS_EXTENSIONS_BETA")
	if value == "" {
		return false
	} else {
		return true
	}
}
