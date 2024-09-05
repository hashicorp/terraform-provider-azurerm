// Copyright © 2024, Oracle and/or its affiliates. All rights reserved

package oracledatabase

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration = Registration{}
)

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		ExadataInfraDataSource{},
		CloudVmClusterDataSource{},
		DBServersDataSource{},
		AutonomousDatabaseRegularDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ExadataInfraResource{},
		CloudVmClusterResource{},
		AdbsRegularResource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "App Service"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"App Service",
	}
}
