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

type SubscriptionAssignmentTestResource struct{}

func TestAccSubscriptionPolicyAssignment_basicWithBuiltInPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_policy_assignment", "test")
	r := SubscriptionAssignmentTestResource{}

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

func TestAccSubscriptionPolicyAssignment_basicWithBuiltInPolicySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_policy_assignment", "test")
	r := SubscriptionAssignmentTestResource{}

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

func TestAccSubscriptionPolicyAssignment_basicWithCustomPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_policy_assignment", "test")
	r := SubscriptionAssignmentTestResource{}

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

func TestAccSubscriptionPolicyAssignment_basicWithCustomPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_policy_assignment", "test")
	r := SubscriptionAssignmentTestResource{}

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

func TestAccSubscriptionPolicyAssignment_basicWithCustomPolicyRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_policy_assignment", "test")
	r := SubscriptionAssignmentTestResource{}

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

func (r SubscriptionAssignmentTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r SubscriptionAssignmentTestResource) withBuiltInPolicyBasic(data acceptance.TestData) string {
	template := r.template()
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_subscription_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = ["%s"]
    }
  })
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r SubscriptionAssignmentTestResource) withBuiltInPolicyUpdated(data acceptance.TestData) string {
	template := r.template()
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_subscription_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = ["%[3]s", "%[4]s"]
    }
  })
}
`, template, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r SubscriptionAssignmentTestResource) withBuiltInPolicySetBasic(data acceptance.TestData) string {
	template := r.template()
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_subscription_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = %[3]q

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r SubscriptionAssignmentTestResource) withBuiltInPolicySetUpdated(data acceptance.TestData) string {
	template := r.template()
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_subscription_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = %[3]q

  identity {
    type = "SystemAssigned"
  }

  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r SubscriptionAssignmentTestResource) withCustomPolicyBasic(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_subscription_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_definition_id = azurerm_policy_definition.test.id
}
`, template, data.RandomInteger)
}

func (r SubscriptionAssignmentTestResource) withCustomPolicyComplete(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	// NOTE: we could include parameters here but it's tested extensively elsewhere
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_subscription_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "This is a policy assignment from an acceptance test"
  display_name         = "AccTest Policy %[2]d"
  enforce              = false
  not_scopes = [
    format("%%s/resourceGroups/blah", data.azurerm_subscription.test.id)
  ]
  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomInteger)
}

func (r SubscriptionAssignmentTestResource) withCustomPolicyRequiresImport(data acceptance.TestData) string {
	template := r.withCustomPolicyBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subscription_policy_assignment" "import" {
  name                 = azurerm_subscription_policy_assignment.test.name
  subscription_id      = azurerm_subscription_policy_assignment.test.subscription_id
  policy_definition_id = azurerm_subscription_policy_assignment.test.policy_definition_id
}
`, template)
}
func (r SubscriptionAssignmentTestResource) withCustomPolicyUpdated(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_subscription_policy_assignment" "test" {
  name                 = "acctestpa-%[2]d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomInteger)
}

func (r SubscriptionAssignmentTestResource) templateWithCustomPolicy(data acceptance.TestData) string {
	template := r.template()
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

func (r SubscriptionAssignmentTestResource) template() string {
	return `data "azurerm_subscription" "test" {}`
}
