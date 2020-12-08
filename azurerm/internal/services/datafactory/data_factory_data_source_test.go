package datafactory_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataFactoryDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_data_factory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryDataSource_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(data.ResourceName),
				),
			},
		},
	})
}

func testAccAzureRMDataFactoryDataSource_basic(data acceptance.TestData) string {
	config := testAccAzureRMDataFactory_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_factory" "test" {
  name                = azurerm_data_factory.test.name
  resource_group_name = azurerm_data_factory.test.resource_group_name
}
`, config)
}

func TestAccAzureRMDataFactoryDataSource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryDataSource_identity(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.type"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
		},
	})
}

func testAccAzureRMDataFactoryDataSource_identity(data acceptance.TestData) string {
	config := testAccAzureRMDataFactory_identity(data)
	return fmt.Sprintf(`
%s

data "azurerm_data_factory" "test" {
  name                = azurerm_data_factory.test.name
  resource_group_name = azurerm_data_factory.test.resource_group_name
}
`, config)
}
