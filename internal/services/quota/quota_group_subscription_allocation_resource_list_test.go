// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package quota_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccQuotaGroupSubscriptionAllocation_list_basic(t *testing.T) {
	r := QuotaGroupSubscriptionAllocationResource{}
	listResourceAddress := "azurerm_quota_group_subscription_allocation.list"

	data := acceptance.BuildTestData(t, "azurerm_quota_group_subscription_allocation", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				// a quota group can only associate the test subscription once, so exactly one allocation is expected
				Config: r.basic(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 1),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryIncludeResource(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 1),
				},
			},
		},
	})
}

func (r QuotaGroupSubscriptionAllocationResource) basicQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_quota_group_subscription_allocation" "list" {
  provider = azurerm
  config {
    quota_group_id = "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Quota/groupQuotas/acctestqg%d"
    location       = "%s"
  }
}
`, data.Client().TenantID, data.RandomInteger, data.Locations.Primary)
}

func (r QuotaGroupSubscriptionAllocationResource) basicQueryIncludeResource(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_quota_group_subscription_allocation" "list" {
  provider         = azurerm
  include_resource = true
  config {
    quota_group_id         = "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Quota/groupQuotas/acctestqg%d"
    location               = "%s"
    resource_provider_name = "Microsoft.Compute"
  }
}
`, data.Client().TenantID, data.RandomInteger, data.Locations.Primary)
}
