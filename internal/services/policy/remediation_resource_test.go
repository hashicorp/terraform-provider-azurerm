// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	policyinsights "github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourcePolicyRemediationResource struct{}

func TestAccAzureRMResourcePolicyRemediation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_remediation", "test")
	r := ResourcePolicyRemediationResource{}

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

func TestAccAzureRMResourcePolicyRemediation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_policy_remediation", "test")
	r := ResourcePolicyRemediationResource{}

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

func (r ResourcePolicyRemediationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := policyinsights.ParseScopedRemediationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Policy.RemediationsClient.GetAtResource(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Policy Remediation %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ResourcePolicyRemediationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-policy-%[1]s"
  location = %[2]q
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%[1]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_resource_policy_assignment" "test" {
  name                 = "acctestpa-res-%[1]s"
  resource_id          = azurerm_virtual_network.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = [azurerm_resource_group.test.location, "%[3]s"]
    }
  })
}
`, data.RandomString, data.Locations.Primary, data.Locations.Secondary)
}

func (r ResourcePolicyRemediationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_policy_remediation" "test" {
  name                 = "acctestremediation-%[2]s"
  resource_id          = azurerm_virtual_network.test.id
  policy_assignment_id = azurerm_resource_policy_assignment.test.id
}
`, r.template(data), data.RandomString)
}

func (r ResourcePolicyRemediationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_policy_remediation" "test" {
  name                    = "acctestremediation-%[2]s"
  resource_id             = azurerm_virtual_network.test.id
  policy_assignment_id    = azurerm_resource_policy_assignment.test.id
  location_filters        = ["westus"]
  resource_discovery_mode = "ReEvaluateCompliance"
  failure_percentage      = 0.5
  parallel_deployments    = 3
  resource_count          = 3
}
`, r.template(data), data.RandomString)
}
