package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDomainService_complete(t *testing.T) {
	resourceName := "azurerm_domain_service.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureDomainServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDomainService_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDomainServiceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "filtered_sync", "Disabled"),
					resource.TestCheckResourceAttr(resourceName, "domain_controller_ip_address.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "ldaps_settings.0.external_access_ip_address"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ldaps_settings.0.pfx_certificate",
					"ldaps_settings.0.pfx_certificate_password",
					"notification_settings",
				},
			},
		},
	})
}

func testCheckAzureDomainServiceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).DomainServices.DomainServicesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_domain_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Domain Service still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMDomainServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Domain Service: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).DomainServices.DomainServicesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Domain Service %q (resource group: %s) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on DomainServicesClient: %s", err)
		}

		return nil
	}
}

func testAccAzureRMDomainService_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestNSG-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                        = "AllowSyncWithAzureAD"
    priority                    = 101
    direction                   = "Inbound"
    access                      = "Allow"
    protocol                    = "Tcp"
    source_port_range           = "*"
    destination_port_range      = "443"
    source_address_prefix       = "AzureActiveDirectoryDomainServices"
    destination_address_prefix  = "*"
  }

  security_rule {
    name                        = "AllowRD"
    priority                    = 201
    direction                   = "Inbound"
    access                      = "Allow"
    protocol                    = "Tcp"
    source_port_range           = "*"
    destination_port_range      = "3389"
    source_address_prefix       = "CorpNetSaw"
    destination_address_prefix  = "*"
  }

  security_rule {
    name                        = "AllowPSRemoting"
    priority                    = 301
    direction                   = "Inbound"
    access                      = "Allow"
    protocol                    = "Tcp"
    source_port_range           = "*"
    destination_port_range      = "5986"
    source_address_prefix       = "AzureActiveDirectoryDomainServices"
    destination_address_prefix  = "*"
  }

  security_rule {
    name                        = "AllowLDAPS"
    priority                    = 401
    direction                   = "Inbound"
    access                      = "Allow"
    protocol                    = "Tcp"
    source_port_range           = "*"
    destination_port_range      = "636"
    source_address_prefix       = "*"
    destination_address_prefix  = "*"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVnet-%d"
  address_space       = ["10.0.1.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  lifecycle {
    ignore_changes = [ dns_servers ]
  }
}

resource "azurerm_subnet" "test" {
  name                 = "acctestSubnet-%d"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_prefix       = "10.0.1.0/24"
  network_security_group_id = "${azurerm_network_security_group.test.id}"
}

resource "azurerm_domain_service" "test" {
  name                  = "test.onmicrosoft.com"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  
  domain_security_settings {
    ntlm_v1 = "Enabled"
    tls_v1 = "Disabled"
    sync_ntlm_passwords = "Enabled"
  }

  ldaps_settings {
    ldaps = "Enabled"
    pfx_certificate = "${filebase64("testdata/domain_service_test.pfx")}"
    pfx_certificate_password = "test"
    external_access = "Enabled"
  }

  subnet_id = "${azurerm_subnet.test.id}"
  filtered_sync = "Disabled"
}
`, rInt, location, rInt, rInt, rInt)
}
