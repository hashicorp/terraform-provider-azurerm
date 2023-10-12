// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	assignments "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-06-01/policyassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

func TestAccManagementGroupPolicyAssignment_basicWithBuiltInPolicyNonComplianceMessage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_assignment", "test")
	r := ManagementGroupAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBuiltInPolicyNonComplianceMessage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("1"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicyNonComplianceMessageUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.0").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicyNonComplianceMessage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("1"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
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

func TestAccManagementGroupPolicyAssignment_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_assignment", "test")
	r := ManagementGroupAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.userAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementGroupPolicyAssignment_basicWithBuiltInPolicySetNonComplianceMessage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_assignment", "test")
	r := ManagementGroupAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBuiltInPolicySetNonComplianceMessage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("1"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicySetNonComplianceMessageUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("2"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
				check.That(data.ResourceName).Key("non_compliance_message.1.content").HasValue("test2"),
				check.That(data.ResourceName).Key("non_compliance_message.1.policy_definition_reference_id").HasValue("AINE_MinimumPasswordLength"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBuiltInPolicySetNonComplianceMessage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("non_compliance_message.#").HasValue("1"),
				check.That(data.ResourceName).Key("non_compliance_message.0.content").HasValue("test"),
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

func TestAccManagementGroupPolicyAssignment_overridesAndResourceSelector(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_assignment", "test")
	r := ManagementGroupAssignmentTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withOverrideAndSelectorsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withOverrideAndSelectorsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withOverrideAndSelectorsBasic(data),
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
	id, err := assignments.ParseScopedPolicyAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	assignment, err := client.Policy.AssignmentsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(assignment.HttpResponse) {
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
  name                 = "acctestpol-mg-%[2]s"
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
  name                 = "acctestpol-mg-%[2]s"
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

func (r ManagementGroupAssignmentTestResource) withBuiltInPolicyNonComplianceMessage(data acceptance.TestData) string {
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
  name                 = "acctestpol-mg-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id

  non_compliance_message {
    content = "test"
  }

  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = ["%[3]s"]
    }
  })
}
`, template, data.RandomString, data.Locations.Primary)
}

func (r ManagementGroupAssignmentTestResource) withBuiltInPolicyNonComplianceMessageUpdated(data acceptance.TestData) string {
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
  name                 = "acctestpol-mg-%[2]s"
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
  name                 = "acctestpol-mg-%[2]s"
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
  name                 = "acctestpol-mg-%[2]s"
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

func (r ManagementGroupAssignmentTestResource) withBuiltInPolicySetNonComplianceMessage(data acceptance.TestData) string {
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
  name                 = "acctestpol-mg-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = %[3]q

  non_compliance_message {
    content = "test"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomString, data.Locations.Primary)
}

func (r ManagementGroupAssignmentTestResource) withBuiltInPolicySetNonComplianceMessageUpdated(data acceptance.TestData) string {
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
  name                 = "acctestpol-mg-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = %[3]q

  non_compliance_message {
    content = "test"
  }

  non_compliance_message {
    content                        = "test2"
    policy_definition_reference_id = "AINE_MinimumPasswordLength"
  }

  identity {
    type = "SystemAssigned"
  }
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
  name                 = "acctestpol-mg-%[2]s"
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
  name                 = "acctestpol-mg-%[2]s"
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
  name                = "acctestpol-mg-%[2]s"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "acctestpol-mg-%[2]s"
  description         = "Description for %[2]s"
  management_group_id = azurerm_management_group.test.id
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
  name                 = "acctestpol-mg-%[2]s"
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
  name                 = "acctestpol-mg-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  metadata = jsonencode({
    "category" : "Testing"
  })
}
`, template, data.RandomString)
}

func (r ManagementGroupAssignmentTestResource) withOverrideAndSelectorsBasic(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-mg-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  metadata = jsonencode({
    "category" : "Testing"
  })

  overrides {
    value = "Disabled"
    selectors {
      in = [data.azurerm_policy_set_definition.test.policy_definition_reference.0.reference_id]
    }
  }

  resource_selectors {
    selectors {
      not_in = ["eastus", "westus"]
      kind   = "resourceLocation"
    }
  }
}
`, template, data.RandomString)
}

func (r ManagementGroupAssignmentTestResource) withOverrideAndSelectorsUpdate(data acceptance.TestData) string {
	template := r.templateWithCustomPolicy(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-mg-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  metadata = jsonencode({
    "category" : "Testing"
  })

  overrides {
    value = "AuditIfNotExists"
    selectors {
      in = [data.azurerm_policy_set_definition.test.policy_definition_reference.0.reference_id]
    }
  }

  resource_selectors {
    name = "selected for policy"
    selectors {
      not_in = ["eastus"]
      kind   = "resourceLocation"
    }
  }
}
`, template, data.RandomString)
}

func (r ManagementGroupAssignmentTestResource) templateWithCustomPolicy(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_policy_definition" "test" {
  name                = "acctestpol-mg-%[2]s"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "acctestpol-mg-%[2]s"
  management_group_id = azurerm_management_group.test.id

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

func (r ManagementGroupAssignmentTestResource) systemAssignedIdentity(data acceptance.TestData) string {
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
  name                 = "acctestpol-mg-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = %[3]q

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomString, data.Locations.Primary)
}

func (r ManagementGroupAssignmentTestResource) userAssignedIdentity(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-pa-%[2]s"
  location = %[3]q
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestua%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpol-mg-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id
  location             = %[3]q

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomString, data.Locations.Primary)
}
