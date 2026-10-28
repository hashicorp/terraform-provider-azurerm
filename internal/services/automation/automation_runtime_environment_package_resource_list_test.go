// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_test

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

func TestAccAutomationRuntimeEnvironmentPackage_listByRuntimeEnvironmentID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_runtime_environment_package", "testlist1")
	r := AutomationRuntimeEnvironmentPackageResource{}

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
					querycheck.ExpectLength("azurerm_automation_runtime_environment_package.list", 3),
					querycheck.ExpectIdentity(
						"azurerm_automation_runtime_environment_package.list",
						map[string]knownvalue.Check{
							"name":                     knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name":      knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"runtime_environment_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":          knownvalue.StringExact(data.Subscriptions.Primary),
							"automation_account_name":  knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
		},
	})
}

func (r AutomationRuntimeEnvironmentPackageResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_runtime_environment_package" "test" {
  count                             = 3
  name                              = "acctest-package-%[2]d-${count.index}"
  automation_runtime_environment_id = azurerm_automation_runtime_environment.test.id
  content_uri                       = "https://www.powershellgallery.com/api/v2/package/Microsoft.Graph.Authentication/2.25.0"
}
`, r.template(data), data.RandomInteger)
}

func (r AutomationRuntimeEnvironmentPackageResource) basicQuery() string {
	return `
list "azurerm_automation_runtime_environment_package" "list" {
  provider = azurerm
  config {
    automation_runtime_environment_id = azurerm_automation_runtime_environment.test.id
  }
}
`
}
