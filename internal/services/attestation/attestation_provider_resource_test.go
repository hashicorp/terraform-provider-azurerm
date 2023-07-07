// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attestation_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AttestationProviderResource struct {
	name string
}

func TestAccAttestationProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{
		name: fmt.Sprintf("acctestap%s", data.RandomStringOfLength(10)),
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAttestationProvider_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{
		name: fmt.Sprintf("acctestap%s", data.RandomStringOfLength(10)),
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAttestationProvider_completeString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{
		name: fmt.Sprintf("acctestap%s", data.RandomStringOfLength(10)),
	}
	testCertificate, err := testGenerateTestCertificate("ENCOM")
	if err != nil {
		t.Fatalf("Test case failed: '%+v'", err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeString(data, testCertificate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// must ignore policy_signing_certificate since the API does not return this value
		data.ImportStep("policy_signing_certificate"),
	})
}

func TestAccAttestationProvider_withPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{
		name: fmt.Sprintf("acctestap%s", data.RandomStringOfLength(10)),
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAttestationProvider_withPolicyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{
		name: fmt.Sprintf("acctestap%s", data.RandomStringOfLength(10)),
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// must ignore policy_signing_certificate since the API does not return these values
		data.ImportStep("policy_signing_certificate"),
	})
}

func TestAccAttestationProvider_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	r := AttestationProviderResource{
		name: fmt.Sprintf("acctestap%s", data.RandomStringOfLength(10)),
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AttestationProviderResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := attestationproviders.ParseAttestationProvidersID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Attestation.ProviderClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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

func (AttestationProviderResource) template(data acceptance.TestData) string {
	// currently only supported in "East US 2", "West Central US" & "UK South"
	data.Locations.Primary = "westus"
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-attestation-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AttestationProviderResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_attestation_provider" "test" {
  name                = %q
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  lifecycle {
    ignore_changes = [
      "open_enclave_policy_base64",
      "sgx_enclave_policy_base64",
      "tpm_policy_base64",
      "sev_snp_policy_base64",
    ]
  }
}
`, r.template(data), r.name)
}

func (r AttestationProviderResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_attestation_provider" "test" {
  name                = %q
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  lifecycle {
    ignore_changes = [
      "open_enclave_policy_base64",
      "sgx_enclave_policy_base64",
      "tpm_policy_base64",
      "sev_snp_policy_base64",
    ]
  }

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), r.name)
}

func (r AttestationProviderResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_attestation_provider" "import" {
  name                = azurerm_attestation_provider.test.name
  resource_group_name = azurerm_attestation_provider.test.resource_group_name
  location            = azurerm_attestation_provider.test.location

  lifecycle {
    ignore_changes = [
      "open_enclave_policy_base64",
      "sgx_enclave_policy_base64",
      "tpm_policy_base64",
      "sev_snp_policy_base64",
    ]
  }
}
`, r.basic(data))
}

func (r AttestationProviderResource) completeString(data acceptance.TestData, testCertificate string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_attestation_provider" "test" {
  name                = %[2]q
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_signing_certificate_data = <<EOT
%[3]s
EOT

  tags = {
    ENV = "Test"
  }

  lifecycle {
    ignore_changes = [
      "open_enclave_policy_base64",
      "sgx_enclave_policy_base64",
      "tpm_policy_base64",
      "sev_snp_policy_base64",
    ]
  }
}
`, r.template(data), r.name, testCertificate)
}

func (r AttestationProviderResource) withPolicy(data acceptance.TestData) string {
	// do not set `policy_signing_certificate_data`, since the policies use `jwt.SigningMethodNone`
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_attestation_provider" "test" {
  name                = %[2]q
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  open_enclave_policy_base64 = %[3]q
  sgx_enclave_policy_base64  = %[3]q
  tpm_policy_base64          = %[3]q
  sev_snp_policy_base64      = %[3]q

  lifecycle {
    ignore_changes = [
      "tpm_policy_base64",
    ]
  }
}
`, r.template(data), r.name, r.genJWT())
}

func (r AttestationProviderResource) genJWT() string {
	// document about create policy: https://learn.microsoft.com/en-us/azure/attestation/author-sign-policy
	policyContent := `version=1.0;
authorizationrules
{
[type=="secureBootEnabled", value==true, issuer=="AttestationService"]=>permit();
};

issuancerules
{
=> issue(type="SecurityLevelValue", value=100);
};`
	b64Encoded := base64.RawURLEncoding.EncodeToString([]byte(policyContent))
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"AttestationPolicy": b64Encoded,
	})
	token.Header["jku"] = fmt.Sprintf("https://%s.uks.attest.azure.net/certs", r.name)
	token.Header["kid"] = "xxx"

	str, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	return str
}
