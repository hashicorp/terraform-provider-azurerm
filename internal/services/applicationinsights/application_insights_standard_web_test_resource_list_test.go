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

func TestAccApplicationInsightsStandardWebTest_list_basic(t *testing.T) {
	r := ApplicationInsightsStandardWebTestResource{}
	listResourceAddress := "azurerm_application_insights_standard_web_test.list"

	data := acceptance.BuildTestData(t, "azurerm_application_insights_standard_web_test", "test1")

	acceptance.RunTest(t, resource.TestCase{
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

func (r ApplicationInsightsStandardWebTestResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_standard_web_test" "test" {
  count                   = 3
  name                    = "acctestwebtest${count.index}-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  geo_locations           = ["us-tx-sn1-azr"]

  request {
    url = "http://microsoft.com"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApplicationInsightsStandardWebTestResource) basicQuery() string {
	return `
list "azurerm_application_insights_standard_web_test" "list" {
  provider = azurerm
  config {}
}
`
}

func (r ApplicationInsightsStandardWebTestResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_application_insights_standard_web_test" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
