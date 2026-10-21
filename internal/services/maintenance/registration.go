// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package maintenance

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/resources/maintenance_assignment_dedicated_host"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/resources/maintenance_assignment_dynamic_scope"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/resources/maintenance_assignment_virtual_machine"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/resources/maintenance_assignment_virtual_machine_scale_set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/resources/maintenance_configuration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/resources/public_maintenance_configurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/maintenance"
}

func (r Registration) Name() string {
	return "Maintenance"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Maintenance",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		maintenance_assignment_dynamic_scope.Resource{},
	}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_maintenance_configuration":         maintenance_configuration.RegisterDataSource(),
		"azurerm_public_maintenance_configurations": public_maintenance_configurations.RegisterDataSource(),
	}
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_maintenance_assignment_dedicated_host":            maintenance_assignment_dedicated_host.RegisterResource(),
		"azurerm_maintenance_assignment_virtual_machine":           maintenance_assignment_virtual_machine.RegisterResource(),
		"azurerm_maintenance_assignment_virtual_machine_scale_set": maintenance_assignment_virtual_machine_scale_set.RegisterResource(),
		"azurerm_maintenance_configuration":                        maintenance_configuration.RegisterResource(),
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
