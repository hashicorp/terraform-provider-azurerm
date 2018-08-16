package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMManagementGroup_basic(t *testing.T) {
	resourceName := "azurerm_management_group.test"

	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRmManagementGroup_basic(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(resourceName),
				),
			},
		},
	})
}

func TestAccAzureRMManagementGroup_withSubscriptions(t *testing.T) {
	resourceName := "azurerm_management_group.test"

	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMManagementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRmManagementGroup_withSubscriptions(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMManagementGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subscription_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "subscription_ids.#", "00000000-1111-2222-3333-444444444444"),
				),
			},
		},
	})
}

func testCheckAzureRMManagementGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).managementGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		recurse := false
		resp, err := client.Get(ctx, name, "", &recurse, "", "no-cache")
		if err != nil {
			return fmt.Errorf("Bad: Get on policyDefinitionsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("policy does not exist: %s", name)
		}

		return nil
	}
}

func testCheckAzureRMManagementGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).managementGroupsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_management_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		recurse := false
		resp, err := client.Get(ctx, name, "", &recurse, "", "no-cache")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("policy still exists:%s", *resp.Name)
		}
	}

	return nil
}

func testAzureRmManagementGroup_basic(ri int) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  name         = "acctestmg-%d"
}
`, ri)
}

func testAzureRmManagementGroup_withSubscriptions(ri int) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  name         = "acctestmg-%d"
  subscription_ids = [
	  "00000000-1111-2222-3333-444444444444"
  ]
}
`, ri)
}
