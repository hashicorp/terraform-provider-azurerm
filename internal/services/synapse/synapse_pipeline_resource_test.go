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

type PipelineResource struct {
}

func TestAccSynapsePipeline_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_pipeline", "test")
	r := PipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapsePipeline_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_pipeline", "test")
	r := PipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.update1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapsePipeline_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_pipeline", "test")
	r := PipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapsePipeline_activities(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_pipeline", "test")
	r := PipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.activities(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.activitiesUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.activities(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapsePipeline_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_pipeline", "test")
	r := PipelineResource{}

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

func (t PipelineResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PipelineID(state.ID)
	if err != nil {
		return nil, err
	}

	environment := client.Account.Environment
	pipelinesClient, err := client.Synapse.PipelinesClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return nil, err
	}
	resp, err := pipelinesClient.GetPipeline(ctx, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r PipelineResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_synapse_pipeline" "test" {
  name                 = "acctest-SynapsePipeline%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r PipelineResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_synapse_pipeline" "test" {
  name                 = "acctest-SynapsePipeline%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id

  parameters = {
    test = "testparameter"
  }

  variables = {
    foo = "test1"
    bar = "test2"
  }

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r PipelineResource) update1(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_synapse_pipeline" "test" {
  name                 = "acctest-SynapsePipeline%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  annotations          = ["test1", "test2", "test3"]
  description          = "test description"

  parameters = {
    test = "testparameter"
  }

  variables = {
    foo = "test1"
    bar = "test2"
  }

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r PipelineResource) update2(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_synapse_pipeline" "test" {
  name                 = "acctest-SynapsePipeline%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  annotations          = ["test1", "test2"]
  description          = "test description2"

  parameters = {
    test  = "testparameter"
    test2 = "testparameter2"
  }

  variables = {
    foo = "test1"
    bar = "test2"
    baz = "test3"
  }

  activities_json = <<JSON
[
  {
    "name": "Append variable1",
    "type": "AppendVariable",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
      "variableName": "foo",
      "value": "something"
    }
  },
  {
    "name": "Append variable2",
    "type": "AppendVariable",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
      "variableName": "bar",
      "value": "something1"
    }
  }
]
JSON

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r PipelineResource) activities(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_synapse_pipeline" "test" {
  name                 = "acctest-SynapsePipeline%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  variables = {
    "bob" = "item1"
  }
  activities_json = <<JSON
[
  {
    "name": "Append variable1",
    "type": "AppendVariable",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
      "variableName": "bob",
      "value": "something"
    }
  }
]
JSON

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r PipelineResource) activitiesUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_synapse_pipeline" "test" {
  name                 = "acctest-SynapsePipeline%d"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  variables = {
    "bob" = "item1"
  }
  activities_json = <<JSON
[
  {
    "name": "Append variable1",
    "type": "AppendVariable",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
      "variableName": "bob",
      "value": "something"
    }
  }
]
JSON

  depends_on = [azurerm_synapse_firewall_rule.test]
}
`, template, data.RandomInteger)
}

func (r PipelineResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
	%s

resource "azurerm_synapse_pipeline" "import" {
  name                 = azurerm_synapse_pipeline.test.name
  synapse_workspace_id = azurerm_synapse_pipeline.test.synapse_workspace_id
}
`, config)
}

func (r PipelineResource) template(data acceptance.TestData) string {
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
  managed_virtual_network_enabled      = true
}

resource "azurerm_synapse_firewall_rule" "test" {
  name                 = "AllowAll"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}
