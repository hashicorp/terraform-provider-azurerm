package attestation_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/attestation/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AttestationProviderResource struct {
}

func TestAccAttestationProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{}
	randStr := strings.ToLower(acctest.RandString(10))

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, randStr),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAttestationProvider_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{}
	randStr := strings.ToLower(acctest.RandString(10))

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, randStr),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(AttestationProviderResource{}.requiresImport),
	})
}

func TestAccAttestationProvider_completeString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{}
	randStr := strings.ToLower(acctest.RandString(10))
	testCertificate, err := testGenerateTestCertificate("ENCOM")
	if err != nil {
		t.Fatalf("Test case failed: '%+v'", err)
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeString(data, randStr, testCertificate),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// must ignore policy_signing_certificate since the API does not return these values
		data.ImportStep("policy_signing_certificate"),
	})
}

func TestAccAttestationProvider_completeFile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{}
	randStr := strings.ToLower(acctest.RandString(10))

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeFile(data, randStr),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// must ignore policy_signing_certificate since the API does not return these values
		data.ImportStep("policy_signing_certificate"),
	})
}

func TestAccAttestationProvider_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{}
	randStr := strings.ToLower(acctest.RandString(10))

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, randStr),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, randStr),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, randStr),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t AttestationProviderResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ProviderID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Attestation.ProviderClient.Get(ctx, id.ResourceGroup, id.AttestationProviderName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Attestation Provider %q (resource group: %q): %+v", id.AttestationProviderName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.StatusResult != nil), nil
}

func testGenerateTestCertificate(organization string) (string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return "", err
	}

	rawCert := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{organization},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &rawCert, &rawCert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", fmt.Errorf("unable to create test certificate: %+v", err)
	}

	encoded := &bytes.Buffer{}
	if err := pem.Encode(encoded, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return "", fmt.Errorf("unable to pem encode test certificate: %+v", err)
	}

	return encoded.String(), nil
}

// currently only supported in "East US 2", "West Central US" & "UK South"
func (AttestationProviderResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
// TODO: switch to using regular regions when this is supported
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-attestation-%d"
  location = "%s"
}
`, data.RandomInteger, "uksouth")
}

func (AttestationProviderResource) basic(data acceptance.TestData, randStr string) string {
	template := AttestationProviderResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation_provider" "test" {
  name                = "acctestap%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, randStr)
}

func (AttestationProviderResource) update(data acceptance.TestData, randStr string) string {
	template := AttestationProviderResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation_provider" "test" {
  name                = "acctestap%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Test"
  }
}
`, template, randStr)
}

func (AttestationProviderResource) requiresImport(data acceptance.TestData) string {
	randStr := strings.ToLower(acctest.RandString(10))
	config := AttestationProviderResource{}.basic(data, randStr)

	return fmt.Sprintf(`
%s

resource "azurerm_attestation_provider" "import" {
  name                = azurerm_attestation_provider.test.name
  resource_group_name = azurerm_attestation_provider.test.resource_group_name
  location            = azurerm_attestation_provider.test.location
}
`, config)
}

func (AttestationProviderResource) completeString(data acceptance.TestData, randStr string, testCertificate string) string {
	template := AttestationProviderResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation_provider" "test" {
  name                = "acctestap%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_signing_certificate_data = <<EOT
%s
EOT

  tags = {
    ENV = "Test"
  }
}
`, template, randStr, testCertificate)
}

func (AttestationProviderResource) completeFile(data acceptance.TestData, randStr string) string {
	template := AttestationProviderResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation_provider" "test" {
  name                = "acctestap%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_signing_certificate_data = file("testdata/cert.pem")

  tags = {
    ENV = "Test"
  }
}
`, template, randStr)
}
