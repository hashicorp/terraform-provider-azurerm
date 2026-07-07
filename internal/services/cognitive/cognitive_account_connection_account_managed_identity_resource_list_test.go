// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
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

func TestAccCognitiveAccountConnectionAccountManagedIdentity_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_account_connection_account_managed_identity", "test")
	r := CognitiveAccountConnectionAccountManagedIdentityResource{}

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
					// Azure allows only one AccountManagedIdentity connection with the AzureKeyVault category per account.
					querycheck.ExpectLengthAtLeast("azurerm_cognitive_account_connection_account_managed_identity.list", 1),
					querycheck.ExpectIdentity(
						"azurerm_cognitive_account_connection_account_managed_identity.list",
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

func (r CognitiveAccountConnectionAccountManagedIdentityResource) basicListQuery(data acceptance.TestData) string {
	return `
list "azurerm_cognitive_account_connection_account_managed_identity" "list" {
  provider = azurerm
  config {
    cognitive_account_id = azurerm_cognitive_account.test.id
  }
}
`
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) basicList(data acceptance.TestData) string {
	return r.basic(data)
}
