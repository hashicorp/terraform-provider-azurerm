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
					resource.TestCheckResourceAttr(data.ResourceName, "description", "Asset description"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMediaAsset_alternate_id(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMediaAssetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMediaAsset_alternate_id(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "alternate_id", "assetalternateid"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMediaAsset_custom_container(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMediaAssetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMediaAsset_custom_container(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "container", "assetcontainer"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMediaAsset_storage_account(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_media_asset", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMediaAssetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMediaAsset_storage_account(data),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "storage_account_name", fmt.Sprintf("acctestsa1%s", data.RandomString)),
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
	name                        = "asset%s"
	description                 = "Asset description"
	resource_group_name         = azurerm_resource_group.test.name
	media_services_account_name = azurerm_media_services_account.test.name
}

`, template, data.RandomString)
}

func testAccAzureRMMediaAsset_alternate_id(data acceptance.TestData) string {
	template := testAccAzureRMMediaAsset_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset" "test" {
	name                        = "asset%s"
	description                 = "Asset description"
	resource_group_name         = azurerm_resource_group.test.name
	media_services_account_name = azurerm_media_services_account.test.name
	alternate_id                = "assetalternateid"
}

`, template, data.RandomString)
}

func testAccAzureRMMediaAsset_storage_account(data acceptance.TestData) string {
	template := testAccAzureRMMediaAsset_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset" "test" {
	name                        = "asset%s"
	description                 = "Asset description"
	resource_group_name         = azurerm_resource_group.test.name
	media_services_account_name = azurerm_media_services_account.test.name
	storage_account_name        = azurerm_storage_account.test.name
}

`, template, data.RandomString)
}

func testAccAzureRMMediaAsset_custom_container(data acceptance.TestData) string {
	template := testAccAzureRMMediaAsset_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_media_asset" "test" {
	name                        = "asset%s"
	description                 = "Asset description"
	resource_group_name         = azurerm_resource_group.test.name
	media_services_account_name = azurerm_media_services_account.test.name
	storage_account_name        = azurerm_storage_account.test.name
	container                   = "assetcontainer"
}

`, template, data.RandomString)
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
