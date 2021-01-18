package datafactory_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type DataFactoryDataSource struct {
}

func TestAccDataFactoryDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_factory", "test")
	r := DataFactoryDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccDataFactoryDataSource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")
	r := DataFactoryDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.identity(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("identity.#").Exists(),
				check.That(data.ResourceName).Key("identity.0.type").Exists(),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
	})
}

func (DataFactoryDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_factory" "test" {
  name                = azurerm_data_factory.test.name
  resource_group_name = azurerm_data_factory.test.resource_group_name
}
`, DataFactoryResource{}.basic(data))
}

func (DataFactoryDataSource) identity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_data_factory" "test" {
  name                = azurerm_data_factory.test.name
  resource_group_name = azurerm_data_factory.test.resource_group_name
}
`, DataFactoryResource{}.identity(data))
}
