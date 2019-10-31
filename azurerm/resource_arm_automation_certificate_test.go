package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Test certificate generated with:
//   openssl req -x509 -nodes -sha256 -days 3650 -subj "/CN=Local" -newkey rsa:2048 -keyout Local.key -out Local.crt
//   openssl pkcs12 -export -in Local.crt -inkey Local.key -CSP "Microsoft Enhanced RSA and AES Cryptographic Provider" -out Local.pfx

var testCertThumbprint = "F3FF6EA8713C72230933B9152329F7F324E6BF67"
var testCertBase64 = "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"

func TestAccAzureRMAutomationCertificate_basic(t *testing.T) {
	resourceName := "azurerm_automation_certificate.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationCertificate_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCertificateExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMAutomationCertificate_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_automation_certificate.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationCertificate_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCertificateExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAutomationCertificate_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_automation_certificate"),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"base64"},
			},
		},
	})
}

func TestAccAzureRMAutomationCertificate_complete(t *testing.T) {
	resourceName := "azurerm_automation_certificate.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationCertificate_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCertificateExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "base64", testCertBase64),
					resource.TestCheckResourceAttr(resourceName, "thumbprint", testCertThumbprint),
					resource.TestCheckResourceAttr(resourceName, "is_exportable", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test certificate for terraform acceptance test"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"base64"},
			},
		},
	})
}

func TestAccAzureRMAutomationCertificate_update(t *testing.T) {
	resourceName := "azurerm_automation_certificate.test"
	ri := tf.AccRandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationCertificate_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d", ri)),
					testCheckAzureRMAutomationCertificateExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMAutomationCertificate_update(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCertificateExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest-%d-updated", ri)),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test certificate for terraform acceptance test"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"base64"},
			},
		},
	})
}

func testCheckAzureRMAutomationCertificateDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).Automation.CertificateClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_certificate" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Certificate: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Certificate still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationCertificateExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Certificate: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).Automation.CertificateClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Certificate '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationCertificateClient: %s\nName: %s, Account name: %s, Resource group: %s OBJECT: %+v", err, name, accName, resourceGroup, rs.Primary)
		}

		return nil
	}
}

func testAccAzureRMAutomationCertificate_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name            = "Basic"
}

resource "azurerm_automation_certificate" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_automation_account.test.name}"
  base64              = "%s"
}
`, rInt, location, rInt, rInt, testCertBase64)
}

func testAccAzureRMAutomationCertificate_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAutomationCertificate_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_certificate" "import" {
  name                = "${azurerm_automation_certificate.test.name}"
  resource_group_name = "${azurerm_automation_certificate.test.resource_group_name}"
  account_name        = "${azurerm_automation_certificate.test.account_name}"
  base64              = "${azurerm_automation_certificate.test.base64}"
  is_exportable       = "${azurerm_automation_certificate.test.is_exportable}"
  thumbprint          = "${azurerm_automation_certificate.test.username}"
}
`, template)
}

func testAccAzureRMAutomationCertificate_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name            = "Basic"
}

resource "azurerm_automation_certificate" "test" {
  name                = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_automation_account.test.name}"
  base64              = "%s"
  description         = "This is a test certificate for terraform acceptance test"
}
`, rInt, location, rInt, rInt, testCertBase64)
}

func testAccAzureRMAutomationCertificate_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku_name            = "Basic"
}

resource "azurerm_automation_certificate" "test" {
  name                = "acctest-%d-updated"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_automation_account.test.name}"
  base64              = "%s"
  description         = "This is a test certificate for terraform acceptance test"
}
`, rInt, location, rInt, rInt, testCertBase64)
}
