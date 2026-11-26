// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package keyvault

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_access_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_certificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_certificate_contacts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_certificate_data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_certificate_issuer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_encrypted_value"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_managed_storage_account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_managed_storage_account_sas_token_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_secret"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/resources/key_vault_secrets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.FrameworkServiceRegistration               = Registration{}
)

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
		"azurerm_key_vault":                    key_vault.RegisterDataSource(),
		"azurerm_key_vault_access_policy":      key_vault_access_policy.RegisterDataSource(),
		"azurerm_key_vault_certificate":        key_vault_certificate.RegisterDataSource(),
		"azurerm_key_vault_certificate_data":   key_vault_certificate_data.RegisterDataSource(),
		"azurerm_key_vault_certificate_issuer": key_vault_certificate_issuer.RegisterDataSource(),
		"azurerm_key_vault_certificates":       key_vault_certificates.RegisterDataSource(),
		"azurerm_key_vault_key":                key_vault_key.RegisterDataSource(),
		"azurerm_key_vault_secret":             key_vault_secret.RegisterDataSource(),
		"azurerm_key_vault_secrets":            key_vault_secrets.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_key_vault":                                              key_vault.RegisterResource(),
		"azurerm_key_vault_access_policy":                                key_vault_access_policy.RegisterResource(),
		"azurerm_key_vault_certificate":                                  key_vault_certificate.RegisterResource(),
		"azurerm_key_vault_certificate_issuer":                           key_vault_certificate_issuer.RegisterResource(),
		"azurerm_key_vault_key":                                          key_vault_key.RegisterResource(),
		"azurerm_key_vault_managed_storage_account":                      key_vault_managed_storage_account.RegisterResource(),
		"azurerm_key_vault_managed_storage_account_sas_token_definition": key_vault_managed_storage_account_sas_token_definition.RegisterResource(),
		"azurerm_key_vault_secret":                                       key_vault_secret.RegisterResource(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		key_vault_encrypted_value.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		key_vault_certificate_contacts.Resource{},
	}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		key_vault_certificate.NewKeyVaultCertificateEphemeralResource,
		key_vault_secret.NewKeyVaultSecretEphemeralResource,
	}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		key_vault.ListResource{},
	}
}
