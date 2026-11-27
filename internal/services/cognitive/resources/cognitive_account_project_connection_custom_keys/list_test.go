// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_custom_keys_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccCognitiveAccountProjectConnectionCustomKeys_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_project_connection_custom_keys", "test")
	r := CognitiveAccountProjectConnectionCustomKeysResource{}

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
				Config: r.basicListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("azurerm_cognitive_account_project_connection_custom_keys.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_cognitive_account_project_connection_custom_keys.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile("acctest-conn")),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
							"account_name":        knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"project_name":        knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
					querycheck.ExpectResourceKnownValues(
						"azurerm_cognitive_account_project_connection_custom_keys.list",
						queryfilter.ByDisplayName(knownvalue.StringRegexp(regexp.MustCompile("acctest-conn2-"))),
						[]querycheck.KnownValueCheck{
							{
								Path:       tfjsonpath.New("authentication_type"),
								KnownValue: knownvalue.StringExact("CustomKeys"),
							},
						},
					),
				},
			},
		},
	})
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) basicListQuery() string {
	return `
list "azurerm_cognitive_account_project_connection_custom_keys" "list" {
  provider         = azurerm
  include_resource = true
  config {
    cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  }
}
`
}

func (r CognitiveAccountProjectConnectionCustomKeysResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_project_connection_custom_keys" "test2" {
  name                         = "acctest-conn2-%[2]d"
  cognitive_account_project_id = azurerm_cognitive_account_project.test.id
  category                     = "CustomKeys"
  target                       = "https://tool.example.com/"

  custom_keys = {
    apiKey = azurerm_cognitive_account.test.primary_access_key
  }
}
`, r.basic(data), data.RandomInteger)
}
