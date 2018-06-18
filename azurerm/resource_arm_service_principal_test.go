package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMServicePrincipal_simple(t *testing.T) {
	resourceName := "azurerm_service_principal.test"
	id := uuid.New().String()
	config := testAccAzureRMServicePrincipal_simple(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMServicePrincipalExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
					resource.TestCheckResourceAttr(resourceName, "display_name", id),
				),
			},
		},
	})
}

func testCheckAzureRMServicePrincipalExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		objectId := rs.Primary.Attributes["object_id"]

		client := testAccProvider.Meta().(*ArmClient).servicePrincipalsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, objectId)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Service Principal %q does not exist", objectId)
			}
			return fmt.Errorf("Bad: Get on servicePrincipalsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMServicePrincipalDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_service_principal" {
			continue
		}

		objectId := rs.Primary.Attributes["object_id"]

		client := testAccProvider.Meta().(*ArmClient).servicePrincipalsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, objectId)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Service Principal still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMServicePrincipal_simple(id string) string {
	return fmt.Sprintf(`
resource "azurerm_ad_application" "test" {
  display_name = "%s"
}

resource "azurerm_service_principal" "test" {
  app_id = "${azurerm_ad_application.test.app_id}"
}
`, id)
}
