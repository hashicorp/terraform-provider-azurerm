// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ResourcesDataSource struct{}

func TestAccDataSourceResources_ByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resources", "test")
	r := ResourcesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			PreConfig: func() { time.Sleep(1 * time.Minute) },
			Config:    r.ByName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resources.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceResources_ByResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resources", "test")
	r := ResourcesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			PreConfig: func() { time.Sleep(1 * time.Minute) },
			Config:    r.ByResourceGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resources.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceResources_ByResourceType(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resources", "test")
	r := ResourcesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			PreConfig: func() { time.Sleep(1 * time.Minute) },
			Config:    r.ByResourceType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resources.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceResources_FilteredByResourceTypeAndSubscriptionIDs(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resources", "test")
	r := ResourcesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			PreConfig: func() { time.Sleep(1 * time.Minute) },
			Config:    r.ByResourceTypeAndSubscriptionIDs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resources.#").HasValue("1"),
			),
		},
	})
}

func TestAccDataSourceResources_FilteredByTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_resources", "test")
	r := ResourcesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.template(data),
		},
		{
			PreConfig: func() { time.Sleep(1 * time.Minute) },
			Config:    r.FilteredByTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resources.#").HasValue("1"),
			),
		},
	})
}

func (r ResourcesDataSource) ByName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  name = azurerm_storage_account.test.name
}
`, r.template(data))
}

func (r ResourcesDataSource) ByResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  resource_group_name = azurerm_storage_account.test.resource_group_name
}
`, r.template(data))
}

func (r ResourcesDataSource) ByResourceType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  resource_group_name = azurerm_storage_account.test.resource_group_name
  type                = "Microsoft.Storage/storageAccounts"
}
`, r.template(data))
}

func (r ResourcesDataSource) ByResourceTypeAndSubscriptionIDs(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

data "azurerm_resources" "test" {
	type                = "Microsoft.Storage/storageAccounts"
	resource_group_name = azurerm_storage_account.test.resource_group_name

	subscription_ids    = [
		data.azurerm_client_config.current.subscription_id
	]
}
`, r.template(data))
}

func (r ResourcesDataSource) FilteredByTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_resources" "test" {
  name                = azurerm_storage_account.test.name
  resource_group_name = azurerm_storage_account.test.resource_group_name

  required_tags = {
    environment = "production"
  }
}
`, r.template(data))
}

func (ResourcesDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "acctestsads%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
