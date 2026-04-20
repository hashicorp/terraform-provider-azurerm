package resource_test

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

func TestAccResourceGroup_list_basic(t *testing.T) {
	r := ResourceGroupResource{}
	listResourceAddress := "azurerm_resource_group.list"

	data := acceptance.BuildTestData(t, "azurerm_resource_group", "test")

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
					querycheck.ExpectLengthAtLeast(listResourceAddress, 4),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryWithFilter(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 3),
				},
			},
		},
	})
}

func (r ResourceGroupResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  count = 3

  name     = "acctestRG${count.index}-%[1]d"
  location = "%[2]s"

  tags = {
    "query" = "test-%[1]d"
  }
}

resource "azurerm_resource_group" "test-no-tags" {
  name     = "acctestRG4-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceGroupResource) basicQuery() string {
	return `
list "azurerm_resource_group" "list" {
  provider = azurerm
  config {}
}
`
}

func (r ResourceGroupResource) basicQueryWithFilter(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_resource_group" "list" {
  provider = azurerm
  config {
    filter = "tagName eq 'query' and tagValue eq 'test-%d'"
  }
}
`, data.RandomInteger)
}
