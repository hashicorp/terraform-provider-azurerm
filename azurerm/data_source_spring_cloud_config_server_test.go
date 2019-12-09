package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMSpringCloudConfigServer_complete(t *testing.T) {
	dataSourceName := "data.azurerm_spring_cloud_config_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpringCloudConfigServer_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(dataSourceName, "label", "config"),
					resource.TestCheckResourceAttr(dataSourceName, "search_paths.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "search_paths.0", "dir1"),
					resource.TestCheckResourceAttr(dataSourceName, "search_paths.1", "dir2"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.0.name", "repo1"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.0.uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.0.label", "config"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.0.search_paths.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.0.search_paths.0", "dir1"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.0.search_paths.1", "dir2"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.1.name", "repo2"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.1.uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.1.label", "config"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.1.search_paths.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.1.search_paths.0", "dir1"),
					resource.TestCheckResourceAttr(dataSourceName, "repositories.1.search_paths.1", "dir2"),
				),
			},
		},
	})
}

func testAccDataSourceSpringCloudConfigServer_complete(rInt int, location string) string {
	config := testAccAzureRMSpringCloudConfigServer_complete(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_config_server" "test" {
    spring_cloud_id  = azurerm_spring_cloud.test.id
}
`, config)
}
