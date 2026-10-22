// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managedhsm

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/resources/key_vault_managed_hardware_security_module"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/resources/key_vault_managed_hardware_security_module_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/resources/key_vault_managed_hardware_security_module_key_rotation_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/resources/key_vault_managed_hardware_security_module_role_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/resources/key_vault_managed_hardware_security_module_role_definition"
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
		"azurerm_key_vault_managed_hardware_security_module": key_vault_managed_hardware_security_module.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_key_vault_managed_hardware_security_module": key_vault_managed_hardware_security_module.RegisterResource(),
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		key_vault_managed_hardware_security_module_key.DataSource{},
		key_vault_managed_hardware_security_module_role_definition.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		key_vault_managed_hardware_security_module_role_assignment.Resource{},
		key_vault_managed_hardware_security_module_key.Resource{},
		key_vault_managed_hardware_security_module_key_rotation_policy.Resource{},
		key_vault_managed_hardware_security_module_role_definition.Resource{},
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
