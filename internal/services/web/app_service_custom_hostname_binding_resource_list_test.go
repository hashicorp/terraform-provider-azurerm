package web_test

import (
	"context"
	"os"
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

func TestAccAppServiceCustomHostnameBinding_listByAppServiceID(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_custom_hostname_binding", "test")
	r := AppServiceCustomHostnameBindingResource{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicConfig(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_app_service_custom_hostname_binding.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_app_service_custom_hostname_binding.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringExact(data.RandomString + "." + os.Getenv("ARM_TEST_DNS_ZONE")),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"site_name":           knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
		},
	})
}

func (r AppServiceCustomHostnameBindingResource) basicQuery() string {
	return `
list "azurerm_app_service_custom_hostname_binding" "list" {
  provider = azurerm
  config {
    app_service_id = azurerm_app_service.test.id
  }
}
`
}
