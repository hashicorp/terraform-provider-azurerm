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

type PolicyDefinitionResource struct{}

func TestAccAzureRMPolicyDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

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

func TestAccAzureRMPolicyDefinition_basicWithDetail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithDetail(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("role_definition_ids.0").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

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

func TestAccAzureRMPolicyDefinition_computedMetadata(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.computedMetadata(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyDefinitionAtMgmtGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

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

func TestAccAzureRMPolicyDefinition_metadata(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

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

func TestAccAzureRMPolicyDefinition_modeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mode(data, "All"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.mode(data, "Indexed"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.mode(data, "All"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyDefinition_removeParameter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PolicyDefinitionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	definitionsClient := client.Policy.DefinitionsClient
	id, err := parse.PolicyDefinitionID(state.ID)
	if err != nil {
		return nil, err
	}

	var resp policy.Definition
	switch scope := id.PolicyScopeId.(type) {
	case parse.ScopeAtSubscription:
		resp, err = definitionsClient.Get(ctx, id.Name)
	case parse.ScopeAtManagementGroup:
		resp, err = definitionsClient.GetAtManagementGroup(ctx, id.Name, scope.ManagementGroupName)
	default:
		return nil, fmt.Errorf("unexpected scope type: %+v", scope)
	}
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Policy Definition %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.DefinitionProperties != nil), nil
}

func (r PolicyDefinitionResource) basic(data acceptance.TestData) string {
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

func (r PolicyDefinitionResource) basicWithDetail(data acceptance.TestData) string {
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
    "field": "type",
    "equals": "Microsoft.ServiceBus/namespaces"
  },
  "then": {
    "effect": "DeployIfNotExists",
    "details": {
      "type": "Microsoft.Insights/diagnosticSettings",
      "name": "acctest",
      "roleDefinitionIds": [
        "/providers/microsoft.authorization/roleDefinitions/749f88d5-cbae-40b8-bcfc-e573ddc772fa",
        "/providers/microsoft.authorization/roleDefinitions/92aaf0da-9dab-42b6-94a3-d43ce8d16293"
      ],
      "deployment": {
        "properties": {
          "mode": "incremental",
          "template": {
            "$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
            "contentVersion": "1.0.0.0",
            "parameters": {}
          }
        }
      }
    }
  }
}
POLICY_RULE
}
`, data.RandomInteger, data.RandomInteger)
}

func (r PolicyDefinitionResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_definition" "import" {
  name         = azurerm_policy_definition.test.name
  policy_type  = azurerm_policy_definition.test.policy_type
  mode         = azurerm_policy_definition.test.mode
  display_name = azurerm_policy_definition.test.display_name
  policy_rule  = azurerm_policy_definition.test.policy_rule
  parameters   = azurerm_policy_definition.test.parameters
}
`, template)
}

func (r PolicyDefinitionResource) computedMetadata(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_policy_definition" "test" {
  name         = "acctest-%d"
  policy_type  = "Custom"
  mode         = "Indexed"
  display_name = "DefaultTags"

  policy_rule = <<POLICY_RULE
    {
  "if": {
    "field": "tags",
    "exists": "false"
  },
  "then": {
    "effect": "append",
    "details": [
      {
        "field": "tags",
        "value": {
          "environment": "D-137",
          "owner": "Rick",
          "application": "Portal",
          "implementor": "Morty"
        }
      }
    ]
  }
  }
POLICY_RULE
}
`, data.RandomInteger)
}

func (r PolicyDefinitionResource) managementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%d"
}

resource "azurerm_policy_definition" "test" {
  name                = "acctestpol-%d"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "acctestpol-%d"
  management_group_id = azurerm_management_group.test.id

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
`, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r PolicyDefinitionResource) metadata(data acceptance.TestData) string {
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

  metadata = <<METADATA
  {
  		"foo": "bar"
  }
METADATA
}
`, data.RandomInteger, data.RandomInteger)
}

func (r PolicyDefinitionResource) mode(data acceptance.TestData, mode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestpol-%d"
  policy_type  = "Custom"
  mode         = "%s"
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
`, data.RandomInteger, mode, data.RandomInteger)
}

func (r PolicyDefinitionResource) additionalParameter(data acceptance.TestData) string {
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
       "effect": "AuditIfNotExists",
        "details": {
          "type": "Microsoft.Insights/diagnosticSettings",
          "existenceCondition": {
            "allOf": [
            {
              "field": "Microsoft.Insights/diagnosticSettings/logs[*].retentionPolicy.enabled",
              "equals": "true"
            },
            {
              "field": "Microsoft.Insights/diagnosticSettings/logs[*].retentionPolicy.days",
              "equals": "[parameters('requiredRetentionDays')]"
            }
          ]
        }
      }
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
    },
    "requiredRetentionDays": {
        "type": "Integer",
        "defaultValue": 365,
        "allowedValues": [
          0,
          30,
          90,
          180,
          365
        ],
        "metadata": {
          "displayName": "Required retention (days)",
          "description": "The required diagnostic logs retention in days"
      }
    }
  }
PARAMETERS
}
`, data.RandomInteger, data.RandomInteger)
}
