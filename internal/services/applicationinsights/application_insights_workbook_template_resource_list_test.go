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

func TestAccApplicationInsightsWorkbookTemplate_list_basic(t *testing.T) {
	r := ApplicationInsightsWorkbookTemplateResource{}
	listResourceAddress := "azurerm_application_insights_workbook_template.list"

	data := acceptance.BuildTestData(t, "azurerm_application_insights_workbook_template", "test0")

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
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r ApplicationInsightsWorkbookTemplateResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights_workbook_template" "test" {
  count               = 3
  name                = "acctest-aiwt${count.index}-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  galleries {
    category = "workbook"
    name     = "test"
  }

  template_data = jsonencode({
    "version" : "Notebook/1.0",
    "items" : [
      {
        "type" : 1,
        "content" : {
          "json" : "## New workbook\n---\n\nWelcome to your new workbook."
        },
        "name" : "text - 2"
      }
    ],
    "styleSettings" : {},
    "$schema" : "https://github.com/Microsoft/Application-Insights-Workbooks/blob/master/schema/workbook.json"
  })
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ApplicationInsightsWorkbookTemplateResource) basicQueryByResourceGroupName() string {
	return `
list "azurerm_application_insights_workbook_template" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
`
}
