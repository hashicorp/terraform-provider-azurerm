// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccStorageActionsTaskDefinition_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_actions_task_definition", "testlist")
	r := StorageActionsTaskDefinitionResource{}
	listResourceAddress := "azurerm_storage_actions_task_definition.list"

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
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile("acctest" + data.RandomString)),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.resourceGroupListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile("acctest" + data.RandomString)),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r StorageActionsTaskDefinitionResource) basicList(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_storage_actions_task_definition" "test" {
  count = 3

  name                = "acctest%s${count.index}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "list test"

  identity {
    type = "SystemAssigned"
  }

  action {
    if {
      condition = "[[equals(AccessTier, 'Cool')]]"

      operation {
        name = "SetBlobTier"

        parameters = {
          tier = "Hot"
        }
      }
    }
  }
}
`, template, data.RandomString)
}

func (r StorageActionsTaskDefinitionResource) subscriptionListQuery() string {
	return `
list "azurerm_storage_actions_task_definition" "list" {
  provider = azurerm
  config {}
}
`
}

func (r StorageActionsTaskDefinitionResource) resourceGroupListQuery() string {
	return `
list "azurerm_storage_actions_task_definition" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
