// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PolicySetDefinitionResourceTest struct{}

func TestAccAzureRMPolicySetDefinition_builtIn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.builtIn(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicySetDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.builtIn(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMPolicySetDefinition_custom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.custom(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicySetDefinition_customNoParameter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customNoParameter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customNoParameterUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicySetDefinition_customUpdateDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.custom(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customUpdateDisplayName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicySetDefinition_customUpdateParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.custom(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customUpdateParameters(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccAzureRMPolicySetDefinition_customUpdateAddNewReference(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.custom(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customUpdateAddNewReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicySetDefinition_customWithPolicyReferenceID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customWithPolicyReferenceID(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicySetDefinition_customWithDefinitionGroups(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customWithDefinitionGroups(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customWithDefinitionGroupsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.customWithDefinitionGroups(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicySetDefinition_customWithGroupsInDefinitionReferenceUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// provision a policy set without group names
			Config: r.customWithDefinitionGroupsNotUsedInPolicyReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("policy_definition_reference.0.policy_group_names").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			// test if group_names were correctly added
			Config: r.customWithDefinitionGroupsUsedInPolicyReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("policy_definition_reference.0.policy_group_names.#").HasValue("3"),
			),
		},
		data.ImportStep(),
		{
			// test if the deletion of the group_names works again
			Config: r.customWithDefinitionGroupsNotUsedInPolicyReference(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("policy_definition_reference.0.policy_group_names.0").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicySetDefinition_managementGroup(t *testing.T) {
	if features.FivePointOh() {
		t.Skip("`skipping test as `management_group_id` has been removed from the `azurerm_policy_set_definition` resource")
	}

	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managementGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicySetDefinition_metadata(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.metadata(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PolicySetDefinitionResourceTest) builtIn(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestpolset-%d"
  policy_type  = "Custom"
  display_name = "acctestpolset-%d"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = "/providers/Microsoft.Authorization/policyDefinitions/e765b5de-1225-4ba3-bd56-1ac6695af988"
    parameter_values     = <<VALUES
    {
      "listOfAllowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
  }
}
`, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "import" {
  name         = azurerm_policy_set_definition.test.name
  policy_type  = azurerm_policy_set_definition.test.policy_type
  display_name = azurerm_policy_set_definition.test.display_name
  parameters   = azurerm_policy_set_definition.test.parameters

  policy_definition_reference {
    policy_definition_id = "/providers/Microsoft.Authorization/policyDefinitions/e765b5de-1225-4ba3-bd56-1ac6695af988"
    parameter_values     = <<VALUES
    {
      "listOfAllowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
  }
}
`, r.builtIn(data))
}

func (r PolicySetDefinitionResourceTest) custom(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = <<VALUES
    {
      "allowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) customUpdateDisplayName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d-updated"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = <<VALUES
    {
      "allowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) customUpdateParameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d-updated"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = <<VALUES
    {
      "allowedLocations": {"value": ["%s"]}
    }
VALUES
  }
}
`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r PolicySetDefinitionResourceTest) customUpdateAddNewReference(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_policy_definition" "allowed_resource_types" {
  display_name = "Allowed resource types"
}

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = <<VALUES
    {
      "allowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
  }

  policy_definition_reference {
    policy_definition_id = data.azurerm_policy_definition.allowed_resource_types.id
    parameter_values     = <<VALUES
    {
      "listOfResourceTypesAllowed": {"value": ["Microsoft.Compute/virtualMachines"]}
    }
VALUES
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) customNoParameter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d"

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
  }
}
`, r.templateNoParameter(data), data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) customNoParameterUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d"

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = "{}"
  }
}
`, r.templateNoParameter(data), data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) managementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%[1]d"
}

resource "azurerm_policy_set_definition" "test" {
  name                = "acctestpolset-%[1]d"
  policy_type         = "Custom"
  display_name        = "acctestpolset-%[1]d"
  management_group_id = azurerm_management_group.test.id

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = "/providers/Microsoft.Authorization/policyDefinitions/e765b5de-1225-4ba3-bd56-1ac6695af988"
    parameter_values     = <<VALUES
    {
      "listOfAllowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
  }
}
`, data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) metadata(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestpolset-%[1]d"
  policy_type  = "Custom"
  display_name = "acctestpolset-%[1]d"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = "/providers/Microsoft.Authorization/policyDefinitions/e765b5de-1225-4ba3-bd56-1ac6695af988"
    parameter_values     = <<VALUES
    {
      "listOfAllowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
  }

  metadata = <<METADATA
    {
        "foo": "bar"
    }
METADATA
}
`, data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) customWithPolicyReferenceID(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = <<VALUES
    {
      "allowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
    reference_id         = "TestRef"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) customWithDefinitionGroups(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = <<VALUES
    {
      "allowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
    policy_group_names   = ["group-1", "group-2"]
  }

  policy_definition_group {
    name = "redundant"
  }

  policy_definition_group {
    name = "Group-1"
  }

  policy_definition_group {
    name = "group-2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) customWithDefinitionGroupsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = <<VALUES
    {
      "allowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
    policy_group_names   = ["group-1", "group-2"]
  }

  policy_definition_group {
    name = "redundant"
  }

  policy_definition_group {
    name         = "Group-1"
    display_name = "Group-Display-1"
    category     = "My Access Control"
    description  = "Controls accesses"
  }

  policy_definition_group {
    name         = "group-2"
    display_name = "group-display-2"
    category     = "My Security Control"
    description  = "Controls security"
  }
}
`, r.template(data), data.RandomInteger)
}

// test adding "group-3" to policy_definition_reference.policy_group_names
func (r PolicySetDefinitionResourceTest) customWithDefinitionGroupsUsedInPolicyReference(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = <<VALUES
    {
      "allowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
VALUES
    policy_group_names   = ["group-1", "group-2", "group-3"]
  }

  policy_definition_group {
    name = "redundant"
  }

  policy_definition_group {
    name         = "group-1"
    display_name = "Group-Display-1"
    category     = "My Access Control"
    description  = "Controls accesses"
  }

  policy_definition_group {
    name         = "group-2"
    display_name = "group-display-2"
    category     = "My Security Control"
    description  = "Controls security"
  }

  policy_definition_group {
    name         = "group-3"
    display_name = "group-display-3"
    category     = "Category-3"
    description  = "Newly added group 3"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

// test adding "group-3" to policy_definition_reference.policy_group_names
func (r PolicySetDefinitionResourceTest) customWithDefinitionGroupsNotUsedInPolicyReference(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%[2]d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%[2]d"

  parameters = <<PARAMETERS
    {
        "allowedLocations": {
            "type": "Array",
            "metadata": {
                "description": "The list of allowed locations for resources.",
                "displayName": "Allowed locations",
                "strongType": "location"
            }
        }
    }
PARAMETERS

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = <<VALUES
    {
        "allowedLocations": {"value": "[parameters('allowedLocations')]"}
    }
    VALUES
  }

  policy_definition_group {
    name = "redundant"
  }

  policy_definition_group {
    name         = "group-1"
    display_name = "Group-Display-1"
    category     = "My Access Control"
    description  = "Controls accesses"
  }

  policy_definition_group {
    name         = "group-2"
    display_name = "group-display-2"
    category     = "My Security Control"
    description  = "Controls security"
  }

  policy_definition_group {
    name         = "group-3"
    display_name = "group-display-3"
    category     = "Category-3"
    description  = "Newly added group 3"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestpol-%[1]d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%[1]d"

  policy_rule = <<POLICY_RULE
    {
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

  parameters = <<PARAMETERS
  {
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
PARAMETERS
}
`, data.RandomInteger)
}

func (r PolicySetDefinitionResourceTest) templateNoParameter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestpol-%[1]d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%[1]d"

  policy_rule = <<POLICY_RULE
    {
    "if": {
      "not": {
        "field": "location",
        "equals": "%s"
      }
    },
    "then": {
      "effect": "deny"
    }
  }
POLICY_RULE
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PolicySetDefinitionResourceTest) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	if !features.FivePointOh() {
		subscriptionId := client.Account.SubscriptionId

		resourceId, err := parse.PolicySetDefinitionID(state.ID)
		if err != nil {
			return nil, err
		}

		if scopeId, ok := resourceId.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
			id := policysetdefinitions.NewProviders2PolicySetDefinitionID(scopeId.ManagementGroupName, resourceId.Name)
			resp, err := client.Policy.PolicySetDefinitionsClient.GetAtManagementGroup(ctx, id, policysetdefinitions.DefaultGetAtManagementGroupOperationOptions())
			if err != nil {
				return nil, fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return pointer.To(resp.Model != nil), nil
		}

		id := policysetdefinitions.NewProviderPolicySetDefinitionID(subscriptionId, resourceId.Name)
		resp, err := client.Policy.PolicySetDefinitionsClient.Get(ctx, id, policysetdefinitions.DefaultGetOperationOptions())
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return pointer.To(resp.Model != nil), nil
	}

	id, err := policysetdefinitions.ParseProviderPolicySetDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Policy.PolicySetDefinitionsClient.Get(ctx, *id, policysetdefinitions.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}
