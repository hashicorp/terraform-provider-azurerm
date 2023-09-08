// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
	id, err := remediations.ParseRemediationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Policy.RemediationsClient.GetAtSubscription(ctx, *id)
	if err != nil || resp.Model == nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Policy Remediation %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model.Properties != nil), nil
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
  name                 = "acctestpa-sub-%[1]d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = ["%[2]s", "%[3]s", "%[4]s"]
    }
  })
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary, data.Locations.Ternary)
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
  resource_discovery_mode = "ReEvaluateCompliance"
  failure_percentage      = 0.5
  parallel_deployments    = 3
  resource_count          = 3
}
`, r.template(data), data.RandomString)
}
