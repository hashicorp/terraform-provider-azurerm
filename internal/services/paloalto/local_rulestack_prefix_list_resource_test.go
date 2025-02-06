// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/prefixlistlocalrulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRuleStackPrefixList struct{}

func TestAccLocalRulestackPrefixList_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_prefix_list", "test")

	r := LocalRuleStackPrefixList{}

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

func TestAccLocalRulestackPrefixList_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_prefix_list", "test")

	r := LocalRuleStackPrefixList{}

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

func TestAccLocalRulestackPrefixList_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_prefix_list", "test")

	r := LocalRuleStackPrefixList{}

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

func TestAccLocalRulestackPrefixList_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_prefix_list", "test")

	r := LocalRuleStackPrefixList{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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

func (r LocalRuleStackPrefixList) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := prefixlistlocalrulestack.ParseLocalRulestackPrefixListID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.Client.PrefixListLocalRulestack.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r LocalRuleStackPrefixList) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack_prefix_list" "test" {
  name         = "testacc-palr-%[2]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  prefix_list = ["10.0.0.0/8", "172.16.0.0/16"]
}
`, r.template(data), data.RandomInteger)
}

func (r LocalRuleStackPrefixList) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_palo_alto_local_rulestack_prefix_list" "import" {
  name         = azurerm_palo_alto_local_rulestack_prefix_list.test.name
  rulestack_id = azurerm_palo_alto_local_rulestack_prefix_list.test.rulestack_id

  prefix_list = azurerm_palo_alto_local_rulestack_prefix_list.test.prefix_list
}
`, r.basic(data))
}

func (r LocalRuleStackPrefixList) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack_prefix_list" "test" {
  name         = "testacc-palr-%[2]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  prefix_list = ["10.0.0.0/8", "172.16.0.0/16"]

  audit_comment = "Updated acceptance test audit comment - %[2]d"
  description   = "Updated acceptance test Desc - %[2]d"

}
`, r.template(data), data.RandomInteger)
}

func (r LocalRuleStackPrefixList) template(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-PAN-%[1]d"
  location = "%[2]s"
}

resource "azurerm_palo_alto_local_rulestack" "test" {
  name                = "testAcc-palrs-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
}

`, data.RandomInteger, data.Locations.Primary)
}
