package keyvault

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_key_vault_access_policy":                    dataSourceKeyVaultAccessPolicy(),
		"azurerm_key_vault_certificate":                      dataSourceKeyVaultCertificate(),
		"azurerm_key_vault_certificate_data":                 dataSourceKeyVaultCertificateData(),
		"azurerm_key_vault_certificate_issuer":               dataSourceKeyVaultCertificateIssuer(),
		"azurerm_key_vault_key":                              dataSourceKeyVaultKey(),
		"azurerm_key_vault_managed_hardware_security_module": dataSourceKeyVaultManagedHardwareSecurityModule(),
		"azurerm_key_vault_secret":                           dataSourceKeyVaultSecret(),
		"azurerm_key_vault_secrets":                          dataSourceKeyVaultSecrets(),
		"azurerm_key_vault":                                  dataSourceKeyVault(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_key_vault_access_policy":                    resourceKeyVaultAccessPolicy(),
		"azurerm_key_vault_certificate":                      resourceKeyVaultCertificate(),
		"azurerm_key_vault_certificate_issuer":               resourceKeyVaultCertificateIssuer(),
		"azurerm_key_vault_key":                              resourceKeyVaultKey(),
		"azurerm_key_vault_managed_hardware_security_module": resourceKeyVaultManagedHardwareSecurityModule(),
		"azurerm_key_vault_secret":                           resourceKeyVaultSecret(),
		"azurerm_key_vault":                                  resourceKeyVault(),
	}
}
