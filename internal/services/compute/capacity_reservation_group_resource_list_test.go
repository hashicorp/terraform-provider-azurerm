package compute_test

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

func TestAccCapacityReservationGroup_listBySubscriptionAndRG(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_capacity_reservation_group", "test1")
	r := CapacityReservationGroupResource{}
	listResourceAddress := "azurerm_capacity_reservation_group.list"

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
					querycheck.ExpectLengthAtLeast(listResourceAddress, 3),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r CapacityReservationGroupResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-compute-%[1]d"
  location = "%[2]s"
}

resource "azurerm_capacity_reservation_group" "test" {
  count               = 3
  name                = "acctest-crg${count.index}-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r CapacityReservationGroupResource) basicQuery() string {
	return `
list "azurerm_capacity_reservation_group" "list" {
  provider = azurerm
  config {}
}
`
}

func (r CapacityReservationGroupResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_capacity_reservation_group" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
