package keyvault

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "KeyVault"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Key Vault",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_key_vault_access_policy":      dataSourceArmKeyVaultAccessPolicy(),
		"azurerm_key_vault_certificate":        dataSourceArmKeyVaultCertificate(),
		"azurerm_key_vault_certificate_issuer": dataSourceArmKeyVaultCertificateIssuer(),
		"azurerm_key_vault_key":                dataSourceArmKeyVaultKey(),
		"azurerm_key_vault_secret":             dataSourceArmKeyVaultSecret(),
		"azurerm_key_vault":                    dataSourceArmKeyVault(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_key_vault_access_policy":      resourceArmKeyVaultAccessPolicy(),
		"azurerm_key_vault_certificate":        resourceArmKeyVaultCertificate(),
		"azurerm_key_vault_certificate_issuer": resourceArmKeyVaultCertificateIssuer(),
		"azurerm_key_vault_key":                resourceArmKeyVaultKey(),
		"azurerm_key_vault_secret":             resourceArmKeyVaultSecret(),
		"azurerm_key_vault":                    resourceArmKeyVault(),
	}
}
