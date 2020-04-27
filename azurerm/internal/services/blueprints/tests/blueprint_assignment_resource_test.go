package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/blueprints/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccBlueprintAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBlueprintAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlueprintAssignment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckBlueprintAssignmentExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, ""),
				),
			},
		},
	})
}

func testCheckBlueprintAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Blueprint Assignment not found: %s", resourceName)
		}
		id, err := parse.AssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Blueprints.AssignmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, id.Scope, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Blueprint Assignment %q (scope %q) was not found", id.Name, id.Scope)
			}
			return fmt.Errorf("Bad: Get on Blueprint Assignment %q (scope %q): %+v", id.Name, id.Scope, err)
		}
		return nil
	}
}

func testCheckAzureRMBlueprintAssignmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Blueprints.AssignmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_blueprint_assignment" {
			continue
		}

		id, err := parse.AssignmentID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.Scope, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Blueprint.AssignmentClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

// This test is intentionally broken as the AccTest environment is not currently capable of supporting this test
func testAccBlueprintAssignment_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_blueprint_definition" "test" {
  name       = "testAcc_basicSubscription"
  scope_type = "subscription"
  scope_name = data.azurerm_client_config.current.subscription_id
}

data "azurerm_blueprint_published_version" "test" {
  subscription_id = data.azurerm_client_config.current.subscription_id
  blueprint_name  = data.azurerm_blueprint_definition.test.name
  version         = "testAcc"
}

resource "azurerm_blueprint_assignment" "test" {
  # name     = "testAccBPAssignment"
  scope_type = "subscription"
  scope      = data.azurerm_client_config.current.subscription_id
  location   = "%s"
  identity {
    type                     = "UserAssigned"
    user_assigned_identities = ["00000000-0000-0000-0000-000000000000"]
  }
  version_id = data.azurerm_blueprint_published_version.test.id
}
`, data.Locations.Primary)
}
