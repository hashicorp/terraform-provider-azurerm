// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LinkedServiceResource struct{}

func TestAccSynapseLinkedService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_linked_service", "test")
	r := LinkedServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("type_properties_json"),
	})
}

func TestAccSynapseLinkedService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_linked_service", "test")
	r := LinkedServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSynapseLinkedService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_linked_service", "test")
	r := LinkedServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("type_properties_json"),
	})
}

func TestAccSynapseLinkedService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_linked_service", "test")
	r := LinkedServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("type_properties_json"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("type_properties_json"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("type_properties_json"),
	})
}

func TestAccSynapseLinkedService_web(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_linked_service", "test")
	r := LinkedServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.web(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("type_properties_json"),
	})
}

func TestAccSynapseLinkedService_search(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_linked_service", "test")
	r := LinkedServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.search(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("type_properties_json"),
	})
}

func (t LinkedServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LinkedServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	suffix, ok := clients.Account.Environment.Synapse.DomainSuffix()
	if !ok {
		return nil, fmt.Errorf("could not determine Synapse domain suffix for environment %q", clients.Account.Environment.Name)
	}

	client, err := clients.Synapse.LinkedServiceClient(id.WorkspaceName, *suffix)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetLinkedService(ctx, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r LinkedServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_linked_service" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  type                 = "AzureBlobStorage"
  type_properties_json = <<JSON
{
  "connectionString": "${azurerm_storage_account.test.primary_connection_string}"
}
JSON

  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r LinkedServiceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_linked_service" "import" {
  name                 = azurerm_synapse_linked_service.test.name
  synapse_workspace_id = azurerm_synapse_linked_service.test.synapse_workspace_id
  type                 = azurerm_synapse_linked_service.test.type
  type_properties_json = azurerm_synapse_linked_service.test.type_properties_json
}
`, r.basic(data))
}

func (r LinkedServiceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_linked_service" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  type                 = "AzureBlobStorage"
  description          = "test description"
  type_properties_json = <<JSON
{
  "connectionString":"${azurerm_storage_account.test.primary_connection_string}"
}
JSON

  additional_properties = {
    foo = "test1"
    bar = "test2"
  }

  annotations = [
    "test1",
    "test2",
    "test3"
  ]

  parameters = {
    "foo" : "bar"
    "Env" : "Test"
  }

  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r LinkedServiceResource) web(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_linked_service" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  type                 = "Web"
  type_properties_json = <<JSON
{
  "authenticationType": "Anonymous",
  "url": "http://www.bing.com"
}
JSON

  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r LinkedServiceResource) search(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_search_service" "test" {
  name                = "acctestsearchservice%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}

resource "azurerm_synapse_linked_service" "test" {
  name                 = "acctestls%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  type                 = "AzureSearch"
  type_properties_json = <<JSON
{
  "url": "https://${azurerm_search_service.test.name}.search.windows.net",
  "key": {
    "type": "SecureString",
    "value": "${azurerm_search_service.test.primary_key}"
  }
}
JSON

  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (LinkedServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
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
  managed_virtual_network_enabled      = true
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_firewall_rule" "test" {
  name                 = "allowAll"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}
