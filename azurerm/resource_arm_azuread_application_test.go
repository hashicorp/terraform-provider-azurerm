package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMActiveDirectoryApplication_basic(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_basic(id)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("https://acctest%s", id)),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
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

func TestAccAzureRMActiveDirectoryApplication_availableToOtherTenants(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_availableToOtherTenants(id)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "available_to_other_tenants", "true"),
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

func TestAccAzureRMActiveDirectoryApplication_complete(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_complete(id)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("https://homepage-%s", id)),
					resource.TestCheckResourceAttr(resourceName, "identifier_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "reply_urls.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
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

func TestAccAzureRMActiveDirectoryApplication_update(t *testing.T) {
	resourceName := "azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryApplication_basic(id)

	updatedId := uuid.New().String()
	updatedConfig := testAccAzureRMActiveDirectoryApplication_complete(updatedId)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("https://acctest%s", id)),
					resource.TestCheckResourceAttr(resourceName, "identifier_uris.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "reply_urls.#", "0"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest%s", updatedId)),
					resource.TestCheckResourceAttr(resourceName, "homepage", fmt.Sprintf("https://homepage-%s", updatedId)),
					resource.TestCheckResourceAttr(resourceName, "identifier_uris.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "reply_urls.#", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMActiveDirectoryApplicationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		client := testAccProvider.Meta().(*ArmClient).graph.ApplicationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, rs.Primary.ID)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure AD Application %q does not exist", rs.Primary.ID)
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

		client := testAccProvider.Meta().(*ArmClient).graph.ApplicationsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, rs.Primary.ID)

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

func testAccAzureRMActiveDirectoryApplication_basic(id string) string {
	return fmt.Sprintf(`
resource "azurerm_azuread_application" "test" {
  name = "acctest%s"
}
`, id)
}

func testAccAzureRMActiveDirectoryApplication_availableToOtherTenants(id string) string {
	return fmt.Sprintf(`
resource "azurerm_azuread_application" "test" {
  name                       = "acctest%s"
  identifier_uris            = ["https://%s.hashicorptest.com"]
  available_to_other_tenants = true
}
`, id, id)
}

func testAccAzureRMActiveDirectoryApplication_complete(id string) string {
	return fmt.Sprintf(`
resource "azurerm_azuread_application" "test" {
  name                       = "acctest%s"
  homepage                   = "https://homepage-%s"
  identifier_uris            = ["http://%s.hashicorptest.com/00000000-0000-0000-0000-00000000"]
  reply_urls                 = ["http://%s.hashicorptest.com"]
  oauth2_allow_implicit_flow = true
}
`, id, id, id, id)
}
