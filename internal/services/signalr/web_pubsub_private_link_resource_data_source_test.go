// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type WebPubsubPrivateLinkedServiceDataSource struct{}

func TestAccWebPubsubPrivateLinkedService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_web_pubsub_private_link_resource", "test")
	r := WebPubsubPrivateLinkedServiceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("shared_private_link_resource_types.#").Exists(),
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (r WebPubsubPrivateLinkedServiceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_web_pubsub" "test" {
  name                = "acctestWebPubsub-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_S1"
  capacity            = 1
}

data "azurerm_web_pubsub_private_link_resource" "test" {
  web_pubsub_id = azurerm_web_pubsub.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
