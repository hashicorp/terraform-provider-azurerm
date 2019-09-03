package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAPIManagementCertificate_basic(t *testing.T) {
	resourceName := "azurerm_api_management_certificate.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementCertificate_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementCertificateExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "subject"),
					resource.TestCheckResourceAttrSet(resourceName, "thumbprint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"data",
					"password",
				},
			},
		},
	})
}

func TestAccAzureRMAPIManagementCertificate_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_certificate.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAPIManagementCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAPIManagementCertificate_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAPIManagementCertificateExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMAPIManagementCertificate_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management_certificate"),
			},
		},
	})
}

func testCheckAzureRMAPIManagementCertificateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).apiManagement.CertificatesClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_certificate" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, serviceName, name)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMAPIManagementCertificateExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := testAccProvider.Meta().(*ArmClient).apiManagement.CertificatesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Certificate %q (Resource Group %q / API Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagementCertificatesClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAPIManagementCertificate_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_certificate" "test" {
  name                = "example-cert"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data                = "${filebase64("testdata/keyvaultcert.pfx")}"
  password            = ""
}
`, rInt, location, rInt)
}

func testAccAzureRMAPIManagementCertificate_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAPIManagementCertificate_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_certificate" "import" {
  name                = "${azurerm_api_management_certificate.test.name}"
  api_management_name = "${azurerm_api_management_certificate.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_certificate.test.resource_group_name}"
  data                = "${azurerm_api_management_certificate.test.data}"
  password            = "${azurerm_api_management_certificate.test.password}"
}
`, template)
}
