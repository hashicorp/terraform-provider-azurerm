package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/avs/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"testing"
)

func TestAccAzureRMavsAuthorization_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_authorization", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsAuthorizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMavsAuthorization_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsAuthorizationExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "express_route_authorization_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "express_route_authorization_key"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMavsAuthorization_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_avs_authorization", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsAuthorizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMavsAuthorization_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsAuthorizationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMavsAuthorization_requiresImport),
		},
	})
}

func testCheckAzureRMavsAuthorizationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Avs.AuthorizationClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("avs Authorization not found: %s", resourceName)
		}
		id, err := parse.AvsAuthorizationID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Avs Authorization %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Avs.AuthorizationClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMavsAuthorizationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Avs.AuthorizationClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_avs_authorization" {
			continue
		}
		id, err := parse.AvsAuthorizationID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.PrivateCloudName, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Avs.AuthorizationClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMavsAuthorization_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-avs-%d"
  location = "%s"
}

resource "azurerm_avs_private_cloud" "test" {
  name                = "acctest-apc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  network_block       = "192.168.48.0/22"
  sku {
    name = "av36"
  }

  management_cluster {
    cluster_size = 3
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMavsAuthorization_basic(data acceptance.TestData) string {
	template := testAccAzureRMavsAuthorization_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_avs_authorization" "test" {
  name             = "acctest-AA-%d"
  private_cloud_id = azurerm_avs_private_cloud.test.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMavsAuthorization_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMavsAuthorization_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_avs_authorization" "import" {
  name             = azurerm_avs_authorization.test.name
  private_cloud_id = azurerm_avs_private_cloud.test.id
}
`, config)
}
