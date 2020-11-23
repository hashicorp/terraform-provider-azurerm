package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PolicyAssignmentResource struct{}

func TestAccAzureRMPolicyAssignment_basicCustom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	r := PolicyAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicCustom(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyAssignment_basicBuiltin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	r := PolicyAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicBuiltin(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyAssignment_basicBuiltInSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	r := PolicyAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicBuiltInSet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	r := PolicyAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicCustom(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMPolicyAssignment_deployIfNotExists_policy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	r := PolicyAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.deployIfNotExistsPolicy(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyAssignment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	r := PolicyAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyAssignment_not_scopes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	r := PolicyAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.notScopes(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyAssignment_enforcement_mode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	r := PolicyAssignmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.enforcementMode(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PolicyAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	assignmentsClient := client.Policy.AssignmentsClient
	resp, err := assignmentsClient.GetByID(ctx, state.ID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Policy Assignment %q: %+v", state.ID, err)
	}
	return utils.Bool(resp.AssignmentProperties != nil), nil
}

func (r PolicyAssignmentResource) basicCustom(data acceptance.TestData) string {
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
        "equals": "%[2]s"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
POLICY_RULE

}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestpa-%[1]d"
  scope                = azurerm_resource_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id

  metadata = <<METADATA
  {
    "category": "General"
  }
METADATA

}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PolicyAssignmentResource) basicBuiltInSet(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_policy_set_definition" "test" {
  display_name = "Audit machines with insecure password security settings"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestpa-%[1]d"
  location             = azurerm_resource_group.test.location
  scope                = azurerm_resource_group.test.id
  policy_definition_id = data.azurerm_policy_set_definition.test.id

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PolicyAssignmentResource) basicBuiltin(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestpa-%[1]d"
  scope                = azurerm_resource_group.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters           = <<PARAMETERS
{
  "listOfAllowedLocations": {
    "value": [ "%[2]s" ]
  }
}
PARAMETERS
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PolicyAssignmentResource) requiresImport(data acceptance.TestData) string {
	template := r.basicCustom(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_assignment" "import" {
  name                 = azurerm_policy_assignment.test.name
  scope                = azurerm_policy_assignment.test.scope
  policy_definition_id = azurerm_policy_assignment.test.policy_definition_id
}
`, template)
}

func (r PolicyAssignmentResource) deployIfNotExistsPolicy(data acceptance.TestData) string {
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
		"equals": "Microsoft.Sql/servers/databases"
	},
	"then": {
		"effect": "DeployIfNotExists",
		"details": {
			"type": "Microsoft.Sql/servers/databases/transparentDataEncryption",
			"name": "current",
			"roleDefinitionIds": [
				"/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c"
			],
			"existenceCondition": {
				"field": "Microsoft.Sql/transparentDataEncryption.status",
				"equals": "Enabled"
			},
			"deployment": {
				"properties": {
					"mode": "incremental",
					"template": {
						"$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
						"contentVersion": "1.0.0.0",
						"parameters": {
							"fullDbName": {
								"type": "string"
							}
						},
						"resources": [{
							"name": "[concat(parameters('fullDbName'), '/current')]",
							"type": "Microsoft.Sql/servers/databases/transparentDataEncryption",
							"apiVersion": "2014-04-01",
							"properties": {
								"status": "Enabled"
							}
						}]
					},
					"parameters": {
						"fullDbName": {
							"value": "[field('fullName')]"
						}
					}
				}
			}
		}
	}
}
POLICY_RULE

}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestpa-%d"
  scope                = azurerm_resource_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id

  identity {
    type = "SystemAssigned"
  }

  location = "%s"
}
`, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (r PolicyAssignmentResource) complete(data acceptance.TestData) string {
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

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestpa-%d"
  scope                = azurerm_resource_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "Policy Assignment created via an Acceptance Test"
  display_name         = "Acceptance Test Run %d"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "%s" ]
  }
}
PARAMETERS

}
`, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (r PolicyAssignmentResource) notScopes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {
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

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestpa-%d"
  scope                = data.azurerm_subscription.current.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "Policy Assignment created via an Acceptance Test"
  not_scopes           = [azurerm_resource_group.test.id]
  display_name         = "Acceptance Test Run %d"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "%s" ]
  }
}
PARAMETERS

}
`, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (r PolicyAssignmentResource) enforcementMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {
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

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_policy_assignment" "test" {
  name                 = "acctestpa-%d"
  scope                = data.azurerm_subscription.current.id
  policy_definition_id = azurerm_policy_definition.test.id
  description          = "Policy Assignment created via an Acceptance Test"
  enforcement_mode     = false
  display_name         = "Acceptance Test Run %d"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "%s" ]
  }
}
PARAMETERS

}
`, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}
