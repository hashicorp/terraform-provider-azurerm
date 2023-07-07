// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssqlmanagedinstance

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
	return "service/mssqlmanagedinstance"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Microsoft SQL Server Managed Instances"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Database",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_mssql_managed_instance_security_alert_policy":       resourceMsSqlManagedInstanceSecurityAlertPolicy(),
		"azurerm_mssql_managed_instance_transparent_data_encryption": resourceMsSqlManagedInstanceTransparentDataEncryption(),
		"azurerm_mssql_managed_instance_vulnerability_assessment":    resourceMsSqlManagedInstanceVulnerabilityAssessment(),
	}
}

// DataSources returns the typed DataSources supported by this service
func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		MsSqlManagedInstanceDataSource{},
	}
}

// Resources returns the typed Resources supported by this service
func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		MsSqlManagedDatabaseResource{},
		MsSqlManagedInstanceActiveDirectoryAdministratorResource{},
		MsSqlManagedInstanceFailoverGroupResource{},
		MsSqlManagedInstanceResource{},
	}
}
