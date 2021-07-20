package appservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/serviceplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/sourcecontrol"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/webapp"
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
			sourcecontrol.AppServiceSourceControlTokenDataSource{},
		}
	}
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	if features.ThreePointOh() {
		return []sdk.Resource{
			sourcecontrol.AppServiceSourceControlResource{},
			sourcecontrol.AppServiceSourceControlTokenResource{},
			webapp.WindowsWebAppResource{},
			webapp.LinuxWebAppResource{},
			serviceplan.AppServicePlanResource{},
		}
	}
	return []sdk.Resource{}
}
