package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPointToSiteVPNGateway_basic(t *testing.T) {
	resourceName := "azurerm_point_to_site_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPointToSiteVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAzureRMPointToSiteVPNGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPointToSiteVPNGatewayExists(resourceName),
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

func TestAccAzureRMPointToSiteVPNGateway_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_point_to_site_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPointToSiteVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAzureRMPointToSiteVPNGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPointToSiteVPNGatewayExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAzureRMPointToSiteVPNGateway_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_point_to_site_vpn_gateway"),
			},
		},
	})
}

func TestAccAzureRMPointToSiteVPNGateway_update(t *testing.T) {
	resourceName := "azurerm_point_to_site_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPointToSiteVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAzureRMPointToSiteVPNGateway_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPointToSiteVPNGatewayExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMAzureRMPointToSiteVPNGateway_updated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPointToSiteVPNGatewayExists(resourceName),
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

func TestAccAzureRMPointToSiteVPNGateway_tags(t *testing.T) {
	resourceName := "azurerm_point_to_site_vpn_gateway.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPointToSiteVPNGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAzureRMPointToSiteVPNGateway_tags(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPointToSiteVPNGatewayExists(resourceName),
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

func testCheckAzureRMPointToSiteVPNGatewayExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("VPN Gateway Server Configuration not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).Network.PointToSiteVpnGatewaysClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: VPN Gateway %q (Resource Group %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on network.PointToSiteVpnGatewaysClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPointToSiteVPNGatewayDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_vpn_server_configuration" {
			continue
		}

		client := testAccProvider.Meta().(*ArmClient).Network.PointToSiteVpnGatewaysClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on network.PointToSiteVpnGatewaysClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMAzureRMPointToSiteVPNGateway_basic(rInt int, location string) string {
	template := testAccAzureRMAzureRMPointToSiteVPNGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_point_to_site_vpn_gateway" "test" {
  name                        = "acctestp2sVPNG-%d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  virtual_hub_id              = azurerm_virtual_hub.test.id
  vpn_server_configuration_id = azurerm_vpn_server_configuration.test.id
  scale_unit                  = 1

  connection_configuration {
    name = "first"
    vpn_client_address_pool {
      address_prefixes = [ "172.100.0.0/14" ]
    }
  }
}
`, template, rInt)
}

func testAccAzureRMAzureRMPointToSiteVPNGateway_updated(rInt int, location string) string {
	template := testAccAzureRMAzureRMPointToSiteVPNGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_point_to_site_vpn_gateway" "test" {
  name                        = "acctestp2sVPNG-%d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  virtual_hub_id              = azurerm_virtual_hub.test.id
  vpn_server_configuration_id = azurerm_vpn_server_configuration.test.id
  scale_unit                  = 2

  connection_configuration {
    name = "first"
    vpn_client_address_pool {
      address_prefixes = [ "172.100.0.0/14", "10.100.0.0/14" ]
    }
  }
}
`, template, rInt)
}

func testAccAzureRMAzureRMPointToSiteVPNGateway_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAzureRMPointToSiteVPNGateway_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_point_to_site_vpn_gateway" "import" {
  name                        = azurerm_point_to_site_vpn_gateway.test.name
  location                    = azurerm_point_to_site_vpn_gateway.test.location
  resource_group_name         = azurerm_point_to_site_vpn_gateway.test.resource_group_name
  virtual_hub_id              = azurerm_point_to_site_vpn_gateway.test.virtual_hub_id
  vpn_server_configuration_id = azurerm_point_to_site_vpn_gateway.test.vpn_server_configuration_id
  scale_unit                  = azurerm_point_to_site_vpn_gateway.test.scale_unit

  connection_configuration {
    name = "first"
    vpn_client_address_pool {
      address_prefixes = [ "172.100.0.0/14" ]
    }
  }
}
`, template)
}

func testAccAzureRMAzureRMPointToSiteVPNGateway_tags(rInt int, location string) string {
	template := testAccAzureRMAzureRMPointToSiteVPNGateway_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_point_to_site_vpn_gateway" "test" {
  name                        = "acctestp2sVPNG-%d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  virtual_hub_id              = azurerm_virtual_hub.test.id
  vpn_server_configuration_id = azurerm_vpn_server_configuration.test.id
  scale_unit                  = 1

  connection_configuration {
    name = "first"
    vpn_client_address_pool {
      address_prefixes = [ "172.100.0.0/14" ]
    }
  }

  tags = {
    Hello = "World"
  }
}
`, template, rInt)
}

func testAccAzureRMAzureRMPointToSiteVPNGateway_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctestvhub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_vpn_server_configuration" "test" {
  name                     = "acctestvpnsc-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  vpn_authentication_types = ["Certificate"]

  client_root_certificate {
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
`, rInt, location, rInt, rInt, rInt)
}
