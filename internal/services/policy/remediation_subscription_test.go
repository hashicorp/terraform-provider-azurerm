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

type SubscriptionPolicyRemediationResource struct{}

func TestAccAzureRMSubscriptionPolicyRemediation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_policy_remediation", "test")
	r := SubscriptionPolicyRemediationResource{}

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

func TestAccAzureRMSubscriptionPolicyRemediation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_policy_remediation", "test")
	r := SubscriptionPolicyRemediationResource{}

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

func (r SubscriptionPolicyRemediationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SubscriptionPolicyRemediationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Policy.RemediationsClient.GetAtSubscription(ctx, id.SubscriptionId, id.RemediationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Policy Remediation %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.RemediationProperties != nil), nil
}

func (r SubscriptionPolicyRemediationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "test" {}

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_subscription_policy_assignment" "test" {
  name                 = "acctestpa-%[1]d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = ["%[2]s", "%[3]s"]
    }
  })
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r SubscriptionPolicyRemediationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subscription_policy_remediation" "test" {
  name                 = "acctestremediation-%[2]s"
  subscription_id      = data.azurerm_subscription.test.id
  policy_assignment_id = azurerm_subscription_policy_assignment.test.id
}
`, r.template(data), data.RandomString)
}

func (r SubscriptionPolicyRemediationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subscription_policy_remediation" "test" {
  name                    = "acctestremediation-%[2]s"
  subscription_id         = data.azurerm_subscription.test.id
  policy_assignment_id    = azurerm_subscription_policy_assignment.test.id
  location_filters        = ["westus"]
  policy_definition_id    = data.azurerm_policy_definition.test.id
  resource_discovery_mode = "ReEvaluateCompliance"
}
`, r.template(data), data.RandomString)
}
