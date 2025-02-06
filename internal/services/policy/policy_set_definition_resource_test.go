// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PolicySetDefinitionResource struct{}

func TestAccAzureRMPolicySetDefinition_builtIn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResource{}

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
	r := PolicySetDefinitionResource{}

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
	r := PolicySetDefinitionResource{}

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
	r := PolicySetDefinitionResource{}

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
	r := PolicySetDefinitionResource{}

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
	r := PolicySetDefinitionResource{}

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
	r := PolicySetDefinitionResource{}

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
	r := PolicySetDefinitionResource{}

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
	r := PolicySetDefinitionResource{}

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
	r := PolicySetDefinitionResource{}

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
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	r := PolicySetDefinitionResource{}

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
	r := PolicySetDefinitionResource{}

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

func (r PolicySetDefinitionResource) builtIn(data acceptance.TestData) string {
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

func (r PolicySetDefinitionResource) requiresImport(data acceptance.TestData) string {
	template := r.builtIn(data)
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
`, template)
}

func (r PolicySetDefinitionResource) custom(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d"

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
`, template, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) customUpdateDisplayName(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d-updated"

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
`, template, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) customUpdateParameters(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d-updated"

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
`, template, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (r PolicySetDefinitionResource) customUpdateAddNewReference(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_definition" "allowed_resource_types" {
  display_name = "Allowed resource types"
}

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d"

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
`, template, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) customNoParameter(data acceptance.TestData) string {
	template := r.templateNoParameter(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d"

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) customNoParameterUpdate(data acceptance.TestData) string {
	template := r.templateNoParameter(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d"

  policy_definition_reference {
    policy_definition_id = azurerm_policy_definition.test.id
    parameter_values     = "{}"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) managementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%d"
}

resource "azurerm_policy_set_definition" "test" {
  name                = "acctestpolset-%d"
  policy_type         = "Custom"
  display_name        = "acctestpolset-%d"
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
`, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) metadata(data acceptance.TestData) string {
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

  metadata = <<METADATA
    {
        "foo": "bar"
    }
METADATA
}
`, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) customWithPolicyReferenceID(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d"

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
`, template, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) customWithDefinitionGroups(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d"

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
`, template, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) customWithDefinitionGroupsUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d"

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
`, template, data.RandomInteger, data.RandomInteger)
}

// test adding "group-3" to policy_definition_reference.policy_group_names
func (r PolicySetDefinitionResource) customWithDefinitionGroupsUsedInPolicyReference(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d"

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
`, template, data.RandomInteger, data.RandomInteger)
}

// test adding "group-3" to policy_definition_reference.policy_group_names
func (r PolicySetDefinitionResource) customWithDefinitionGroupsNotUsedInPolicyReference(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "test" {
  name         = "acctestPolSet-%d"
  policy_type  = "Custom"
  display_name = "acctestPolSet-display-%d"

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
`, template, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestpol-%d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%d"

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
`, data.RandomInteger, data.RandomInteger)
}

func (r PolicySetDefinitionResource) templateNoParameter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestpol-%d"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-%d"

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
`, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (r PolicySetDefinitionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.PolicySetDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	var resp policy.SetDefinition
	if mgmtGroupID, ok := id.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
		resp, err = client.Policy.SetDefinitionsClient.GetAtManagementGroup(ctx, id.Name, mgmtGroupID.ManagementGroupName)
	} else {
		resp, err = client.Policy.SetDefinitionsClient.Get(ctx, id.Name)
	}
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Policy Set Definition %q: %+v", id.Name, err)
	}

	return utils.Bool(true), nil
}
