package authorization_test

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

func TestAccRoleAssignment(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning a couple at a time
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"emptyName":      testAccRoleAssignment_emptyName,
			"roleName":       testAccRoleAssignment_roleName,
			"dataActions":    testAccRoleAssignment_dataActions,
			"builtin":        testAccRoleAssignment_builtin,
			"custom":         testAccRoleAssignment_custom,
			"requiresImport": testAccRoleAssignment_requiresImport,
		},
		"assignment": {
			"sp":     testAccActiveDirectoryServicePrincipal_servicePrincipal,
			"spType": testAccActiveDirectoryServicePrincipal_servicePrincipalWithType,
			"group":  testAccActiveDirectoryServicePrincipal_group,
		},
		"management": {
			"assign": testAccRoleAssignment_managementGroup,
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

func testAccRoleAssignment_emptyName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignment_emptyNameConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"skip_service_principal_aad_check",
				},
			},
		},
	})
}

func testAccRoleAssignment_roleName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	id := uuid.New().String()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignment_roleNameConfig(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_definition_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_definition_name", "Log Analytics Reader"),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"skip_service_principal_aad_check",
				},
			},
		},
	})
}

func testAccRoleAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignment_roleNameConfig(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_definition_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_definition_name", "Log Analytics Reader"),
				),
			},
			{
				Config:      testAccRoleAssignment_requiresImportConfig(id),
				ExpectError: acceptance.RequiresImportError("azurerm_role_assignment"),
			},
		},
	})
}

func testAccRoleAssignment_dataActions(t *testing.T) {
	id := uuid.New().String()
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignment_dataActionsConfig(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_definition_id"),
				),
			},
			data.ImportStep("skip_service_principal_aad_check"),
		},
	})
}

func testAccRoleAssignment_builtin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	id := uuid.New().String()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignment_builtinConfig(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep("skip_service_principal_aad_check"),
		},
	})
}

func testAccRoleAssignment_custom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	roleDefinitionId := uuid.New().String()
	roleAssignmentId := uuid.New().String()
	rInt := acceptance.RandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignment_customConfig(roleDefinitionId, roleAssignmentId, rInt),
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleAssignmentExists(data.ResourceName),
				),
			},
			{
				ResourceName:      data.ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"skip_service_principal_aad_check",
				},
			},
		},
	})
}

func testAccActiveDirectoryServicePrincipal_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignment_servicePrincipal(ri, id),
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleAssignmentExists("azurerm_role_assignment.test"),
					resource.TestCheckResourceAttr(data.ResourceName, "principal_type", "ServicePrincipal"),
				),
			},
		},
	})
}

func testAccActiveDirectoryServicePrincipal_servicePrincipalWithType(t *testing.T) {
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignment_servicePrincipalWithType(ri, id),
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleAssignmentExists("azurerm_role_assignment.test"),
				),
			},
		},
	})
}

func testAccActiveDirectoryServicePrincipal_group(t *testing.T) {
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignment_group(ri, id),
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleAssignmentExists("azurerm_role_assignment.test"),
				),
			},
		},
	})
}

func testCheckRoleAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Authorization.RoleAssignmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		scope := rs.Primary.Attributes["scope"]
		roleAssignmentName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, scope, roleAssignmentName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Role Assignment %q (Scope: %q) does not exist", roleAssignmentName, scope)
			}
			return fmt.Errorf("Bad: Get on roleDefinitionsClient: %+v", err)
		}

		return nil
	}
}

func testCheckRoleAssignmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Authorization.RoleAssignmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_role_assignment" {
			continue
		}

		scope := rs.Primary.Attributes["scope"]
		roleAssignmentName := rs.Primary.Attributes["name"]

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

// TODO - "real" management group with appropriate required for testing
func testAccRoleAssignment_managementGroup(t *testing.T) {
	groupId := uuid.New().String()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleAssignment_managementGroupConfig(groupId),
				Check: resource.ComposeTestCheckFunc(
					testCheckRoleAssignmentExists("azurerm_role_assignment.test"),
				),
			},
		},
	})
}

func testAccRoleAssignment_emptyNameConfig() string {
	return `
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_role_assignment" "test" {
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.object_id}"
}
`
}

func testAccRoleAssignment_roleNameConfig(id string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Log Analytics Reader"
  principal_id         = data.azurerm_client_config.test.object_id
}
`, id)
}

func testAccRoleAssignment_requiresImportConfig(id string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_role_assignment" "import" {
  name                 = azurerm_role_assignment.test.name
  scope                = azurerm_role_assignment.test.id
  role_definition_name = azurerm_role_assignment.test.role_definition_name
  principal_id         = azurerm_role_assignment.test.principal_id
}
`, testAccRoleAssignment_roleNameConfig(id))
}

func testAccRoleAssignment_dataActionsConfig(id string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Virtual Machine User Login"
  principal_id         = data.azurerm_client_config.test.object_id
}
`, id)
}

func testAccRoleAssignment_builtinConfig(id string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

data "azurerm_role_definition" "test" {
  name = "Site Recovery Reader"
}

resource "azurerm_role_assignment" "test" {
  name               = "%s"
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_role_definition.test.id}"
  principal_id       = data.azurerm_client_config.test.object_id
}
`, id)
}

func testAccRoleAssignment_customConfig(roleDefinitionId string, roleAssignmentId string, rInt int) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = data.azurerm_subscription.primary.id
  description        = "Created by the Role Assignment Acceptance Test"

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id,
  ]
}

resource "azurerm_role_assignment" "test" {
  name               = "%s"
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = azurerm_role_definition.test.role_definition_resource_id
  principal_id       = data.azurerm_client_config.test.object_id
}
`, roleDefinitionId, rInt, roleAssignmentId)
}

func testAccRoleAssignment_servicePrincipal(rInt int, roleAssignmentID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {
}

resource "azuread_application" "test" {
  name = "acctestspa-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = data.azurerm_subscription.current.id
  role_definition_name = "Reader"
  principal_id         = azuread_service_principal.test.id
}
`, rInt, roleAssignmentID)
}

func testAccRoleAssignment_servicePrincipalWithType(rInt int, roleAssignmentID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {
}

resource "azuread_application" "test" {
  name = "acctestspa-%d"
}

resource "azuread_service_principal" "test" {
  application_id = azuread_application.test.application_id
}

resource "azurerm_role_assignment" "test" {
  name                             = "%s"
  scope                            = data.azurerm_subscription.current.id
  role_definition_name             = "Reader"
  principal_id                     = azuread_service_principal.test.id
  skip_service_principal_aad_check = true
}
`, rInt, roleAssignmentID)
}

func testAccRoleAssignment_group(rInt int, roleAssignmentID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {
}

resource "azuread_group" "test" {
  name = "acctestspa-%d"
}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = data.azurerm_subscription.current.id
  role_definition_name = "Reader"
  principal_id         = azuread_group.test.id
}
`, rInt, roleAssignmentID)
}

func testAccRoleAssignment_managementGroupConfig(groupId string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "test" {
}

data "azurerm_role_definition" "test" {
  name = "Monitoring Reader"
}

resource "azurerm_management_group" "test" {
  group_id = "%s"
}

resource "azurerm_role_assignment" "test" {
  scope              = azurerm_management_group.test.id
  role_definition_id = data.azurerm_role_definition.test.id
  principal_id       = data.azurerm_client_config.test.object_id
}
`, groupId)
}
