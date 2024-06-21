// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arcresourcebridge_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ArcResourceBridgeApplianceDataSource struct{}

func TestAccArcResourceBridgeApplianceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_arc_resource_bridge_appliance", "test")
	d := ArcResourceBridgeApplianceDataSource{}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").IsNotEmpty(),
				check.That(data.ResourceName).Key("distro").IsNotEmpty(),
				check.That(data.ResourceName).Key("infrastructure_provider").IsNotEmpty(),
				check.That(data.ResourceName).Key("public_key_base64").IsNotEmpty(),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func (d ArcResourceBridgeApplianceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
data "azurerm_arc_resource_bridge_appliance" "test" {
  name                = azurerm_arc_resource_bridge_appliance.test.name
  resource_group_name = azurerm_arc_resource_bridge_appliance.test.resource_group_name
}
`, ArcResourceBridgeApplianceResource{}.complete(data))
}
