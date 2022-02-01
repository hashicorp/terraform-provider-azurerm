package postgresqlhsc

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

func (r Registration) PackagePath() string {
	return "TODO: Not implemented yet"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"PostgreSQL HyperScale",
	}
}

func (r Registration) Name() string {
	return "PostgreSQL HyperScale"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		PostgreSQLHyperScaleServerGroupResource{},
	}
}
