// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_adbs_character_sets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_adbs_national_character_sets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_autonomous_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_autonomous_database_backup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_autonomous_database_backups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_autonomous_database_clone_from_backup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_autonomous_database_clone_from_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_autonomous_database_cross_region_disaster_recovery"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_cloud_vm_cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_database_system_versions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_db_nodes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_db_servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_db_system_shapes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_exadata_infrastructure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_exascale_database_storage_vault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_gi_versions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/resources/oracle_resource_anchor"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration             = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
)

func (r Registration) AssociatedGitHubLabel() string {
	return "service/oracle"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		oracle_adbs_character_sets.DataSource{},
		oracle_adbs_national_character_sets.DataSource{},
		oracle_autonomous_database_backup.DataSource{},
		oracle_autonomous_database_backups.DataSource{},
		oracle_autonomous_database_clone_from_backup.DataSource{},
		oracle_autonomous_database_clone_from_database.DataSource{},
		oracle_autonomous_database_cross_region_disaster_recovery.DataSource{},
		oracle_autonomous_database.DataSource{},
		oracle_cloud_vm_cluster.DataSource{},
		oracle_database_system_versions.DataSource{},
		oracle_db_nodes.DataSource{},
		oracle_db_servers.DataSource{},
		oracle_db_system_shapes.DataSource{},
		oracle_exadata_infrastructure.DataSource{},
		oracle_exascale_database_storage_vault.DataSource{},
		oracle_gi_versions.DataSource{},
		oracle_resource_anchor.DataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		oracle_autonomous_database_backup.Resource{},
		oracle_autonomous_database_clone_from_backup.Resource{},
		oracle_autonomous_database_clone_from_database.Resource{},
		oracle_autonomous_database_cross_region_disaster_recovery.Resource{},
		oracle_autonomous_database.Resource{},
		oracle_cloud_vm_cluster.Resource{},
		oracle_exadata_infrastructure.Resource{},
		oracle_exascale_database_storage_vault.Resource{},
		oracle_resource_anchor.Resource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Oracle"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Oracle",
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
