// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountassemblies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppIntegrationAccountAssemblyResource struct{}

func TestAccLogicAppIntegrationAccountAssembly_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_assembly", "test")
	r := LogicAppIntegrationAccountAssemblyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content"), // not returned from the API
	})
}

func TestAccLogicAppIntegrationAccountAssembly_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_assembly", "test")
	r := LogicAppIntegrationAccountAssemblyResource{}

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

func TestAccLogicAppIntegrationAccountAssembly_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_assembly", "test")
	r := LogicAppIntegrationAccountAssemblyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content"), // not returned from the API
	})
}

func TestAccLogicAppIntegrationAccountAssembly_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_assembly", "test")
	r := LogicAppIntegrationAccountAssemblyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content"), // not returned from the API
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("content_link_uri"), // expected value isn't returned from the API
	})
}

func (r LogicAppIntegrationAccountAssemblyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := integrationaccountassemblies.ParseAssemblyID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Logic.IntegrationAccountAssemblyClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r LogicAppIntegrationAccountAssemblyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-ia-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogicAppIntegrationAccountAssemblyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_assembly" "test" {
  name                     = "acctest-assembly-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  assembly_name            = "TestAssembly"
  content                  = filebase64("testdata/log4net.dll")
}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppIntegrationAccountAssemblyResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_assembly" "import" {
  name                     = azurerm_logic_app_integration_account_assembly.test.name
  resource_group_name      = azurerm_logic_app_integration_account_assembly.test.resource_group_name
  integration_account_name = azurerm_logic_app_integration_account_assembly.test.integration_account_name
  assembly_name            = azurerm_logic_app_integration_account_assembly.test.assembly_name
  content                  = azurerm_logic_app_integration_account_assembly.test.content
}
`, config)
}

func (r LogicAppIntegrationAccountAssemblyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_assembly" "test" {
  name                     = "acctest-assembly-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  assembly_name            = "TestAssembly"
  assembly_version         = "1.1.1.1"
  content                  = filebase64("testdata/log4net.dll")

  metadata = {
    foo = "bar"
  }
}
`, template, data.RandomInteger)
}

func (r LogicAppIntegrationAccountAssemblyResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_account" "test" {
  name                            = "acctestsa%s"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestsc%s"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "test" {
  name                   = "log4net.dll"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source                 = "testdata/log4net.dll"
}

resource "azurerm_logic_app_integration_account_assembly" "test" {
  name                     = "acctest-assembly-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  assembly_name            = "TestAssembly2"
  assembly_version         = "2.2.2.2"
  content_link_uri         = azurerm_storage_blob.test.id

  metadata = {
    foo = "bar2"
  }
}
`, template, data.RandomString, data.RandomString, data.RandomInteger)
}
