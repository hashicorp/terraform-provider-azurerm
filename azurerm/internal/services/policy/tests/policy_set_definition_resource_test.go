package tests

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2020-03-01-preview/policy"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPolicySetDefinition_builtInDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_builtInDeprecated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

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

func TestAccAzureRMPolicySetDefinition_customDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_customDeprecated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
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

func TestAccAzureRMPolicySetDefinition_customNoParameter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_customNoParameter(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicySetDefinition_customUpdateDisplayName(t *testing.T) {
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
			{
				Config: testAccAzureRMPolicySetDefinition_customUpdateDisplayName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicySetDefinition_customUpdateParameters(t *testing.T) {
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
			{
				Config: testAccAzureRMPolicySetDefinition_customUpdateParameters(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMPolicySetDefinition_customUpdateAddNewReference(t *testing.T) {
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
			{
				Config: testAccAzureRMPolicySetDefinition_customUpdateAddNewReference(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
		},
	})
}

func TestAccAzureRMPolicySetDefinition_customWithPolicyReferenceID(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_customWithPolicyReferenceID(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicySetDefinition_managementGroupDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_managementGroupDeprecated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicySetDefinition_managementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_managementGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicySetDefinition_metadataDeprecated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_metadataDeprecated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicySetDefinition_metadata(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_set_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicySetDefinition_metadata(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicySetDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAzureRMPolicySetDefinition_builtInDeprecated(data acceptance.TestData) string {
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

func testAzureRMPolicySetDefinition_builtIn(data acceptance.TestData) string {
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

func testAzureRMPolicySetDefinition_requiresImport(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_builtInDeprecated(data)
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

func testAzureRMPolicySetDefinition_customDeprecated(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_template(data)
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
`, template, data.RandomInteger, data.RandomInteger)
}

func testAzureRMPolicySetDefinition_custom(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_template(data)
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

func testAccAzureRMPolicySetDefinition_customUpdateDisplayName(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_template(data)
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

func testAccAzureRMPolicySetDefinition_customUpdateParameters(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_template(data)
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

func testAccAzureRMPolicySetDefinition_customUpdateAddNewReference(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_template(data)
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

func testAzureRMPolicySetDefinition_customNoParameter(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_templateNoParameter(data)
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

func testAzureRMPolicySetDefinition_managementGroupDeprecated(data acceptance.TestData) string {
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
  management_group_id = azurerm_management_group.test.group_id

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

func testAzureRMPolicySetDefinition_managementGroup(data acceptance.TestData) string {
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
  management_group_id = azurerm_management_group.test.group_id

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

func testAzureRMPolicySetDefinition_metadataDeprecated(data acceptance.TestData) string {
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

  metadata = <<METADATA
    {
        "foo": "bar"
    }
METADATA
}
`, data.RandomInteger, data.RandomInteger)
}

func testAzureRMPolicySetDefinition_metadata(data acceptance.TestData) string {
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

func testAzureRMPolicySetDefinition_customWithPolicyReferenceID(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_template(data)
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

func testAzureRMPolicySetDefinition_template(data acceptance.TestData) string {
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

func testAzureRMPolicySetDefinition_templateNoParameter(data acceptance.TestData) string {
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

func testCheckAzureRMPolicySetDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.SetDefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		subscriptionId := acceptance.AzureProvider.Meta().(*clients.Client).Account.SubscriptionId

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		id, err := parse.PolicySetDefinitionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		var resp policy.SetDefinition
		if mgmtGroupID, ok := id.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
			resp, err = client.GetAtManagementGroup(ctx, id.Name, mgmtGroupID.ManagementGroupName)
		} else {
			resp, err = client.Get(ctx, id.Name, subscriptionId)
		}

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("policy set definition does not exist: %s", id.Name)
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
	subscriptionId := acceptance.AzureProvider.Meta().(*clients.Client).Account.SubscriptionId

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_policy_set_definition" {
			continue
		}

		id, err := parse.PolicySetDefinitionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		var resp policy.SetDefinition
		if mgmtGroupID, ok := id.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
			resp, err = client.GetAtManagementGroup(ctx, id.Name, mgmtGroupID.ManagementGroupName)
		} else {
			resp, err = client.Get(ctx, id.Name, subscriptionId)
		}

		if err == nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Policy.SetDefinitionsClient: %+v", err)
			}
		}
	}

	return nil
}
