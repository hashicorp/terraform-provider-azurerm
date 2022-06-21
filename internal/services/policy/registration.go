package policy

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
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
		"azurerm_policy_definition":                               dataSourceArmPolicyDefinition(),
		"azurerm_policy_set_definition":                           dataSourceArmPolicySetDefinition(),
		"azurerm_policy_virtual_machine_configuration_assignment": dataSourcePolicyVirtualMachineConfigurationAssignment(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_policy_definition":                               resourceArmPolicyDefinition(),
		"azurerm_policy_set_definition":                           resourceArmPolicySetDefinition(),
		"azurerm_management_group_policy_remediation":             resourceArmManagementGroupPolicyRemediation(),
		"azurerm_resource_policy_remediation":                     resourceArmResourcePolicyRemediation(),
		"azurerm_management_group_policy_exemption":               resourceArmManagementGroupPolicyExemption(),
		"azurerm_resource_policy_exemption":                       resourceArmResourcePolicyExemption(),
		"azurerm_resource_group_policy_exemption":                 resourceArmResourceGroupPolicyExemption(),
		"azurerm_subscription_policy_exemption":                   resourceArmSubscriptionPolicyExemption(),
		"azurerm_resource_group_policy_remediation":               resourceArmResourceGroupPolicyRemediation(),
		"azurerm_subscription_policy_remediation":                 resourceArmSubscriptionPolicyRemediation(),
		"azurerm_policy_virtual_machine_configuration_assignment": resourcePolicyVirtualMachineConfigurationAssignment(),
	}
}
