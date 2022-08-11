package containerapps_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerAppEnvironmentDaprComponentDataSource struct{}

func TestAccContainerAppEnvironmentDaprComponentDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_environment_dapr_component", "test")
	r := ContainerAppEnvironmentDaprComponentDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("type").HasValue("state.azure.blobstorage"),
				check.That(data.ResourceName).Key("version").HasValue("v1"),
			),
		},
	})
}

func (d ContainerAppEnvironmentDaprComponentDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_environment_dapr_component" "test" {
  name 						   = azurerm_container_app_environment_dapr_component.test.name
  container_app_environment_id = azurerm_container_app_environment_dapr_component.test.container_app_environment_id   
}

`, ContainerAppEnvironmentDaprComponentResource{}.basic(data))
}
