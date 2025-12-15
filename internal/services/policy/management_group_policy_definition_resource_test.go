// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policydefinitions"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagementGroupPolicyDefinitionResourceTest struct{}

func TestAccManagementGroupPolicyDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_definition", "test")
	r := ManagementGroupPolicyDefinitionResourceTest{}

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

func TestAccManagementGroupPolicyDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_definition", "test")
	r := ManagementGroupPolicyDefinitionResourceTest{}

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

func TestAccManagementGroupPolicyDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_definition", "test")
	r := ManagementGroupPolicyDefinitionResourceTest{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementGroupPolicyDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_definition", "test")
	r := ManagementGroupPolicyDefinitionResourceTest{}

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

func TestAccManagementGroupPolicyDefinition_removeParameter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_definition", "test")
	r := ManagementGroupPolicyDefinitionResourceTest{}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withAdditionalParameter(data),
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

func (ManagementGroupPolicyDefinitionResourceTest) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := policydefinitions.ParseProviders2PolicyDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Policy.PolicyDefinitionsClient.GetAtManagementGroup(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagementGroupPolicyDefinitionResourceTest) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  name = "acctestmg-%[1]d"
}

resource "azurerm_management_group_policy_definition" "test" {
  name                = "acctestpol-%[1]d"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "acctestpol-%[1]d"
  management_group_id = azurerm_management_group.test.id
}
`, data.RandomInteger)
}

func (r ManagementGroupPolicyDefinitionResourceTest) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_management_group_policy_definition" "import" {
  name                = azurerm_management_group_policy_definition.test.name
  policy_type         = azurerm_management_group_policy_definition.test.policy_type
  mode                = azurerm_management_group_policy_definition.test.mode
  display_name        = azurerm_management_group_policy_definition.test.display_name
  management_group_id = azurerm_management_group_policy_definition.test.management_group_id
}
`, template)
}

func (r ManagementGroupPolicyDefinitionResourceTest) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  name = "acctestmg-%[1]d"
}

resource "azurerm_management_group_policy_definition" "test" {
  name                = "acctestpol-%[1]d"
  policy_type         = "Custom"
  mode                = "Indexed"
  display_name        = "acctestpol-%[1]d"
  management_group_id = azurerm_management_group.test.id
  description         = "Test Policy Definition %[1]d"

  metadata = <<METADATA
    {
      "category": "General"
    }
METADATA

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

  policy_rule = <<POLICY_RULE
    {
      "if": {
        "not": {
          "field": "location",
          "in": "[parameters('allowedLocations')]"
        }
      },
      "then": {
        "effect": "Deny"
      }
    }
POLICY_RULE
}
`, data.RandomInteger)
}

func (r ManagementGroupPolicyDefinitionResourceTest) withAdditionalParameter(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  name = "acctestmg-%[1]d"
}

resource "azurerm_management_group_policy_definition" "test" {
  name                = "acctestpol-%[1]d"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "acctestpol-%[1]d"
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
      "additionalParameter": {
        "type": "String",
        "metadata": {
          "description": "An additional parameter",
          "displayName": "Additional Parameter"
        }
      }
    }
PARAMETERS

  policy_rule = <<POLICY_RULE
    {
      "if": {
        "not": {
          "field": "location",
          "in": "[parameters('allowedLocations')]"
        }
      },
      "then": {
        "effect": "Deny"
      }
    }
POLICY_RULE
}
`, data.RandomInteger)
}
