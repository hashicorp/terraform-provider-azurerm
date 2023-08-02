// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StreamAnalyticsManagedPrivateEndpointResource struct{}

func TestAccStreamAnalyticsManagedPrivateEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stream_analytics_managed_private_endpoint", "test")
	r := StreamAnalyticsManagedPrivateEndpointResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
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
	id, err := privateendpoints.ParsePrivateEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.StreamAnalytics.EndpointsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func (r StreamAnalyticsManagedPrivateEndpointResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_managed_private_endpoint" "test" {
  name                          = "acctestprivate_endpoint-%d"
  resource_group_name           = azurerm_resource_group.test.name
  stream_analytics_cluster_name = azurerm_stream_analytics_cluster.test.name
  target_resource_id            = azurerm_storage_account.test.id
  subresource_name              = "blob"
}
`, template, data.RandomInteger)
}

func (r StreamAnalyticsManagedPrivateEndpointResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_stream_analytics_managed_private_endpoint" "import" {
  name                          = azurerm_stream_analytics_managed_private_endpoint.test.name
  resource_group_name           = azurerm_stream_analytics_managed_private_endpoint.test.resource_group_name
  stream_analytics_cluster_name = azurerm_stream_analytics_managed_private_endpoint.test.stream_analytics_cluster_name
  target_resource_id            = azurerm_stream_analytics_managed_private_endpoint.test.target_resource_id
  subresource_name              = azurerm_stream_analytics_managed_private_endpoint.test.subresource_name
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

resource "azurerm_storage_account" "test" {
  name                     = "accteststorageacc%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_stream_analytics_cluster" "test" {
  name                = "acctestcluster-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  streaming_capacity  = 36
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
