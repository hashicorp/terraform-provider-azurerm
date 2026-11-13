// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/actions/mssql_execute_job"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_database_extended_auditing_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_database_vulnerability_assessment_rule_baseline"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_elasticpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_failover_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_firewall_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_job"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_job_agent"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_job_credential"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_job_schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_job_step"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_job_target_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_outbound_firewall_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_server_dns_alias"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_server_extended_auditing_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_server_microsoft_support_auditing_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_server_security_alert_policy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_server_transparent_data_encryption"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_server_vulnerability_assessment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_virtual_machine"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_virtual_machine_availability_group_listener"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_virtual_machine_group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/resources/mssql_virtual_network_rule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
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
		"azurerm_mssql_database":    mssql_database.RegisterDataSource(),
		"azurerm_mssql_elasticpool": mssql_elasticpool.RegisterDataSource(),
		"azurerm_mssql_server":      mssql_server.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_mssql_database":                                        mssql_database.RegisterResource(),
		"azurerm_mssql_database_extended_auditing_policy":               mssql_database_extended_auditing_policy.RegisterResource(),
		"azurerm_mssql_database_vulnerability_assessment_rule_baseline": mssql_database_vulnerability_assessment_rule_baseline.RegisterResource(),
		"azurerm_mssql_elasticpool":                                     mssql_elasticpool.RegisterResource(),
		"azurerm_mssql_firewall_rule":                                   mssql_firewall_rule.RegisterResource(),
		"azurerm_mssql_job_agent":                                       mssql_job_agent.RegisterResource(),
		"azurerm_mssql_job_credential":                                  mssql_job_credential.RegisterResource(),
		"azurerm_mssql_outbound_firewall_rule":                          mssql_outbound_firewall_rule.RegisterResource(),
		"azurerm_mssql_server":                                          mssql_server.RegisterResource(),
		"azurerm_mssql_server_extended_auditing_policy":                 mssql_server_extended_auditing_policy.RegisterResource(),
		"azurerm_mssql_server_microsoft_support_auditing_policy":        mssql_server_microsoft_support_auditing_policy.RegisterResource(),
		"azurerm_mssql_server_security_alert_policy":                    mssql_server_security_alert_policy.RegisterResource(),
		"azurerm_mssql_server_transparent_data_encryption":              mssql_server_transparent_data_encryption.RegisterResource(),
		"azurerm_mssql_server_vulnerability_assessment":                 mssql_server_vulnerability_assessment.RegisterResource(),
		"azurerm_mssql_virtual_machine":                                 mssql_virtual_machine.RegisterResource(),
		"azurerm_mssql_virtual_network_rule":                            mssql_virtual_network_rule.RegisterResource(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		mssql_failover_group.DataSource{},
	}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		mssql_failover_group.Resource{},
		mssql_job.Resource{},
		mssql_job_schedule.Resource{},
		mssql_job_step.Resource{},
		mssql_job_target_group.Resource{},
		mssql_virtual_machine_availability_group_listener.Resource{},
		mssql_virtual_machine_group.Resource{},
		mssql_server_dns_alias.Resource{},
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{
		mssql_execute_job.RegisterAction,
	}
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
	return []sdk.FrameworkListWrappedResource{
		mssql_database.ListResource{},
		mssql_elasticpool.ListResource{},
		mssql_job_agent.ListResource{},
		mssql_server.ListResource{},
		mssql_virtual_machine.ListResource{},
	}
}
