package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSynapseFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_firewall_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseFirewallRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseFirewallRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSynapseFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_firewall_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseFirewallRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseFirewallRuleExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSynapseFirewallRule_requiresImport),
		},
	})
}

func TestAccAzureRMSynapseFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_firewall_rule", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSynapseFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSynapseFirewallRule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseFirewallRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSynapseFirewallRule_withUpdates(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSynapseFirewallRuleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSynapseFirewallRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Synapse.FirewallRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("synapse Firewall Rule not found: %s", resourceName)
		}
		id, err := parse.SynapseFirewallRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.Workspace.ResourceGroup, id.Workspace.Name, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Synapse Firewall Rule %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Synapse.FirewallRulesClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMSynapseFirewallRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Synapse.FirewallRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_synapse_firewall_rule" {
			continue
		}
		id, err := parse.SynapseFirewallRuleID(rs.Primary.ID)
		if err != nil {
			return err
		}
		resp, err := client.Get(ctx, id.Workspace.ResourceGroup, id.Workspace.Name, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Synapse.FirewallRulesClient: %+v", err)
			}
			return nil
		}
		return fmt.Errorf("expected no Firewall Rule but found %+v", resp)
	}
	return nil
}

func testAccAzureRMSynapseFirewallRule_basic(data acceptance.TestData) string {
	template := testAccAzureRMSynapseFirewallRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_firewall_rule" "test" {
  name                 = "FirewallRule%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}
`, template, data.RandomInteger)
}

func testAccAzureRMSynapseFirewallRule_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMSynapseFirewallRule_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_firewall_rule" "import" {
  name                 = azurerm_synapse_firewall_rule.test.name
  synapse_workspace_id = azurerm_synapse_firewall_rule.test.synapse_workspace_id
  start_ip_address     = azurerm_synapse_firewall_rule.test.start_ip_address
  end_ip_address       = azurerm_synapse_firewall_rule.test.end_ip_address
}
`, config)
}

func testAccAzureRMSynapseFirewallRule_withUpdates(data acceptance.TestData) string {
	template := testAccAzureRMSynapseFirewallRule_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_firewall_rule" "test" {
  name                 = "FirewallRule%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "10.0.17.62"
  end_ip_address       = "10.0.17.62"
}
`, template, data.RandomInteger)
}

func testAccAzureRMSynapseFirewallRule_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-Synapse-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}
