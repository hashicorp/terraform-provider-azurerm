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

func TestAccAzureRMVirtualHubSecurityPartnerProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_security_partner_provider", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubSecurityPartnerProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubSecurityPartnerProvider_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubSecurityPartnerProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubSecurityPartnerProvider_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_security_partner_provider", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubSecurityPartnerProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubSecurityPartnerProvider_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubSecurityPartnerProviderExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMVirtualHubSecurityPartnerProvider_requiresImport),
		},
	})
}

func TestAccAzureRMVirtualHubSecurityPartnerProvider_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_security_partner_provider", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubSecurityPartnerProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubSecurityPartnerProvider_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubSecurityPartnerProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMVirtualHubSecurityPartnerProvider_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_security_partner_provider", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMVirtualHubSecurityPartnerProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualHubSecurityPartnerProvider_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubSecurityPartnerProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMVirtualHubSecurityPartnerProvider_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualHubSecurityPartnerProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMVirtualHubSecurityPartnerProviderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityPartnerProviderClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Security Partner Provider not found: %s", resourceName)
		}

		id, err := parse.SecurityPartnerProviderID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Security Partner Provider %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Network.VirtualHubSecurityPartnerProviderClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMVirtualHubSecurityPartnerProviderDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityPartnerProviderClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_hub_security_partner_provider" {
			continue
		}

		id, err := parse.SecurityPartnerProviderID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Network.SecurityPartnerProviderClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMVirtualHubSecurityPartnerProvider_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vhub-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-vwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.2.0/24"
}

resource "azurerm_vpn_gateway" "test" {
  name                = "acctest-VPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMVirtualHubSecurityPartnerProvider_basic(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubSecurityPartnerProvider_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_security_partner_provider" "test" {
  name                   = "acctest-SPP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  security_provider_name = "ZScaler"

  depends_on = [azurerm_vpn_gateway.test]
}
`, template, data.RandomInteger)
}

func testAccAzureRMVirtualHubSecurityPartnerProvider_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMVirtualHubSecurityPartnerProvider_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_security_partner_provider" "import" {
  name                   = azurerm_virtual_hub_security_partner_provider.test.name
  resource_group_name    = azurerm_virtual_hub_security_partner_provider.test.resource_group_name
  location               = azurerm_virtual_hub_security_partner_provider.test.location
  security_provider_name = azurerm_virtual_hub_security_partner_provider.test.security_provider_name
}
`, config)
}

func testAccAzureRMVirtualHubSecurityPartnerProvider_complete(data acceptance.TestData) string {
	template := testAccAzureRMVirtualHubSecurityPartnerProvider_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_security_partner_provider" "test" {
  name                   = "acctest-SPP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  virtual_hub_id         = azurerm_virtual_hub.test.id
  security_provider_name = "ZScaler"

  tags = {
    ENv = "Test"
  }

  depends_on = [azurerm_vpn_gateway.test]
}
`, template, data.RandomInteger)
}
