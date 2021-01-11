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
		"azurerm_key_vault_access_policy":      dataSourceKeyVaultAccessPolicy(),
		"azurerm_key_vault_certificate":        dataSourceKeyVaultCertificate(),
		"azurerm_key_vault_certificate_issuer": dataSourceKeyVaultCertificateIssuer(),
		"azurerm_key_vault_key":                dataSourceKeyVaultKey(),
		"azurerm_key_vault_secret":             dataSourceKeyVaultSecret(),
		"azurerm_key_vault":                    dataSourceKeyVault(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_key_vault_access_policy":      resourceKeyVaultAccessPolicy(),
		"azurerm_key_vault_certificate":        resourceKeyVaultCertificate(),
		"azurerm_key_vault_certificate_issuer": resourceKeyVaultCertificateIssuer(),
		"azurerm_key_vault_key":                resourceKeyVaultKey(),
		"azurerm_key_vault_secret":             resourceKeyVaultSecret(),
		"azurerm_key_vault":                    resourceKeyVault(),
	}
}
