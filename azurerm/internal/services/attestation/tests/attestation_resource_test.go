package tests

import (
	"fmt"
	"testing"

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

func TestAccAzureRMAttestation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_attestation", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAttestationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAttestation_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAttestationExists(data.ResourceName),
				),
			},
			data.ImportStep(),
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
				Config: testAccAzureRMAttestation_complete(data),
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

// hard coded location because attestation resources are only available in the 'eastus2,centralus,uksouth' regions.
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

func testAccAzureRMAttestation_requiresImport(data acceptance.TestData) string {
	config := testAccAzureRMAttestation_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation" "import" {
  name                = azurerm_attestation.test.name
  resource_group_name = azurerm_attestation.test.resource_group_name
  location            = azurerm_attestation.test.location
  attest_uri          = azurerm_attestation.test.attest_uri
  trust_model         = azurerm_attestation.test.trust_model
  type                = azurerm_attestation.test.type
}
`, config)
}

func testAccAzureRMAttestation_complete(data acceptance.TestData) string {
	template := testAccAzureRMAttestation_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_attestation" "test" {
  name                = "ap%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  policy_signing_certificate {
    key {
      kty  = "RSA"
      x5c = ["LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJ6akNDQVRHZ0F3SUJBZ0lCQVRBS0JnZ3Foa2pPUFFRREJEQVNNUkF3RGdZRFZRUUtFd2RCWTIxbElFTnYKTUI0WERUQTVNVEV4TURJek1EQXdNRm9YRFRFd01EVXdPVEl6TURBd01Gb3dFakVRTUE0R0ExVUVDaE1IUVdOdApaU0JEYnpDQm16QVFCZ2NxaGtqT1BRSUJCZ1VyZ1FRQUl3T0JoZ0FFQUpuVnlBdkoyOW1DUlJoVTQ4UVdrUDJKCmJ6K0ZCRmFET0xTMXBhVW9KUGJ6R0NZd05velg4NjBZTXdVOGtzd1BZTzkrUDE0c08vRnZkRVQ3ZURFZlN2VUMKQUFocEtJT01jSHdjcDE3WG5NazFPSmIwZExWK2wrUVViL1pDR3pkYTdRc2J3MjA3Zlp2OERYOHhWMFYwQkVnOApyUW1IM0htNCtCK0xsaVU4bFFXU3lFcG5velV3TXpBT0JnTlZIUThCQWY4RUJBTUNCNEF3RXdZRFZSMGxCQXd3CkNnWUlLd1lCQlFVSEF3RXdEQVlEVlIwVEFRSC9CQUl3QURBS0JnZ3Foa2pPUFFRREJBT0JpZ0F3Z1lZQ1FYNkQKYUVqVjhQS0V5SWJyNUFBd1pyOTc5QkEzQndtcGUxNUd4ODlUbk5ScDZ1MjhZbSt0eFB6RjBFK3ZXS3NuQjF6aQpTdURhYmVRMHpMbGRCaVhRVHdkV0FrRURGUEdjb2owUjZvS3Y1N0lZZjRjd3hjdklmMjI1SEhrVXVDdXAzQit4ClFZcXRxWW44OTdJSDRWUXdrbVpCWnBoVXVjY2xEaWxIQlhQeWE3RTZZb2dBMmc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="]
    }
  }

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
