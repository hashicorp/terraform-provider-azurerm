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

func TestAccDataboxEdgeOrder_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_order", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeOrder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataboxEdgeOrder_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_order", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeOrder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccDataboxEdgeOrder_requiresImport),
		},
	})
}

func TestAccDataboxEdgeOrder_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_order", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeOrder_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataboxEdgeOrder_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_order", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeOrder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDataboxEdgeOrder_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDataboxEdgeOrder_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataboxEdgeOrder_updateContactInformation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_order", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeOrder_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDataboxEdgeOrder_updateContactInformation(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataboxEdgeOrder_updateCurrentStatus(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_order", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeOrder_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDataboxEdgeOrder_updateCurrentStatus(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccDataboxEdgeOrder_updateShippingAddress(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databox_edge_order", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckDataboxEdgeOrderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataboxEdgeOrder_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccDataboxEdgeOrder_updateShippingAddress(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckDataboxEdgeOrderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckDataboxEdgeOrderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataboxEdge.OrderClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Databox Edge Order not found: %s", resourceName)
		}
		id, err := parse.DataboxEdgeOrderID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.DeviceName, id.ResourceGroup); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Databox Edge Order does not exist")
			}
			return fmt.Errorf("bad: Get on DataboxEdge.OrderClient: %+v", err)
		}
		return nil
	}
}

func testCheckDataboxEdgeOrderDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataboxEdge.OrderClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_databox_edge_order" {
			continue
		}
		id, err := parse.DataboxEdgeOrderID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.DeviceName, id.ResourceGroup); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on DataboxEdge.OrderClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccDataboxEdgeOrder_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-databoxedge-%d"
  location = "%s"
}

resource "azurerm_databox_edge_device" "test" {
  name                = "acctest-dd-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataboxEdgeOrder_basic(data acceptance.TestData) string {
	template := testAccDataboxEdgeOrder_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_order" "test" {
  resource_group_name = azurerm_resource_group.test.name
  device_name         = azurerm_databox_edge_device.test.name
}
`, template)
}

func testAccDataboxEdgeOrder_requiresImport(data acceptance.TestData) string {
	config := testAccDataboxEdgeOrder_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_order" "import" {
	name                = azurerm_databox_edge_order.test.device_name
	resource_group_name = azurerm_databox_edge_order.test.resource_group_name
}
`, config)
}

func testAccDataboxEdgeOrder_complete(data acceptance.TestData) string {
	template := testAccDataboxEdgeOrder_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_order" "test" {
  name                = azurerm_databox_edge_device.test.name
  resource_group_name = azurerm_resource_group.test.name

  contact_information {
    name           = "TerraForm Test"
    emails         = ["test@foo.com"]
    company_name   = "Microsoft"
    phone_number   = "(800) 555-1234"
  }

  current_status {
    comments = ""
  }

  shipping_address {
    address_line1 = "One Microsoft Way"
    city          = "Redmond"
    postal_code   = "98052"
    state         = "WA"
    country       = "USA"
  }
}
`, template)
}

func testAccDataboxEdgeOrder_updateContactInformation(data acceptance.TestData) string {
	template := testAccDataboxEdgeOrder_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_order" "test" {
  resource_group_name = azurerm_resource_group.test.name
  device_name         = azurerm_databox_edge_device.test.name

  contact_information {
    company_name   = "Microsoft"
    contact_person = "John Mcclane"
    email_lists    = ["john@microsoft.com"]
    phone          = "(800) 426-9400"
  }

  current_status {
    comments = ""
  }

  shipping_address {
    address_line1 = "Microsoft Corporation"
    city          = "WA"
    country       = "USA"
    postal_code   = "98052"
    state         = "WA"
    address_line2 = "One Microsoft Way"
    address_line3 = "Redmond"
  }
}
`, template)
}

func testAccDataboxEdgeOrder_updateCurrentStatus(data acceptance.TestData) string {
	template := testAccDataboxEdgeOrder_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_order" "test" {
  resource_group_name = azurerm_resource_group.test.name
  device_name = azurerm_databox_edge_device.test.name

  contact_information {
    company_name   = "Microsoft"
    contact_person = "John Mcclane"
    email_lists    = ["john@microsoft.com"]
    phone          = "(800) 426-9400"
  }

  current_status {
    comments = ""
  }

  shipping_address {
    address_line1 = "Microsoft Corporation"
    city          = "WA"
    country       = "USA"
    postal_code   = "98052"
    state         = "WA"
    address_line2 = "One Microsoft Way"
    address_line3 = "Redmond"
  }
}
`, template)
}

func testAccDataboxEdgeOrder_updateShippingAddress(data acceptance.TestData) string {
	template := testAccDataboxEdgeOrder_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databox_edge_order" "test" {
  resource_group_name = azurerm_resource_group.test.name
  device_name         = azurerm_databox_edge_device.test.name

  contact_information {
    company_name   = "Microsoft"
    contact_person = "John Mcclane"
    email_lists    = ["john@microsoft.com"]
    phone          = "(800) 426-9400"
  }

  current_status {
    comments = ""
  }

  shipping_address {
    address_line1 = "Microsoft Corporation"
    city          = "WA"
    country       = "USA"
    postal_code   = "98052"
    state         = "WA"
    address_line2 = "One Microsoft Way"
    address_line3 = "Redmond"
  }
}
`, template)
}
