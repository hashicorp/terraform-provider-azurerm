package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMPolicyAssignment_basicCustom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_basicCustom(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyAssignment_basicBuiltin(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_basicBuiltin(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyAssignment_basicBuiltInSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_basicBuiltInSet(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyAssignment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_basicCustom(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAzureRMPolicyAssignment_requiresImport),
		},
	})
}

func TestAccAzureRMPolicyAssignment_deployIfNotExists_policy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_deployIfNotExists_policy(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyAssignment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyAssignment_not_scopes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_not_scopes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyAssignment_enforcement_mode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_assignment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_enforcement_mode(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMPolicyAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.AssignmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		id := rs.Primary.ID
		resp, err := client.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("Bad: Get on policyAssignmentsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Policy Assignment does not exist: %s", id)
		}

		return nil
	}
}

func testCheckAzureRMPolicyAssignmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.AssignmentsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_policy_definition" {
			continue
		}

		id := rs.Primary.ID
		resp, err := client.GetByID(ctx, id)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Policy Assignment still exists:%s", *resp.Name)
		}
	}

	return nil
}

func testAzureRMPolicyAssignment_basicCustom(data acceptance.TestData) string {
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

func testAzureRMPolicyAssignment_basicBuiltInSet(data acceptance.TestData) string {
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

func testAzureRMPolicyAssignment_basicBuiltin(data acceptance.TestData) string {
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

func testAzureRMPolicyAssignment_requiresImport(data acceptance.TestData) string {
	template := testAzureRMPolicyAssignment_basicCustom(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_assignment" "import" {
  name                 = azurerm_policy_assignment.test.name
  scope                = azurerm_policy_assignment.test.scope
  policy_definition_id = azurerm_policy_assignment.test.policy_definition_id
}
`, template)
}

func testAzureRMPolicyAssignment_deployIfNotExists_policy(data acceptance.TestData) string {
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

func testAzureRMPolicyAssignment_complete(data acceptance.TestData) string {
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

func testAzureRMPolicyAssignment_not_scopes(data acceptance.TestData) string {
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

func testAzureRMPolicyAssignment_enforcement_mode(data acceptance.TestData) string {
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
