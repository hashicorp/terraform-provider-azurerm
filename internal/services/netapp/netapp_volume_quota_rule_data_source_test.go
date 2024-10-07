// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NetAppVolumeQuotaRuleDataSource struct{}

func TestAccNetAppVolumeQuotaRuleDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_netapp_volume_quota_rule", "test")
	d := NetAppVolumeQuotaRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("quota_size_in_kib").Exists(),
			),
		},
	})
}

func (d NetAppVolumeQuotaRuleDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_netapp_volume_quota_rule" "test" {
  name      = azurerm_netapp_volume_quota_rule.test.name
  volume_id = azurerm_netapp_volume_quota_rule.test.volume_id
}
`, NetAppVolumeQuotaRuleResource{}.individualUserQuotaType(data))
}
