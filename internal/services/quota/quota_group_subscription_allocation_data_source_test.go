// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package quota_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type QuotaGroupSubscriptionAllocationDataSource struct{}

func TestAccQuotaGroupSubscriptionAllocationDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_quota_group_subscription_allocation", "test")
	d := QuotaGroupSubscriptionAllocationDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				acceptance.TestCheckResourceAttrSet(data.ResourceName, "id"),
				acceptance.TestCheckResourceAttrSet(data.ResourceName, "allocation.#"),
			),
		},
	})
}

func (d QuotaGroupSubscriptionAllocationDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_quota_group_subscription_allocation" "test" {
  quota_group_id  = azurerm_quota_group_subscription_allocation.test.quota_group_id
  subscription_id = azurerm_quota_group_subscription_allocation.test.subscription_id
  location        = azurerm_quota_group_subscription_allocation.test.location
}
`, QuotaGroupSubscriptionAllocationResource{}.basic(data))
}
