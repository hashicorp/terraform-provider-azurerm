package tests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMActiveDirectoryApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azuread_application", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryApplication_basic(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(data.ResourceName, "homepage", fmt.Sprintf("https://acctest%s", id)),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMActiveDirectoryApplication_availableToOtherTenants(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azuread_application", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryApplication_availableToOtherTenants(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "available_to_other_tenants", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMActiveDirectoryApplication_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azuread_application", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryApplication_complete(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(data.ResourceName, "homepage", fmt.Sprintf("https://homepage-%s", id)),
					resource.TestCheckResourceAttr(data.ResourceName, "identifier_uris.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "reply_urls.#", "1"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "application_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMActiveDirectoryApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_azuread_application", "test")
	id := uuid.New().String()

	updatedId := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryApplication_basic(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(data.ResourceName, "homepage", fmt.Sprintf("https://acctest%s", id)),
					resource.TestCheckResourceAttr(data.ResourceName, "identifier_uris.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "reply_urls.#", "0"),
				),
			},
			{
				Config: testAccAzureRMActiveDirectoryApplication_complete(updatedId),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctest%s", updatedId)),
					resource.TestCheckResourceAttr(data.ResourceName, "homepage", fmt.Sprintf("https://homepage-%s", updatedId)),
					resource.TestCheckResourceAttr(data.ResourceName, "identifier_uris.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "reply_urls.#", "1"),
				),
			},
		},
	})
}

func testCheckAzureRMActiveDirectoryApplicationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Graph.ApplicationsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

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
	client := acceptance.AzureProvider.Meta().(*clients.Client).Graph.ApplicationsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_azuread_application" {
			continue
		}

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
