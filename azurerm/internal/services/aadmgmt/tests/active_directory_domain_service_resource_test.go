package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMActiveDirectoryDomainService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_active_directory_domain_service", "test")
	dnsSuffix := acctest.RandString(7)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryDomainServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryDomainService_basic(data, dnsSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryDomainServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "domain_controller_ip_addresses.#", "2"),
				),
			},
			data.ImportStep("ldaps.0.pfx_certificate", "ldaps.0.pfx_certificate_password"),
			{
				Config:      testAccAzureRMActiveDirectoryDomainService_requiresImport(data, dnsSuffix),
				ExpectError: acceptance.RequiresImportError("azurerm_active_directory_domain_service"),
			},
		},
	})
}

func TestAccAzureRMActiveDirectoryDomainService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_active_directory_domain_service", "test")
	dnsSuffix := acctest.RandString(7)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryDomainServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryDomainService_complete(data, dnsSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryDomainServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "domain_controller_ip_address.#", "2"),
				),
			},
			data.ImportStep("ldaps.0.pfx_certificate", "ldaps.0.pfx_certificate_password"),
		},
	})
}

func TestAccAzureRMActiveDirectoryDomainService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_active_directory_domain_service", "test")
	dnsSuffix := acctest.RandString(7)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryDomainServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryDomainService_basic(data, dnsSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryDomainServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "domain_controller_ip_addresses.#", "2"),
				),
			},
			data.ImportStep("ldaps.0.pfx_certificate", "ldaps.0.pfx_certificate_password"),
			{
				Config: testAccAzureRMActiveDirectoryDomainService_complete(data, dnsSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryDomainServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "domain_controller_ip_address.#", "2"),
				),
			},
			data.ImportStep("ldaps.0.pfx_certificate", "ldaps.0.pfx_certificate_password"),
			{
				Config: testAccAzureRMActiveDirectoryDomainService_basic(data, dnsSuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryDomainServiceExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "domain_controller_ip_addresses.#", "2"),
				),
			},
			data.ImportStep("ldaps.0.pfx_certificate", "ldaps.0.pfx_certificate_password"),
		},
	})
}

func testCheckAzureRMActiveDirectoryDomainServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).AadMgmt.DomainServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on DomainServicesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Domain Service %q (Resource Group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMActiveDirectoryDomainServiceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).AadMgmt.DomainServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_active_directory_domain_service" {
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

func testAccAzureRMActiveDirectoryDomainService_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aadds-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestVnet-aadds-%[1]d"
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  lifecycle {
    ignore_changes = [dns_servers]
  }
}

resource "azurerm_subnet" "test" {
  name                 = "acctestSubnet-aadds-%[1]d"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestNSG-aadds-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                       = "AllowSyncWithAzureAD"
    priority                   = 101
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "443"
    source_address_prefix      = "AzureActiveDirectoryDomainServices"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowRD"
    priority                   = 201
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "3389"
    source_address_prefix      = "CorpNetSaw"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowPSRemoting"
    priority                   = 301
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "5986"
    source_address_prefix      = "AzureActiveDirectoryDomainServices"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "AllowLDAPS"
    priority                   = 401
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "636"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource azurerm_subnet_network_security_group_association "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMActiveDirectoryDomainService_basic(data acceptance.TestData, dnsSuffix string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_active_directory_domain_service" "test" {
  name                = "acctest-%s.onmicrosoft.com"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
}
`, testAccAzureRMActiveDirectoryDomainService_template(data), dnsSuffix)
}

func testAccAzureRMActiveDirectoryDomainService_requiresImport(data acceptance.TestData, dnsSuffix string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_active_directory_domain_service" "test2" {
  name                = azurerm_active_directory_domain_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
  filtered_sync       = false
}
`, testAccAzureRMActiveDirectoryDomainService_basic(data, dnsSuffix))
}

func testAccAzureRMActiveDirectoryDomainService_complete(data acceptance.TestData, dnsSuffix string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_active_directory_domain_service" "test" {
  name                = "acctest-%s.onmicrosoft.com"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test.id
  filtered_sync       = false

  security {
    ntlm_v1             = true
    tls_v1              = true
    sync_ntlm_passwords = true
  }

  ldaps {
    external_access          = true
    ldaps                    = true
    pfx_certificate          = "TODO Generate a dummy pfx key+cert (https://docs.microsoft.com/en-us/azure/active-directory-domain-services/tutorial-configure-ldaps)"
    pfx_certificate_password = "test"
  }
}
`, testAccAzureRMActiveDirectoryDomainService_template(data), dnsSuffix)
}
