package appservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

func (r Registration) PackagePath() string {
	return "TODO: Not implemented yet"
}

func (r Registration) WebsiteCategories() []string {
	return nil
}

func (r Registration) Name() string {
	return "AppService"
}

func (r Registration) DataSources() []sdk.DataSource {
	if features.ThreePointOhBetaResources() {
		return []sdk.DataSource{
			AppServiceSourceControlTokenDataSource{},
			LinuxWebAppDataSource{},
			ServicePlanDataSource{},
			WindowsWebAppDataSource{},
		}
	}
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	if features.ThreePointOhBetaResources() {
		return []sdk.Resource{
			AppServiceSourceControlResource{},
			AppServiceSourceControlTokenResource{},
			LinuxFunctionAppResource{},
			LinuxWebAppResource{},
			ServicePlanResource{},
			WindowsWebAppResource{},
			WindowsFunctionAppResource{},
		}
	}
	return []sdk.Resource{}
}
