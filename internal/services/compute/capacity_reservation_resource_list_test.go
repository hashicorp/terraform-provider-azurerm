// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

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

func TestAccCapacityReservation_listByCapacityReservationGroupID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation", "testlist1")
	r := CapacityReservationResource{}
	listResourceAddress := "azurerm_capacity_reservation.list"

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
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
					querycheck.ExpectIdentity(
						listResourceAddress,
						map[string]knownvalue.Check{
							"name":                            knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name":             knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"capacity_reservation_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":                 knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r CapacityReservationResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-compute-%[1]d"
  location = "%[2]s"
}

resource "azurerm_capacity_reservation_group" "test" {
  name                = "acctest-ccrg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  zones               = ["1", "2"]
}

resource "azurerm_capacity_reservation" "test" {
  count                         = 2
  name                          = "acctest-ccr${count.index}-%[1]d"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  zone                          = tostring(count.index + 1)
  sku {
    name     = "Standard_F2"
    capacity = 1
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CapacityReservationResource) basicQuery() string {
	return `
list "azurerm_capacity_reservation" "list" {
  provider = azurerm
  config {
    capacity_reservation_group_id = azurerm_capacity_reservation_group.test.id
  }
}
`
}
