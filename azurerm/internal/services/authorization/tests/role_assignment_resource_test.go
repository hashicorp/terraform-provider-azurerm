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

func TestAccAzureRMRoleAssignment(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being happy about provisioning a couple at a time
	testCases := map[string]map[string]func(t *testing.T){
		"basic": {
			"emptyName":      testAccAzureRMRoleAssignment_emptyName,
			"roleName":       testAccAzureRMRoleAssignment_roleName,
			"dataActions":    testAccAzureRMRoleAssignment_dataActions,
			"builtin":        testAccAzureRMRoleAssignment_builtin,
			"custom":         testAccAzureRMRoleAssignment_custom,
			"requiresImport": testAccAzureRMRoleAssignment_requiresImport,
		},
		"assignment": {
			"sp":     testAccAzureRMActiveDirectoryServicePrincipal_servicePrincipal,
			"spType": testAccAzureRMActiveDirectoryServicePrincipal_servicePrincipalWithType,
			"group":  testAccAzureRMActiveDirectoryServicePrincipal_group,
		},
		"management": {
			"assign": testAccAzureRMRoleAssignment_managementGroup,
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
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleAssignment_emptyNameConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(data.ResourceName),
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

func testAccAzureRMRoleAssignment_roleName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	id := uuid.New().String()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleAssignment_roleNameConfig(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(data.ResourceName),
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

func testAccAzureRMRoleAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleAssignment_roleNameConfig(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_definition_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "role_definition_name", "Log Analytics Reader"),
				),
			},
			{
				Config:      testAccAzureRMRoleAssignment_requiresImportConfig(id),
				ExpectError: acceptance.RequiresImportError("azurerm_role_assignment"),
			},
		},
	})
}

func testAccAzureRMRoleAssignment_dataActions(t *testing.T) {
	id := uuid.New().String()
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleAssignment_dataActionsConfig(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "role_definition_id"),
				),
			},
			data.ImportStep("skip_service_principal_aad_check"),
		},
	})
}

func testAccAzureRMRoleAssignment_builtin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	id := uuid.New().String()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleAssignment_builtinConfig(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep("skip_service_principal_aad_check"),
		},
	})
}

func testAccAzureRMRoleAssignment_custom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	roleDefinitionId := uuid.New().String()
	roleAssignmentId := uuid.New().String()
	rInt := acceptance.RandTimeInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleAssignment_customConfig(roleDefinitionId, roleAssignmentId, rInt),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(data.ResourceName),
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

func testAccAzureRMActiveDirectoryServicePrincipal_servicePrincipal(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_role_assignment", "test")
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleAssignment_servicePrincipal(ri, id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists("azurerm_role_assignment.test"),
					resource.TestCheckResourceAttr(data.ResourceName, "principal_type", "ServicePrincipal"),
				),
			},
		},
	})
}

func testAccAzureRMActiveDirectoryServicePrincipal_servicePrincipalWithType(t *testing.T) {
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleAssignment_servicePrincipalWithType(ri, id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists("azurerm_role_assignment.test"),
				),
			},
		},
	})
}

func testAccAzureRMActiveDirectoryServicePrincipal_group(t *testing.T) {
	ri := acceptance.RandTimeInt()
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleAssignment_group(ri, id),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists("azurerm_role_assignment.test"),
				),
			},
		},
	})
}

func testCheckAzureRMRoleAssignmentExists(resourceName string) resource.TestCheckFunc {
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

func testCheckAzureRMRoleAssignmentDestroy(s *terraform.State) error {
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
func testAccAzureRMRoleAssignment_managementGroup(t *testing.T) {
	groupId := uuid.New().String()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleAssignment_managementGroupConfig(groupId),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists("azurerm_role_assignment.test"),
				),
			},
		},
	})
}

func testAccAzureRMRoleAssignment_emptyNameConfig() string {
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

func testAccAzureRMRoleAssignment_roleNameConfig(id string) string {
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

func testAccAzureRMRoleAssignment_requiresImportConfig(id string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_role_assignment" "import" {
  name                 = azurerm_role_assignment.test.name
  scope                = azurerm_role_assignment.test.id
  role_definition_name = azurerm_role_assignment.test.role_definition_name
  principal_id         = azurerm_role_assignment.test.principal_id
}
`, testAccAzureRMRoleAssignment_roleNameConfig(id))
}

func testAccAzureRMRoleAssignment_dataActionsConfig(id string) string {
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

func testAccAzureRMRoleAssignment_builtinConfig(id string) string {
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

func testAccAzureRMRoleAssignment_customConfig(roleDefinitionId string, roleAssignmentId string, rInt int) string {
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

func testAccAzureRMRoleAssignment_servicePrincipal(rInt int, roleAssignmentID string) string {
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

func testAccAzureRMRoleAssignment_servicePrincipalWithType(rInt int, roleAssignmentID string) string {
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

func testAccAzureRMRoleAssignment_group(rInt int, roleAssignmentID string) string {
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

func testAccAzureRMRoleAssignment_managementGroupConfig(groupId string) string {
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
