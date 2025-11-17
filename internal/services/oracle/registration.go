// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistration = Registration{}

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
		DatabaseVersionsDataSource{},
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
