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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourcePolicyExemptionResource struct{}

func TestAccAzureRMResourcePolicyExemption_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_exemption", "test")
	r := ResourcePolicyExemptionResource{}

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

func TestAccAzureRMResourcePolicyExemption_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_exemption", "test")
	r := ResourcePolicyExemptionResource{}
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

func TestAccAzureRMResourcePolicyExemption_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_exemption", "test")
	r := ResourcePolicyExemptionResource{}
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

func (r ResourcePolicyExemptionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ResourcePolicyExemptionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Policy.ExemptionsClient.Get(ctx, id.ResourceId, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}
	return utils.Bool(true), nil
}

func (r ResourcePolicyExemptionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_resource_policy_exemption" "test" {
  name                 = "acctest-exemption-%d"
  resource_id          = azurerm_resource_policy_assignment.test.resource_id
  policy_assignment_id = azurerm_resource_policy_assignment.test.id
  exemption_category   = "Mitigated"
}
`, ResourceAssignmentTestResource{}.withBuiltInPolicySetBasic(data), data.RandomInteger)
}

func (r ResourcePolicyExemptionResource) complete(data acceptance.TestData, endDate string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_policy_exemption" "test" {
  name                 = "acctest-exemption-%d"
  resource_id          = azurerm_resource_policy_assignment.test.resource_id
  policy_assignment_id = azurerm_resource_policy_assignment.test.id
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
`, ResourceAssignmentTestResource{}.withBuiltInPolicySetBasic(data), data.RandomInteger, endDate)
}
