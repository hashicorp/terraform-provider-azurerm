// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package codesigning_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type TrustedSigningAccountDataSource struct{}

func TestAccTrustedSigningAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_trusted_signing_account", "test")
	r := TrustedSigningAccountDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("sku_name").HasValue("Basic"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("test"),
			),
		},
	})
}

func (TrustedSigningAccountDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[3]s"
}

resource "azurerm_trusted_signing_account" "test" {
  name                = "acctest-%[2]s"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"

  tags = {
    env = "test"
  }
}

data "azurerm_trusted_signing_account" "test" {
  name                = azurerm_trusted_signing_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.RandomString, data.Locations.Primary)
}
