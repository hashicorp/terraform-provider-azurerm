package codesigning_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2024-09-30-preview/certificateprofiles"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type TrustedSigningCertificateProfileResource struct{}

func TestAccTrustedSigningCertificateProfile_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID") == "" {
		t.Skip("Skipping as ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID is not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificates.0.created_date").Exists(),
				check.That(data.ResourceName).Key("certificates.0.enhanced_key_usage").Exists(),
				check.That(data.ResourceName).Key("certificates.0.expiry_date").Exists(),
				check.That(data.ResourceName).Key("certificates.0.serial_number").Exists(),
				check.That(data.ResourceName).Key("certificates.0.status").Exists(),
				check.That(data.ResourceName).Key("certificates.0.subject_name").Exists(),
				check.That(data.ResourceName).Key("certificates.0.thumbprint").Exists(),
				check.That(data.ResourceName).Key("status").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTrustedSigningCertificateProfile_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID") == "" {
		t.Skip("Skipping as ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID is not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificates.0.created_date").Exists(),
				check.That(data.ResourceName).Key("certificates.0.enhanced_key_usage").Exists(),
				check.That(data.ResourceName).Key("certificates.0.expiry_date").Exists(),
				check.That(data.ResourceName).Key("certificates.0.serial_number").Exists(),
				check.That(data.ResourceName).Key("certificates.0.status").Exists(),
				check.That(data.ResourceName).Key("certificates.0.subject_name").Exists(),
				check.That(data.ResourceName).Key("certificates.0.thumbprint").Exists(),
				check.That(data.ResourceName).Key("status").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAcccodesigningCertificateProfile_complete(t *testing.T) {
	if os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID") == "" {
		t.Skip("Skipping as ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID is not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificates.0.created_date").Exists(),
				check.That(data.ResourceName).Key("certificates.0.enhanced_key_usage").Exists(),
				check.That(data.ResourceName).Key("certificates.0.expiry_date").Exists(),
				check.That(data.ResourceName).Key("certificates.0.serial_number").Exists(),
				check.That(data.ResourceName).Key("certificates.0.status").Exists(),
				check.That(data.ResourceName).Key("certificates.0.subject_name").Exists(),
				check.That(data.ResourceName).Key("certificates.0.thumbprint").Exists(),
				check.That(data.ResourceName).Key("status").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccTrustedSigningCertificateProfile_update(t *testing.T) {
	if os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID") == "" {
		t.Skip("Skipping as ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID is not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_trusted_signing_certificate_profile", "test")
	r := TrustedSigningCertificateProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificates.0.created_date").Exists(),
				check.That(data.ResourceName).Key("certificates.0.enhanced_key_usage").Exists(),
				check.That(data.ResourceName).Key("certificates.0.expiry_date").Exists(),
				check.That(data.ResourceName).Key("certificates.0.serial_number").Exists(),
				check.That(data.ResourceName).Key("certificates.0.status").Exists(),
				check.That(data.ResourceName).Key("certificates.0.subject_name").Exists(),
				check.That(data.ResourceName).Key("certificates.0.thumbprint").Exists(),
				check.That(data.ResourceName).Key("status").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("certificates.0.created_date").Exists(),
				check.That(data.ResourceName).Key("certificates.0.enhanced_key_usage").Exists(),
				check.That(data.ResourceName).Key("certificates.0.expiry_date").Exists(),
				check.That(data.ResourceName).Key("certificates.0.serial_number").Exists(),
				check.That(data.ResourceName).Key("certificates.0.status").Exists(),
				check.That(data.ResourceName).Key("certificates.0.subject_name").Exists(),
				check.That(data.ResourceName).Key("certificates.0.thumbprint").Exists(),
				check.That(data.ResourceName).Key("status").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (r TrustedSigningCertificateProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := certificateprofiles.ParseCertificateProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.CodeSigning.Client.CertificateProfiles
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r TrustedSigningCertificateProfileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}
resource "azurerm_trusted_signing_account" "test" {
  name                = "acctest-%[3]s"
  location            = "%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r TrustedSigningCertificateProfileResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	identityId := os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID")
	return fmt.Sprintf(`
				%s

resource "azurerm_trusted_signing_certificate_profile" "test" {
  name                       = "acctest-ccp-%d"
  trusted_signing_account_id = azurerm_trusted_signing_account.test.id
  identity_validation_id     = "%s"
  profile_type               = "PrivateTrust"
}
`, template, data.RandomInteger, identityId)
}

func (r TrustedSigningCertificateProfileResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_trusted_signing_certificate_profile" "import" {
  name                       = azurerm_trusted_signing_certificate_profile.test.name
  trusted_signing_account_id = azurerm_trusted_signing_certificate_profile.test.trusted_signing_account_id
  identity_validation_id     = azurerm_trusted_signing_certificate_profile.test.identity_validation_id
  profile_type               = azurerm_trusted_signing_certificate_profile.test.profile_type
}
`, config)
}

func (r TrustedSigningCertificateProfileResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	identityId := os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID")
	return fmt.Sprintf(`
			%s

resource "azurerm_trusted_signing_certificate_profile" "test" {
  name                       = "acctest-ccp-%d"
  trusted_signing_account_id = azurerm_trusted_signing_account.test.id
  identity_validation_id     = "%s"
  include_city               = true
  include_country            = true
  include_postal_code        = true
  include_state              = true
  include_street_address     = true
  profile_type               = "PrivateTrust"

}
`, template, data.RandomInteger, identityId)
}

func (r TrustedSigningCertificateProfileResource) update(data acceptance.TestData) string {
	template := r.template(data)
	identityId := os.Getenv("ARM_TEST_TRUSTED_SIGNING_IDENTITY_ID")
	return fmt.Sprintf(`
			%s

resource "azurerm_trusted_signing_certificate_profile" "test" {
  name                       = "acctest-ccp-%d"
  trusted_signing_account_id = azurerm_trusted_signing_account.test.id
  identity_validation_id     = "%s"
  include_city               = false
  include_country            = false
  include_postal_code        = false
  include_state              = false
  include_street_address     = false
  profile_type               = "PrivateTrust"

}
`, template, data.RandomInteger, identityId)
}
