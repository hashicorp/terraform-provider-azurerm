// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package policy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagementGroupPolicyRemediationResource struct{}

func TestAccAzureRMManagementGroupPolicyRemediation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_remediation", "test")
	r := ManagementGroupPolicyRemediationResource{}

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

func TestAccAzureRMManagementGroupPolicyRemediation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_remediation", "test")
	r := ManagementGroupPolicyRemediationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMManagementGroupPolicyRemediation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_policy_remediation", "test")
	r := ManagementGroupPolicyRemediationResource{}

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

func (r ManagementGroupPolicyRemediationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	parsed, err := parse.ParseManagementGroupRemediationID(state.ID)
	if err != nil {
		return nil, err
	}
	id := parsed.ToRemediationID()

	resp, err := client.Policy.RemediationsClient.GetAtResource(ctx, id)
	if err != nil || resp.Model == nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving Policy Remediation %q: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagementGroupPolicyRemediationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_management_group" "test" {
  display_name = "Acceptance Test MgmtGroup %[1]d"
}

data "azurerm_policy_definition" "test" {
  display_name = "Allowed locations"
}

resource "azurerm_management_group_policy_assignment" "test" {
  name                 = "acctestpa-mg-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_definition_id = data.azurerm_policy_definition.test.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = ["%[3]s"]
    }
  })
}
`, data.RandomInteger, data.Locations.Primary, data.Locations.Secondary)
}

func (r ManagementGroupPolicyRemediationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_management_group_policy_remediation" "test" {
  name                 = "acctestremediation-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_assignment_id = azurerm_management_group_policy_assignment.test.id
}
`, r.template(data), data.RandomString)
}

func (r ManagementGroupPolicyRemediationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_management_group_policy_remediation" "import" {
  name                 = azurerm_management_group_policy_remediation.test.id
  management_group_id  = azurerm_management_group_policy_remediation.test.management_group_id
  policy_assignment_id = azurerm_management_group_policy_remediation.test.policy_assignment_id
}
`, r.basic(data))
}

func (r ManagementGroupPolicyRemediationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_management_group_policy_remediation" "test" {
  name                 = "acctestremediation-%[2]s"
  management_group_id  = azurerm_management_group.test.id
  policy_assignment_id = azurerm_management_group_policy_assignment.test.id
  location_filters     = [%[3]q]
}
`, r.template(data), data.RandomString, data.Locations.Secondary)
}
