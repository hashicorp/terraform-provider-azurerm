// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/key-vault"
}

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
		"azurerm_key_vault_certificates":                     dataSourceKeyVaultCertificates(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_key_vault_access_policy":                                resourceKeyVaultAccessPolicy(),
		"azurerm_key_vault_certificate":                                  resourceKeyVaultCertificate(),
		"azurerm_key_vault_certificate_issuer":                           resourceKeyVaultCertificateIssuer(),
		"azurerm_key_vault_key":                                          resourceKeyVaultKey(),
		"azurerm_key_vault_managed_hardware_security_module":             resourceKeyVaultManagedHardwareSecurityModule(),
		"azurerm_key_vault_secret":                                       resourceKeyVaultSecret(),
		"azurerm_key_vault":                                              resourceKeyVault(),
		"azurerm_key_vault_managed_storage_account":                      resourceKeyVaultManagedStorageAccount(),
		"azurerm_key_vault_managed_storage_account_sas_token_definition": resourceKeyVaultManagedStorageAccountSasTokenDefinition(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		EncryptedValueDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		KeyVaultCertificateContactsResource{},
	}
}
