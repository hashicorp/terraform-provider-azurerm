package azurerm

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var testCertThumbprintRaw, _ = ioutil.ReadFile(filepath.Join("testdata", "automation_certificate_test.thumb"))
var testCertRaw, _ = ioutil.ReadFile(filepath.Join("testdata", "automation_certificate_test.pfx"))

var testCertBase64 = base64.StdEncoding.EncodeToString(testCertRaw)

func TestAccAzureRMAutomationCertificate_basic(t *testing.T) {
	resourceName := "azurerm_automation_certificate.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationCertificate_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCertificateExists(resourceName),
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

func TestAccAzureRMAutomationCertificate_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_automation_certificate.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
				ExpectError: acceptance.RequiresImportError("azurerm_automation_certificate"),
			},
		},
	})
}

func TestAccAzureRMAutomationCertificate_complete(t *testing.T) {
	resourceName := "azurerm_automation_certificate.test"
	ri := tf.AccRandTimeInt()
	testCertThumbprint := strings.TrimSpace(string(testCertThumbprintRaw))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationCertificate_complete(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCertificateExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "base64", testCertBase64),
					resource.TestCheckResourceAttr(resourceName, "thumbprint", testCertThumbprint),
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
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationCertificate_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCertificateExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"base64"},
			},
			{
				Config: testAccAzureRMAutomationCertificate_update(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCertificateExists(resourceName),
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
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.CertificateClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_certificate" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

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
		accName := rs.Primary.Attributes["automation_account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.CertificateClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_certificate" "test" {
  name                    = "acctest-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  base64                  = "%s"
}
`, rInt, location, rInt, rInt, testCertBase64)
}

func testAccAzureRMAutomationCertificate_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAutomationCertificate_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_certificate" "import" {
  name                    = azurerm_automation_certificate.test.name
  resource_group_name     = azurerm_automation_certificate.test.resource_group_name
  automation_account_name = azurerm_automation_certificate.test.account_name
  base64                  = azurerm_automation_certificate.test.base64
  thumbprint              = azurerm_automation_certificate.test.thumbprint
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_certificate" "test" {
  name                    = "acctest-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  base64                  = "%s"
  description             = "This is a test certificate for terraform acceptance test"
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
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_certificate" "test" {
  name                    = "acctest-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  base64                  = "%s"
  description             = "This is a test certificate for terraform acceptance test"
}
`, rInt, location, rInt, rInt, testCertBase64)
}
