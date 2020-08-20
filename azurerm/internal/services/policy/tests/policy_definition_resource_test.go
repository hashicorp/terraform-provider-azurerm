package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMPolicyDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyDefinition_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyDefinitionExists(data.ResourceName, "All"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyDefinition_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyDefinitionExists(data.ResourceName, "All"),
				),
			},
			data.RequiresImportErrorStep(testAzureRMPolicyDefinition_requiresImport),
		},
	})
}

func TestAccAzureRMPolicyDefinition_computedMetadata(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyDefinition_computedMetadata(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyDefinitionExists(data.ResourceName, "Indexed"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyDefinitionAtMgmtGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyDefinitionDestroyInMgmtGroup,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyDefinition_managementGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyDefinitionExistsInMgmtGroup(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyDefinition_metadata(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyDefinition_metadata(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyDefinitionExists(data.ResourceName, "All"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyDefinition_mode_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	number := data.RandomInteger
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyDefinition_mode(number, "All"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyDefinitionExists(data.ResourceName, "All"),
				),
			},
			{
				Config: testAzureRMPolicyDefinition_mode(number, "Indexed"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyDefinitionExists(data.ResourceName, "Indexed"),
				),
			},
			{
				Config: testAzureRMPolicyDefinition_mode(number, "All"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyDefinitionExists(data.ResourceName, "All"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMPolicyDefinitionExistsInMgmtGroup(policyName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.DefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[policyName]
		if !ok {
			return fmt.Errorf("not found: %s", policyName)
		}

		id, err := parse.PolicyDefinitionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		scopeId, ok := id.PolicyScopeId.(parse.ScopeAtManagementGroup)
		if !ok {
			return fmt.Errorf("Bad: cannot get the management group from Policy Definition %q", id.Name)
		}

		if resp, err := client.GetAtManagementGroup(ctx, id.Name, scopeId.ManagementGroupName); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Policy Definition %q does not exist", id.Name)
			}
			return fmt.Errorf("Bad: GetAtManagementGroup on Policy.DefinitionsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMPolicyDefinitionDestroyInMgmtGroup(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.DefinitionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_policy_definition" {
			continue
		}

		id, err := parse.PolicyDefinitionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		scopeId, ok := id.PolicyScopeId.(parse.ScopeAtManagementGroup)
		if !ok {
			return fmt.Errorf("Bad: cannot get the management group from Policy Definition %q", id.Name)
		}

		if resp, err := client.GetAtManagementGroup(ctx, id.Name, scopeId.ManagementGroupName); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Policy.DefinitionsClient: %+v", err)
			}
		}
	}

	return nil
}

func testCheckAzureRMPolicyDefinitionExists(resourceName string, mode string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.DefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		id, err := parse.PolicyDefinitionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := client.Get(ctx, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Policy Definition %q does not exist", id.Name)
			}
			return fmt.Errorf("Bad: Get on Policy.DefinitionsClient: %+v", err)
		}

		if mode != *resp.DefinitionProperties.Mode {
			return fmt.Errorf("Bad: Policy Definition Mode is different. Expected: %s, Actual: %s", mode, *resp.DefinitionProperties.Mode)
		}

		return nil
	}
}

func testCheckAzureRMPolicyDefinitionDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.DefinitionsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_policy_definition" {
			continue
		}

		id, err := parse.PolicyDefinitionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		if resp, err := client.Get(ctx, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Policy.DefinitionsClient: %+v", err)
			}
		}
	}

	return nil
}

func testAzureRMPolicyDefinition_basic(data acceptance.TestData) string {
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

func testAzureRMPolicyDefinition_requiresImport(data acceptance.TestData) string {
	template := testAzureRMPolicyDefinition_basic(data)
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

func testAzureRMPolicyDefinition_computedMetadata(data acceptance.TestData) string {
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

func testAzureRMPolicyDefinition_managementGroup(data acceptance.TestData) string {
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
  management_group_id = azurerm_management_group.test.group_id

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

func testAzureRMPolicyDefinition_metadata(data acceptance.TestData) string {
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

func testAzureRMPolicyDefinition_mode(number int, mode string) string {
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
`, number, mode, number)
}
