// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagementGroupPolicySetDefinitionResourceTest struct{}

func TestAccManagementGroupPolicySetDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_set_definition", "test")
	r := ManagementGroupPolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementGroupPolicySetDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_set_definition", "test")
	r := ManagementGroupPolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccManagementGroupPolicySetDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_set_definition", "test")
	r := ManagementGroupPolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementGroupPolicySetDefinition_policyDefinitionVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_set_definition", "test")
	r := ManagementGroupPolicySetDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.policyDefinitionVersion(data, "9.*"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.policyDefinitionVersion(data, "9.3.*"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.policyDefinitionVersion(data, "9.*.*"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementGroupPolicySetDefinition_removeParameter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_set_definition", "test")
	r := ManagementGroupPolicySetDefinitionResourceTest{}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.additionalParameter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionReplace),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ManagementGroupPolicySetDefinitionResourceTest) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := policysetdefinitions.ParseProviders2PolicySetDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Policy.PolicySetDefinitionsClient.GetAtManagementGroup(ctx, *id, policysetdefinitions.DefaultGetAtManagementGroupOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagementGroupPolicySetDefinitionResourceTest) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_management_group_policy_set_definition" "test" {
  name                = "acctestpolset-%[2]d"
  policy_type         = "Custom"
  display_name        = "acctestpolset-%[2]d"
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

func (r ManagementGroupPolicySetDefinitionResourceTest) additionalParameter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_management_group_policy_set_definition" "test" {
  name                = "acctestpolset-%[2]d"
  policy_type         = "Custom"
  display_name        = "acctestpolset-%[2]d"
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
       },
       "allowedResourceTypes": {
           "type": "Array",
           "defaultValue": [
                "Microsoft.Compute/virtualMachines"
            ],
           "metadata": {
               "description": "The list of allowed resource types.",
               "displayName": "Allowed resource types",
               "strongType": "resourceType"
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
    policy_definition_id = "/providers/Microsoft.Authorization/policyDefinitions/a08ec900-254a-4555-9bf5-e42af04b5c5c"
    parameter_values     = <<VALUES
  {
      "listOfResourceTypesAllowed": {"value": "[parameters('allowedResourceTypes')]"}
  }
VALUES
  }
}`, r.template(data), data.RandomInteger)
}

func (r ManagementGroupPolicySetDefinitionResourceTest) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_management_group_policy_set_definition" "import" {
  name                = azurerm_management_group_policy_set_definition.test.name
  policy_type         = azurerm_management_group_policy_set_definition.test.policy_type
  display_name        = azurerm_management_group_policy_set_definition.test.display_name
  management_group_id = azurerm_management_group.test.id

  parameters = azurerm_management_group_policy_set_definition.test.parameters

  policy_definition_reference {
    policy_definition_id = azurerm_management_group_policy_set_definition.test.policy_definition_reference.0.policy_definition_id
    parameter_values     = azurerm_management_group_policy_set_definition.test.policy_definition_reference.0.parameter_values
  }
}`, r.basic(data))
}

func (r ManagementGroupPolicySetDefinitionResourceTest) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_management_group_policy_set_definition" "test" {
  name                = "acctestpolset-%[2]d"
  display_name        = "acctestpolset-%[2]d"
  management_group_id = azurerm_management_group.test.id
  policy_type         = "Custom"

  description = "A description for this policy set definition"
  metadata    = <<METADATA
    {
        "foo": "bar"
    }
METADATA

  policy_definition_group {
    name         = "Group-1"
    category     = "My Access Control"
    description  = "Controls accesses"
    display_name = "Group-Display-1"
  }

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
    policy_group_names   = ["Group-1"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ManagementGroupPolicySetDefinitionResourceTest) policyDefinitionVersion(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_management_group_policy_set_definition" "test" {
  name                = "acctestpolset-%[2]d"
  display_name        = "acctestpolset-%[2]d"
  management_group_id = azurerm_management_group.test.id
  policy_type         = "Custom"

  policy_definition_reference {
    version              = %q
    policy_definition_id = "/providers/Microsoft.Authorization/policyDefinitions/e345eecc-fa47-480f-9e88-67dcc122b164"
    parameter_values     = <<VALUES
    {
       "cpuLimit": {"value": "100m"},
       "memoryLimit": {"value": "100Mi"}
    }
VALUES
  }
}
`, r.template(data), data.RandomInteger, version)
}

func (r ManagementGroupPolicySetDefinitionResourceTest) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%d"
}

resource "azurerm_policy_definition" "test" {
  name                = "acctestpol-%[1]d"
  display_name        = "acctestpol-%[1]d"
  management_group_id = azurerm_management_group.test.id
  mode                = "All"
  policy_type         = "Custom"

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
