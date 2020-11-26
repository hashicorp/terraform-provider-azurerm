package batch_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccBatchCertificate_Pfx(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_certificate", "test")
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	certificateID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-42c107874fd0e4a9583292a2f1098e8fe4b2edda", subscriptionID, data.RandomInteger, data.RandomString)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchCertificatePfx(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", certificateID),
					resource.TestCheckResourceAttr(data.ResourceName, "format", "Pfx"),
					resource.TestCheckResourceAttr(data.ResourceName, "thumbprint", "42c107874fd0e4a9583292a2f1098e8fe4b2edda"),
					resource.TestCheckResourceAttr(data.ResourceName, "thumbprint_algorithm", "sha1"),
				),
			},
			{
				ResourceName:            data.ResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "password"},
			},
		},
	})
}

func TestAccBatchCertificate_PfxWithoutPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccBatchCertificatePfxWithoutPassword(data),
				ExpectError: regexp.MustCompile("Password is required"),
			},
		},
	})
}

func TestAccBatchCertificate_Cer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_certificate", "test")
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	certificateID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-312d31a79fa0cef49c00f769afc2b73e9f4edf34", subscriptionID, data.RandomInteger, data.RandomString)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchCertificateCer(data),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr(data.ResourceName, "id", certificateID),
					resource.TestCheckResourceAttr(data.ResourceName, "format", "Cer"),
					resource.TestCheckResourceAttr(data.ResourceName, "thumbprint", "312d31a79fa0cef49c00f769afc2b73e9f4edf34"),
					resource.TestCheckResourceAttr(data.ResourceName, "thumbprint_algorithm", "sha1"),
				),
			},
			data.ImportStep("certificate"),
		},
	})
}

func TestAccBatchCertificate_CerWithPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_batch_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckBatchCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccBatchCertificateCerWithPassword(data),
				ExpectError: regexp.MustCompile("Password must not be specified"),
			},
		},
	})
}

func testAccBatchCertificatePfx(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.pfx")
  format               = "Pfx"
  password             = "terraform"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccBatchCertificatePfxWithoutPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.pfx")
  format               = "Pfx"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccBatchCertificateCer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.cer")
  format               = "Cer"
  thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34"
  thumbprint_algorithm = "SHA1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testAccBatchCertificateCerWithPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  account_name         = azurerm_batch_account.test.name
  certificate          = filebase64("testdata/batch_certificate.cer")
  format               = "Cer"
  password             = "should not have a password for Cer"
  thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34"
  thumbprint_algorithm = "SHA1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testCheckBatchCertificateDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Batch.CertificateClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_batch_certificate" {
			continue
		}

		id, err := parse.CertificateID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}
