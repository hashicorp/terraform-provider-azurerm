package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SynapseFirewallRuleResource struct{}

func TestAccSynapseFirewallRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_firewall_rule", "test")
	r := SynapseFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseFirewallRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_firewall_rule", "test")
	r := SynapseFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSynapseFirewallRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_firewall_rule", "test")
	r := SynapseFirewallRuleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withUpdates(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SynapseFirewallRuleResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.FirewallRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Synapse.FirewallRulesClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Synapse Firewall Rule %q (Workspace %q / Resource Group %q): %+v", id.Name, id.WorkspaceName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r SynapseFirewallRuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
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

func (r SynapseFirewallRuleResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
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

func (r SynapseFirewallRuleResource) withUpdates(data acceptance.TestData) string {
	template := r.template(data)
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

func (r SynapseFirewallRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%d"
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
