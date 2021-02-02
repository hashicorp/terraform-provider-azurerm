package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databoxedge/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccDataboxEdgeDevice_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_device", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeDevice_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeDeviceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataboxEdgeDevice_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_device", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeDevice_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeDeviceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccDataboxEdgeDevice_requiresImport),
		},
	})
}

func TestAccDataboxEdgeDevice_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_device", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeDevice_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeDeviceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataboxEdgeDevice_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_device", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeDevice_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeDeviceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDataboxEdgeDevice_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeDeviceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDataboxEdgeDevice_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeDeviceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckDataboxEdgeDeviceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataboxEdge.DeviceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Databox Edge Device not found: %s", resourceName)
		}
		id, err := parse.DataboxEdgeDeviceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.Name, id.ResourceGroup); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Databox Edge Device %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on DataboxEdge.DeviceClient: %+v", err)
		}
		return nil
	}
}

func testCheckDataboxEdgeDeviceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataboxEdge.DeviceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_databox_edge_device" {
			continue
		}
		id, err := parse.DataboxEdgeDeviceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.Name, id.ResourceGroup); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on DataboxEdge.DeviceClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

// Location has to be hard coded due to limited support of locations for this resource
func testAccDataboxEdgeDevice_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databoxedge-%d"
  location = "%s"
}
`, data.RandomInteger, "eastus")
}

func testAccDataboxEdgeDevice_basic(data acceptance.TestData) string {
	template := testAccDataboxEdgeDevice_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_device" "test" {
  name                = "acctest-dd-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "Edge-Standard"
}
`, template, data.RandomString)
}

func testAccDataboxEdgeDevice_requiresImport(data acceptance.TestData) string {
	config := testAccDataboxEdgeDevice_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_device" "import" {
  name                = azurerm_databox_edge_device.test.name
  resource_group_name = azurerm_databox_edge_device.test.resource_group_name
  location            = azurerm_databox_edge_device.test.location

  sku_name = "Edge-Standard"
}
`, config)
}

func testAccDataboxEdgeDevice_complete(data acceptance.TestData) string {
	template := testAccDataboxEdgeDevice_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_device" "test" {
  name                = "acctest-dd-%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "Edge-Standard"

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomString)
}
