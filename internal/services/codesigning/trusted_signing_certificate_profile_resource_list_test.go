// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package codesigning_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/codesigningaccounts"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccTrustedSigningCertificateProfile_list(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	r.preCheck(t)
	accountId := codesigningaccounts.NewCodeSigningAccountID(data.Subscriptions.Primary, fmt.Sprintf("acctestRG-%d", data.RandomInteger), "acctest-"+data.RandomString)
	name := "acctest-" + data.RandomString

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm"),
		Steps: []resource.TestStep{
			{Config: r.complete(data)},
			{
				Query: true,
				Config: fmt.Sprintf(`
list "azurerm_trusted_signing_certificate_profile" "list" {
  provider         = azurerm
  include_resource = true
  config {
    trusted_signing_account_id = "%s"
  }
}
`, accountId.ID()),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("azurerm_trusted_signing_certificate_profile.list", 1),
					querycheck.ExpectIdentity("azurerm_trusted_signing_certificate_profile.list", map[string]knownvalue.Check{
						"name":                      knownvalue.StringExact(name),
						"subscription_id":           knownvalue.StringExact(accountId.SubscriptionId),
						"resource_group_name":       knownvalue.StringExact(accountId.ResourceGroupName),
						"code_signing_account_name": knownvalue.StringExact(accountId.CodeSigningAccountName),
					}),
					querycheck.ExpectResourceKnownValues("azurerm_trusted_signing_certificate_profile.list", queryfilter.ByDisplayName(knownvalue.StringExact(name)), []querycheck.KnownValueCheck{
						{Path: tfjsonpath.New("trusted_signing_account_id"), KnownValue: knownvalue.StringExact(accountId.ID())},
						{Path: tfjsonpath.New("identity_validation_id"), KnownValue: knownvalue.StringExact(os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID"))},
						{Path: tfjsonpath.New("profile_type"), KnownValue: knownvalue.StringExact("PrivateTrust")},
						{Path: tfjsonpath.New("include_city"), KnownValue: knownvalue.Bool(true)},
						{Path: tfjsonpath.New("include_country"), KnownValue: knownvalue.Bool(true)},
						{Path: tfjsonpath.New("include_postal_code"), KnownValue: knownvalue.Bool(true)},
						{Path: tfjsonpath.New("include_state"), KnownValue: knownvalue.Bool(true)},
						{Path: tfjsonpath.New("include_street_address"), KnownValue: knownvalue.Bool(true)},
					}),
				},
			},
		},
	})
}
