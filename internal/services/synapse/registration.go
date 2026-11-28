// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_firewall_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_integration_runtime_azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_integration_runtime_self_hosted"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_linked_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_managed_private_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_private_link_hub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_role_assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_spark_pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_sql_pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_sql_pool_extended_auditing_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_sql_pool_security_alert_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_sql_pool_vulnerability_assessment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_sql_pool_vulnerability_assessment_baseline"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_sql_pool_workload_classifier"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_sql_pool_workload_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_workspace_aad_admin"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_workspace_extended_auditing_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_workspace_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_workspace_security_alert_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_workspace_sql_aad_admin"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/resources/synapse_workspace_vulnerability_assessment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/synapse"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Synapse"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Synapse",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_synapse_workspace": synapse_workspace.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_synapse_firewall_rule":                              synapse_firewall_rule.RegisterResource(),
		"azurerm_synapse_integration_runtime_azure":                  synapse_integration_runtime_azure.RegisterResource(),
		"azurerm_synapse_integration_runtime_self_hosted":            synapse_integration_runtime_self_hosted.RegisterResource(),
		"azurerm_synapse_linked_service":                             synapse_linked_service.RegisterResource(),
		"azurerm_synapse_managed_private_endpoint":                   synapse_managed_private_endpoint.RegisterResource(),
		"azurerm_synapse_private_link_hub":                           synapse_private_link_hub.RegisterResource(),
		"azurerm_synapse_role_assignment":                            synapse_role_assignment.RegisterResource(),
		"azurerm_synapse_spark_pool":                                 synapse_spark_pool.RegisterResource(),
		"azurerm_synapse_sql_pool":                                   synapse_sql_pool.RegisterResource(),
		"azurerm_synapse_sql_pool_extended_auditing_policy":          synapse_sql_pool_extended_auditing_policy.RegisterResource(),
		"azurerm_synapse_sql_pool_security_alert_policy":             synapse_sql_pool_security_alert_policy.RegisterResource(),
		"azurerm_synapse_sql_pool_vulnerability_assessment":          synapse_sql_pool_vulnerability_assessment.RegisterResource(),
		"azurerm_synapse_sql_pool_vulnerability_assessment_baseline": synapse_sql_pool_vulnerability_assessment_baseline.RegisterResource(),
		"azurerm_synapse_sql_pool_workload_classifier":               synapse_sql_pool_workload_classifier.RegisterResource(),
		"azurerm_synapse_sql_pool_workload_group":                    synapse_sql_pool_workload_group.RegisterResource(),
		"azurerm_synapse_workspace":                                  synapse_workspace.RegisterResource(),
		"azurerm_synapse_workspace_aad_admin":                        synapse_workspace_aad_admin.RegisterResource(),
		"azurerm_synapse_workspace_extended_auditing_policy":         synapse_workspace_extended_auditing_policy.RegisterResource(),
		"azurerm_synapse_workspace_key":                              synapse_workspace_key.RegisterResource(),
		"azurerm_synapse_workspace_security_alert_policy":            synapse_workspace_security_alert_policy.RegisterResource(),
		"azurerm_synapse_workspace_sql_aad_admin":                    synapse_workspace_sql_aad_admin.RegisterResource(),
		"azurerm_synapse_workspace_vulnerability_assessment":         synapse_workspace_vulnerability_assessment.RegisterResource(),
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
