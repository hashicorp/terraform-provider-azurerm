// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SubscriptionPolicyExemptionResource struct{}

func TestAccAzureRMSubscriptionPolicyExemption_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_policy_exemption", "test")
	r := SubscriptionPolicyExemptionResource{}

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

func TestAccAzureRMSubscriptionPolicyExemption_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_policy_exemption", "test")
	r := SubscriptionPolicyExemptionResource{}
	endDate := time.Now().UTC().Add(time.Hour * 24).Format(time.RFC3339)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, endDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMSubscriptionPolicyExemption_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_policy_exemption", "test")
	r := SubscriptionPolicyExemptionResource{}
	endDate := time.Now().UTC().Add(time.Hour * 24).Format(time.RFC3339)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, endDate),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SubscriptionPolicyExemptionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SubscriptionPolicyExemptionID(state.ID)
	if err != nil {
		return nil, err
	}

	subscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	resp, err := client.Policy.ExemptionsClient.Get(ctx, subscriptionId.ID(), id.PolicyExemptionName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}
	return utils.Bool(true), nil
}

func (r SubscriptionPolicyExemptionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_subscription_policy_exemption" "test" {
  name                 = "acctest-exemption-%d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_assignment_id = azurerm_subscription_policy_assignment.test.id
  exemption_category   = "Mitigated"
}
`, SubscriptionAssignmentTestResource{}.withBuiltInPolicySetBasic(data), data.RandomInteger)
}

func (r SubscriptionPolicyExemptionResource) complete(data acceptance.TestData, endDate string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subscription_policy_exemption" "test" {
  name                 = "acctest-exemption-%d"
  subscription_id      = data.azurerm_subscription.test.id
  policy_assignment_id = azurerm_subscription_policy_assignment.test.id
  exemption_category   = "Waiver"

  display_name = "Policy Exemption for acceptance test"
  description  = "Policy Exemption created in an acceptance test"
  expires_on   = "%[3]s"

  metadata = <<METADATA
    {
        "foo": "bar"
    }
METADATA
}
`, SubscriptionAssignmentTestResource{}.withBuiltInPolicySetBasic(data), data.RandomInteger, endDate)
}
