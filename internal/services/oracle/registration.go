// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracle

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration = Registration{}
)

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		AdbsCharSetsDataSource{},
		AdbsNCharSetsDataSource{},
		AutonomousDatabaseRegularDataSource{},
		CloudVmClusterDataSource{},
		DBNodesDataSource{},
		DBServersDataSource{},
		DbSystemShapesDataSource{},
		ExadataInfraDataSource{},
		GiVersionsDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AutonomousDatabaseRegularResource{},
		CloudVmClusterResource{},
		ExadataInfraResource{},
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
