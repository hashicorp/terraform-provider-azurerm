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

func TestAccStorageActionsTask_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task", "testlist")
	r := StorageActionsTaskResource{}
	listResourceAddress := "azurerm_storage_actions_task.list"

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
				Config: r.subscriptionListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.resourceGroupListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r StorageActionsTaskResource) basicList(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task" "test" {
  count = 3

  name                = "acctest%s${count.index}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "list test"
  enabled             = true

  identity {
    type = "SystemAssigned"
  }

  action {
    if {
      condition = "[[equals(AccessTier, 'Cool')]]"

      operation {
        name       = "SetBlobTier"
        on_failure = "break"
        on_success = "continue"

        parameters = {
          tier = "Hot"
        }
      }
    }
  }
}
`, template, data.RandomString)
}

func (r StorageActionsTaskResource) subscriptionListQuery() string {
	return `
list "azurerm_storage_actions_task" "list" {
  provider = azurerm
  config {}
}
`
}

func (r StorageActionsTaskResource) resourceGroupListQuery() string {
	return `
list "azurerm_storage_actions_task" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
