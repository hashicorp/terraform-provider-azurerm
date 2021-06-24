package apimanagement_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementApiSchemaResource struct {
}

func TestAccApiManagementApiSchema_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_schema", "test")
	r := ApiManagementApiSchemaResource{}
	schema, _ := ioutil.ReadFile("testdata/api_management_api_schema.xml")

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
	schema, _ := ioutil.ReadFile("testdata/api_management_api_schema_swagger.json")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicSwagger(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("value").HasValue(strings.TrimRight(string(schema), "\r\n")),
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

func (ApiManagementApiSchemaResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiName := id.Path["apis"]
	schemaID := id.Path["schemas"]

	resp, err := clients.ApiManagement.ApiSchemasClient.Get(ctx, resourceGroup, serviceName, apiName, schemaID)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagementApi Schema (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
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
  content_type        = "application/vnd.ms-azure-apim.swagger.definitions+json"
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
