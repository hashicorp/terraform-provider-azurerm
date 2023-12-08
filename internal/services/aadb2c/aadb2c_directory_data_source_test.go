// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package aadb2c_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AadB2cDirectoryDataSource struct{}

func TestAccAadB2cDirectoryDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_aadb2c_directory", "test")
	d := AadB2cDirectoryDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("data_residency_location").HasValue("United States"),
				check.That(data.ResourceName).Key("sku_name").HasValue("PremiumP1"),
			),
		},
	})
}

func (d AadB2cDirectoryDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

data "azurerm_aadb2c_directory" "test" {
  domain_name         = azurerm_aadb2c_directory.test.domain_name
  resource_group_name = azurerm_aadb2c_directory.test.resource_group_name
}
`, AadB2cDirectoryResource{}.basic(data))
}
