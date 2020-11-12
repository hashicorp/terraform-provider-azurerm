package appconfiguration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAppConfigurationDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAppConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAppConfigurationResource_standard(data),
			},
			{
				Config: testAppConfigurationDataSource_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAppConfigurationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "location"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "sku"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_read_key.0.connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_read_key.0.id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_read_key.0.secret"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_write_key.0.connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_write_key.0.id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_write_key.0.secret"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_read_key.0.connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_read_key.0.id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_read_key.0.secret"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_write_key.0.connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_write_key.0.id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_write_key.0.secret"),
				),
			},
		},
	})
}

func testAppConfigurationDataSource_basic(data acceptance.TestData) string {
	template := testAppConfigurationResource_standard(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_configuration" "test" {
  name                = azurerm_app_configuration.test.name
  resource_group_name = azurerm_app_configuration.test.resource_group_name
}
`, template)
}
