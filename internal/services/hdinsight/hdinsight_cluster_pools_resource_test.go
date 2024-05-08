// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2024-05-01/hdinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"testing"
)

type ClusterPoolResource struct{}

func (r ClusterPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := hdinsights.ParseClusterPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.HDInsight2024.Hdinsights.ClusterPoolsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retreiving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func TestAccHDInsightClusterPools_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hdinsight_cluster_pools", "test")
	r := ClusterPoolResource{}

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

func (r ClusterPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_hdinsight_cluster_pools" "test" {
  name 			  = "acctestpool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location 		  = azurerm_resource_group.test.location
  
  compute_profile {
	vm_size = "Standard_D3_v2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ClusterPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
