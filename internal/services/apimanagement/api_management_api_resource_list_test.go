package apimanagement_test

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

func TestAccApiManagementApi_listByServiceID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api", "testlist1")
	r := ApiManagementApiResource{}

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
					querycheck.ExpectLengthAtLeast("azurerm_api_management_api.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_api_management_api.list",
						map[string]knownvalue.Check{
							"api_id":              knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"service_name":        knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r ApiManagementApiResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api" "test" {
  count               = 3
  name                = "acctestapi-%d-${count.index}"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api${count.index}"
  path                = "api${count.index}"
  protocols           = ["https"]
  revision            = "1"
}
`, r.template(data, SkuNameConsumption), data.RandomInteger)
}

func (r ApiManagementApiResource) basicQuery() string {
	return `
list "azurerm_api_management_api" "list" {
  provider = azurerm
  config {
    service_id = azurerm_api_management.test.id
  }
}
`
}
