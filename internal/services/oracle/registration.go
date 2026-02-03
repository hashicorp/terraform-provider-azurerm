// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration = Registration{}
	_ sdk.TypedServiceRegistration     = Registration{}
)

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		AdbsCharSetsDataSource{},
		AdbsNCharSetsDataSource{},
		AutonomousDatabaseBackupDataSource{},
		AutonomousDatabaseBackupsDataSource{},
		AutonomousDatabaseCloneFromBackupDataSource{},
		AutonomousDatabaseCloneFromDatabaseDataSource{},
		AutonomousDatabaseRegularDataSource{},
		CloudVmClusterDataSource{},
		DBNodesDataSource{},
		DBServersDataSource{},
		DbSystemShapesDataSource{},
		ExadataInfraDataSource{},
		ExascaleDatabaseStorageVaultDataSource{},
		GiVersionsDataSource{},
		ResourceAnchorDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AutonomousDatabaseBackupResource{},
		AutonomousDatabaseCloneFromBackupResource{},
		AutonomousDatabaseCloneFromDatabaseResource{},
		AutonomousDatabaseRegularResource{},
		CloudVmClusterResource{},
		ExadataInfraResource{},
		ExascaleDatabaseStorageVaultResource{},
		ResourceAnchorResource{},
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
