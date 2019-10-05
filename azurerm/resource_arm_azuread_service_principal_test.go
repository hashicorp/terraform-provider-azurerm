package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMActiveDirectoryServicePrincipal_basic(t *testing.T) {
	resourceName := "azurerm_azuread_service_principal.test"
	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryServicePrincipal_basic(id)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryServicePrincipalExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
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

func TestAccAzureRMActiveDirectoryServicePrincipal_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_azuread_service_principal.test"
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryServicePrincipal_basic(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryServicePrincipalExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMActiveDirectoryServicePrincipal_requiresImport(id),
				ExpectError: testRequiresImportError("azurerm_azuread_service_principal"),
			},
		},
	})
}

func testCheckAzureRMActiveDirectoryServicePrincipalExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		client := testAccProvider.Meta().(*ArmClient).graph.ServicePrincipalsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, rs.Primary.ID)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure AD Service Principal %q does not exist", rs.Primary.ID)
			}
			return fmt.Errorf("Bad: Get on Azure AD servicePrincipalsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMActiveDirectoryServicePrincipalDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_azuread_service_principal" {
			continue
		}

		client := testAccProvider.Meta().(*ArmClient).graph.ServicePrincipalsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, rs.Primary.ID)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Azure AD Service Principal still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMActiveDirectoryServicePrincipal_basic(id string) string {
	return fmt.Sprintf(`
resource "azurerm_azuread_application" "test" {
  name = "acctestspa%s"
}

resource "azurerm_azuread_service_principal" "test" {
  application_id = "${azurerm_azuread_application.test.application_id}"
}
`, id)
}

func testAccAzureRMActiveDirectoryServicePrincipal_requiresImport(id string) string {
	template := testAccAzureRMActiveDirectoryServicePrincipal_basic(id)
	return fmt.Sprintf(`
%s

resource "azurerm_azuread_service_principal" "import" {
  application_id = "${azurerm_azuread_service_principal.test.application_id}"
}
`, template)
}
