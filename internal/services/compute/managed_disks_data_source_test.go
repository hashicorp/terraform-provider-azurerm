// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagedDisksDataSource struct{}

func TestAccDataSourceManagedDisks_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_disks", "test")
	r := ManagedDisksDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("disk.#").HasValue("2"),
			),
		},
	})
}

func (ManagedDisksDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestmanageddisk-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
  zone                 = "2"

  tags = {
    environment = "acctest"
  }
}

resource "azurerm_managed_disk" "test2" {
  name                 = "acctestmanageddisk2-%[2]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Premium_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
  zone                 = "2"

  tags = {
    environment = "acctest"
  }
}

data "azurerm_managed_disks" "test" {
  resource_group_name = azurerm_resource_group.test.name

  // Force Terraform to create the resources before evaluating this Data Source
  depends_on = [azurerm_managed_disk.test, azurerm_managed_disk.test2]
}
`, data.Locations.Primary, data.RandomInteger)
}
