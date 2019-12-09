package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMSpringCloudConfigServer_basic(t *testing.T) {
	resourceName := "azurerm_spring_cloud_config_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudConfigServer_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(resourceName, "label", "config"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSpringCloudConfigServer_update(t *testing.T) {
	resourceName := "azurerm_spring_cloud_config_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudConfigServer_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(resourceName, "label", "config"),
					resource.TestCheckResourceAttr(resourceName, "search_paths.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "repositories.#", "0"),
				),
			},
			{
				Config: testAccAzureRMSpringCloudConfigServer_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(resourceName, "label", "config"),
					resource.TestCheckResourceAttr(resourceName, "search_paths.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "search_paths.0", "dir1"),
					resource.TestCheckResourceAttr(resourceName, "search_paths.1", "dir2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.name", "repo1"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.label", "config"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.search_paths.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.search_paths.0", "dir1"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.search_paths.1", "dir2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.name", "repo2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.label", "config"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.search_paths.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.search_paths.0", "dir1"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.search_paths.1", "dir2"),
				),
			},
			{
				Config: testAccAzureRMSpringCloudConfigServer_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "search_paths.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "repositories.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMSpringCloudConfigServer_complete(t *testing.T) {
	resourceName := "azurerm_spring_cloud_config_server.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSpringCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudConfigServer_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(resourceName, "label", "config"),
					resource.TestCheckResourceAttr(resourceName, "search_paths.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "search_paths.0", "dir1"),
					resource.TestCheckResourceAttr(resourceName, "search_paths.1", "dir2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.name", "repo1"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.label", "config"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.search_paths.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.search_paths.0", "dir1"),
					resource.TestCheckResourceAttr(resourceName, "repositories.0.search_paths.1", "dir2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.name", "repo2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.uri", "https://github.com/Azure-Samples/piggymetrics"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.label", "config"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.search_paths.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.search_paths.0", "dir1"),
					resource.TestCheckResourceAttr(resourceName, "repositories.1.search_paths.1", "dir2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAzureRMSpringCloudConfigServer_basic(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_config_server" "test" {
  spring_cloud_id = azurerm_spring_cloud.test.id

  uri = "https://github.com/Azure-Samples/piggymetrics"
  label = "config"
}

`, testAccAzureRMSpringCloud_basic(rInt, location))
}

func testAccAzureRMSpringCloudConfigServer_complete(rInt int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_config_server" "test" {
  spring_cloud_id = azurerm_spring_cloud.test.id
  uri = "https://github.com/Azure-Samples/piggymetrics"
  label = "config"
  search_paths = ["dir1", "dir2"]
  
  repositories {
     name = "repo1"
     uri  = "https://github.com/Azure-Samples/piggymetrics"
     label = "config"
     search_paths =  ["dir1", "dir2"]
  }

  repositories {
     name = "repo2"
     uri  = "https://github.com/Azure-Samples/piggymetrics"
     label = "config"
     search_paths =  ["dir1", "dir2"]
  }
}
`, testAccAzureRMSpringCloud_basic(rInt, location))
}
