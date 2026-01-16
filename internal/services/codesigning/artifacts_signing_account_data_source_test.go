// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package codesigning_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ArtifactsSigningAccountDataSource struct{}

func TestAccArtifactsSigningAccountDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_artifacts_signing_account", "test")
	r := ArtifactsSigningAccountDataSource{}

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

func (ArtifactsSigningAccountDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[3]s"
}

resource "azurerm_artifacts_signing_account" "test" {
  name                = "acctest-%[2]s"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"

  tags = {
    env = "test"
  }
}

data "azurerm_artifacts_signing_account" "test" {
  name                = azurerm_artifacts_signing_account.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.RandomString, data.Locations.Primary)
}
