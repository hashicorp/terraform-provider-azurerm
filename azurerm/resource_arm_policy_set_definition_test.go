package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform/terraform"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMPolicySetDefinition_builtIn(t *testing.T) {
	resourceName := "azurerm_policy_set_definition.test"

	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_builtIn(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(resourceName),
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

func TestAccAzureRMPolicySetDefinition_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_policy_set_definition.test"

	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_builtIn(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(resourceName),
				),
			},
			{
				Config:      testAzureRMPolicySetDefinition_requiresImport(ri),
				ExpectError: testRequiresImportError("azurerm_policy_set_definition"),
			},
		},
	})
}

func TestAccAzureRMPolicySetDefinition_custom(t *testing.T) {
	resourceName := "azurerm_policy_set_definition.test"

	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_custom(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(resourceName),
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

func TestAccAzureRMPolicySetDefinition_ManagementGroup(t *testing.T) {
	resourceName := "azurerm_policy_set_definition.test"

	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_ManagementGroup(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(resourceName),
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

func testAzureRMPolicySetDefinition_builtIn(ri int) string {
	return fmt.Sprintf(`
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

  policy_definitions = <<POLICY_DEFINITIONS
    [
        {
            "parameters": {
                "listOfAllowedLocations": {
                    "value": "[parameters('allowedLocations')]"
                }
            },
            "policyDefinitionId": "/providers/Microsoft.Authorization/policyDefinitions/e765b5de-1225-4ba3-bd56-1ac6695af988"
        }
    ]
POLICY_DEFINITIONS
}
`, ri, ri)
}

func testAzureRMPolicySetDefinition_requiresImport(ri int) string {
	return fmt.Sprintf(`
%s 

resource "azurerm_policy_set_definition" "import" {
  name         = "${azurerm_policy_set_definition.test.name}"
  policy_type  = "${azurerm_policy_set_definition.test.policy_type}"
  display_name = "${azurerm_policy_set_definition.test.display_name}"
  parameters   = "${azurerm_policy_set_definition.test.parameters}"
}
`, testAzureRMPolicySetDefinition_builtIn(ri))
}

func testAzureRMPolicySetDefinition_custom(ri int) string {
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

  policy_definitions = <<POLICY_DEFINITIONS
    [
        {
            "parameters": {
                "allowedLocations": {
                    "value": "[parameters('allowedLocations')]"
                }
            },
            "policyDefinitionId": "${azurerm_policy_definition.test.id}"
        }
    ]
POLICY_DEFINITIONS
}
`, ri, ri, ri, ri)
}

func testAzureRMPolicySetDefinition_ManagementGroup(ri int) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%d"
}

resource "azurerm_policy_set_definition" "test" {
  name                = "acctestpolset-%d"
  policy_type         = "Custom"
  display_name        = "acctestpolset-%d"
  management_group_id = "${azurerm_management_group.test.group_id}"

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

  policy_definitions = <<POLICY_DEFINITIONS
    [
        {
            "parameters": {
                "listOfAllowedLocations": {
                    "value": "[parameters('allowedLocations')]"
                }
            },
            "policyDefinitionId": "/providers/Microsoft.Authorization/policyDefinitions/e765b5de-1225-4ba3-bd56-1ac6695af988"
        }
    ]
POLICY_DEFINITIONS
}
`, ri, ri, ri)
}

func testCheckAzureRMPolicySetDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		policySetName := rs.Primary.Attributes["name"]
		managementGroupId := rs.Primary.Attributes["management_group_id"]

		client := testAccProvider.Meta().(*ArmClient).policy.SetDefinitionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		var err error
		var resp policy.SetDefinition
		if managementGroupId != "" {
			resp, err = client.GetAtManagementGroup(ctx, policySetName, managementGroupId)
		} else {
			resp, err = client.Get(ctx, policySetName)
		}

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("policy set definition does not exist: %s", policySetName)
			} else {
				return fmt.Errorf("Bad: Get on policySetDefinitionsClient: %s", err)
			}
		}

		return nil
	}
}

func testCheckAzureRMPolicySetDefinitionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).policy.SetDefinitionsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_policy_set_definition" {
			continue
		}

		policySetName := rs.Primary.Attributes["name"]
		managementGroupId := rs.Primary.Attributes["management_group_id"]

		var err error
		var resp policy.SetDefinition
		if managementGroupId != "" {
			resp, err = client.GetAtManagementGroup(ctx, policySetName, managementGroupId)
		} else {
			resp, err = client.Get(ctx, policySetName)
		}

		if err == nil {
			return fmt.Errorf("policy set definition still exists: %s", *resp.Name)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		} else {
			return err
		}
	}

	return nil
}
