package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PolicyExemptionResource struct{}

func TestAccAzureRMPolicyExemption_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_exemption", "test")
	r := PolicyExemptionResource{}

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

func TestAccAzureRMPolicyExemption_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_exemption", "test")
	r := PolicyExemptionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMPolicyExemption_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_policy_exemption", "test")
	r := PolicyExemptionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r PolicyExemptionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.PolicyExemptionID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Policy.ExemptionsClient.Get(ctx, id.ScopeId(), id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Policy Exemption %q (Scope %q): %+v", id.Name, id.ScopeId(), err)
	}
	return utils.Bool(true), nil
}

func (r PolicyExemptionResource) basic(data acceptance.TestData) string {
	template := PolicyAssignmentResource{}.basicCustom(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_exemption" "test" {
  name                 = "acctest-exemption-%d"
  scope                = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id

  exemption_category = "Mitigated"
}
`, template, data.RandomInteger)
}

func (r PolicyExemptionResource) complete(data acceptance.TestData) string {
	template := PolicyAssignmentResource{}.basicCustom(data)
	return fmt.Sprintf(`
%s

resource "azurerm_policy_exemption" "test" {
  name                 = "acctest-exemption-%d"
  scope                = azurerm_policy_assignment.test.scope
  policy_assignment_id = azurerm_policy_assignment.test.id

  exemption_category = "Waiver"

  display_name = "Policy Exemption for acceptance test"
  description  = "Policy Exemption created in an acceptance test"
}
`, template, data.RandomInteger)
}
