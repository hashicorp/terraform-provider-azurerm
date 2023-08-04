// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communication_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CommunicationServiceDataSource struct{}

func TestAccCommunicationServiceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_communication_service", "test")
	d := CommunicationServiceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("data_location").Exists(),
				check.That(data.ResourceName).Key("primary_connection_string").Exists(),
				check.That(data.ResourceName).Key("secondary_connection_string").Exists(),
				check.That(data.ResourceName).Key("primary_key").Exists(),
				check.That(data.ResourceName).Key("secondary_key").Exists(),
			),
		},
	})
}

func TestAccCommunicationServiceDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_communication_service", "test")
	d := CommunicationServiceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.env").HasValue("Test"),
			),
		},
	})
}

func (d CommunicationServiceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_communication_service" "test" {
  name                = azurerm_communication_service.test.name
  resource_group_name = azurerm_communication_service.test.resource_group_name
}
`, CommunicationServiceTestResource{}.basic(data))
}

func (d CommunicationServiceDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_communication_service" "test" {
  name                = azurerm_communication_service.test.name
  resource_group_name = azurerm_communication_service.test.resource_group_name
}
`, CommunicationServiceTestResource{}.complete(data))
}
