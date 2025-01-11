// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoDatabaseResource struct{}

func TestAccKustoDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database", "test")
	r := KustoDatabaseResource{}

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

func TestAccKustoDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database", "test")
	r := KustoDatabaseResource{}

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

func TestAccKustoDatabase_name(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database", "test")
	r := KustoDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.name(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKustoDatabase_softDeletePeriod(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database", "test")
	r := KustoDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.softDeletePeriod(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("soft_delete_period").HasValue("P7D"),
			),
		},
		{
			Config: r.softDeletePeriodUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("soft_delete_period").HasValue("P31D"),
			),
		},
	})
}

func TestAccKustoDatabase_hotCachePeriod(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kusto_database", "test")
	r := KustoDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hotCachePeriod(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("hot_cache_period").HasValue("P7D"),
			),
		},
		{
			Config: r.hotCachePeriodUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("hot_cache_period").HasValue("P14DT12H"),
			),
		},
	})
}

func (KustoDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (KustoDatabaseResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name
  soft_delete_period  = "P31D"
  hot_cache_period    = "P7D"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (KustoDatabaseResource) name(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctest_kd_%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (KustoDatabaseResource) softDeletePeriod(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name

  soft_delete_period = "P7D"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (KustoDatabaseResource) softDeletePeriodUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name

  soft_delete_period = "P31D"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (KustoDatabaseResource) hotCachePeriod(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name

  hot_cache_period = "P7D"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (KustoDatabaseResource) hotCachePeriodUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_kusto_cluster" "cluster" {
  name                = "acctestkc%s"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  sku {
    name     = "Dev(No SLA)_Standard_D11_v2"
    capacity = 1
  }
}

resource "azurerm_kusto_database" "test" {
  name                = "acctestkd-%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  cluster_name        = azurerm_kusto_cluster.cluster.name

  hot_cache_period = "P14DT12H"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (KustoDatabaseResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseKustoDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Kusto.DatabasesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	if resp.Model != nil {
		value, ok := resp.Model.(databases.ReadWriteDatabase)
		if !ok {
			return nil, fmt.Errorf("%s is not a ReadWriteDatabase", id.String())
		}

		exists := value.Properties != nil

		return &exists, nil
	} else {
		return nil, fmt.Errorf("response model is empty")
	}
}
