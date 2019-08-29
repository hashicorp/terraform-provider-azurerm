package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMPolicyAssignment_basic(t *testing.T) {
	resourceName := "azurerm_policy_assignment.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPolicyAssignment_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_policy_assignment.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(resourceName),
				),
			},
			{
				Config:      testAzureRMPolicyAssignment_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_policy_assignment"),
			},
		},
	})
}

func TestAccAzureRMPolicyAssignment_deployIfNotExists_policy(t *testing.T) {
	resourceName := "azurerm_policy_assignment.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_deployIfNotExists_policy(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPolicyAssignment_complete(t *testing.T) {
	resourceName := "azurerm_policy_assignment.test"

	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMPolicyAssignment_not_scopes(t *testing.T) {
	resourceName := "azurerm_policy_assignment.test"

	ri := tf.AccRandTimeInt()

	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyAssignment_not_scopes(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyAssignmentExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMPolicyAssignmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		client := testAccProvider.Meta().(*ArmClient).policy.AssignmentsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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
	client := testAccProvider.Meta().(*ArmClient).policy.AssignmentsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAzureRMPolicyAssignment_basic(ri int, location string) string {
	return fmt.Sprintf(`
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
      "effect": "audit"
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
  scope                = "${azurerm_resource_group.test.id}"
  policy_definition_id = "${azurerm_policy_definition.test.id}"
}
`, ri, ri, location, ri, location, ri)
}

func testAzureRMPolicyAssignment_requiresImport(ri int, location string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_policy_assignment" "import" {
  name                 = "${azurerm_policy_assignment.test.name}"
  scope                = "${azurerm_policy_assignment.test.scope}"
  policy_definition_id = "${azurerm_policy_assignment.test.policy_definition_id}"
}
`, testAzureRMPolicyAssignment_basic(ri, location))
}

func testAzureRMPolicyAssignment_deployIfNotExists_policy(ri int, location string) string {
	return fmt.Sprintf(`
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
  scope                = "${azurerm_resource_group.test.id}"
  policy_definition_id = "${azurerm_policy_definition.test.id}"

  identity {
    type = "SystemAssigned"
  }

  location = "%s"
}
`, ri, ri, ri, location, ri, location)
}

func testAzureRMPolicyAssignment_complete(ri int, location string) string {
	return fmt.Sprintf(`
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
  scope                = "${azurerm_resource_group.test.id}"
  policy_definition_id = "${azurerm_policy_definition.test.id}"
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
`, ri, ri, ri, location, ri, ri, location)
}

func testAzureRMPolicyAssignment_not_scopes(ri int, location string) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "current" {}

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
  scope                = "${data.azurerm_subscription.current.id}"
  policy_definition_id = "${azurerm_policy_definition.test.id}"
  description          = "Policy Assignment created via an Acceptance Test"
  not_scopes           = ["${azurerm_resource_group.test.id}"]
  display_name         = "Acceptance Test Run %d"

  parameters = <<PARAMETERS
{
  "allowedLocations": {
    "value": [ "%s" ]
  }
}
PARAMETERS
}
`, ri, ri, ri, location, ri, ri, location)
}
