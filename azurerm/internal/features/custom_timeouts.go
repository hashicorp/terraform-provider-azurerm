package features

import (
	"os"
	"strings"
)

// SupportsCustomTimeouts returns whether Custom Timeouts are supported
//
// This feature allows Resources to define Custom Timeouts for Creation, Updating and Deletion
// which helps work with Azure resources that take longer to provision/delete.
// When this feature is disabled, all resources have a hard-coded timeout of 3 hours.
//
// This feature-toggle defaults to off in 1.x versions of the Azure Provider, however this will
// become the default behaviour in version 2.0 of the Azure Provider. As outlined in the announcement
// for v2.0 of the Azure Provider: https://github.com/terraform-providers/terraform-provider-azurerm/issues/2807
//
// Operators wishing to adopt this behaviour can opt-into this behaviour in 1.x versions of the
// Azure Provider by setting the Environment Variable 'ARM_PROVIDER_CUSTOM_TIMEOUTS' to 'true'
func SupportsCustomTimeouts() bool {
	return strings.EqualFold(os.Getenv("ARM_PROVIDER_CUSTOM_TIMEOUTS"), "true")
}
