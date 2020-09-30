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

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/attestation/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAttestationProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	randStr := strings.ToLower(acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestationProvider_basic(data, randStr),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAttestationProvider_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	randStr := strings.ToLower(acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestationProvider_basic(data, randStr),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationProviderExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAttestationProvider_requiresImport),
		},
	})
}

func TestAccAzureRMAttestationProvider_completeString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	randStr := strings.ToLower(acctest.RandString(10))
	testCertificate, err := testAzureRMGenerateTestCertificate("ENCOM")
	if err != nil {
		t.Fatalf("Test case failed: '%+v'", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestationProvider_completeString(data, randStr, testCertificate),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationProviderExists(data.ResourceName),
				),
			},
			// must ignore policy_signing_certificate since the API does not return these values
			data.ImportStep("policy_signing_certificate"),
		},
	})
}

func TestAccAzureRMAttestationProvider_completeFile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	randStr := strings.ToLower(acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestationProvider_completeFile(data, randStr),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationProviderExists(data.ResourceName),
				),
			},
			// must ignore policy_signing_certificate since the API does not return these values
			data.ImportStep("policy_signing_certificate"),
		},
	})
}

func TestAccAzureRMAttestationProvider_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation_provider", "test")
	randStr := strings.ToLower(acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestationProvider_basic(data, randStr),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAttestationProvider_update(data, randStr),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAttestationProvider_basic(data, randStr),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationProviderExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAttestationProviderExists(resourceName string) resource.TestCheckFunc {
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
				return fmt.Errorf("bad: Attestation Provider %q does not exist", id.Name)
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
		return "", fmt.Errorf("unable to create test certificate: %+v", err)
	}

	encoded := &bytes.Buffer{}
	if err := pem.Encode(encoded, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return "", fmt.Errorf("unable to pem encode test certificate: %+v", err)
	}

	return encoded.String(), nil
}

func testCheckAzureRMAttestationProviderDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Attestation.ProviderClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_attestation_provider" {
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
func testAccAzureRMAttestationProvider_template(data acceptance.TestData) string {
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

func testAccAzureRMAttestationProvider_basic(data acceptance.TestData, randStr string) string {
	template := testAccAzureRMAttestationProvider_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation_provider" "test" {
  name                = "acctestap%s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, randStr)
}

func testAccAzureRMAttestationProvider_update(data acceptance.TestData, randStr string) string {
	template := testAccAzureRMAttestationProvider_template(data)
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

func testAccAzureRMAttestationProvider_requiresImport(data acceptance.TestData) string {
	randStr := strings.ToLower(acctest.RandString(10))
	config := testAccAzureRMAttestationProvider_basic(data, randStr)

	return fmt.Sprintf(`
%s

resource "azurerm_attestation_provider" "import" {
  name                = azurerm_attestation_provider.test.name
  resource_group_name = azurerm_attestation_provider.test.resource_group_name
  location            = azurerm_attestation_provider.test.location
}
`, config)
}

func testAccAzureRMAttestationProvider_completeString(data acceptance.TestData, randStr string, testCertificate string) string {
	template := testAccAzureRMAttestationProvider_template(data)
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

func testAccAzureRMAttestationProvider_completeFile(data acceptance.TestData, randStr string) string {
	template := testAccAzureRMAttestationProvider_template(data)
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
