package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LinkedCustomServiceResource struct {
}

func TestAccDataFactoryLinkedCustomService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_custom_service", "test")
	r := LinkedCustomServiceResource{}

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

func TestAccDataFactoryLinkedCustomService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_custom_service", "test")
	r := LinkedCustomServiceResource{}

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

func TestAccDataFactoryLinkedCustomService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_custom_service", "test")
	r := LinkedCustomServiceResource{}

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

func TestAccDataFactoryLinkedCustomService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_custom_service", "test")
	r := LinkedCustomServiceResource{}

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

func TestAccDataFactoryLinkedCustomService_web(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_custom_service", "test")
	r := LinkedCustomServiceResource{}

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

func TestAccDataFactoryLinkedCustomService_search(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_linked_custom_service", "test")
	r := LinkedCustomServiceResource{}

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

func (t LinkedCustomServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.LinkedServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataFactory.LinkedServiceClient.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r LinkedCustomServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_linked_custom_service" "test" {
  name                 = "acctestls%d"
  data_factory_id      = azurerm_data_factory.test.id
  type                 = "AzureBlobStorage"
  type_properties_json = <<JSON
{
  "connectionString": "${azurerm_storage_account.test.primary_connection_string}"
}
JSON
}
`, r.template(data), data.RandomInteger)
}

func (r LinkedCustomServiceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_linked_custom_service" "import" {
  name                 = azurerm_data_factory_linked_custom_service.test.name
  data_factory_id      = azurerm_data_factory_linked_custom_service.test.data_factory_id
  type                 = azurerm_data_factory_linked_custom_service.test.type
  type_properties_json = azurerm_data_factory_linked_custom_service.test.type_properties_json
}
`, r.basic(data))
}

func (r LinkedCustomServiceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_factory_linked_custom_service" "test" {
  name                 = "acctestls%d"
  data_factory_id      = azurerm_data_factory.test.id
  type                 = "AzureBlobStorage"
  description          = "test description"
  type_properties_json = <<JSON
{
  "connectionString":"${azurerm_storage_account.test.primary_connection_string}"
}
JSON

  integration_runtime {
    name = azurerm_data_factory_integration_runtime_managed.test.name
    parameters = {
      "Key" : "value"
    }
  }

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
}
`, r.template(data), data.RandomInteger)
}

func (r LinkedCustomServiceResource) web(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_factory_linked_custom_service" "test" {
  name                 = "acctestls%d"
  data_factory_id      = azurerm_data_factory.test.id
  type                 = "Web"
  type_properties_json = <<JSON
{
  "authenticationType": "Anonymous",
  "url": "http://www.bing.com"
}
JSON
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r LinkedCustomServiceResource) search(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_search_service" "test" {
  name                = "acctestsearchservice%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}

resource "azurerm_data_factory_linked_custom_service" "test" {
  name                 = "acctestls%d"
  data_factory_id      = azurerm_data_factory.test.id
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
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (LinkedCustomServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name                = "acctest-irm%d"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  node_size                        = "Standard_D8_v3"
  number_of_nodes                  = 2
  max_parallel_executions_per_node = 8
  edition                          = "Standard"
  license_type                     = "LicenseIncluded"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString, data.RandomInteger)
}
