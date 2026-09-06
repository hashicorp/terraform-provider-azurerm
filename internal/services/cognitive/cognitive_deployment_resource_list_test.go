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

func TestAccCognitiveDeployment_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_deployment", "test")
	r := CognitiveDeploymentResource{}

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
				Config: r.basicListQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_cognitive_deployment.list", 2),
					querycheck.ExpectIdentity(
						"azurerm_cognitive_deployment.list",
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

func (r CognitiveDeploymentResource) basicListQuery(data acceptance.TestData) string {
	return `
list "azurerm_cognitive_deployment" "list" {
  provider = azurerm
  config {
    cognitive_account_id = azurerm_cognitive_account.test.id
  }
}
`
}

func (r CognitiveDeploymentResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_cognitive_deployment" "test2" {
  name                 = "acctest-cd2-%[2]d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  model {
    format = "OpenAI"
    name   = "text-embedding-ada-002"
  }
  sku {
    name = "Standard"
  }
  lifecycle {
    ignore_changes = [model[0].version]
  }
}
`, r.basic(data), data.RandomInteger)
}
