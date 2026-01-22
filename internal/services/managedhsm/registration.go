// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/managed-hsm"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Managed HSM"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		// Managed HSM is grouped under Key Vault
		"Key Vault",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_key_vault_managed_hardware_security_module": dataSourceKeyVaultManagedHardwareSecurityModule(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_key_vault_managed_hardware_security_module": resourceKeyVaultManagedHardwareSecurityModule(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		KeyvaultMHSMRoleDefinitionDataSource{},
		KeyvaultMHSMKeyDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		KeyVaultMHSMKeyResource{},
		KeyVaultMHSMRoleDefinitionResource{},
		KeyVaultManagedHSMRoleAssignmentResource{},
		KeyVaultMHSMKeyRotationPolicyResource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
