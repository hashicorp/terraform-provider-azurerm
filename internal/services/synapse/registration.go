// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

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
		"azurerm_synapse_workspace": dataSourceSynapseWorkspace(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_synapse_firewall_rule":                              resourceSynapseFirewallRule(),
		"azurerm_synapse_integration_runtime_azure":                  resourceSynapseIntegrationRuntimeAzure(),
		"azurerm_synapse_integration_runtime_self_hosted":            resourceSynapseIntegrationRuntimeSelfHosted(),
		"azurerm_synapse_linked_service":                             resourceSynapseLinkedService(),
		"azurerm_synapse_managed_private_endpoint":                   resourceSynapseManagedPrivateEndpoint(),
		"azurerm_synapse_private_link_hub":                           resourceSynapsePrivateLinkHub(),
		"azurerm_synapse_role_assignment":                            resourceSynapseRoleAssignment(),
		"azurerm_synapse_spark_pool":                                 resourceSynapseSparkPool(),
		"azurerm_synapse_sql_pool":                                   resourceSynapseSqlPool(),
		"azurerm_synapse_sql_pool_extended_auditing_policy":          resourceSynapseSqlPoolExtendedAuditingPolicy(),
		"azurerm_synapse_sql_pool_security_alert_policy":             resourceSynapseSqlPoolSecurityAlertPolicy(),
		"azurerm_synapse_sql_pool_vulnerability_assessment":          resourceSynapseSqlPoolVulnerabilityAssessment(),
		"azurerm_synapse_sql_pool_vulnerability_assessment_baseline": resourceSynapseSqlPoolVulnerabilityAssessmentBaseline(),
		"azurerm_synapse_sql_pool_workload_classifier":               resourceSynapseSQLPoolWorkloadClassifier(),
		"azurerm_synapse_sql_pool_workload_group":                    resourceSynapseSQLPoolWorkloadGroup(),
		"azurerm_synapse_workspace":                                  resourceSynapseWorkspace(),
		"azurerm_synapse_workspace_aad_admin":                        resourceSynapseWorkspaceAADAdmin(),
		"azurerm_synapse_workspace_extended_auditing_policy":         resourceSynapseWorkspaceExtendedAuditingPolicy(),
		"azurerm_synapse_workspace_key":                              resourceSynapseWorkspaceKey(),
		"azurerm_synapse_workspace_security_alert_policy":            resourceSynapseWorkspaceSecurityAlertPolicy(),
		"azurerm_synapse_workspace_sql_aad_admin":                    resourceSynapseWorkspaceSqlAADAdmin(),
		"azurerm_synapse_workspace_vulnerability_assessment":         resourceSynapseWorkspaceVulnerabilityAssessment(),
	}
}
