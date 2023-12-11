// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration                   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/mssql"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Microsoft SQL Server / Azure SQL"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Database",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_mssql_database":    dataSourceMsSqlDatabase(),
		"azurerm_mssql_elasticpool": dataSourceMsSqlElasticpool(),
		"azurerm_mssql_server":      dataSourceMsSqlServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_mssql_database":                                        resourceMsSqlDatabase(),
		"azurerm_mssql_database_extended_auditing_policy":               resourceMsSqlDatabaseExtendedAuditingPolicy(),
		"azurerm_mssql_database_vulnerability_assessment_rule_baseline": resourceMsSqlDatabaseVulnerabilityAssessmentRuleBaseline(),
		"azurerm_mssql_elasticpool":                                     resourceMsSqlElasticPool(),
		"azurerm_mssql_firewall_rule":                                   resourceMsSqlFirewallRule(),
		"azurerm_mssql_job_agent":                                       resourceMsSqlJobAgent(),
		"azurerm_mssql_job_credential":                                  resourceMsSqlJobCredential(),
		"azurerm_mssql_outbound_firewall_rule":                          resourceMsSqlOutboundFirewallRule(),
		"azurerm_mssql_server":                                          resourceMsSqlServer(),
		"azurerm_mssql_server_extended_auditing_policy":                 resourceMsSqlServerExtendedAuditingPolicy(),
		"azurerm_mssql_server_microsoft_support_auditing_policy":        resourceMsSqlServerMicrosoftSupportAuditingPolicy(),
		"azurerm_mssql_server_security_alert_policy":                    resourceMsSqlServerSecurityAlertPolicy(),
		"azurerm_mssql_server_transparent_data_encryption":              resourceMsSqlTransparentDataEncryption(),
		"azurerm_mssql_server_vulnerability_assessment":                 resourceMsSqlServerVulnerabilityAssessment(),
		"azurerm_mssql_virtual_machine":                                 resourceMsSqlVirtualMachine(),
		"azurerm_mssql_virtual_network_rule":                            resourceMsSqlVirtualNetworkRule(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		MsSqlFailoverGroupResource{},
		MsSqlVirtualMachineAvailabilityGroupListenerResource{},
		MsSqlVirtualMachineGroupResource{},
		ServerDNSAliasResource{},
	}
}
