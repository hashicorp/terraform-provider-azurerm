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
	if features.ThreePointOh() {
		return []sdk.DataSource{
			AppServiceSourceControlTokenDataSource{},
		}
	}
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	if features.ThreePointOh() {
		return []sdk.Resource{
			AppServiceSourceControlResource{},
			AppServiceSourceControlTokenResource{},
			WindowsWebAppResource{},
			LinuxWebAppResource{},
			AppServicePlanResource{},
		}
	}
	return []sdk.Resource{}
}
