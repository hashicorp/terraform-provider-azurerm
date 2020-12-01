package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
)

func TestAccAzureRMMediaAsset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMediaAssetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMediaAsset_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", "Asset-Content1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMediaAsset_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMediaAssetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMediaAsset_complete(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "alternate_id", "Asset-alternateid"),
					resource.TestCheckResourceAttr(data.ResourceName, "storage_account_name", fmt.Sprintf("acctestsa1%s", data.RandomString)),
					resource.TestCheckResourceAttr(data.ResourceName, "container", "asset-container"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Asset description"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMediaAsset_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMediaAssetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMediaAsset_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", "Asset-Content1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMediaAsset_complete(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "alternate_id", "Asset-alternateid"),
					resource.TestCheckResourceAttr(data.ResourceName, "container", "asset-container"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Asset description"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMediaAsset_basic(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "alternate_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "Asset-Content1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMMediaAssetDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Media.AssetsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_media_asset" {
			continue
		}

		id, err := parse.MediaAssetsID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := conn.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Media Asset still exists:\n%#v", resp)
		}
	}

	return nil
}

func testAccAzureRMMediaAsset_basic(data acceptance.TestData) string {
	template := testAccAzureRMMediaAsset_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset" "test" {
  name                        = "Asset-Content1"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
}

`, template)
}

func testAccAzureRMMediaAsset_complete(data acceptance.TestData) string {
	template := testAccAzureRMMediaAsset_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset" "test" {
  name                        = "Asset-Content1"
  description                 = "Asset description"
  resource_group_name         = azurerm_resource_group.test.name
  media_services_account_name = azurerm_media_services_account.test.name
  storage_account_name        = azurerm_storage_account.test.name
  alternate_id                = "Asset-alternateid"
  container                   = "asset-container"
}

`, template)
}

func testAccAzureRMMediaAsset_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-media-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa1%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_media_services_account" "test" {
  name                = "acctestmsa%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  storage_account {
    id         = azurerm_storage_account.test.id
    is_primary = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}
