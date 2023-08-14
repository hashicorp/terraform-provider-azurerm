// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetAppSnapshotPolicyDataSource struct{}

func TestAccDataSourceNetAppSnapshotPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_snapshot_policy", "test")
	r := NetAppSnapshotPolicyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("account_name").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("enabled").Exists(),
			),
		},
	})
}

func (NetAppSnapshotPolicyDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_snapshot_policy" "test" {
  name                = azurerm_netapp_snapshot_policy.test.name
  resource_group_name = azurerm_netapp_snapshot_policy.test.resource_group_name
  account_name        = azurerm_netapp_snapshot_policy.test.account_name
}
`, NetAppSnapshotPolicyResource{}.basic(data))
}
