package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ResourceAssignmentTestResource struct{}

func TestAccResourcePolicyAssignment_basicWithBuiltInPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBuiltInPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicyUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourcePolicyAssignment_basicWithBuiltInPolicySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBuiltInPolicySetBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicySetUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicySetBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourcePolicyAssignment_basicWithCustomPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCustomPolicyUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCustomPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourcePolicyAssignment_basicWithCustomPolicyRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		data.RequiresImportErrorStep(r.withCustomPolicyRequiresImport),
	})
}

func TestAccResourcePolicyAssignment_basicWithCustomPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_assignment", "test")
	r := ResourceAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomPolicyComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCustomPolicyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCustomPolicyComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ResourceAssignmentTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PolicyAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	assignment, err := client.Policy.AssignmentsClient.Get(ctx, id.Scope, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(assignment.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r ResourceAssignmentTestResource) withBuiltInPolicyBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = [azurerm_resource_group.test.location]
    }
  })
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withBuiltInPolicyUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = [azurerm_resource_group.test.location, "%[3]s"]
    }
  })
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func (r ResourceAssignmentTestResource) withBuiltInPolicySetBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withBuiltInPolicySetUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }

  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withCustomPolicyBasic(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = azurerm_policy_definition.test.id
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withCustomPolicyComplete(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	// NOTE: we could include parameters here but it's tested extensively elsewhere
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "This is a policy assignment from an acceptance test"
  display_name         = "AccTest Policy %[2]d"
  enforce              = false
  not_scopes = [
    format("%%s/subnets/subnet1", azurerm_virtual_network.test.id)
  ]
  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) withCustomPolicyRequiresImport(data acceptance.TestData) string {
	template := r.withCustomPolicyBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_policy_assignment" "import" {
  name                 = azurerm_resource_policy_assignment.test.name
  resource_id          = azurerm_resource_policy_assignment.test.resource_id
  policy_definition_id = azurerm_resource_policy_assignment.test.policy_definition_id
}
`, template)
}

func (r ResourceAssignmentTestResource) withCustomPolicyUpdated(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) templateWithCustomPolicy(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_policy_definition" "test" {
  name         = "acctestpol-%[2]d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%[2]d"

  policy_rule = <<POLICY_RULE
	{
    "if": {
      "not": {
        "field": "name",
        "equals": "bob"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE
}
`, template, data.RandomInteger)
}

func (r ResourceAssignmentTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest%[1]d"
  location = %[2]q
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}
`, data.RandomInteger, data.Locations.Primary)
}
