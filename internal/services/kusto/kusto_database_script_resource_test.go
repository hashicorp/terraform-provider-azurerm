package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KustoScriptResource struct{}

func TestAccKustoScript_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_script", "test")
	r := KustoScriptResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sas_token"),
	})
}

func TestAccKustoScript_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_script", "test")
	r := KustoScriptResource{}
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

func TestAccKustoScript_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_script", "test")
	r := KustoScriptResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sas_token"),
	})
}

func TestAccKustoScript_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_script", "test")
	r := KustoScriptResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sas_token"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sas_token"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sas_token"),
	})
}

func (r KustoScriptResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ScriptID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Kusto.ScriptsClient.Get(ctx, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r KustoScriptResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-kusto-%[1]d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "setup-files"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name                   = "script.txt"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source_content         = ".create table MyTable (Level:string, Timestamp:datetime, UserId:string, TraceId:string, Message:string, ProcessId:int32)"
}

data "azurerm_storage_account_blob_container_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  container_name    = azurerm_storage_container.test.name
  https_only        = true

  start  = "2017-03-21"
  expiry = "2022-03-21"

  permissions {
    read   = true
    add    = false
    create = false
    write  = true
    delete = false
    list   = true
  }
}
`, data.RandomIntOfLength(12), data.Locations.Primary)
}

func (r KustoScriptResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_script" "test" {
  name        = "acctest-ks-%d"
  database_id = azurerm_kusto_database.test.id
  url         = azurerm_storage_blob.test.id
  sas_token   = data.azurerm_storage_account_blob_container_sas.test.sas
}
`, template, data.RandomInteger)
}

func (r KustoScriptResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_script" "import" {
  name        = azurerm_kusto_script.test.name
  database_id = azurerm_kusto_script.test.database_id
  url         = azurerm_kusto_script.test.url
  sas_token   = azurerm_kusto_script.test.sas_token
}
`, config)
}

func (r KustoScriptResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_script" "test" {
  name                               = "acctest-ks-%d"
  database_id                        = azurerm_kusto_database.test.id
  url                                = azurerm_storage_blob.test.id
  sas_token                          = data.azurerm_storage_account_blob_container_sas.test.sas
  continue_on_errors_enabled         = true
  force_an_update_when_value_changed = "first"
}
`, template, data.RandomInteger)
}
