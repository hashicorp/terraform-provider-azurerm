package mssql

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Microsoft SQL Server / SQL Azure"
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
		"azurerm_mssql_database":    dataSourceArmMsSqlDatabase(),
		"azurerm_mssql_elasticpool": dataSourceArmMsSqlElasticpool(),
		"azurerm_mssql_server":      dataSourceMsSqlServer(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_mssql_database":                                        resourceArmMsSqlDatabase(),
		"azurerm_mssql_database_extended_auditing_policy":               resourceArmMsSqlDatabaseExtendedAuditingPolicy(),
		"azurerm_mssql_database_vulnerability_assessment_rule_baseline": resourceArmMssqlDatabaseVulnerabilityAssessmentRuleBaseline(),
		"azurerm_mssql_elasticpool":                                     resourceArmMsSqlElasticPool(),
		"azurerm_mssql_server":                                          resourceArmMsSqlServer(),
		"azurerm_mssql_server_extended_auditing_policy":                 resourceArmMsSqlServerExtendedAuditingPolicy(),
		"azurerm_mssql_server_security_alert_policy":                    resourceArmMssqlServerSecurityAlertPolicy(),
		"azurerm_mssql_server_vulnerability_assessment":                 resourceArmMssqlServerVulnerabilityAssessment(),
		"azurerm_mssql_virtual_machine":                                 resourceArmMsSqlVirtualMachine(),
	}
}
