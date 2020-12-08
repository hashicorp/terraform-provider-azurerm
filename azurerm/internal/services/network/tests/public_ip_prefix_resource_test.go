package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func testCheckAzureRMPublicIPPrefixExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PublicIPPrefixesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		publicIpPrefixName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for public ip prefix: %s", publicIpPrefixName)
		}

		resp, err := client.Get(ctx, resourceGroup, publicIpPrefixName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on publicIPPrefixClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Public IP Prefix %q (resource group: %q) does not exist", publicIpPrefixName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPublicIPPrefixDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PublicIPPrefixesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		publicIpPrefixName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for public ip prefix: %s", publicIpPrefixName)
		}

		future, err := client.Delete(ctx, resourceGroup, publicIpPrefixName)
		if err != nil {
			return fmt.Errorf("Error deleting Public IP Prefix %q (Resource Group %q): %+v", publicIpPrefixName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Public IP Prefix %q (Resource Group %q): %+v", publicIpPrefixName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMPublicIPPrefixDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PublicIPPrefixesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_public_ip_prefix" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Public IP Prefix still exists:\n%#v", resp.PublicIPPrefixPropertiesFormat)
		}
	}

	return nil
}

func TestAccAzureRMPublicIpPrefix_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPublicIPPrefix_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_prefix"),
					resource.TestCheckResourceAttr(data.ResourceName, "prefix_length", "28"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPublicIpPrefix_prefixLength31(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPublicIPPrefix_prefixLength31(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_prefix"),
					resource.TestCheckResourceAttr(data.ResourceName, "prefix_length", "31"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPublicIpPrefix_prefixLength24(t *testing.T) {
	// NOTE: This test will fail unless the subscription is updated
	//        to accept a minimum PrefixLength of 24
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPublicIPPrefix_prefixLength24(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_prefix"),
					resource.TestCheckResourceAttr(data.ResourceName, "prefix_length", "24"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPublicIpPrefix_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPublicIPPrefix_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMPublicIPPrefix_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMPublicIpPrefix_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPublicIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMPublicIPPrefix_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(data.ResourceName),
					testCheckAzureRMPublicIPPrefixDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAzureRMPublicIPPrefix_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMPublicIPPrefix_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMPublicIPPrefix_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMPublicIPPrefix_prefixLength31(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  prefix_length = 31
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMPublicIPPrefix_prefixLength24(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  prefix_length = 24
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
