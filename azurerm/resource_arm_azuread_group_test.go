package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMActiveDirectoryGroup_basic(t *testing.T) {
	resourceName := "azurerm_azuread_group.test"
	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryGroup(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest%s", id)),
				),
			},
		},
	})
}

func TestAccAzureRMActiveDirectoryGroup_complete(t *testing.T) {
	resourceName := "azurerm_azuread_group.test"
	id := uuid.New().String()
	config := testAccAzureRMActiveDirectoryGroup(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctest%s", id)),
				),
			},
		},
	})
}

func testCheckAzureRMActiveDirectoryGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		client := testAccProvider.Meta().(*ArmClient).groupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, rs.Primary.ID)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure AD Group %q does not exist", rs.Primary.ID)
			}
			return fmt.Errorf("Bad: Get on Azure AD groupsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMActiveDirectoryGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_azuread_group" {
			continue
		}

		client := testAccProvider.Meta().(*ArmClient).groupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, rs.Primary.ID)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Azure AD group still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMActiveDirectoryGroup(id string) string {
	return fmt.Sprintf(`
resource "azurerm_azuread_group" "test" {
  name = "acctest%s"
}
`, id)
}
