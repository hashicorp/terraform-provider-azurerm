package tests

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMEventHubNamespaceNetworkRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventHubNamespaceNetworkRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceNetworkRule_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventHubNamespaceNetworkRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMEventHubNamespaceNetworkRule_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMEventHubNamespaceNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMEventHubNamespaceNetworkRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceNetworkRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceNetworkRuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMEventHubNamespaceNetworkRule_requiresImport),
		},
	})
}

func TestAccAzureRMEventHubNamespaceNetworkRule_TrustedServiceAccessEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_network_rule_set", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMEventHubNamespaceNetworkRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMEventHubNamespaceNetworkRuleTrustedServiceAccessEnabled(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMEventHubNamespaceNetworkRuleExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "trusted_service_access_enabled", "true"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMEventHubNamespaceNetworkRule_requiresImport),
		},
	})
}

func testCheckAzureRMEventHubNamespaceNetworkRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.NamespacesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Event Hub Namespace Network Rule Set not found: %s", resourceName)
		}

		id, err := parse.EventHubNamespaceNetworkRuleSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.GetNetworkRuleSet(ctx, id.ResourceGroup, id.NamespaceName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Event Hub Namespace Network Rule Set (Namespace %q / Resource Group %q) does not exist", id.NamespaceName, id.ResourceGroup)
			}
			return fmt.Errorf("failed to GetNetworkRuleSet on EventHub.NamespacesClientPreview: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMEventHubNamespaceNetworkRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Eventhub.NamespacesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_eventhub_namespace_network_rule_set" {
			continue
		}

		id, err := parse.EventHubNamespaceNetworkRuleSetID(rs.Primary.ID)
		if err != nil {
			return err
		}

		// this resource cannot be deleted, instead, we check if this setting was set back to empty
		resp, err := client.GetNetworkRuleSet(ctx, id.ResourceGroup, id.NamespaceName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("failed to GetNetworkRuleSet on EventHub.NamespacesClientPreview: %+v", err)
			}
			return nil
		}

		if !eventhub.CheckNetworkRuleNullified(resp) {
			return fmt.Errorf("the Event Hub Namespace Network Rule Set (Namespace %q / Resource Group %q) still exists", id.NamespaceName, id.ResourceGroup)
		}
	}

	return nil
}

func testAccAzureRMEventHubNamespaceNetworkRule_basic(data acceptance.TestData) string {
	template := testAccAzureRMEventHubNamespaceNetworkRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace_network_rule_set" "test" {
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  default_action = "Deny"

  network_rules {
    subnet_id                            = azurerm_subnet.test.id
    ignore_missing_vnet_service_endpoint = false
  }
}
`, template)
}

func testAccAzureRMEventHubNamespaceNetworkRuleTrustedServiceAccessEnabled(data acceptance.TestData) string {
	template := testAccAzureRMEventHubNamespaceNetworkRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace_network_rule_set" "test" {
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  trusted_service_access_enabled = true
  default_action                 = "Deny"

  network_rules {
    subnet_id                            = azurerm_subnet.test.id
    ignore_missing_vnet_service_endpoint = false
  }
}
`, template)
}

func testAccAzureRMEventHubNamespaceNetworkRule_complete(data acceptance.TestData) string {
	template := testAccAzureRMEventHubNamespaceNetworkRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace_network_rule_set" "test" {
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  default_action = "Deny"

  network_rules {
    subnet_id                            = azurerm_subnet.test.id
    ignore_missing_vnet_service_endpoint = false
  }

  ip_rules = ["1.1.1.1"]
}
`, template)
}

func testAccAzureRMEventHubNamespaceNetworkRule_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sb-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-sb-namespace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"

  capacity = 1
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["172.17.0.0/16"]
  dns_servers         = ["10.0.0.4", "10.0.0.5"]
}

resource "azurerm_subnet" "test" {
  name                   = "${azurerm_virtual_network.test.name}-default"
  resource_group_name    = azurerm_resource_group.test.name
  virtual_network_name   = azurerm_virtual_network.test.name
  address_prefixes       = ["172.17.0.0/24"]

  service_endpoints = ["Microsoft.EventHub"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMEventHubNamespaceNetworkRule_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMEventHubNamespaceNetworkRule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub_namespace_network_rule_set" "import" {
  namespace_name      = azurerm_eventhub_namespace_network_rule_set.test.namespace_name
  resource_group_name = azurerm_eventhub_namespace_network_rule_set.test.resource_group_name
}
`, template)
}
