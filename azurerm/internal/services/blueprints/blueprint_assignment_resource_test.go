package blueprints_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/blueprints/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type BlueprintAssignmentResource struct {
}

// Scenario: Basic BP, no artefacts etc.  Stored and applied at Subscription.
func TestAccBlueprintAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")
	r := BlueprintAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "testAcc_basicSubscription", "v0.1_testAcc"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBlueprintAssignment_basicUpdated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")
	r := BlueprintAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "testAcc_basicSubscription", "v0.1_testAcc"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "testAcc_basicSubscription", "v0.2_testAcc"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccBlueprintAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")
	r := BlueprintAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "testAcc_basicSubscription", "v0.1_testAcc"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, "testAcc_basicSubscription", "v0.1_testAcc"),
			ExpectError: acceptance.RequiresImportError("azurerm_blueprint_assignment"),
		},
	})
}

// Scenario: BP with RG's, locking and parameters/policies stored at Subscription, applied to subscription
func TestAccBlueprintAssignment_subscriptionComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")
	r := BlueprintAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subscriptionComplete(data, "testAcc_subscriptionComplete", "v0.1_testAcc"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

// Scenario: BP stored at Root Management Group, applied to Subscription
func TestAccBlueprintAssignment_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_blueprint_assignment", "test")
	r := BlueprintAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.rootManagementGroup(data, "testAcc_basicRootManagementGroup", "v0.1_testAcc"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t BlueprintAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
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

func (BlueprintAssignmentResource) basic(data acceptance.TestData, bpName string, version string) string {
	subscription := data.Client().SubscriptionIDAlt
	return fmt.Sprintf(`
provider "azurerm" {
  subscription_id = "%s"
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_subscription" "test" {}

data "azurerm_blueprint_definition" "test" {
  name     = "%s"
  scope_id = data.azurerm_subscription.test.id
}

data "azurerm_blueprint_published_version" "test" {
  scope_id       = data.azurerm_blueprint_definition.test.scope_id
  blueprint_name = data.azurerm_blueprint_definition.test.name
  version        = "%s"
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
  name                   = "testAccBPAssignment"
  target_subscription_id = data.azurerm_subscription.test.id
  version_id             = data.azurerm_blueprint_published_version.test.id
  location               = "%s"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, subscription, bpName, version, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

// This test config creates a UM-MSI and assigns Owner to the target subscription.  This is necessary due to the changes
// the referenced Blueprint Version needs to make to successfully apply.  If the test does not exit cleanly, "dangling"
// resources can include the Role Assignment(s) at the Subscription, which will need to be removed
func (BlueprintAssignmentResource) subscriptionComplete(data acceptance.TestData, bpName string, version string) string {
	subscription := data.Client().SubscriptionIDAlt

	return fmt.Sprintf(`
provider "azurerm" {
  subscription_id = "%s"
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_subscription" "test" {}

data "azurerm_blueprint_definition" "test" {
  name     = "%s"
  scope_id = data.azurerm_subscription.test.id
}

data "azurerm_blueprint_published_version" "test" {
  scope_id       = data.azurerm_blueprint_definition.test.scope_id
  blueprint_name = data.azurerm_blueprint_definition.test.name
  version        = "%s"
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

resource "azurerm_role_assignment" "operator" {
  scope                = data.azurerm_subscription.test.id
  role_definition_name = "Blueprint Operator"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "owner" {
  scope                = data.azurerm_subscription.test.id
  role_definition_name = "Owner"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_blueprint_assignment" "test" {
  name                   = "testAccBPAssignment"
  target_subscription_id = data.azurerm_subscription.test.id
  version_id             = data.azurerm_blueprint_published_version.test.id
  location               = "%s"

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
    azurerm_role_assignment.operator,
    azurerm_role_assignment.owner
  ]
}
`, subscription, bpName, version, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (BlueprintAssignmentResource) rootManagementGroup(data acceptance.TestData, bpName string, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_subscription" "test" {}

data "azurerm_management_group" "root" {
  group_id = data.azurerm_client_config.current.tenant_id
}

data "azurerm_blueprint_definition" "test" {
  name     = "%s"
  scope_id = data.azurerm_management_group.root.id
}

data "azurerm_blueprint_published_version" "test" {
  scope_id       = data.azurerm_blueprint_definition.test.scope_id
  blueprint_name = data.azurerm_blueprint_definition.test.name
  version        = "%s"
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

resource "azurerm_role_assignment" "operator" {
  scope                = data.azurerm_subscription.test.id
  role_definition_name = "Blueprint Operator"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_role_assignment" "owner" {
  scope                = data.azurerm_subscription.test.id
  role_definition_name = "Owner"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_blueprint_assignment" "test" {
  name                   = "testAccBPAssignment"
  target_subscription_id = data.azurerm_subscription.test.id
  version_id             = data.azurerm_blueprint_published_version.test.id
  location               = "%s"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  depends_on = [
    azurerm_role_assignment.operator,
    azurerm_role_assignment.owner
  ]
}
`, bpName, version, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (BlueprintAssignmentResource) requiresImport(data acceptance.TestData, bpName string, version string) string {
	template := BlueprintAssignmentResource{}.basic(data, bpName, version)

	return fmt.Sprintf(`
%s

resource "azurerm_blueprint_assignment" "import" {
  name                   = azurerm_blueprint_assignment.test.name
  target_subscription_id = azurerm_blueprint_assignment.test.target_subscription_id
  version_id             = azurerm_blueprint_assignment.test.version_id
  location               = azurerm_blueprint_assignment.test.location

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  depends_on = [
    azurerm_role_assignment.test
  ]
}


`, template)
}
