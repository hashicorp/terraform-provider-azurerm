// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package codesigning_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccTrustedSigningCertificateProfile_resourceIdentity(t *testing.T) {
	for _, v := range []string{"ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID"} {
		if os.Getenv(v) == "" {
			t.Skipf("Skipping as `%s` was not specified", v)
		}
	}

	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}

	checkedFields := map[string]struct{}{
		"name":                      {},
		"code_signing_account_name": {},
		"resource_group_name":       {},
		"subscription_id":           {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_trusted_signing_certificate_profile.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_trusted_signing_certificate_profile.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_trusted_signing_certificate_profile.test", tfjsonpath.New("code_signing_account_name"), tfjsonpath.New("trusted_signing_account_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_trusted_signing_certificate_profile.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("trusted_signing_account_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_trusted_signing_certificate_profile.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("trusted_signing_account_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
