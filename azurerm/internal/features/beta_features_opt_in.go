package features

import (
	"os"
	"strings"
)

// VMDataDiskBeta returns whether or not the beta for VM Data Disks for Linux and Windows VM resources is
// enabled.
//
// Set the Environment Variable `ARM_PROVIDER_VM_DATADISKS_BETA` to `true`
func VMDataDiskBeta() bool {
	return strings.EqualFold(os.Getenv("ARM_PROVIDER_VM_DATADISKS_BETA"), "true")
}
