// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connections_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ManagedApiTestDataSource struct{}

func TestAccDataSourceManagedApi_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_managed_api", "test")
	r := ManagedApiTestDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func (ManagedApiTestDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_managed_api" "test" {
  name     = "servicebus"
  location = %q
}
`, data.Locations.Primary)
}
