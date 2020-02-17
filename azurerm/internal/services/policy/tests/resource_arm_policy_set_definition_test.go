package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAzureRMPolicySetDefinition_builtIn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_builtIn(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicySetDefinition_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_builtIn(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAzureRMPolicySetDefinition_requiresImport),
		},
	})
}

func TestAccAzureRMPolicySetDefinition_custom(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_custom(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicySetDefinition_ManagementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_ManagementGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAzureRMPolicySetDefinition_builtIn(data acceptance.TestData) string {
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
`, data.RandomInteger, data.RandomInteger)
}

func testAzureRMPolicySetDefinition_requiresImport(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_builtIn(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_set_definition" "import" {
  name         = "${azurerm_policy_set_definition.test.name}"
  policy_type  = "${azurerm_policy_set_definition.test.policy_type}"
  display_name = "${azurerm_policy_set_definition.test.display_name}"
  parameters   = "${azurerm_policy_set_definition.test.parameters}"
}
`, template)
}

func testAzureRMPolicySetDefinition_custom(data acceptance.TestData) string {
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
`, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAzureRMPolicySetDefinition_ManagementGroup(data acceptance.TestData) string {
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
`, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testCheckAzureRMPolicySetDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.SetDefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
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
	client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.SetDefinitionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
