package networkfunction_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/azuretrafficcollectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkFunctionAzureTrafficCollectorResource struct{}

func TestAccNetworkFunctionAzureTrafficCollector_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_azure_traffic_collector", "test")
	r := NetworkFunctionAzureTrafficCollectorResource{}
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

func TestAccNetworkFunctionAzureTrafficCollector_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_azure_traffic_collector", "test")
	r := NetworkFunctionAzureTrafficCollectorResource{}
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

func TestAccNetworkFunctionAzureTrafficCollector_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_azure_traffic_collector", "test")
	r := NetworkFunctionAzureTrafficCollectorResource{}
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

func TestAccNetworkFunctionAzureTrafficCollector_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_azure_traffic_collector", "test")
	r := NetworkFunctionAzureTrafficCollectorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r NetworkFunctionAzureTrafficCollectorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := azuretrafficcollectors.ParseAzureTrafficCollectorID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.NetworkFunction.AzureTrafficCollectorsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NetworkFunctionAzureTrafficCollectorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkFunctionAzureTrafficCollectorResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_function_azure_traffic_collector" "test" {
  name                = "acctest-nfatc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkFunctionAzureTrafficCollectorResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_azure_traffic_collector" "import" {
  name                = azurerm_network_function_azure_traffic_collector.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
}
`, config, data.Locations.Primary)
}

func (r NetworkFunctionAzureTrafficCollectorResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_azure_traffic_collector" "test" {
  name                = "acctest-nfatc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkFunctionAzureTrafficCollectorResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_azure_traffic_collector" "test" {
  name                = "acctest-nfatc-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  tags = {
    key = "value2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
