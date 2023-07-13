// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datashare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ShareKustoClusterDataSetResource struct{}

func TestAccDataShareKustoClusterDataSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_kusto_cluster", "test")
	r := ShareKustoClusterDataSetResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("kusto_cluster_location").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataShareKustoClusterDataSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_kusto_cluster", "test")
	r := ShareKustoClusterDataSetResource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t ShareKustoClusterDataSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dataset.ParseDataSetID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataShare.DataSetClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		ds := *model
		if _, ok := ds.(dataset.KustoClusterDataSet); ok {
			return utils.Bool(true), nil
		}
	}

	return nil, fmt.Errorf("%s is not a kusto cluster dataset", *id)
}

func (ShareKustoClusterDataSetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datashare-%[1]d"
  location = "%[2]s"
}

resource "azurerm_data_share_account" "test" {
  name                = "acctest-DSA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_share" "test" {
  name       = "acctest_DS_%[1]d"
  account_id = azurerm_data_share_account.test.id
  kind       = "InPlace"
}

resource "azurerm_kusto_cluster" "test" {
  name                = "acctestkc%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_kusto_cluster.test.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_data_share_account.test.identity.0.principal_id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r ShareKustoClusterDataSetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_kusto_cluster" "test" {
  name             = "acctest-DSKC-%d"
  share_id         = azurerm_data_share.test.id
  kusto_cluster_id = azurerm_kusto_cluster.test.id
  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r ShareKustoClusterDataSetResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_data_share_dataset_kusto_cluster" "import" {
  name             = azurerm_data_share_dataset_kusto_cluster.test.name
  share_id         = azurerm_data_share.test.id
  kusto_cluster_id = azurerm_kusto_cluster.test.id
}
`, r.basic(data))
}
