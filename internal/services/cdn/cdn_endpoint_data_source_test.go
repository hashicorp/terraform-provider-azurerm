// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn"
)

type CdnEndpointDataSource struct{}

func TestAccCdnEndpointDataSource_basic(t *testing.T) {
	if cdn.IsCdnDeprecatedForCreation() {
		t.Skipf("skipping as CDN (Classic) endpoint creation is deprecated")
	}

	data := acceptance.BuildTestData(t, "data.azurerm_cdn_endpoint", "test")
	d := CdnEndpointDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("is_http_allowed").HasValue("true"),
				check.That(data.ResourceName).Key("is_https_allowed").HasValue("true"),
				check.That(data.ResourceName).Key("origin.#").HasValue("1"),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
	})
}

func (d CdnEndpointDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acctestcdnend%d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name      = "acceptanceTestCdnOrigin1"
    host_name = "www.contoso.com"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}

data "azurerm_cdn_endpoint" "test" {
  name                = azurerm_cdn_endpoint.test.name
  profile_name        = azurerm_cdn_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
