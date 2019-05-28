package flags

import (
	"os"
	"strings"
)

// NOTE: we'll need to add an infobox to MySQL|PostgreSQL Configuration when this goes live
// since these resources can't support import
// in addition the virtual resources will need adjusting

// This file contains feature flags for functionality which will prove more challenging to implement en-mass
var RequireResourcesToBeImported = strings.EqualFold(os.Getenv("ARM_PROVIDER_STRICT"), "true")
