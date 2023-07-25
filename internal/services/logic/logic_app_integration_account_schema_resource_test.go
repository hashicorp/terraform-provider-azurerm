// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountschemas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppIntegrationAccountSchemaResource struct{}

func TestAccLogicAppIntegrationAccountSchema_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_schema", "test")
	r := LogicAppIntegrationAccountSchemaResource{}

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

func TestAccLogicAppIntegrationAccountSchema_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_schema", "test")
	r := LogicAppIntegrationAccountSchemaResource{}

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

func TestAccLogicAppIntegrationAccountSchema_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_schema", "test")
	r := LogicAppIntegrationAccountSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("file_name", "content"), // not returned from the API
	})
}

func TestAccLogicAppIntegrationAccountSchema_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_schema", "test")
	r := LogicAppIntegrationAccountSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("file_name", "content"), // not returned from the API
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("file_name", "content"), // not returned from the API
	})
}

func (LogicAppIntegrationAccountSchemaResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := integrationaccountschemas.ParseSchemaID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Logic.IntegrationAccountSchemaClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r LogicAppIntegrationAccountSchemaResource) template(data acceptance.TestData) string {
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
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogicAppIntegrationAccountSchemaResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_schema" "test" {
  name                     = "acctest-iaschema-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  content                  = file("testdata/integration_account_schema_content.xsd")
}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppIntegrationAccountSchemaResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_schema" "test" {
  name                     = "acctest-iaschema-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  content                  = file("testdata/integration_account_schema_content.xsd")
  file_name                = "TestFile.xsd"

  metadata = <<METADATA
    {
        "foo": "bar"
    }
METADATA
}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppIntegrationAccountSchemaResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_schema" "test" {
  name                     = "acctest-iaschema-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  content                  = file("testdata/integration_account_schema_content2.xsd")
  file_name                = "TestFile2.xsd"

  metadata = <<METADATA
    {
        "foo": "bar2"
    }
METADATA
}
`, r.template(data), data.RandomInteger)
}

func (r LogicAppIntegrationAccountSchemaResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_schema" "import" {
  name                     = azurerm_logic_app_integration_account_schema.test.name
  resource_group_name      = azurerm_logic_app_integration_account_schema.test.resource_group_name
  integration_account_name = azurerm_logic_app_integration_account_schema.test.integration_account_name
  content                  = azurerm_logic_app_integration_account_schema.test.content
}
`, r.basic(data))
}
