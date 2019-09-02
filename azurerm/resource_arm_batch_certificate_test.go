package azurerm

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMBatchCertificate_Pfx(t *testing.T) {
	resourceName := "azurerm_batch_certificate.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	certificateID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-42c107874fd0e4a9583292a2f1098e8fe4b2edda", subscriptionID, ri, rs)

	config := testAccAzureRMBatchCertificatePfx(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", certificateID),
					resource.TestCheckResourceAttr(resourceName, "format", "Pfx"),
					resource.TestCheckResourceAttr(resourceName, "thumbprint", "42c107874fd0e4a9583292a2f1098e8fe4b2edda"),
					resource.TestCheckResourceAttr(resourceName, "thumbprint_algorithm", "sha1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate", "password"},
			},
		},
	})
}

func TestAccAzureRMBatchCertificate_PfxWithoutPassword(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	config := testAccAzureRMBatchCertificatePfxWithoutPassword(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Password is required"),
			},
		},
	})
}

func TestAccAzureRMBatchCertificate_Cer(t *testing.T) {
	resourceName := "azurerm_batch_certificate.test"
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")
	certificateID := fmt.Sprintf("/subscriptions/%s/resourceGroups/testaccbatch%d/providers/Microsoft.Batch/batchAccounts/testaccbatch%s/certificates/sha1-312d31a79fa0cef49c00f769afc2b73e9f4edf34", subscriptionID, ri, rs)

	config := testAccAzureRMBatchCertificateCer(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr(resourceName, "id", certificateID),
					resource.TestCheckResourceAttr(resourceName, "format", "Cer"),
					resource.TestCheckResourceAttr(resourceName, "thumbprint", "312d31a79fa0cef49c00f769afc2b73e9f4edf34"),
					resource.TestCheckResourceAttr(resourceName, "thumbprint_algorithm", "sha1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"certificate"},
			},
		},
	})
}

func TestAccAzureRMBatchCertificate_CerWithPassword(t *testing.T) {
	ri := tf.AccRandTimeInt()
	rs := acctest.RandString(4)
	location := testLocation()

	config := testAccAzureRMBatchCertificateCerWithPassword(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("Password must not be specified"),
			},
		},
	})
}

func testAccAzureRMBatchCertificatePfx(rInt int, batchAccountSuffix string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = "${azurerm_resource_group.test.name}"
  account_name         = "${azurerm_batch_account.test.name}"
  certificate          = "${filebase64("testdata/batch_certificate.pfx")}"
  format               = "Pfx"
  password             = "terraform"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}
`, rInt, location, batchAccountSuffix)
}

func testAccAzureRMBatchCertificatePfxWithoutPassword(rInt int, batchAccountSuffix string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = "${azurerm_resource_group.test.name}"
  account_name         = "${azurerm_batch_account.test.name}"
  certificate          = "${filebase64("testdata/batch_certificate.pfx")}"
  format               = "Pfx"
  thumbprint           = "42c107874fd0e4a9583292a2f1098e8fe4b2edda"
  thumbprint_algorithm = "SHA1"
}
`, rInt, location, batchAccountSuffix)
}
func testAccAzureRMBatchCertificateCer(rInt int, batchAccountSuffix string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = "${azurerm_resource_group.test.name}"
  account_name         = "${azurerm_batch_account.test.name}"
  certificate          = "${filebase64("testdata/batch_certificate.cer")}"
  format               = "Cer"
  thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34"
  thumbprint_algorithm = "SHA1"
}
`, rInt, location, batchAccountSuffix)
}
func testAccAzureRMBatchCertificateCerWithPassword(rInt int, batchAccountSuffix string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = "${azurerm_resource_group.test.name}"
  account_name         = "${azurerm_batch_account.test.name}"
  certificate          = "${filebase64("testdata/batch_certificate.cer")}"
  format               = "Cer"
  password             = "should not have a password for Cer"
  thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34"
  thumbprint_algorithm = "SHA1"
}
`, rInt, location, batchAccountSuffix)
}

func testCheckAzureRMBatchCertificateDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).batch.CertificateClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_batch_certificate" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accountName := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}
