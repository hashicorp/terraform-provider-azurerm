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

type ManagementGroupAssignmentTestResource struct{}

func TestAccManagementGroupPolicyAssignment_basicWithBuiltInPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_assignment", "test")
	r := ManagementGroupAssignmentTestResource{}

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

func TestAccManagementGroupPolicyAssignment_basicWithBuiltInPolicySet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_assignment", "test")
	r := ManagementGroupAssignmentTestResource{}

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

func TestAccManagementGroupPolicyAssignment_basicWithCustomPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_assignment", "test")
	r := ManagementGroupAssignmentTestResource{}

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

func TestAccManagementGroupPolicyAssignment_basicWithCustomPolicyRequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_assignment", "test")
	r := ManagementGroupAssignmentTestResource{}

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

func TestAccManagementGroupPolicyAssignment_basicWithCustomPolicyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_assignment", "test")
	r := ManagementGroupAssignmentTestResource{}

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

func TestAccManagementGroupPolicyAssignment_basicComplexWithCustomPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_assignment", "test")
	r := ManagementGroupAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withComplexCustomPolicyAndParameters(data, "Acceptance"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withComplexCustomPolicyAndParameters(data, "Testing"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withComplexCustomPolicyAndParameters(data, "Acceptance Testing"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagementGroupAssignmentTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r ManagementGroupAssignmentTestResource) withBuiltInPolicyBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = ["%[3]s"]
    }
  })
}
`, template, data.RandomString, data.Locations.Primary)
}

func (r ManagementGroupAssignmentTestResource) withBuiltInPolicyUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = ["%[3]s", "%[4]s"]
    }
  })
}
`, template, data.RandomString, data.Locations.Primary, data.Locations.Secondary)
}

func (r ManagementGroupAssignmentTestResource) withBuiltInPolicySetBasic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = %[3]q

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomString, data.Locations.Primary)
}

func (r ManagementGroupAssignmentTestResource) withBuiltInPolicySetUpdated(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = %[3]q

  identity {
    type = "SystemAssigned"
  }

  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomString, data.Locations.Primary)
}

func (r ManagementGroupAssignmentTestResource) withCustomPolicyBasic(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id
}
`, template, data.RandomString)
}

func (r ManagementGroupAssignmentTestResource) withCustomPolicyComplete(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	// NOTE: we could include parameters here but it's tested extensively elsewhere
	// NOTE: intentionally avoiding NotScopes for Management Groups since it's awkward to test, but
	// this is covered in the other resources
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "This is a policy assignment from an acceptance test"
  display_name         = "AccTest Policy %[2]s"
  enforce              = false
  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomString)
}

func (r ManagementGroupAssignmentTestResource) withCustomPolicyRequiresImport(data acceptance.TestData) string {
	template := r.withCustomPolicyBasic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_management_group_policy_assignment" "import" {
  name                 = azurerm_management_group_policy_assignment.test.name
  management_group_id  = azurerm_management_group_policy_assignment.test.management_group_id
  policy_definition_id = azurerm_management_group_policy_assignment.test.policy_definition_id
}
`, template)
}

func (r ManagementGroupAssignmentTestResource) withComplexCustomPolicyAndParameters(data acceptance.TestData, metadataValue string) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_policy_definition" "test" {
  name                = "acctestpol-%[2]s"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "acctestpol-%[2]s"
  description         = "Description for %[2]s"
  management_group_id = azurerm_management_group.test.group_id
  metadata            = <<METADATA
  {
    "category": "%[3]s"
  }
METADATA

  parameters = <<PARAMETERS
  {
    "effect": {
      "type": "String",
      "metadata": {
        "displayName": "Effect",
        "description": "Enable or disable the execution of the policy"
      },
      "allowedValues": [
        "Deny",
        "Audit",
        "Disabled"
      ],
      "defaultValue": "Deny"
    },
    "sourceIp": {
      "type": "Array",
      "metadata": {
        "displayName": "Source IP ranges",
        "description": "The inbound IP range to deny. Default to *, ANY, Internet"
      },
      "defaultValue": [
        "*",
        "Any",
        "Internet",
        "0.0.0.0"
      ]
    }
  }
PARAMETERS

  policy_rule = <<POLICY_RULE
  {
    "if": {
      "allOf": [
        {
          "field": "type",
          "equals": "Microsoft.Network/networkSecurityGroups/securityRules"
        },
        {
          "allOf": [
            {
              "field": "Microsoft.Network/networkSecurityGroups/securityRules/access",
              "equals": "Allow"
            },
            {
              "field": "Microsoft.Network/networkSecurityGroups/securityRules/direction",
              "equals": "Inbound"
            },
            {
              "anyOf": [
                {
                  "field": "Microsoft.Network/networkSecurityGroups/securityRules/destinationPortRange",
                  "equals": "*"
                },
                {
                  "not": {
                    "field": "Microsoft.Network/networkSecurityGroups/securityRules/destinationPortRanges[*]",
                    "notEquals": "*"
                  }
                },
                {
                  "field": "Microsoft.Network/networkSecurityGroups/securityRules/sourceAddressPrefix",
                  "in": "[parameters('sourceIP')]"
                },
                {
                  "not": {
                    "field": "Microsoft.Network/networkSecurityGroups/securityRules/sourceAddressPrefixes[*]",
                    "notIn": "[parameters('sourceIP')]"
                  }
                }
              ]
            }
          ]
        }
      ]
    },
    "then": {
      "effect": "[parameters('effect')]"
    }
  }
POLICY_RULE
}

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id
}
`, template, data.RandomString, metadataValue)
}

func (r ManagementGroupAssignmentTestResource) withCustomPolicyUpdated(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomString)
}

func (r ManagementGroupAssignmentTestResource) templateWithCustomPolicy(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_policy_definition" "test" {
  name                = "acctestpol-%[2]s"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "acctestpol-%[2]s"
  management_group_id = azurerm_management_group.test.group_id

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
`, template, data.RandomString)
}

func (r ManagementGroupAssignmentTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  display_name = "Acceptance Test MgmtGroup %[1]d"
}
`, data.RandomInteger)
}
