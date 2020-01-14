package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
					testCheckAzureRMPolicyDefinitionExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMPolicyDefinition_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyDefinition_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyDefinitionExists(data.ResourceName),
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
					testCheckAzureRMPolicyDefinitionExists(data.ResourceName),
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
		CheckDestroy: testCheckAzureRMPolicyDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAzureRMPolicyDefinition_ManagementGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPolicyDefinitionExistsInMgmtGroup(data.ResourceName),
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

		policyName := rs.Primary.Attributes["name"]
		managementGroupID := rs.Primary.Attributes["management_group_id"]

		resp, err := client.GetAtManagementGroup(ctx, policyName, managementGroupID)
		if err != nil {
			return fmt.Errorf("Bad: GetAtManagementGroup on policyDefinitionsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("policy does not exist: %s", policyName)
		}

		return nil
	}
}

func testCheckAzureRMPolicyDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Policy.DefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		policyName := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, policyName)
		if err != nil {
			return fmt.Errorf("Bad: Get on policyDefinitionsClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("policy does not exist: %s", policyName)
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

		name := rs.Primary.Attributes["name"]

		resp, err := client.Get(ctx, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("policy still exists:%s", *resp.Name)
		}
	}

	return nil
}

func testAzureRMPolicyDefinition_basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.RandomInteger)
}

func testAzureRMPolicyDefinition_requiresImport(data acceptance.TestData) string {
	template := testAzureRMPolicyDefinition_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_definition" "import" {
  name         = "${azurerm_policy_definition.test.name}"
  policy_type  = "${azurerm_policy_definition.test.policy_type}"
  mode         = "${azurerm_policy_definition.test.mode}"
  display_name = "${azurerm_policy_definition.test.display_name}"
  policy_rule  = "${azurerm_policy_definition.test.policy_rule}"
  parameters   = "${azurerm_policy_definition.test.parameters}"
}
`, template)
}

func testAzureRMPolicyDefinition_computedMetadata(data acceptance.TestData) string {
	return fmt.Sprintf(`
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

func testAzureRMPolicyDefinition_ManagementGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%d"
}

resource "azurerm_policy_definition" "test" {
  name                = "acctestpol-%d"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "acctestpol-%d"
  management_group_id = "${azurerm_management_group.test.group_id}"

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
