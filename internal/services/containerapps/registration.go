package containerapps

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Container Apps",
	}
}

func (r Registration) Name() string {
	return "Container Apps"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		ContainerAppEnvironmentDataSource{},
		ContainerAppEnvironmentCertificateDataSource{},
		ContainerAppEnvironmentDaprComponentDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ContainerAppEnvironmentCertificateResource{},
		ContainerAppEnvironmentDaprComponentResource{},
		ContainerAppEnvironmentResource{},
		ContainerAppEnvironmentStorageResource{},
		ContainerAppResource{},
	}
}
