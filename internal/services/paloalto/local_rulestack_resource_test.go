// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRulestackResource struct{}

func TestAccPaloAltoLocalRulestack_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack", "test")

	r := LocalRulestackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		data.ImportStep(),
	})
}

func TestAccPaloAltoLocalRulestack_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack", "test")

	r := LocalRulestackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccPaloAltoLocalRulestack_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack", "test")

	r := LocalRulestackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		data.ImportStep(),
	})
}

func TestAccPaloAltoLocalRulestack_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack", "test")

	r := LocalRulestackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check:  acceptance.ComposeTestCheckFunc(),
		},
		data.ImportStep(),
	})
}

func (r LocalRulestackResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := localrulestacks.ParseLocalRulestackID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.Client.LocalRulestacks.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r LocalRulestackResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack" "test" {
  name                = "testAcc-palrs-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
}

`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r LocalRulestackResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_palo_alto_local_rulestack" "test" {
  name                = "testAcc-palrs-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  anti_spyware_profile  = "BestPractice"
  anti_virus_profile    = "BestPractice"
  url_filtering_profile = "BestPractice"
  file_blocking_profile = "BestPractice"
  dns_subscription      = "BestPractice"
  vulnerability_profile = "BestPractice"

  description = "Acceptance Test Desc - %[2]d"
}


`, r.template(data), data.RandomInteger, data.Locations.Primary)
}

func (r LocalRulestackResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_palo_alto_local_rulestack" "import" {
  name                = azurerm_palo_alto_local_rulestack.test.name
  resource_group_name = azurerm_palo_alto_local_rulestack.test.resource_group_name
  location            = azurerm_palo_alto_local_rulestack.test.location
}

`, r.basic(data))
}

func (r LocalRulestackResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-PALRS-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}
