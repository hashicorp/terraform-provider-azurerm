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

func TestAccAzureRMContainerRegistryNetworkRuleset_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryNetworkRuleset_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryNetworkRulesetExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistryNetworkRuleset_updateRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryNetworkRuleset_basicTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryNetworkRulesetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Deny"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMContainerRegistryNetworkRuleset_updateTemplate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryNetworkRulesetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Allow"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistryNetworkRuleset_addIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryNetworkRuleset_addIp(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryNetworkRulesetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.ip_rule.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMContainerRegistryNetworkRuleset_addVnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMContainerRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMContainerRegistryNetworkRuleset_addVnet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMContainerRegistryNetworkRulesetExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.default_action", "Deny"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_rule_set.0.virtual_network.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMContainerRegistry_common(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_container_registry" "test" {
  name                = "acctestacr%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestRG-%[1]d-network"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctestRG-%[1]d-subnet"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.1.0/24"]
  service_endpoints    = ["Microsoft.ContainerRegistry"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMContainerRegistryNetworkRuleset_basicTemplate(data acceptance.TestData) string {
	template := testAccAzureRMContainerRegistry_common(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_network_rule_set" "test" {
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  container_registry_name = azurerm_container_registry.test.name
  depends_on              = [azurerm_virtual_network.test, azurerm_subnet.test]
  network_rule_set {
    default_action = "Deny"
  }
}
`, template)
}

func testAccAzureRMContainerRegistryNetworkRuleset_updateTemplate(data acceptance.TestData) string {
	template := testAccAzureRMContainerRegistry_common(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_network_rule_set" "test" {
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  container_registry_name = azurerm_container_registry.test.name
  depends_on              = [azurerm_virtual_network.test, azurerm_subnet.test]
  network_rule_set {
    default_action = "Allow"
  }
}
`, template)
}

func testAccAzureRMContainerRegistryNetworkRuleset_addIp(data acceptance.TestData) string {
	template := testAccAzureRMContainerRegistry_common(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_network_rule_set" "test" {
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  container_registry_name = azurerm_container_registry.test.name
  depends_on              = [azurerm_virtual_network.test, azurerm_subnet.test]
  network_rule_set {
    default_action = "Deny"
    ip_rule {
      action   = "Allow"
      ip_range = "43.0.0.0/24"
    }
  }
}
`, template)
}

func testAccAzureRMContainerRegistryNetworkRuleset_addVnet(data acceptance.TestData) string {
	template := testAccAzureRMContainerRegistry_common(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_network_rule_set" "test" {
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  container_registry_name = azurerm_container_registry.test.name
  depends_on              = [azurerm_virtual_network.test, azurerm_subnet.test]
  network_rule_set {
    default_action = "Deny"
    virtual_network {
      action    = "Allow"
      subnet_id = azurerm_subnet.test.id
    }
  }
}
`, template)
}

func testCheckAzureRMContainerRegistryNetworkRulesetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Containers.RegistriesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := azure.ParseAzureResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}

		resourceGroup := id.ResourceGroup
		containerRegistryName := id.Path["registries"]

		containerRegistry, err := client.Get(ctx, resourceGroup, containerRegistryName)
		if err != nil {
			if utils.ResponseWasNotFound(containerRegistry.Response) {
				return fmt.Errorf("Azure Container Registry %q (Resource Group %q) was not found", containerRegistryName, resourceGroup)
			}

			return fmt.Errorf("Error retrieving Azure Container Registry %q (Resource Group %q): %+v", containerRegistryName, resourceGroup, err)
		}

		if rules := containerRegistry.NetworkRuleSet; rules == nil {
			return fmt.Errorf("Network rule set for Azure Container Registry %q (Resource Group %q): %+v does not exist", containerRegistryName, resourceGroup, err)
		}

		return nil
	}
}
