// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestStorageAccount_list_basic(t *testing.T) {
	r := StorageAccountResource{}

	data := acceptance.BuildTestData(t, "azurerm_storage_account", "list1")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:             true,
				Config:            r.basic_query(data),
				ConfigQueryChecks: []querycheck.QueryCheck{}, // TODO
			},
			{
				Query:             true,
				Config:            r.basic_queryByResourceGroup(data),
				ConfigQueryChecks: []querycheck.QueryCheck{}, // TODO
			},
		},
	})
}

func (r StorageAccountResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
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

func (r StorageAccountResource) basic_query(_ acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_storage_account" "test" {
  provider = azurerm
  config {}
}`)
}

func (r StorageAccountResource) basic_queryByResourceGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_storage_account" "test" {
  provider = azurerm
  config {
	resource_group_name = "acctestRG-storage-%d"
  }
}
`, data.RandomInteger)
}
