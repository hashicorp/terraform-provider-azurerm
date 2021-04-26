package mssql

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_mssql_database":    dataSourceMsSqlDatabase(),
		"azurerm_mssql_elasticpool": dataSourceMsSqlElasticpool(),
		"azurerm_mssql_server":      dataSourceMsSqlServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_mssql_database":                                        resourceMsSqlDatabase(),
		"azurerm_mssql_database_extended_auditing_policy":               resourceMsSqlDatabaseExtendedAuditingPolicy(),
		"azurerm_mssql_database_vulnerability_assessment_rule_baseline": resourceMsSqlDatabaseVulnerabilityAssessmentRuleBaseline(),
		"azurerm_mssql_elasticpool":                                     resourceMsSqlElasticPool(),
		"azurerm_mssql_job_agent":                                       resourceMsSqlJobAgent(),
		"azurerm_mssql_job_credential":                                  resourceMsSqlJobCredential(),
		"azurerm_mssql_firewall_rule":                                   resourceMsSqlFirewallRule(),
		"azurerm_mssql_server":                                          resourceMsSqlServer(),
		"azurerm_mssql_server_extended_auditing_policy":                 resourceMsSqlServerExtendedAuditingPolicy(),
		"azurerm_mssql_server_security_alert_policy":                    resourceMsSqlServerSecurityAlertPolicy(),
		"azurerm_mssql_server_vulnerability_assessment":                 resourceMsSqlServerVulnerabilityAssessment(),
		"azurerm_mssql_virtual_machine":                                 resourceMsSqlVirtualMachine(),
		"azurerm_mssql_virtual_network_rule":                            resourceMsSqlVirtualNetworkRule(),
		"azurerm_mssql_server_transparent_data_encryption":              resourceMsSqlTransparentDataEncryption(),
	}
}
