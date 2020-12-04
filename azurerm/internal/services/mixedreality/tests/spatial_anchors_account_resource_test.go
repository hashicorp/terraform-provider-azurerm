package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mixedreality/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSpatialAnchorsAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spatial_anchors_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpatialAnchorsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpatialAnchorsAccount_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpatialAnchorsAccountExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSpatialAnchorsAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spatial_anchors_account", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpatialAnchorsAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpatialAnchorsAccount_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpatialAnchorsAccountExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Environment", "Production"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSpatialAnchorsAccountExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).MixedReality.SpatialAnchorsAccountClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		id, err := parse.SpatialAnchorsAccountID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on spatialAnchorsAccountClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Spatial Anchors Account %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMSpatialAnchorsAccountDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).MixedReality.SpatialAnchorsAccountClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_spatial_anchors_account" {
			continue
		}

		id, err := parse.SpatialAnchorsAccountID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Spatial Anchors Account still exists: %q", id.Name)
		}
	}

	return nil
}

func testAccAzureRMSpatialAnchorsAccount_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mr-%d"
  location = "%s"
}

resource "azurerm_spatial_anchors_account" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

func testAccAzureRMSpatialAnchorsAccount_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mr-%d"
  location = "%s"
}

resource "azurerm_spatial_anchors_account" "test" {
  name                = "acCTestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    Environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}
