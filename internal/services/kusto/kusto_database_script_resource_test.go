// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/scripts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
		data.ImportStep("sas_token", "script_content"),
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
		data.ImportStep("sas_token", "script_content"),
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
		data.ImportStep("sas_token", "script_content"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sas_token", "script_content"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sas_token", "script_content"),
	})
}

func TestAccKustoScript_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_script", "test")
	r := KustoScriptResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(fmt.Sprintf("%s%d", data.ResourceName, 2)).ExistsInAzure(r),
				check.That(fmt.Sprintf("%s%d", data.ResourceName, 3)).ExistsInAzure(r),
				check.That(fmt.Sprintf("%s%d", data.ResourceName, 4)).ExistsInAzure(r),
			),
		},
		data.ImportStep("sas_token", "script_content"),
	})
}

func TestAccKustoScript_scriptContent(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_script", "test")
	r := KustoScriptResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scriptContent(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sas_token", "script_content"),
	})
}

func (r KustoScriptResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := scripts.ParseScriptID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Kusto.ScriptsClient.Get(ctx, *id)
	exists := true

	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			exists = false
		} else {
			return nil, fmt.Errorf("retrieving %q: %+v", id, err)
		}
	}

	return &exists, nil
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
  storage_account_id    = azurerm_storage_account.test.id
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

  start  = "2022-03-21"
  expiry = "2027-03-21"

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

func (r KustoScriptResource) multiple(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kusto_database" "test2" {
  name                = "acctest-kd-2-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_kusto_database" "test3" {
  name                = "acctest-kd-3-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_kusto_database" "test4" {
  name                = "acctest-kd-4-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  cluster_name        = azurerm_kusto_cluster.test.name
}

resource "azurerm_kusto_script" "test" {
  name        = "acctest-ks-%d"
  database_id = azurerm_kusto_database.test.id
  url         = azurerm_storage_blob.test.id
  sas_token   = data.azurerm_storage_account_blob_container_sas.test.sas
}

resource "azurerm_kusto_script" "test2" {
  name        = "acctest-ks-2-%d"
  database_id = azurerm_kusto_database.test2.id
  url         = azurerm_storage_blob.test.id
  sas_token   = data.azurerm_storage_account_blob_container_sas.test.sas
}

resource "azurerm_kusto_script" "test3" {
  name        = "acctest-ks-3-%d"
  database_id = azurerm_kusto_database.test3.id
  url         = azurerm_storage_blob.test.id
  sas_token   = data.azurerm_storage_account_blob_container_sas.test.sas
}

resource "azurerm_kusto_script" "test4" {
  name        = "acctest-ks-4-%d"
  database_id = azurerm_kusto_database.test4.id
  url         = azurerm_storage_blob.test.id
  sas_token   = data.azurerm_storage_account_blob_container_sas.test.sas
}
`, template, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger,
		data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r KustoScriptResource) scriptContent(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_kusto_script" "test" {
  name                               = "acctest-ks-%d"
  database_id                        = azurerm_kusto_database.test.id
  continue_on_errors_enabled         = true
  force_an_update_when_value_changed = "first"
  script_content                     = ".create table MyTable (Level:string, Timestamp:datetime, UserId:string, TraceId:string, Message:string, ProcessId:int32)"
}
`, template, data.RandomInteger)
}
