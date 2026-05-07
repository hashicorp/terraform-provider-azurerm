// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

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

func TestAccStorageActionsTaskAssignment_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_assignment", "testlist")
	r := StorageActionsTaskAssignmentResource{}
	listResourceAddress := "azurerm_storage_actions_task_assignment.list"

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.storageAccountListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r StorageActionsTaskAssignmentResource) basicList(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_assignment" "test" {
  count = 3

  name               = "acctest%s${count.index}"
  storage_account_id = azurerm_storage_account.test.id
  task_id            = azurerm_storage_actions_task.test.id
  description        = "list test"
  enabled            = true

  execution_context {
    trigger {
      type       = "OnSchedule"
      interval   = 1
      start_from = "2030-01-01T00:00:00Z"
      end_by     = "2031-01-01T00:00:00Z"
    }
  }

  report_prefix = "report"
}
`, template, data.RandomString)
}

func (r StorageActionsTaskAssignmentResource) storageAccountListQuery() string {
	return `
list "azurerm_storage_actions_task_assignment" "list" {
  provider = azurerm
  config {
    storage_account_id = azurerm_storage_account.test.id
  }
}
`
}
