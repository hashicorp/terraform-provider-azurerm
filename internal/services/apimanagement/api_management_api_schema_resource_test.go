// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apischema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiSchemaResource struct{}

func TestAccApiManagementApiSchema_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_schema", "test")
	r := ApiManagementApiSchemaResource{}
	schema, _ := os.ReadFile("testdata/api_management_api_schema.xml")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue(string(schema)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiSchema_basicSwagger(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_schema", "test")
	r := ApiManagementApiSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSwagger(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiSchema_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_schema", "test")
	r := ApiManagementApiSchemaResource{}

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

func TestAccApiManagementApiSchema_components(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_schema", "test")
	r := ApiManagementApiSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.components(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("components").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementApiSchema_definitions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_schema", "test")
	r := ApiManagementApiSchemaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.definitionsJson(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementApiSchemaResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apischema.ParseApiSchemaID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiSchemasClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (r ApiManagementApiSchemaResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_schema" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management_api.test.api_management_name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  schema_id           = "acctestSchema%d"
  content_type        = "application/vnd.ms-azure-apim.xsd+xml"
  value               = file("testdata/api_management_api_schema.xml")
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiSchemaResource) basicSwagger(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_schema" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management_api.test.api_management_name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  schema_id           = "acctestSchema%d"
  content_type        = "application/json"
  value               = file("testdata/api_management_api_schema_swagger.json")
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiSchemaResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_schema" "import" {
  api_name            = azurerm_api_management_api_schema.test.api_name
  api_management_name = azurerm_api_management_api_schema.test.api_management_name
  resource_group_name = azurerm_api_management_api_schema.test.resource_group_name
  schema_id           = azurerm_api_management_api_schema.test.schema_id
  content_type        = azurerm_api_management_api_schema.test.content_type
  value               = azurerm_api_management_api_schema.test.value
}
`, r.basic(data))
}

func (ApiManagementApiSchemaResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiSchemaResource) components(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_schema" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management_api.test.api_management_name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  schema_id           = "acctestSchema%d"
  content_type        = "application/vnd.oai.openapi.components+json"
  components          = file("testdata/api_management_api_schema_swagger.json")
}
`, r.template(data), data.RandomInteger)
}

func (r ApiManagementApiSchemaResource) definitionsJson(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_schema" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management_api.test.api_management_name
  resource_group_name = azurerm_api_management_api.test.resource_group_name
  schema_id           = "acctestSchema%d"
  content_type        = "application/vnd.ms-azure-apim.swagger.definitions+json"
  definitions         = file("testdata/api_management_api_swagger_definitions.json")
}
`, r.template(data), data.RandomInteger)
}
