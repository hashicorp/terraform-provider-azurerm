// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.UntypedServiceRegistration               = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/policy"
}

type Registration struct{}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		AssignmentDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ManagementGroupAssignmentResource{},
		ManagementGroupPolicySetDefinitionResource{},
		PolicySetDefinitionResource{},
		ResourceAssignmentResource{},
		ResourceGroupAssignmentResource{},
		SubscriptionAssignmentResource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Policy"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Policy",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_policy_definition":                               dataSourcePolicyDefinition(),
		"azurerm_policy_definition_built_in":                      dataSourcePolicyDefinitionBuiltIn(),
		"azurerm_policy_set_definition":                           dataSourcePolicySetDefinition(),
		"azurerm_policy_virtual_machine_configuration_assignment": dataSourcePolicyVirtualMachineConfigurationAssignment(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_management_group_policy_exemption":               resourceManagementGroupPolicyExemption(),
		"azurerm_management_group_policy_remediation":             resourceManagementGroupPolicyRemediation(),
		"azurerm_policy_definition":                               resourcePolicyDefinition(),
		"azurerm_policy_virtual_machine_configuration_assignment": resourcePolicyVirtualMachineConfigurationAssignment(),
		"azurerm_resource_group_policy_exemption":                 resourceResourceGroupPolicyExemption(),
		"azurerm_resource_group_policy_remediation":               resourceResourceGroupPolicyRemediation(),
		"azurerm_resource_policy_exemption":                       resourceResourcePolicyExemption(),
		"azurerm_resource_policy_remediation":                     resourceResourcePolicyRemediation(),
		"azurerm_subscription_policy_exemption":                   resourceSubscriptionPolicyExemption(),
		"azurerm_subscription_policy_remediation":                 resourceSubscriptionPolicyRemediation(),
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
