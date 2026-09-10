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

func TestAccApiManagementNamedValue_listByServiceID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_named_value", "testlist1")
	r := ApiManagementNamedValueResource{}

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
					querycheck.ExpectLengthAtLeast("azurerm_api_management_named_value.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_api_management_named_value.list",
						map[string]knownvalue.Check{
							"named_value_id":      knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
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

func (r ApiManagementNamedValueResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_named_value" "test" {
  count = 3 
  name                = "acctestAMProperty-%d-${count.index}"
  resource_group_name = azurerm_api_management.test.resource_group_name
  api_management_name = azurerm_api_management.test.name
  display_name        = "TestProperty%d${count.index}"
  value               = "Test Value"
  tags                = ["tag1", "tag2"]
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementNamedValueResource) basicQuery() string {
	return `
list "azurerm_api_management_named_value" "list" {
  provider = azurerm
  config {
    service_id = azurerm_api_management.test.id
  }
}
`
}
