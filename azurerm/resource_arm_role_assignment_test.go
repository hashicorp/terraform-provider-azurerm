package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRoleAssignment(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning a couple at a time
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"emptyName":   testAccAzureRMRoleAssignment_emptyName,
			"roleName":    testAccAzureRMRoleAssignment_roleName,
			"dataActions": testAccAzureRMRoleAssignment_dataActions,
			"builtin":     testAccAzureRMRoleAssignment_builtin,
			"custom":      testAccAzureRMRoleAssignment_custom,
		},
		"import": {
			"basic":  testAccAzureRMRoleAssignment_importBasic,
			"custom": testAccAzureRMRoleAssignment_importCustom,
		},
	}

	for group, m := range testCases {
		m := m
		t.Run(group, func(t *testing.T) {
			for name, tc := range m {
				tc := tc
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccAzureRMRoleAssignment_emptyName(t *testing.T) {
	resourceName := "azurerm_role_assignment.test"
	config := testAccAzureRMRoleAssignment_emptyNameConfig()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
		},
	})
}

func testAccAzureRMRoleAssignment_roleName(t *testing.T) {
	id := uuid.New().String()
	resourceName := "azurerm_role_assignment.test"
	config := testAccAzureRMRoleAssignment_roleNameConfig(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "role_definition_id"),
				),
			},
		},
	})
}

func testAccAzureRMRoleAssignment_dataActions(t *testing.T) {
	id := uuid.New().String()
	resourceName := "azurerm_role_assignment.test"
	config := testAccAzureRMRoleAssignment_dataActionsConfig(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "role_definition_id"),
				),
			},
		},
	})
}

func testAccAzureRMRoleAssignment_builtin(t *testing.T) {
	id := uuid.New().String()
	config := testAccAzureRMRoleAssignment_builtinConfig(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists("azurerm_role_assignment.test"),
				),
			},
		},
	})
}

func testAccAzureRMRoleAssignment_custom(t *testing.T) {
	roleDefinitionId := uuid.New().String()
	roleAssignmentId := uuid.New().String()
	rInt := acctest.RandInt()
	config := testAccAzureRMRoleAssignment_customConfig(roleDefinitionId, roleAssignmentId, rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists("azurerm_role_assignment.test"),
				),
			},
		},
	})
}

func testCheckAzureRMRoleAssignmentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		scope := rs.Primary.Attributes["scope"]
		roleAssignmentName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).roleAssignmentsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, scope, roleAssignmentName)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Role Assignment %q (Scope: %q) does not exist", name, scope)
			}
			return fmt.Errorf("Bad: Get on roleDefinitionsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRoleAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_role_assignment" {
			continue
		}

		scope := rs.Primary.Attributes["scope"]
		roleAssignmentName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).roleAssignmentsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, scope, roleAssignmentName)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Role Definition still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMRoleAssignment_emptyNameConfig() string {
	return `
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_builtin_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_role_assignment" "test" {
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_builtin_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.service_principal_object_id}"
}
`
}

func testAccAzureRMRoleAssignment_roleNameConfig(id string) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = "${data.azurerm_subscription.primary.id}"
  role_definition_name = "Log Analytics Reader"
  principal_id         = "${data.azurerm_client_config.test.service_principal_object_id}"
}
`, id)
}

func testAccAzureRMRoleAssignment_dataActionsConfig(id string) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = "${data.azurerm_subscription.primary.id}"
  role_definition_name = "Virtual Machine User Login"
  principal_id         = "${data.azurerm_client_config.test.service_principal_object_id}"
}
`, id)
}

func testAccAzureRMRoleAssignment_builtinConfig(id string) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_builtin_role_definition" "test" {
  name = "Site Recovery Reader"
}

resource "azurerm_role_assignment" "test" {
  name               = "%s"
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_builtin_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.service_principal_object_id}"
}
`, id)
}

func testAccAzureRMRoleAssignment_customConfig(roleDefinitionId string, roleAssignmentId string, rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = "${data.azurerm_subscription.primary.id}"
  description        = "Created by the Role Assignment Acceptance Test"

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}

resource "azurerm_role_assignment" "test" {
  name               = "%s"
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${azurerm_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.service_principal_object_id}"
}
`, roleDefinitionId, rInt, roleAssignmentId)
}
