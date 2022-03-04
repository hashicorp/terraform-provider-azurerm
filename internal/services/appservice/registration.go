package appservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/app-service"
}

func (r Registration) WebsiteCategories() []string {
	return nil
}

func (r Registration) Name() string {
	return "AppService"
}

func (r Registration) DataSources() []sdk.DataSource {
	if features.ThreePointOhAppServiceResources() {
		return []sdk.DataSource{
			AppServiceSourceControlTokenDataSource{},
			LinuxFunctionAppDataSource{},
			LinuxWebAppDataSource{},
			ServicePlanDataSource{},
			WindowsFunctionAppDataSource{},
			WindowsWebAppDataSource{},
		}
	}
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	if features.ThreePointOhAppServiceResources() {
		return []sdk.Resource{
			AppServiceSourceControlTokenResource{},
			FunctionAppActiveSlotResource{},
			LinuxFunctionAppResource{},
			LinuxFunctionAppSlotResource{},
			LinuxWebAppResource{},
			LinuxWebAppSlotResource{},
			ServicePlanResource{},
			SourceControlResource{},
			SourceControlSlotResource{},
			WebAppActiveSlotResource{},
			WindowsWebAppResource{},
			WindowsFunctionAppResource{},
			WindowsWebAppSlotResource{},
			WindowsFunctionAppSlotResource{},
		}
	}
	return []sdk.Resource{}
}
