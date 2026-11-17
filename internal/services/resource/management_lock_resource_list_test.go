// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_test

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

func TestAccManagementLock_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "list")
	r := ManagementLockResource{}

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
				Config: r.basicListQueryByScope(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_management_lock.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_management_lock.list",
						map[string]knownvalue.Check{
							"name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"scope": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResourceGroup(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_management_lock.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_management_lock.list",
						map[string]knownvalue.Check{
							"name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"scope": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryByResource(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_management_lock.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_management_lock.list",
						map[string]knownvalue.Check{
							"name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"scope": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.basicListQueryBySubscription(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_management_lock.list", 0),
				},
			},
		},
	})
}

func (ManagementLockResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_management_lock" "test1" {
  name       = "acctestlock1-%[1]d"
  scope      = azurerm_resource_group.test.id
  lock_level = "ReadOnly"
}

resource "azurerm_management_lock" "test2" {
  name       = "acctestlock2-%[1]d"
  scope      = azurerm_resource_group.test.id
  lock_level = "CanNotDelete"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ManagementLockResource) basicListQueryByScope() string {
	return `
list "azurerm_management_lock" "list" {
  provider = azurerm
  config {
    scope = azurerm_resource_group.test.id
  }
}
`
}

func (ManagementLockResource) basicListQueryByResourceGroup() string {
	return `
list "azurerm_management_lock" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}

func (ManagementLockResource) basicListQueryByResource() string {
	return `
list "azurerm_management_lock" "list" {
  provider = azurerm
  config {
    resource_id = azurerm_resource_group.test.id
  }
}
`
}

func (ManagementLockResource) basicListQueryBySubscription() string {
	return `
list "azurerm_management_lock" "list" {
  provider = azurerm
  config {
  }
}
`
}
