// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/management_group_policy_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/management_group_policy_exemption"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/management_group_policy_remediation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/management_group_policy_set_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/policy_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/policy_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/policy_definition_built_in"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/policy_set_definition"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/policy_virtual_machine_configuration_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/resource_group_policy_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/resource_group_policy_exemption"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/resource_group_policy_remediation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/resource_policy_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/resource_policy_exemption"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/resource_policy_remediation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/subscription_policy_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/subscription_policy_exemption"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/resources/subscription_policy_remediation"
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
		policy_assignment.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		management_group_policy_assignment.Resource{},
		management_group_policy_set_definition.Resource{},
		policy_set_definition.Resource{},
		resource_policy_assignment.Resource{},
		resource_group_policy_assignment.Resource{},
		subscription_policy_assignment.Resource{},
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
		"azurerm_policy_definition":                               policy_definition.RegisterDataSource(),
		"azurerm_policy_definition_built_in":                      policy_definition_built_in.RegisterDataSource(),
		"azurerm_policy_set_definition":                           policy_set_definition.RegisterDataSource(),
		"azurerm_policy_virtual_machine_configuration_assignment": policy_virtual_machine_configuration_assignment.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_management_group_policy_exemption":               management_group_policy_exemption.RegisterResource(),
		"azurerm_management_group_policy_remediation":             management_group_policy_remediation.RegisterResource(),
		"azurerm_policy_definition":                               policy_definition.RegisterResource(),
		"azurerm_policy_virtual_machine_configuration_assignment": policy_virtual_machine_configuration_assignment.RegisterResource(),
		"azurerm_resource_group_policy_exemption":                 resource_group_policy_exemption.RegisterResource(),
		"azurerm_resource_group_policy_remediation":               resource_group_policy_remediation.RegisterResource(),
		"azurerm_resource_policy_exemption":                       resource_policy_exemption.RegisterResource(),
		"azurerm_resource_policy_remediation":                     resource_policy_remediation.RegisterResource(),
		"azurerm_subscription_policy_exemption":                   subscription_policy_exemption.RegisterResource(),
		"azurerm_subscription_policy_remediation":                 subscription_policy_remediation.RegisterResource(),
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
