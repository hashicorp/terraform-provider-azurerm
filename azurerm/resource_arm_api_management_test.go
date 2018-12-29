package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMApiManagement_basic(t *testing.T) {
	resourceName := "azurerm_api_management.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApiManagement_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApiManagement_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management.test"
	ri := acctest.RandInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagement_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagement_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management"),
			},
		},
	})
}

func TestAccAzureRMApiManagement_customProps(t *testing.T) {
	resourceName := "azurerm_api_management.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApiManagement_customProps(ri, testAltLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMApiManagement_complete(t *testing.T) {
	resourceName := "azurerm_api_management.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApiManagement_complete(ri, testLocation(), testAltLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.Acceptance", "Test"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip_addresses.#"),
				),
			},
			{
				ResourceName:      resourceName,
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

func testCheckAzureRMApiManagementDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).apiManagementServiceClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testCheckAzureRMApiManagementExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		apiMangementName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Api Management: %s", apiMangementName)
		}

		conn := testAccProvider.Meta().(*ArmClient).apiManagementServiceClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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

func testAccAzureRMApiManagement_basic(rInt int, location string) string {
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
`, rInt, location, rInt)
}

func testAccAzureRMApiManagement_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagement_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management" "import" {
  name                = "${azurerm_api_management.test.name}"
  location            = "${azurerm_api_management.test.location}"
  resource_group_name = "${azurerm_api_management.test.resource_group_name}"
  publisher_name      = "${azurerm_api_management.test.publisher_name}"
  publisher_email     = "${azurerm_api_management.test.publisher_email}"

  sku {
    name     = "Developer"
    capacity = 1
  }
}
`, template)
}

func testAccAzureRMApiManagement_customProps(rInt int, location string) string {
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

  security {
    disable_frontend_tls10     = true
    disable_triple_des_chipers = true
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMApiManagement_complete(rInt int, location string, altLocation string) string {
	template := testAccAzureRMApiManagement_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-%d"
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

  certificate {
    encoded_certificate  = "${base64encode(file("testdata/api_management_api_test.pfx"))}"
    certificate_password = "terraform"
    store_name           = "CertificateAuthority"
  }

  certificate {
    encoded_certificate  = "${base64encode(file("testdata/api_management_api_test.pfx"))}"
    certificate_password = "terraform"
    store_name           = "Root"
  }

  security {
    disable_backend_tls11 = true
  }

  hostname_configuration {
    proxy {
      host_name                    = "api.terraform.io"
      certificate                  = "${base64encode(file("testdata/api_management_api_test.pfx"))}"
      certificate_password         = "terraform"
      default_ssl_binding          = true
      negotiate_client_certificate = false
    }

    proxy {
      host_name                    = "api2.terraform.io"
      certificate                  = "${base64encode(file("testdata/api_management_api2_test.pfx"))}"
      certificate_password         = "terraform"
      negotiate_client_certificate = true
    }

    portal {
      host_name            = "portal.terraform.io"
      certificate          = "${base64encode(file("testdata/api_management_portal_test.pfx"))}"
      certificate_password = "terraform"
    }
	}
	
	virtual_network_type	= "External"

	virtual_network_configuration {
		subnet_id	= "${azurerm_subnet.test.id}"
	}

  sku {
    name     = "Premium"
    capacity = 1
  }

  tags {
    "Acceptance" = "Test"
  }

  location            = "${azurerm_resource_group.test1.location}"
  resource_group_name = "${azurerm_resource_group.test1.name}"
}
`, template, rInt, location, rInt, altLocation, rInt)
}

func testAccAzureRMApiManagement_template(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet-%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.0.0/24"
}
`, rInt, location, rInt, rInt)
}
