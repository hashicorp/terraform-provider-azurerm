package features

import (
	"os"
	"strings"
)

// NOTE: we'll need to add an infobox to MySQL|PostgreSQL Configuration when this goes live
// since these resources can't support import
// in addition the virtual resources will need adjusting
func ShouldResourcesBeImported() bool {
	return strings.EqualFold(os.Getenv("ARM_PROVIDER_STRICT"), "true")
}
