package appservice_test

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

func TestAccServicePlan_list_basic(t *testing.T) {
	r := ServicePlanResource{}
	listResourceAddress := "azurerm_service_plan.list"

	data := acceptance.BuildTestData(t, "azurerm_service_plan", "test")

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
				Config: r.basicQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r ServicePlanResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appserviceplan-%[1]d"
  location = "%[2]s"
}

resource "azurerm_service_plan" "test" {
  count = 3

  name                = "acctest-${count.index}-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "B1"
  os_type             = "Windows"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ServicePlanResource) basicQuery() string {
	return `
list "azurerm_service_plan" "list" {
  provider = azurerm
  config {}
}
`
}

func (r ServicePlanResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_service_plan" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestRG-appserviceplan-%[1]d"
  }
}
`, data.RandomInteger)
}
