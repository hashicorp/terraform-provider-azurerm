package networkfunction_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/networkfunction/2022-11-01/collectorpolicies"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetworkFunctionCollectorPolicyResource struct{}

func TestAccNetworkFunctionCollectorPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
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

func TestAccNetworkFunctionCollectorPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
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

func TestAccNetworkFunctionCollectorPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
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

func TestAccNetworkFunctionCollectorPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_function_collector_policy", "test")
	r := NetworkFunctionCollectorPolicyResource{}
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

func (r NetworkFunctionCollectorPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := collectorpolicies.ParseCollectorPolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.NetworkFunction.CollectorPoliciesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NetworkFunctionCollectorPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_network_function_azure_traffic_collector" "test" {
  name                = "acctest-nfatc-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r NetworkFunctionCollectorPolicyResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_network_function_collector_policy" "test" {
  name                                        = "acctest-nfcp-%d"
  network_function_azure_traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id
  location                                    = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkFunctionCollectorPolicyResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_collector_policy" "import" {
  name                                        = azurerm_network_function_collector_policy.test.name
  network_function_azure_traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id
  location                                    = "%s"
}
`, config, data.Locations.Primary)
}

func (r NetworkFunctionCollectorPolicyResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_collector_policy" "test" {
  name                                        = "acctest-nfcp-%d"
  network_function_azure_traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id
  location                                    = "%s"
  ingestion_policy {
    ingestion_sources {
      resource_id = ""
    }
  }
  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkFunctionCollectorPolicyResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_network_function_collector_policy" "test" {
  name                                        = "acctest-nfcp-%d"
  network_function_azure_traffic_collector_id = azurerm_network_function_azure_traffic_collector.test.id
  location                                    = "%s"
  ingestion_policy {
    ingestion_sources {
      resource_id = ""
    }
  }
  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
