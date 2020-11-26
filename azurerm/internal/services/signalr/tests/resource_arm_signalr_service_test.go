package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr/parse"
)

func TestAccAzureRMSignalRService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSignalRService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSignalRService_requiresImport),
		},
	})
}

func TestAccAzureRMSignalRService_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_standardWithCapacity(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSignalRService_standardWithCap2(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_standardWithCapacity(data, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSignalRService_skuUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			{
				Config: testAccAzureRMSignalRService_standardWithCapacity(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			{
				Config: testAccAzureRMSignalRService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func TestAccAzureRMSignalRService_capacityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_standardWithCapacity(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			{
				Config: testAccAzureRMSignalRService_standardWithCapacity(data, 5),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "5"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			{
				Config: testAccAzureRMSignalRService_standardWithCapacity(data, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func TestAccAzureRMSignalRService_skuAndCapacityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			{
				Config: testAccAzureRMSignalRService_standardWithCapacity(data, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Standard_S1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			{
				Config: testAccAzureRMSignalRService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.name", "Free_F1"),
					resource.TestCheckResourceAttr(data.ResourceName, "sku.0.capacity", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
		},
	})
}

func TestAccAzureRMSignalRService_serviceMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_withServiceMode(data, "Serverless"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSignalRService_cors(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSignalRServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSignalRService_withCors(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSignalRServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "cors.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "cors.0.allowed_origins.#", "2"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hostname"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "server_port"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "primary_connection_string"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_access_key"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "secondary_connection_string"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMSignalRService_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Free_F1"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMSignalRService_requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_signalr_service" "import" {
  name                = azurerm_signalr_service.test.name
  location            = azurerm_signalr_service.test.location
  resource_group_name = azurerm_signalr_service.test.resource_group_name

  sku {
    name     = "Free_F1"
    capacity = 1
  }
}
`, testAccAzureRMSignalRService_basic(data))
}

func testAccAzureRMSignalRService_standardWithCapacity(data acceptance.TestData, capacity int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Standard_S1"
    capacity = %d
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, capacity)
}

func testAccAzureRMSignalRService_withCors(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Free_F1"
    capacity = 1
  }

  cors {
    allowed_origins = [
      "https://example.com",
      "https://contoso.com",
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMSignalRService_withServiceMode(data acceptance.TestData, serviceMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Free_F1"
    capacity = 1
  }

  features {
    flag  = "ServiceMode"
    value = "%s"
  }

  features {
    flag  = "EnableConnectivityLogs"
    value = "False"
  }

  features {
    flag  = "EnableMessagingLogs"
    value = "False"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, serviceMode)
}

func testCheckAzureRMSignalRServiceDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).SignalR.Client
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_signalr_service" {
			continue
		}

		id, err := parse.SignalRServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil
		}
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("SignalR service still exists:\n%#v", resp)
		}
	}
	return nil
}

func testCheckAzureRMSignalRServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).SignalR.Client
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parse.SignalRServiceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("Bad: Get on signalRClient: %+v", err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: SignalR service %q (resource group: %q) does not exist", id.Name, id.ResourceGroup)
		}

		return nil
	}
}
