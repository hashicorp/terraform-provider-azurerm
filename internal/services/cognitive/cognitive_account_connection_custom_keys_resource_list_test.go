// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

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

func TestAccCognitiveAccountConnectionCustomKeys_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_custom_keys", "test")
	r := CognitiveAccountConnectionCustomKeysResource{}

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
				Config: r.basicListQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_cognitive_account_connection_custom_keys.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_cognitive_account_connection_custom_keys.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
							"account_name":        knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
						},
					),
				},
			},
		},
	})
}

func (r CognitiveAccountConnectionCustomKeysResource) basicListQuery(data acceptance.TestData) string {
	return `
list "azurerm_cognitive_account_connection_custom_keys" "list" {
  provider = azurerm
  config {
    cognitive_account_id = azurerm_cognitive_account.test.id
  }
}
`
}

func (r CognitiveAccountConnectionCustomKeysResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_account_connection_custom_keys" "test2" {
  name                 = "acctest-conn2-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  category             = "CustomKeys"
  target               = azurerm_cognitive_account.openai.endpoint

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.openai.primary_access_key
    secondaryKey = azurerm_cognitive_account.openai.secondary_access_key
  }
}
`, r.basic(data), data.RandomInteger)
}
