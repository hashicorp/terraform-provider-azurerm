package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataFactoryDataSource_basic(t *testing.T) {
	dsn := "azurerm_data_factory.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryDataSource_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryExists(dsn),
				),
			},
		},
	})
}

func testAccAzureRMDataFactoryDataSource_basic(rInt int, location string) string {
	config := testAccAzureRMDataFactory_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_data_factory" "test" {
  name                = "${azurerm_data_factory.test.name}"
  resource_group_name = "${azurerm_data_factory.test.resource_group_name}"
}
`, config)
}
