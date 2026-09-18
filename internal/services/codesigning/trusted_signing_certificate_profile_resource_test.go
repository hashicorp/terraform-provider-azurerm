// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package codesigning_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/certificateprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	privateTrustIdentityEnvVar = "ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID"
	publicTrustIdentityEnvVar  = "ARM_TEST_TRUSTED_SIGNING_PUBLIC_IDENTITY_ID"
)

type TrustedSigningCertificateProfileResource struct{}

func TestAccTrustedSigningCertificateProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	r.preCheck(t, privateTrustIdentityEnvVar)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func TestAccTrustedSigningCertificateProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	r.preCheck(t, privateTrustIdentityEnvVar)

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
	r.preCheck(t, privateTrustIdentityEnvVar)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func TestAccTrustedSigningCertificateProfile_privateTrustCIPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	r.preCheck(t, privateTrustIdentityEnvVar)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.profile(data, os.Getenv(privateTrustIdentityEnvVar), "PrivateTrustCIPolicy", ""),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func TestAccTrustedSigningCertificateProfile_publicTrust(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	r.preCheck(t, publicTrustIdentityEnvVar)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.profile(data, os.Getenv(publicTrustIdentityEnvVar), "PublicTrust", ""),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep(),
	})
}

func (TrustedSigningCertificateProfileResource) preCheck(t *testing.T, identityEnvVar string) {
	// Identity validation requires the Azure Portal and can be shared by accounts in
	// the same subscription. The test creates its own resource group and account.
	// https://learn.microsoft.com/azure/artifact-signing/quickstart#create-an-artifact-signing-account
	if os.Getenv(identityEnvVar) == "" {
		t.Skipf("%s is required: complete a matching identity validation in the test subscription using the Azure Portal", identityEnvVar)
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
	return r.profile(data, os.Getenv(privateTrustIdentityEnvVar), "PrivateTrust", "")
}

func (r TrustedSigningCertificateProfileResource) complete(data acceptance.TestData) string {
	return r.profile(data, os.Getenv(privateTrustIdentityEnvVar), "PrivateTrust", `
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
  identity_validation_id     = "%[3]s"
  profile_type               = "%[4]s"
  trusted_signing_account_id = azurerm_trusted_signing_account.test.id
  %[5]s
}
`, TrustedSigningAccountResource{}.basic(data), data.RandomString, identityId, profileType, options)
}

func (r TrustedSigningCertificateProfileResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_trusted_signing_certificate_profile" "import" {
  name                       = azurerm_trusted_signing_certificate_profile.test.name
  identity_validation_id     = azurerm_trusted_signing_certificate_profile.test.identity_validation_id
  profile_type               = azurerm_trusted_signing_certificate_profile.test.profile_type
  trusted_signing_account_id = azurerm_trusted_signing_certificate_profile.test.trusted_signing_account_id
}
`, r.basic(data))
}
