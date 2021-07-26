package azuread

import (
	"os"
	"strings"
)

// This file contains feature flags for functionality which will prove more challenging to implement en-mass
var requireResourcesToBeImported = strings.EqualFold(os.Getenv("ARM_PROVIDER_STRICT"), "true")
