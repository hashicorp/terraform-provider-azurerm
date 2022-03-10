package policy

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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
	return []sdk.DataSource{}
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
	resources := map[string]*pluginsdk.Resource{
		"azurerm_policy_definition":                               resourceArmPolicyDefinition(),
		"azurerm_policy_set_definition":                           resourceArmPolicySetDefinition(),
		"azurerm_policy_remediation":                              resourceArmPolicyRemediation(),
		"azurerm_policy_virtual_machine_configuration_assignment": resourcePolicyVirtualMachineConfigurationAssignment(),
	}

	if !features.ThreePointOhBeta() {
		resources["azurerm_policy_assignment"] = resourceArmPolicyAssignment()
		resources["azurerm_virtual_machine_configuration_policy_assignment"] = resourceVirtualMachineConfigurationPolicyAssignment()
	}

	return resources
}
