package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVpnSite_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_site", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVpnSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVpnSite_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnSiteExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVpnSite_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_site", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVpnSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVpnSite_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnSiteExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVpnSite_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_site", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVpnSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVpnSite_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnSiteExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVpnSite_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnSiteExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVpnSite_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnSiteExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVpnSite_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vpn_site", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVpnSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVpnSite_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVpnSiteExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMVpnSite_requiresImport),
		},
	})
}

func testCheckAzureRMVpnSiteExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VpnSitesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Vpn Site not found: %s", resourceName)
		}

		id, err := parse.VpnSiteID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Vpn Site %q (Resource Group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Getting on Network.VpnSites: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVpnSiteDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.VpnSitesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_vpn_site" {
			continue
		}

		id, err := parse.VpnSiteID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err == nil {
			return fmt.Errorf("Network.VpnSites still exists")
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Getting on Network.VpnSites: %+v", err)
		}
		return nil
	}

	return nil
}

func testAccAzureRMVpnSite_basic(data acceptance.TestData) string {
	template := testAccAzureRMVpnSite_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_site" "test" {
  name                = "acctest-VpnSite-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_wan_id      = azurerm_virtual_wan.test.id
  link {
    name       = "link1"
    ip_address = "10.0.0.1"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVpnSite_complete(data acceptance.TestData) string {
	template := testAccAzureRMVpnSite_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_site" "test" {
  name                = "acctest-VpnSite-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_cidrs       = ["10.0.0.0/24", "10.0.1.0/24"]

  device_vendor = "Cisco"
  device_model  = "foobar"

  link {
    name          = "link1"
    provider_name = "Verizon"
    speed_in_mbps = 50
    ip_address    = "10.0.0.1"
    bgp {
      asn             = 12345
      peering_address = "10.0.0.1"
    }
  }

  link {
    name = "link2"
    fqdn = "foo.com"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVpnSite_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMVpnSite_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_vpn_site" "import" {
  name                = "acctest-VpnSite-%d"
  location            = azurerm_vpn_site.test.location
  resource_group_name = azurerm_vpn_site.test.resource_group_name
  virtual_wan_id      = azurerm_vpn_site.test.virtual_wan_id
  link {
    name       = "link1"
    ip_address = "10.0.0.1"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMVpnSite_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}


resource "azurerm_virtual_wan" "test" {
  name                = "acctest-vwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
