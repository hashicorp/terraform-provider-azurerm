// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountagreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogicAppIntegrationAccountAgreementResource struct{}

func TestAccLogicAppIntegrationAccountAgreement_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_agreement", "test")
	r := LogicAppIntegrationAccountAgreementResource{}

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

func TestAccLogicAppIntegrationAccountAgreement_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_agreement", "test")
	r := LogicAppIntegrationAccountAgreementResource{}

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

func TestAccLogicAppIntegrationAccountAgreement_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_agreement", "test")
	r := LogicAppIntegrationAccountAgreementResource{}

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

func TestAccLogicAppIntegrationAccountAgreement_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_integration_account_agreement", "test")
	r := LogicAppIntegrationAccountAgreementResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LogicAppIntegrationAccountAgreementResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := integrationaccountagreements.ParseAgreementID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Logic.IntegrationAccountAgreementClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r LogicAppIntegrationAccountAgreementResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-ia-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogicAppIntegrationAccountAgreementResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_partner" "host" {
  name                     = "acctest-hostpartner-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  business_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamNY"
  }
}

resource "azurerm_logic_app_integration_account_partner" "guest" {
  name                     = "acctest-guestpartner-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  business_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamDC"
  }
}

resource "azurerm_logic_app_integration_account_agreement" "test" {
  name                     = "acctest-agreement-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  agreement_type           = "AS2"
  host_partner_name        = azurerm_logic_app_integration_account_partner.host.name
  guest_partner_name       = azurerm_logic_app_integration_account_partner.guest.name
  content                  = file("testdata/integration_account_agreement_content_as2.json")

  host_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamNY"
  }

  guest_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamDC"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LogicAppIntegrationAccountAgreementResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_agreement" "import" {
  name                     = azurerm_logic_app_integration_account_agreement.test.name
  resource_group_name      = azurerm_logic_app_integration_account_agreement.test.resource_group_name
  integration_account_name = azurerm_logic_app_integration_account_agreement.test.integration_account_name
  agreement_type           = azurerm_logic_app_integration_account_agreement.test.agreement_type
  host_partner_name        = azurerm_logic_app_integration_account_agreement.test.host_partner_name
  guest_partner_name       = azurerm_logic_app_integration_account_agreement.test.guest_partner_name
  content                  = azurerm_logic_app_integration_account_agreement.test.content

  host_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamNY"
  }

  guest_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamDC"
  }
}
`, r.basic(data))
}

func (r LogicAppIntegrationAccountAgreementResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_partner" "host" {
  name                     = "acctest-hostpartner-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  business_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamNY"
  }
}

resource "azurerm_logic_app_integration_account_partner" "guest" {
  name                     = "acctest-guestpartner-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  business_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamDC"
  }
}

resource "azurerm_logic_app_integration_account_agreement" "test" {
  name                     = "acctest-agreement-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  agreement_type           = "AS2"
  host_partner_name        = azurerm_logic_app_integration_account_partner.host.name
  guest_partner_name       = azurerm_logic_app_integration_account_partner.guest.name
  content                  = file("testdata/integration_account_agreement_content_as2.json")

  host_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamNY"
  }

  guest_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamDC"
  }

  metadata = {
    foo = "bar"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r LogicAppIntegrationAccountAgreementResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_integration_account_partner" "host" {
  name                     = "acctest-hostpartner-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  business_identity {
    qualifier = "2"
    value     = "FabrikamNY"
  }
}

resource "azurerm_logic_app_integration_account_partner" "guest" {
  name                     = "acctest-guestpartner-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  business_identity {
    qualifier = "2"
    value     = "FabrikamDC"
  }
}

resource "azurerm_logic_app_integration_account_agreement" "test" {
  name                     = "acctest-agreement-%d"
  resource_group_name      = azurerm_resource_group.test.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  agreement_type           = "X12"
  host_partner_name        = azurerm_logic_app_integration_account_partner.host.name
  guest_partner_name       = azurerm_logic_app_integration_account_partner.guest.name
  content                  = file("testdata/integration_account_agreement_content_x12.json")

  host_identity {
    qualifier = "2"
    value     = "FabrikamNY"
  }

  guest_identity {
    qualifier = "2"
    value     = "FabrikamDC"
  }

  metadata = {
    foo = "bar2"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
