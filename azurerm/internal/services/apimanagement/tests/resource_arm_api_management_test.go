package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagement_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

// Remove in 2.0
func TestAccAzureRMApiManagement_basicClassic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_basicClassic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

// Remove in 2.0
func TestAccAzureRMApiManagement_basicNotDefined(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMApiManagement_basicNotDefined(data),
				ExpectError: regexp.MustCompile("either 'sku_name' or 'sku' must be defined in the configuration file"),
			},
		},
	})
}

func TestAccAzureRMApiManagement_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMApiManagement_requiresImport),
		},
	})
}

func TestAccAzureRMApiManagement_customProps(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_customProps(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.Acceptance", "Test"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "public_ip_addresses.#"),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate", // not returned from API, sensitive
					"hostname_configuration.0.portal.0.certificate",          // not returned from API, sensitive
					"hostname_configuration.0.portal.0.certificate_password", // not returned from API, sensitive
					"hostname_configuration.0.proxy.0.certificate",           // not returned from API, sensitive
					"hostname_configuration.0.proxy.0.certificate_password",  // not returned from API, sensitive
					"hostname_configuration.0.proxy.1.certificate",           // not returned from API, sensitive
					"hostname_configuration.0.proxy.1.certificate_password",  // not returned from API, sensitive
				},
			},
		},
	})
}

func TestAccAzureRMApiManagement_signInSignUpSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_signInSignUpSettings(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMApiManagement_policy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_policyXmlContent(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMApiManagement_policyXmlLink(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"policy.0.xml_link",
				},
			},
			{
				Config: testAccAzureRMApiManagement_policyRemoved(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMApiManagementDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ServiceClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMApiManagementExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		apiMangementName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Api Management: %s", apiMangementName)
		}

		conn := acceptance.AzureProvider.Meta().(*clients.Client).ApiManagement.ServiceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := conn.Get(ctx, resourceGroup, apiMangementName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Api Management %q (resource group: %q) does not exist", apiMangementName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on apiManagementClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMApiManagement_basic(data acceptance.TestData) string {
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

  sku_name = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

// Remove in 2.0
func testAccAzureRMApiManagement_basicClassic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

// Remove in 2.0
func testAccAzureRMApiManagement_basicNotDefined(data acceptance.TestData) string {
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
}
`, data.RandomInteger, data.Locations, data.RandomInteger)
}

func testAccAzureRMApiManagement_policyXmlContent(data acceptance.TestData) string {
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

  sku_name = "Developer_1"

  policy {
    xml_content = <<XML
<policies>
  <inbound>
    <find-and-replace from="xyz" to="abc" />
  </inbound>
</policies>
XML
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagement_policyXmlLink(data acceptance.TestData) string {
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

  sku_name = "Developer_1"

  policy {
    xml_link = "https://gist.githubusercontent.com/tombuildsstuff/4f58581599d2c9f64b236f505a361a67/raw/0d29dcb0167af1e5afe4bd52a6d7f69ba1e05e1f/example.xml"
  }
}
`, data.RandomInteger, data.Locations, data.RandomInteger)
}

func testAccAzureRMApiManagement_policyRemoved(data acceptance.TestData) string {
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

  sku_name = "Developer_1"

  policy = []
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagement_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMApiManagement_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "import" {
  name                = "${azurerm_api_management.test.name}"
  location            = "${azurerm_api_management.test.location}"
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  publisher_name      = "${azurerm_api_management.test.publisher_name}"
  publisher_email     = "${azurerm_api_management.test.publisher_email}"

  sku_name = "Developer_1"
}
`, template)
}

func testAccAzureRMApiManagement_customProps(data acceptance.TestData) string {
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

  sku_name = "Developer_1"

  security {
    enable_frontend_tls10     = true
    enable_triple_des_ciphers = true
  }
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

func testAccAzureRMApiManagement_signInSignUpSettings(data acceptance.TestData) string {
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

  sku_name = "Developer_1"

  sign_in {
    enabled = true
  }

  sign_up {
    enabled = true

    terms_of_service {
      enabled          = true
      consent_required = false
      text             = "Lorem Ipsum Dolor Morty"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMApiManagement_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-api1-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-api1-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test3" {
  name     = "acctestRG-api1-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                      = "acctestAM-%d"
  publisher_name            = "pub1"
  publisher_email           = "pub1@email.com"
  notification_sender_email = "notification@email.com"

  additional_location {
    location = "${azurerm_resource_group.test2.location}"
  }

  additional_location {
    location = "${azurerm_resource_group.test3.location}"
  }

  certificate {
    encoded_certificate  = "${filebase64("testdata/api_management_api_test.pfx")}"
    certificate_password = "terraform"
    store_name           = "CertificateAuthority"
  }

  certificate {
    encoded_certificate  = "${filebase64("testdata/api_management_api_test.pfx")}"
    certificate_password = "terraform"
    store_name           = "Root"
  }

  security {
    enable_backend_tls11      = true
    enable_backend_ssl30      = true
    enable_backend_tls10      = true
    enable_frontend_ssl30     = true
    enable_frontend_tls10     = true
    enable_frontend_tls11     = true
    enable_triple_des_ciphers = true
  }

  hostname_configuration {
    proxy {
      host_name                    = "api.terraform.io"
      certificate                  = "${filebase64("testdata/api_management_api_test.pfx")}"
      certificate_password         = "terraform"
      default_ssl_binding          = true
      negotiate_client_certificate = false
    }

    proxy {
      host_name                    = "api2.terraform.io"
      certificate                  = "${filebase64("testdata/api_management_api2_test.pfx")}"
      certificate_password         = "terraform"
      negotiate_client_certificate = true
    }

    portal {
      host_name            = "portal.terraform.io"
      certificate          = "${filebase64("testdata/api_management_portal_test.pfx")}"
      certificate_password = "terraform"
    }
  }

  sku_name = "Premium_1"

  tags = {
    "Acceptance" = "Test"
  }

  location            = "${azurerm_resource_group.test1.location}"
  resource_group_name = "${azurerm_resource_group.test1.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary, data.RandomInteger, data.Locations.Ternary, data.RandomInteger)
}
