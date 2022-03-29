package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceGroupPolicyRemediationResource struct{}

func TestAccAzureRMResourceGroupPolicyRemediation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_policy_remediation", "test")
	r := ResourceGroupPolicyRemediationResource{}

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

func TestAccAzureRMResourceGroupPolicyRemediation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_policy_remediation", "test")
	r := ResourceGroupPolicyRemediationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ResourceGroupPolicyRemediationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ResourceGroupPolicyRemediationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Policy.RemediationsClient.GetAtResourceGroup(ctx, id.SubscriptionId, id.ResourceGroup, id.RemediationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Policy Remediation %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.RemediationProperties != nil), nil
}

func (r ResourceGroupPolicyRemediationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-policy-%[1]s"
  location = "%[2]s"
}

resource "azurerm_policy_definition" "test" {
  name         = "acctestDef-%[1]s"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "my-policy-definition"

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

resource "azurerm_resource_group_policy_assignment" "test" {
  name                 = "acctestpa-%[1]s"
  resource_group_id    = azurerm_resource_group.test.id
  policy_definition_id = azurerm_policy_definition.test.id

  non_compliance_message {
    content = "test"
  }

  parameters = jsonencode({
    "allowedLocations" = {
      "value" = [azurerm_resource_group.test.location]
    }
  })
}
`, data.RandomString, data.Locations.Primary)
}

func (r ResourceGroupPolicyRemediationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group_policy_remediation" "test" {
  name                 = "acctestremediation-%[2]s"
  resource_group_id    = azurerm_resource_group_policy_assignment.test.resource_group_id
  policy_assignment_id = azurerm_resource_group_policy_assignment.test.id
}
`, r.template(data), data.RandomString)
}

func (r ResourceGroupPolicyRemediationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group_policy_remediation" "test" {
  name                    = "acctestremediation-%[2]s"
  resource_group_id       = azurerm_resource_group_policy_assignment.test.resource_group_id
  policy_assignment_id    = azurerm_resource_group_policy_assignment.test.id
  location_filters        = ["westus"]
  policy_definition_id    = azurerm_policy_definition.test.id
  resource_discovery_mode = "ReEvaluateCompliance"
}
`, r.template(data), data.RandomString)
}
