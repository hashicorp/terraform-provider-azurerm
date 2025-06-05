// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/fqdnlistlocalrulestack"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRulestackFQDNList struct{}

func TestAccPaloAltoLocalRulestackFQDNList_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_fqdn_list", "test")

	r := LocalRulestackFQDNList{}

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

func TestAccPaloAltoLocalRulestackFQDNList_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_fqdn_list", "test")

	r := LocalRulestackFQDNList{}

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

func TestAccPaloAltoLocalRulestackFQDNList_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_fqdn_list", "test")

	r := LocalRulestackFQDNList{}

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
	})
}

func TestAccPaloAltoLocalRulestackFQDNList_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_palo_alto_local_rulestack_fqdn_list", "test")

	r := LocalRulestackFQDNList{}

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

func (r LocalRulestackFQDNList) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fqdnlistlocalrulestack.ParseLocalRulestackFqdnListID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.PaloAlto.Client.FqdnListLocalRulestack.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r LocalRulestackFQDNList) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_palo_alto_local_rulestack_fqdn_list" "test" {
  name         = "testacc-pafqdn-%[2]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  fully_qualified_domain_names = ["contoso.com", "test.example.com"]
}


`, r.template(data), data.RandomInteger)
}

func (r LocalRulestackFQDNList) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_palo_alto_local_rulestack_fqdn_list" "test" {
  name         = "testacc-pafqdn-%[2]d"
  rulestack_id = azurerm_palo_alto_local_rulestack.test.id

  fully_qualified_domain_names = ["contoso.com", "test.example.com", "anothertest.example.com"]

  audit_comment = "Acc Test Audit Comment - %[2]d"
  description   = "Acc Test Description - %[2]d"
}


`, r.template(data), data.RandomInteger)
}

func (r LocalRulestackFQDNList) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_palo_alto_local_rulestack_fqdn_list" "import" {
  name         = azurerm_palo_alto_local_rulestack_fqdn_list.test.name
  rulestack_id = azurerm_palo_alto_local_rulestack_fqdn_list.test.rulestack_id

  fully_qualified_domain_names = azurerm_palo_alto_local_rulestack_fqdn_list.test.fully_qualified_domain_names
}


`, r.basic(data))
}

func (r LocalRulestackFQDNList) template(data acceptance.TestData) string {
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
