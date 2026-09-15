// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package codesigning_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/certificateprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type TrustedSigningCertificateProfileResource struct{}

func TestAccTrustedSigningCertificateProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("profile_type").HasValue("PrivateTrust"),
				check.That(data.ResourceName).Key("include_city").HasValue("false"),
				check.That(data.ResourceName).Key("include_country").HasValue("false"),
				check.That(data.ResourceName).Key("include_postal_code").HasValue("false"),
				check.That(data.ResourceName).Key("include_state").HasValue("false"),
				check.That(data.ResourceName).Key("include_street_address").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTrustedSigningCertificateProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccTrustedSigningCertificateProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("include_city").HasValue("true"),
				check.That(data.ResourceName).Key("include_country").HasValue("true"),
				check.That(data.ResourceName).Key("include_postal_code").HasValue("true"),
				check.That(data.ResourceName).Key("include_state").HasValue("true"),
				check.That(data.ResourceName).Key("include_street_address").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTrustedSigningCertificateProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.basic(data)},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("include_city").HasValue("true"),
				check.That(data.ResourceName).Key("include_country").HasValue("true"),
				check.That(data.ResourceName).Key("include_postal_code").HasValue("true"),
				check.That(data.ResourceName).Key("include_state").HasValue("true"),
				check.That(data.ResourceName).Key("include_street_address").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("include_city").HasValue("false"),
				check.That(data.ResourceName).Key("include_country").HasValue("false"),
				check.That(data.ResourceName).Key("include_postal_code").HasValue("false"),
				check.That(data.ResourceName).Key("include_state").HasValue("false"),
				check.That(data.ResourceName).Key("include_street_address").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTrustedSigningCertificateProfile_privateTrustCIPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: strings.ReplaceAll(r.basic(data), "PrivateTrust", "PrivateTrustCIPolicy"),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func TestAccTrustedSigningCertificateProfile_publicTrust(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	identityId := os.Getenv("ARM_TEST_TRUSTED_SIGNING_PUBLIC_IDENTITY_ID")
	if identityId == "" {
		t.Skip("ARM_TEST_TRUSTED_SIGNING_PUBLIC_IDENTITY_ID is required: complete a Public Trust identity validation in the test subscription using the Azure Portal")
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.profile(data, identityId, "PublicTrust", ""),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func (TrustedSigningCertificateProfileResource) preCheck(t *testing.T) {
	// Identity validation requires the Azure Portal and can be shared by accounts in
	// the same subscription. The test creates its own resource group and account.
	// https://learn.microsoft.com/azure/artifact-signing/quickstart#create-an-artifact-signing-account
	if os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID") == "" {
		t.Skip("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID is required: complete a Private Trust identity validation in the test subscription using the Azure Portal")
	}
}

func (TrustedSigningCertificateProfileResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := certificateprofiles.ParseCertificateProfileID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.CodeSigning.Client.CertificateProfiles.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r TrustedSigningCertificateProfileResource) basic(data acceptance.TestData) string {
	return r.profile(data, os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID"), "PrivateTrust", "")
}

func (r TrustedSigningCertificateProfileResource) complete(data acceptance.TestData) string {
	return r.profile(data, os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID"), "PrivateTrust", `
  include_city           = true
  include_country        = true
  include_postal_code    = true
  include_state          = true
  include_street_address = true
`)
}

func (TrustedSigningCertificateProfileResource) profile(data acceptance.TestData, identityId, profileType, options string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_trusted_signing_certificate_profile" "test" {
  name                       = "acctest-%[2]s"
  trusted_signing_account_id = azurerm_trusted_signing_account.test.id
  identity_validation_id     = "%[3]s"
  profile_type               = "%[4]s"
  %[5]s
}
`, TrustedSigningAccountResource{}.basic(data), data.RandomString, identityId, profileType, options)
}

func (r TrustedSigningCertificateProfileResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_trusted_signing_certificate_profile" "import" {
  name                       = azurerm_trusted_signing_certificate_profile.test.name
  trusted_signing_account_id = azurerm_trusted_signing_certificate_profile.test.trusted_signing_account_id
  identity_validation_id     = azurerm_trusted_signing_certificate_profile.test.identity_validation_id
  profile_type               = azurerm_trusted_signing_certificate_profile.test.profile_type
}
`, r.basic(data))
}
