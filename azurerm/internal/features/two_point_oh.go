package features

// SupportsTwoPointZeroResources returns whether the new VM and VMSS resources from 2.0
// should be supported
//
// There's 5 new resources coming as a part of 2.0, which are intentionally feature-flagged off
// until all 5 are supported:
//  * `azurerm_linux_virtual_machine`
//  * `azurerm_linux_virtual_machine_scale_set`
//  * `azurerm_windows_virtual_machine`
//  * `azurerm_windows_virtual_machine_scale_set`
//  * `azurerm_virtual_machine_scale_set_extension`
//
// This feature-toggle defaults to off in 1.x versions of the Azure Provider, however this will
// become enabled by default in version 2.0 of the Azure Provider (where this toggle will be removed).
// As outlined in the announcement for v2.0 of the Azure Provider:
// https://github.com/terraform-providers/terraform-provider-azurerm/issues/2807
//
// Operators wishing to beta-test these resources can opt-into them in 1.x versions of the
// Azure Provider by setting the Environment Variable 'ARM_PROVIDER_TWOPOINTZERO_RESOURCES' to 'true'
func SupportsTwoPointZeroResources() bool {
	return true
}
