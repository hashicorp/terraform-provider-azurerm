package healthcare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type HealthCareWorkspaceIotConnectorDataSource struct{}

func TestAccHealthCareWorkspaceIotConnectorDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_healthcare_iot_connector", "test")
	r := HealthCareWorkspaceIotConnectorDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists()),
		},
	})
}

func (HealthCareWorkspaceIotConnectorDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_healthcare_iot_connector" "test" {
  name                = azurerm_healthcare_iot_connector.test.name
  resource_group_name = azurerm_healthcare_iot_connector.test.resource_group_name
}
`, HealthCareWorkspaceIotConnectorResource{}.basic(data))
}
