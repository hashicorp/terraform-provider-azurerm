package postgresqlhsc

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.UntypedServiceRegistration = Registration{}

type Registration struct {
}

func (r Registration) Name() string {
	return "Postgresql HSC"
}

func (r Registration) WebsiteCategories() []string {
	return []string{}
}

func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}

func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{}
}
