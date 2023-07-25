// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devtestlabs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/policies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DevTestPolicyResource struct{}

func TestAccDevTestPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_policy", "test")
	r := DevTestPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDevTestPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_policy", "test")
	r := DevTestPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_dev_test_policy"),
		},
	})
}

func TestAccDevTestPolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dev_test_policy", "test")
	r := DevTestPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Acceptance").HasValue("Test"),
			),
		},
		data.ImportStep(),
	})
}

func (DevTestPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := policies.ParsePolicyID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DevTestLabs.PoliciesClient.Get(ctx, *id, policies.GetOperationOptions{})
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (DevTestPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_policy" "test" {
  name                = "LabVmCount"
  policy_set_name     = "default"
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DevTestPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_policy" "import" {
  name                = azurerm_dev_test_policy.test.name
  policy_set_name     = azurerm_dev_test_policy.test.policy_set_name
  lab_name            = azurerm_dev_test_policy.test.lab_name
  resource_group_name = azurerm_dev_test_policy.test.resource_group_name
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
}
`, r.basic(data))
}

func (DevTestPolicyResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_test_policy" "test" {
  name                = "LabVmCount"
  policy_set_name     = "default"
  lab_name            = azurerm_dev_test_lab.test.name
  resource_group_name = azurerm_resource_group.test.name
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"
  description         = "Aloha this is the max number of VM's'"

  tags = {
    "Acceptance" = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
