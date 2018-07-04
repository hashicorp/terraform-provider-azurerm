package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMActiveDirectoryApplication_simple(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_simple(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", id),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("http://%s", id)),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccAzureRMActiveDirectoryApplication_advanced(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_advanced(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", id),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("http://homepage-%s", id)),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccAzureRMActiveDirectoryApplication_updateAdvanced(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_simple(id)

	updatedId := uuid.New().String()
	updatedConfig := testAccAzureRMActiveDirectoryApplication_advanced(updatedId)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", id),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("http://%s", id)),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedId),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("http://homepage-%s", updatedId)),
				),
			},
		},
	})
}

func testCheckAzureRMActiveDirectoryApplicationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		id := rs.Primary.Attributes["id"]

		client := testAccProvider.Meta().(*ArmClient).applicationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, id)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure AD Application %q does not exist", id)
			}
			return fmt.Errorf("Bad: Get on Azure AD applicationsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMActiveDirectoryApplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_azuread_application" {
			continue
		}

		id := rs.Primary.Attributes["id"]

		client := testAccProvider.Meta().(*ArmClient).applicationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, id)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Azure AD Application still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMActiveDirectoryApplication_simple(id string) string {
	return fmt.Sprintf(`
resource "azurerm_azuread_application" "test" {
  name = "%s"
}
`, id)
}

func testAccAzureRMActiveDirectoryApplication_advanced(id string) string {
	return fmt.Sprintf(`
resource "azurerm_azuread_application" "test" {
  name = "%s"
  homepage = "http://homepage-%s"
  identifier_uris = ["http://uri-%s"]
  reply_urls = ["http://replyurl-%s"]
  available_to_other_tenants = false
  oauth2_allow_implicit_flow = true
}
`, id, id, id, id)
}
