// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask_test

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

type DurableTaskRetentionPolicyListResource struct{}

func TestAccDurableTaskRetentionPolicyList_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_durable_task_retention_policy", "test")
	r := DurableTaskRetentionPolicyListResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
			{
				Query:  true,
				Config: r.basicListQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_durable_task_retention_policy.list", 1),
				},
			},
		},
	})
}

func (r DurableTaskRetentionPolicyListResource) basicListQuery(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_durable_task_retention_policy" "list" {
  provider = azurerm
  config {
    scheduler_id = "/subscriptions/%s/resourceGroups/acctestRG-durabletask-%d/providers/Microsoft.DurableTask/schedulers/acctestdts%s"
  }
}
`, data.Subscriptions.Primary, data.RandomInteger, data.RandomString)
}
