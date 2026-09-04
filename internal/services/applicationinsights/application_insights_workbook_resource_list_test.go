// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package applicationinsights_test

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

func TestAccApplicationInsightsWorkbook_listBySubscriptionAndRG(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook", "test1")
	r := ApplicationInsightsWorkbookResource{}
	listResourceAddress := "azurerm_application_insights_workbook.list"

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

func (r ApplicationInsightsWorkbookResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights_workbook" "test1" {
  count               = 3
  name                = "99c8b2b4-740d-4510-9a69-8b4d7cbb2fe${count.index}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  display_name        = "acctest-amw-${count.index}-%[1]d"
  data_json = jsonencode({
    "version" = "Notebook/1.0",
    "items" = [
      {
        "type"    = 1,
        "content" = { "json" = "Test" },
        "name"    = "text - 0"
      }
    ],
    "isLocked"            = false,
    "fallbackResourceIds" = ["Azure Monitor"]
  })
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApplicationInsightsWorkbookResource) basicQuery() string {
	return `
list "azurerm_application_insights_workbook" "list" {
  provider = azurerm
  config {}
}
`
}

func (r ApplicationInsightsWorkbookResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_application_insights_workbook" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
