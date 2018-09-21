package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAzureRMApiManagementName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "a",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "api1",
			ErrCount: 0,
		},
		{
			Value:    "company-api",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloworld21!",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := azure.ValidateApiManagementName(tc.Value, "azurerm_api_management")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Api Management Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAccAzureRMApiManagement_basic(t *testing.T) {
	resourceName := "azurerm_api_management.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApiManagement_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists("azurerm_api_management.test"),
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

func TestAccAzureRMApiManagement_customProps(t *testing.T) {
	resourceName := "azurerm_api_management.test"
	ri := acctest.RandInt()
	config := testAccAzureRMApiManagement_customProps(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists("azurerm_api_management.test"),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementExists("azurerm_api_management.test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate",                                             // not returned from API, sensitive
					"hostname_configurations.0.portal.0.certificate",          // not returned from API, sensitive
					"hostname_configurations.0.portal.0.certificate_password", // not returned from API, sensitive
					"hostname_configurations.0.proxy.0.certificate",           // not returned from API, sensitive
					"hostname_configurations.0.proxy.0.certificate_password",  // not returned from API, sensitive
					"hostname_configurations.0.proxy.1.certificate",           // not returned from API, sensitive
					"hostname_configurations.0.proxy.1.certificate_password",  // not returned from API, sensitive
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
  name     = "amtestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku {
    name = "Developer"
  }
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMApiManagement_customProps(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "amtestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name            = "acctestAM-%d"
  publisher_name  = "pub1"
  publisher_email = "pub1@email.com"

  sku {
    name = "Developer"
  }

  security {
    disable_frontend_tls10     = true
    disable_triple_des_chipers = true
  }

	location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMApiManagement_complete(rInt int, location string, altLocation string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test1" {
  name     = "amtestRG1-%d"
  location = "%s"
}

resource "azurerm_resource_group" "test2" {
  name     = "amtestRG2-%d"
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

  hostname_configurations {
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

  sku {
    name = "Premium"
  }

  tags {
    test = "true"
  }

  location            = "${azurerm_resource_group.test1.location}"
  resource_group_name = "${azurerm_resource_group.test1.name}"
}
`, rInt, location, rInt, altLocation, rInt)
}
