package policy

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

type PolicyDefinitionResource struct{}

func TestAccAzureRMPolicyDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMPolicyDefinition_computedMetadata(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.computedMetadata(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyDefinitionAtMgmtGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.managementGroup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyDefinition_metadata(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.metadata(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyDefinition_modeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_definition", "test")
	r := PolicyDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.mode(data, "All"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.mode(data, "Indexed"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.mode(data, "All"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PolicyDefinitionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resp, err := getPolicyDefinitionByName(ctx, client.Policy.DefinitionsClient, state.Attributes["name"], state.Attributes["management_group_id"])
	if err != nil {
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
