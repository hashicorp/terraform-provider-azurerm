package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StreamAnalyticsManagedPrivateEndpointResource struct{}

func TestAccStreamAnalyticsManagedPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_managed_private_endpoint", "test")
	r := StreamAnalyticsManagedPrivateEndpointResource{}

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

func TestAccStreamAnalyticsManagedPrivateEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_managed_private_endpoint", "test")
	r := StreamAnalyticsManagedPrivateEndpointResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStreamAnalyticsManagedPrivateEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_managed_private_endpoint", "test")
	r := StreamAnalyticsManagedPrivateEndpointResource{}

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

func (r StreamAnalyticsManagedPrivateEndpointResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PrivateEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.StreamAnalytics.EndpointsClient.Get(ctx, id.ResourceGroup, id.ClusterName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r StreamAnalyticsManagedPrivateEndpointResource) basic(data acceptance.TestData) string {
	//template := r.template(data)
	return fmt.Sprintf(`

resource "azurerm_stream_analytics_managed_private_endpoint" "test" {
  name                          = "acctestprivate_endpoint-%d"
  resource_group_name           = "acctest-sw-31121"
  stream_analytics_cluster_name = "acctestcluster-31121"
  target_resource_id = "/subscriptions/1a6092a6-137e-4025-9a7c-ef77f76f2c02/resourceGroups/acctest-sw-31121/providers/Microsoft.Storage/storageAccounts/examplestorageacc31121"
  subresource_name  = "blob"
}

`, data.RandomInteger)
}

func (r StreamAnalyticsManagedPrivateEndpointResource) updated(data acceptance.TestData) string {
	//template := r.template(data)
	return fmt.Sprintf(`
resource "azurerm_stream_analytics_managed_private_endpoint" "test2" {
  name                          = "acctestprivate_endpoint31121-%d"
  resource_group_name           = "acctest-sw-31121"
  stream_analytics_cluster_name = "acctestcluster-31121"
  target_resource_id = "/subscriptions/1a6092a6-137e-4025-9a7c-ef77f76f2c02/resourceGroups/acctest-sw-31121/providers/Microsoft.EventHub/namespaces/acceptanceTestEventHubNamespace31121"
  subresource_name  = "namespace"
}

`, data.RandomInteger)
}

func (r StreamAnalyticsManagedPrivateEndpointResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_managed_private_endpoint" "import" {
  name                = azurerm_stream_analytics_managed_private_endpoint.test.name
  resource_group_name = azurerm_stream_analytics_managed_private_endpoint.test.resource_group_name
  location            = azurerm_stream_analytics_managed_private_endpoint.test.location
  streaming_capacity  = azurerm_stream_analytics_managed_private_endpoint.test.streaming_capacity
}
`, template)
}

func (r StreamAnalyticsManagedPrivateEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
