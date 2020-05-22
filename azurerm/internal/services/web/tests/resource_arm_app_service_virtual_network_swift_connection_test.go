package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO: requires import

func TestAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceVirtualNetworkSwiftConnection_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
				),
			},
			{
				Config: testAccAzureRMAppServiceVirtualNetworkSwiftConnection_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
				),
			},
		},
	})
}

func TestAccAzureRMAppServiceVirtualNetworkSwiftConnection_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_virtual_network_swift_connection", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "subnet_id"),
					testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id := rs.Primary.Attributes["id"]
		parsedID, err := azure.ParseAzureResourceID(id)
		if err != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", id)
		}
		name := parsedID.Path["sites"]
		resourceGroup := parsedID.ResourceGroup

		resp, err := client.GetSwiftVirtualNetworkConnection(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service Virtual Network Association %q (Resource Group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServicesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id := rs.Primary.Attributes["id"]
		parsedID, err := azure.ParseAzureResourceID(id)
		if err != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", id)
		}
		name := parsedID.Path["sites"]
		resourceGroup := parsedID.ResourceGroup

		resp, err := client.DeleteSwiftVirtualNetwork(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp) {
				return fmt.Errorf("Bad: Delete on appServicesClient: %+v", err)
			}
		}

		return nil
	}
}

func testCheckAzureRMAppServiceVirtualNetworkSwiftConnectionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_virtual_network_swift_connection" {
			continue
		}

		id := rs.Primary.Attributes["id"]
		parsedID, err := azure.ParseAzureResourceID(id)
		if err != nil {
			return fmt.Errorf("Error parsing Azure Resource ID %q", id)
		}
		name := parsedID.Path["sites"]
		resourceGroup := parsedID.ResourceGroup

		resp, err := client.GetSwiftVirtualNetworkConnection(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testAccAzureRMAppServiceVirtualNetworkSwiftConnection_base(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appservice-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VNET-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "%s"
  resource_group_name = azurerm_resource_group.test.name
  lifecycle {
    ignore_changes = [ddos_protection_plan]
  }
}

resource "azurerm_subnet" "test" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest-IP-%d"
  location            = azurerm_virtual_network.test.location
  resource_group_name = azurerm_resource_group.test.name

  allocation_method = "Dynamic"
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctest-VNETGW-%d"
  location            = azurerm_virtual_network.test.location
  resource_group_name = azurerm_resource_group.test.name

  type     = "Vpn"
  vpn_type = "RouteBased"

  active_active = false
  enable_bgp    = false
  sku           = "VpnGw1"

  ip_configuration {
    name                          = "vnetGatewayConfig"
    public_ip_address_id          = azurerm_public_ip.test.id
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.test.id
  }

  vpn_client_configuration {
    address_space = ["192.168.0.96/29"]

    vpn_client_protocols = ["SSTP"]

    root_certificate {
      name = "DigiCert-Federated-ID-Root-CA"

      public_cert_data = <<EOF
	MIIDuzCCAqOgAwIBAgIQCHTZWCM+IlfFIRXIvyKSrjANBgkqhkiG9w0BAQsFADBn
	MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
	d3cuZGlnaWNlcnQuY29tMSYwJAYDVQQDEx1EaWdpQ2VydCBGZWRlcmF0ZWQgSUQg
	Um9vdCBDQTAeFw0xMzAxMTUxMjAwMDBaFw0zMzAxMTUxMjAwMDBaMGcxCzAJBgNV
	BAYTAlVTMRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdp
	Y2VydC5jb20xJjAkBgNVBAMTHURpZ2lDZXJ0IEZlZGVyYXRlZCBJRCBSb290IENB
	MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvAEB4pcCqnNNOWE6Ur5j
	QPUH+1y1F9KdHTRSza6k5iDlXq1kGS1qAkuKtw9JsiNRrjltmFnzMZRBbX8Tlfl8
	zAhBmb6dDduDGED01kBsTkgywYPxXVTKec0WxYEEF0oMn4wSYNl0lt2eJAKHXjNf
	GTwiibdP8CUR2ghSM2sUTI8Nt1Omfc4SMHhGhYD64uJMbX98THQ/4LMGuYegou+d
	GTiahfHtjn7AboSEknwAMJHCh5RlYZZ6B1O4QbKJ+34Q0eKgnI3X6Vc9u0zf6DH8
	Dk+4zQDYRRTqTnVO3VT8jzqDlCRuNtq6YvryOWN74/dq8LQhUnXHvFyrsdMaE1X2
	DwIDAQABo2MwYTAPBgNVHRMBAf8EBTADAQH/MA4GA1UdDwEB/wQEAwIBhjAdBgNV
	HQ4EFgQUGRdkFnbGt1EWjKwbUne+5OaZvRYwHwYDVR0jBBgwFoAUGRdkFnbGt1EW
	jKwbUne+5OaZvRYwDQYJKoZIhvcNAQELBQADggEBAHcqsHkrjpESqfuVTRiptJfP
	9JbdtWqRTmOf6uJi2c8YVqI6XlKXsD8C1dUUaaHKLUJzvKiazibVuBwMIT84AyqR
	QELn3e0BtgEymEygMU569b01ZPxoFSnNXc7qDZBDef8WfqAV/sxkTi8L9BkmFYfL
	uGLOhRJOFprPdoDIUBB+tmCl3oDcBy3vnUeOEioz8zAkprcb3GHwHAK+vHmmfgcn
	WsfMLH4JCLa/tRYL+Rw/N3ybCkDp00s0WUZ+AoDywSl0Q/ZEnNY0MsFiw6LyIdbq
	M/s/1JRtO3bDSzD9TazRVzn2oBqzSa8VgIo5C1nOnoAKJTlsClJKvIhnRlaLQqk=
	EOF

    }
  }
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctest-ASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctest-AS-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMAppServiceVirtualNetworkSwiftConnection_basic(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceVirtualNetworkSwiftConnection_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_virtual_network_swift_connection" "test" {
  app_service_id = azurerm_app_service.test.id
  subnet_id      = azurerm_subnet.test1.id
}
`, template)
}

func testAccAzureRMAppServiceVirtualNetworkSwiftConnection_update(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceVirtualNetworkSwiftConnection_base(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_virtual_network_swift_connection" "test" {
  app_service_id = azurerm_app_service.test.id
  subnet_id      = azurerm_subnet.test2.id
}
`, template)
}
