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

func TestAccAzureRMApiManagementOpenIDConnectProvider_basic(t *testing.T) {
	resourceName := "azurerm_api_management_openid_connect_provider.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementOpenIDConnectProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementOpenIDConnectProvider_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementOpenIDConnectProviderExists(resourceName),
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

func TestAccAzureRMApiManagementOpenIDConnectProvider_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_api_management_openid_connect_provider.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementOpenIDConnectProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementOpenIDConnectProvider_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementOpenIDConnectProviderExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMApiManagementOpenIDConnectProvider_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_api_management_openid_connect_provider"),
			},
		},
	})
}

func TestAccAzureRMApiManagementOpenIDConnectProvider_update(t *testing.T) {
	resourceName := "azurerm_api_management_openid_connect_provider.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApiManagementOpenIDConnectProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMApiManagementOpenIDConnectProvider_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementOpenIDConnectProviderExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMApiManagementOpenIDConnectProvider_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMApiManagementOpenIDConnectProviderExists(resourceName),
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

func testCheckAzureRMApiManagementOpenIDConnectProviderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("API Management OpenID Connect Provider not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		client := testAccProvider.Meta().(*ArmClient).apiManagement.OpenIdConnectClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, serviceName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: OpenID Connect Provider %q (Resource Group %q / API Management Service %q) does not exist", name, resourceGroup, serviceName)
			}
			return fmt.Errorf("Bad: Get on apiManagement.OpenIdConnectClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMApiManagementOpenIDConnectProviderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).apiManagement.OpenIdConnectClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_api_management_openid_connect_provider" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["api_management_name"]

		if resp, err := client.Get(ctx, resourceGroup, serviceName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on apiManagement.OpenIdConnectClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMApiManagementOpenIDConnectProvider_basic(rInt int, location string) string {
	template := testAccAzureRMApiManagementOpenIDConnectProvider_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  client_id           = "00001111-2222-3333-%d"
  client_secret       = "%d-cwdavsxbacsaxZX-%d"
  display_name        = "Initial Name"
  metadata_endpoint   = "https://azacctest.hashicorptest.com/example/foo"
}
`, template, rInt, rInt, rInt, rInt)
}

func testAccAzureRMApiManagementOpenIDConnectProvider_requiresImport(rInt int, location string) string {
	template := testAccAzureRMApiManagementOpenIDConnectProvider_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "import" {
  name                = "${azurerm_api_management_openid_connect_provider.test.name}"
  api_management_name = "${azurerm_api_management_openid_connect_provider.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_openid_connect_provider.test.resource_group_name}"
  client_id           = "${azurerm_api_management_openid_connect_provider.test.client_id}"
  client_secret       = "${azurerm_api_management_openid_connect_provider.test.client_secret}"
  display_name        = "${azurerm_api_management_openid_connect_provider.test.display_name}"
  metadata_endpoint   = "${azurerm_api_management_openid_connect_provider.test.metadata_endpoint}"
}
`, template)
}

func testAccAzureRMApiManagementOpenIDConnectProvider_complete(rInt int, location string) string {
	template := testAccAzureRMApiManagementOpenIDConnectProvider_template(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_openid_connect_provider" "test" {
  name                = "acctest-%d"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  client_id           = "00001111-3333-2222-%d"
  client_secret       = "%d-423egvwdcsjx-%d"
  display_name        = "Updated Name"
  description         = "Example description"
  metadata_endpoint   = "https://azacctest.hashicorptest.com/example/updated"
}
`, template, rInt, rInt, rInt, rInt)
}

func testAccAzureRMApiManagementOpenIDConnectProvider_template(rInt int, location string) string {
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
