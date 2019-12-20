package mssql

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Microsoft SQL Server / SQL Azure"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_mssql_elasticpool": dataSourceArmMsSqlElasticpool(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_mssql_elasticpool":                                     resourceArmMsSqlElasticPool(),
		"azurerm_mssql_database_vulnerability_assessment_rule_baseline": resourceArmMssqlDatabaseVulnerabilityAssessmentRuleBaseline(),
		"azurerm_mssql_server_security_alert_policy":                    resourceArmMssqlServerSecurityAlertPolicy(),
		"azurerm_mssql_server_vulnerability_assessment":                 resourceArmMssqlServerVulnerabilityAssessment(),
	}
}
