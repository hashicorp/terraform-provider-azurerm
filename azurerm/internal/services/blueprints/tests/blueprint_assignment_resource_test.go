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

// Scenario: Basic BP, no artifacts etc.  Stored and applied at Subscription.
func TestAccBlueprintAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBlueprintAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlueprintAssignment_basic(data, "testAcc_basicSubscription"),
				Check: resource.ComposeTestCheckFunc(
					testCheckBlueprintAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

// Scenario: BP with RG's, locking and parameters/policies stored at Subscription, applied to subscription
func TestAccBlueprintAssignment_subscriptionComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBlueprintAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlueprintAssignment_subscriptionComplete(data, "testAcc_subscriptionComplete"),
				Check: resource.ComposeTestCheckFunc(
					testCheckBlueprintAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

// Scenario: BP stored at Root Management Group, applied to Subscription
func TestAccBlueprintAssignment_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMBlueprintAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlueprintAssignment_rootManagementGroup(data, "testAcc_basicRootManagementGroup"),
				Check: resource.ComposeTestCheckFunc(
					testCheckBlueprintAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
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

func testAccBlueprintAssignment_basic(data acceptance.TestData, bpName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_subscription" "test" {}

data "azurerm_blueprint_definition" "test" {
  name       = "%s"
  scope_type = "subscriptions"
  scope_name = data.azurerm_client_config.current.subscription_id
}

data "azurerm_blueprint_published_version" "test" {
  subscription_id = data.azurerm_client_config.current.subscription_id
  blueprint_name  = data.azurerm_blueprint_definition.test.name
  version         = "v0.1_testAcc"
}

resource "azurerm_resource_group" "test" {
  name     = "accTestRG-bp-%d"
  location = "westeurope"
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "bp-user-%d"
}

resource "azurerm_role_assignment" "test" {
  scope                = data.azurerm_subscription.test.id
  role_definition_name = "Blueprint Operator"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_blueprint_assignment" "test" {
  name       = "testAccBPAssignment"
  scope_type = "subscriptions"
  scope      = data.azurerm_client_config.current.subscription_id
  location   = "%s"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  version_id = data.azurerm_blueprint_published_version.test.id

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, bpName, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

// This test config creates a UM-MSI and assigns Owner to the target subscription.  This is necessary due to the changes
// the referenced Blueprint Version needs to make to successfully apply.  If the test panics or otherwise fails,
// Dangling resources can include the Role Assignment(s) at the Subscription, which will need to be removed
func testAccBlueprintAssignment_subscriptionComplete(data acceptance.TestData, bpName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_subscription" "test" {}

data "azurerm_blueprint_definition" "test" {
  name       = "%s"
  scope_type = "subscriptions"
  scope_name = data.azurerm_client_config.current.subscription_id
}

data "azurerm_blueprint_published_version" "test" {
  subscription_id = data.azurerm_client_config.current.subscription_id
  blueprint_name  = data.azurerm_blueprint_definition.test.name
  version         = "v0.1_testAcc"
}

resource "azurerm_resource_group" "test" {
  name     = "accTestRG-bp-%d"
  location = "%s"

  tags = {
    testAcc = "true"
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "bp-user-%d"
}

resource "azurerm_role_assignment" "test" {
  scope                = data.azurerm_subscription.test.id
  role_definition_name = "Blueprint Operator"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "test2" {
  scope                = data.azurerm_subscription.test.id
  role_definition_name = "Owner"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_blueprint_assignment" "test" {
  name       = "testAccBPAssignment"
  scope_type = "subscriptions"
  scope      = data.azurerm_client_config.current.subscription_id
  version_id = data.azurerm_blueprint_published_version.test.id
  location   = "%s"

  lock_mode = "AllResourcesDoNotDelete"

  lock_exclude_principals = [
    data.azurerm_client_config.current.object_id,
  ]

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  resource_groups = <<GROUPS
    {
      "ResourceGroup": {
        "name": "accTestRG-BP-%d"
      }
    }
  GROUPS

  parameter_values = <<VALUES
    {
      "allowedlocationsforresourcegroups_listOfAllowedLocations": {
        "value": ["westus", "westus2", "eastus", "centralus", "centraluseuap", "southcentralus", "northcentralus", "westcentralus", "eastus2", "eastus2euap", "brazilsouth", "brazilus", "northeurope", "westeurope", "eastasia", "southeastasia", "japanwest", "japaneast", "koreacentral", "koreasouth", "indiasouth", "indiawest", "indiacentral", "australiaeast", "australiasoutheast", "canadacentral", "canadaeast", "uknorth", "uksouth2", "uksouth", "ukwest", "francecentral", "francesouth", "australiacentral", "australiacentral2", "uaecentral", "uaenorth", "southafricanorth", "southafricawest", "switzerlandnorth", "switzerlandwest", "germanynorth", "germanywestcentral", "norwayeast", "norwaywest"]
      }
    }
  VALUES

  depends_on = [
    azurerm_role_assignment.test,
    azurerm_role_assignment.test2
  ]
}
`, bpName, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccBlueprintAssignment_rootManagementGroup(data acceptance.TestData, bpName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_subscription" "test" {}

data "azurerm_blueprint_definition" "test" {
  name       = "%s"
  scope_type = "managementGroup"
  scope_name = data.azurerm_client_config.current.tenant_id
}

data "azurerm_blueprint_published_version" "test" {
  management_group = data.azurerm_client_config.current.tenant_id
  blueprint_name   = data.azurerm_blueprint_definition.test.name
  version          = "v0.1_testAcc"
}

resource "azurerm_resource_group" "test" {
  name     = "accTestRG-bp-%d"
  location = "%s"

  tags = {
    testAcc = "true"
  }
}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "bp-user-%d"
}

resource "azurerm_role_assignment" "test" {
  scope                = data.azurerm_subscription.test.id
  role_definition_name = "Blueprint Operator"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "test2" {
  scope                = data.azurerm_subscription.test.id
  role_definition_name = "Owner"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_blueprint_assignment" "test" {
  name       = "testAccBPAssignment"
  scope_type = "subscriptions"
  scope      = data.azurerm_client_config.current.subscription_id
  version_id = data.azurerm_blueprint_published_version.test.id
  location   = "%s"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  depends_on = [
    azurerm_role_assignment.test,
    azurerm_role_assignment.test2
  ]
}
`, bpName, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}
