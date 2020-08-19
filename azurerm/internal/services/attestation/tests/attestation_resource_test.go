package tests

import (
	"bytes"
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/attestation/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAttestation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAttestation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAttestation_requiresImport),
		},
	})
}

func TestAccAzureRMAttestation_twoKeys(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation", "test")

	testCertificate1, err := testAzureRMGenerateTestCertificate("FLYNN'S ARCADE")
	if err != nil {
		t.Fatalf("Test case failed: '%s'", err)
	}

	testCertificate2, err := testAzureRMGenerateTestCertificate("SPACE PARANOIDS")
	if err != nil {
		t.Fatalf("Test case failed: '%s'", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestation_twoKeys(data, testCertificate1, testCertificate2),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationExists(data.ResourceName),
				),
			},
			// must ignore policy_signing_certificate since the API does not return these values
			data.ImportStep("policy_signing_certificate"),
		},
	})
}

func TestAccAzureRMAttestation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation", "test")
	testCertificate, err := testAzureRMGenerateTestCertificate("ENCOM")
	if err != nil {
		t.Fatalf("Test case failed: '%s'", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestation_complete(data, testCertificate),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationExists(data.ResourceName),
				),
			},
			// must ignore policy_signing_certificate since the API does not return these values
			data.ImportStep("policy_signing_certificate"),
		},
	})
}
func TestAccAzureRMAttestation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAttestation_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAttestation_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAttestationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Attestation.ProviderClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("attestation AttestationProvider not found: %s", resourceName)
		}
		id, err := parse.AttestationId(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Attestation AttestationProvider %q does not exist", id.Name)
			}
			return fmt.Errorf("bad: Get on Attestation.ProviderClient: %+v", err)
		}
		return nil
	}
}

func testAzureRMGenerateTestCertificate(organization string) (string, error) {
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
		return "", fmt.Errorf("unable to create test certificate: %s", err)
	}

	encoded := &bytes.Buffer{}
	if err := pem.Encode(encoded, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return "", fmt.Errorf("unable to pem encode test certificate: %s", err)
	}

	base64Cert, err := testAzureRMConvertCertificateToBase64String(encoded.String())
	if err != nil {
		return "", fmt.Errorf("unable to convert test certificate into single base64 string: %s", err)
	}

	return base64Cert, nil
}

func testAzureRMConvertCertificateToBase64String(input string) (string, error) {
	cleanInput := strings.TrimSpace(input)

	if len(cleanInput) == 0 {
		return "", fmt.Errorf("unable to convert empty test certificate")
	}

	cleanCert := bytes.Split([]byte(cleanInput), []byte{10})
	concatCert := []byte{}

	for i, v := range cleanCert {
		if i == 0 || len(v) == 0 || i == (len(cleanCert)-1) {
			continue
		}

		concatCert = append(concatCert, cleanCert[i]...)
	}

	return string(concatCert), nil
}

func testCheckAzureRMAttestationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Attestation.ProviderClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_attestation" {
			continue
		}
		id, err := parse.AttestationId(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on Attestation.ProviderClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

// currently only supported in "East US 2", "West Central US" & "UK South"
func testAccAzureRMAttestation_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-attestation-%d"
  location = "%s"
}
`, data.RandomInteger, "uksouth")
}

func testAccAzureRMAttestation_basic(data acceptance.TestData) string {
	template := testAccAzureRMAttestation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation" "test" {
  name                = "ap%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func testAccAzureRMAttestation_update(data acceptance.TestData) string {
	template := testAccAzureRMAttestation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation" "test" {
  name                = "ap%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMAttestation_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMAttestation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation" "import" {
  name                = azurerm_attestation.test.name
  resource_group_name = azurerm_attestation.test.resource_group_name
  location            = azurerm_attestation.test.location
}
`, config)
}

func testAccAzureRMAttestation_complete(data acceptance.TestData, testCertificate string) string {
	template := testAccAzureRMAttestation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation" "test" {
  name                = "ap%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_signing_certificate {
    key {
      kty = "RSA"
      x5c = ["%s"]
    }
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, testCertificate)
}

func testAccAzureRMAttestation_twoKeys(data acceptance.TestData, testCertificate1 string, testCertificate2 string) string {
	template := testAccAzureRMAttestation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation" "test" {
  name                = "ap%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_signing_certificate {
    key {
      kty = "RSA"
      x5c = ["%s"]
    }
    key {
      kty = "RSA"
      x5c = ["%s"]
    }
  }
}
`, template, data.RandomInteger, testCertificate1, testCertificate2)
}
