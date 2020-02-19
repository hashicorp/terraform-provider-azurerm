package features

// ShouldResourcesBeImported returns whether the feature Requiring Resources to be Imported
// should be enabled.
//
// This feature prevents Terraform from 'adopting' resources which already exist, which is the
// behaviour used by ARM Templates which will update these resources rather than overwriting them
// Instead existing resources will need to be imported using `terraform import`, as is the case
// with other Terraform Providers.
//
// This feature-toggle defaults to off in 1.x versions of the Azure Provider, however this will
// become the default behaviour in version 2.0 of the Azure Provider. As outlined in the announcement
// for v2.0 of the Azure Provider: https://github.com/terraform-providers/terraform-provider-azurerm/issues/2807
//
// Operators wishing to adopt this behaviour can opt-into this behaviour in 1.x versions of the
// Azure Provider by setting the Environment Variable 'ARM_PROVIDER_STRICT' to 'true'
func ShouldResourcesBeImported() bool {
	// NOTE: we'll need to add an infobox to the following resources when this goes live:
	// * App Service Source Control Token
	// * MySQL|PostgreSQL Configuration
	// since these resources can't support import
	// in addition the virtual resources will need adjusting
	return true
}
