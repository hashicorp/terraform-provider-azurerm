// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	resourceParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceGroupPolicyExemptionResource struct{}

func TestAccAzureRMResourceGroupPolicyExemption_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_policy_exemption", "test")
	r := ResourceGroupPolicyExemptionResource{}

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

func TestAccAzureRMResourceGroupPolicyExemption_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_policy_exemption", "test")
	r := ResourceGroupPolicyExemptionResource{}
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

func TestAccAzureRMResourceGroupPolicyExemption_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_policy_exemption", "test")
	r := ResourceGroupPolicyExemptionResource{}
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

func (r ResourceGroupPolicyExemptionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ResourceGroupPolicyExemptionID(state.ID)
	if err != nil {
		return nil, err
	}

	resourceGroupId := resourceParse.NewResourceGroupID(id.SubscriptionId, id.ResourceGroup)

	resp, err := client.Policy.ExemptionsClient.Get(ctx, resourceGroupId.ID(), id.PolicyExemptionName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}
	return utils.Bool(true), nil
}

func (r ResourceGroupPolicyExemptionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_resource_group_policy_exemption" "test" {
  name                 = "acctest-exemption-%d"
  resource_group_id    = azurerm_resource_group.test.id
  policy_assignment_id = azurerm_resource_group_policy_assignment.test.id
  exemption_category   = "Mitigated"
}
`, ResourceGroupAssignmentTestResource{}.withBuiltInPolicySetBasic(data), data.RandomInteger)
}

func (r ResourceGroupPolicyExemptionResource) complete(data acceptance.TestData, endDate string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group_policy_exemption" "test" {
  name                 = "acctest-exemption-%d"
  resource_group_id    = azurerm_resource_group.test.id
  policy_assignment_id = azurerm_resource_group_policy_assignment.test.id
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
`, ResourceGroupAssignmentTestResource{}.withBuiltInPolicySetBasic(data), data.RandomInteger, endDate)
}
