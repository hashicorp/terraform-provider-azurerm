// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MonitorDataCollectionEndpointDataSource struct{}

func TestAccMonitorDataCollectionEndpointDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_data_collection_endpoint", "test")
	d := MonitorDataCollectionEndpointDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("kind").HasValue("Windows"),
				check.That(data.ResourceName).Key("public_network_access_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("configuration_access_endpoint").Exists(),
				check.That(data.ResourceName).Key("logs_ingestion_endpoint").Exists(),
				check.That(data.ResourceName).Key("metrics_ingestion_endpoint").Exists(),
				check.That(data.ResourceName).Key("immutable_id").Exists(),
			),
		},
	})
}

func (d MonitorDataCollectionEndpointDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_monitor_data_collection_endpoint" "test" {
  name                = azurerm_monitor_data_collection_endpoint.test.name
  resource_group_name = azurerm_monitor_data_collection_endpoint.test.resource_group_name
}
`, MonitorDataCollectionEndpointResource{}.complete(data))
}
