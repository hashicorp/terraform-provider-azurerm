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
		AutonomousDatabaseRegularDataSource{},
		CloudVmClusterDataSource{},
		DBNodesDataSource{},
		DBServersDataSource{},
		DbSystemShapesDataSource{},
		ExadataInfraDataSource{},
		ExascaleDatabaseStorageVaultDataSource{},
		GiVersionsDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AutonomousDatabaseBackupResource{},
		AutonomousDatabaseRegularResource{},
		CloudVmClusterResource{},
		ExadataInfraResource{},
		ExascaleDatabaseStorageVaultResource{},
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
