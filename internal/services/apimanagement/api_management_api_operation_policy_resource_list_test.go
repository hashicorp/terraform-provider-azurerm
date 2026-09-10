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

func TestAccApiManagementApiOperationPolicy_listByOperationID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_operation_policy", "testlist1")
	r := ApiManagementApiOperationPolicyResource{}

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
					querycheck.ExpectLengthAtLeast("azurerm_api_management_api_operation_policy.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_api_management_api_operation_policy.list",
						map[string]knownvalue.Check{
							"api_id":              knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"operation_id":        knownvalue.StringExact("acctest-operation-0"),
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

func (r ApiManagementApiOperationPolicyResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_operation" "test" {
  count               = 3
  operation_id        = "acctest-operation-${count.index}"
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  display_name        = "DELETE Resource ${count.index}"
  method              = "DELETE"
  url_template        = "/resource-${count.index}"
}

resource "azurerm_api_management_api_operation_policy" "test" {
  count               = 3
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  operation_id        = azurerm_api_management_api_operation.test[count.index].operation_id
  xml_link            = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/refs/heads/main/internal/services/apimanagement/testdata/api_management_api_operation_policy.xml"
}
`, ApiManagementApiOperationResource{}.template(data))
}

func (r ApiManagementApiOperationPolicyResource) basicQuery() string {
	return `
list "azurerm_api_management_api_operation_policy" "list" {
  provider = azurerm
  config {
    operation_id = azurerm_api_management_api_operation.test[0].id
  }
}
`
}
