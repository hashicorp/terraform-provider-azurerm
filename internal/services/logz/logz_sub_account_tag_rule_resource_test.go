// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logz_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogzSubAccountTagRuleResource struct{}

func TestAccLogzSubAccountTagRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_sub_account_tag_rule", "test")
	r := LogzSubAccountTagRuleResource{}
	email := "72776074-85d8-47fc-a1c2-cd23693522da@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogzSubAccountTagRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_sub_account_tag_rule", "test")
	r := LogzSubAccountTagRuleResource{}
	email := "59cb83fc-bb49-414b-8762-e5cd463f1463@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, email),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccLogzSubAccountTagRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_sub_account_tag_rule", "test")
	r := LogzSubAccountTagRuleResource{}
	email := "77e97659-200c-41f1-824d-f596c9566aee@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogzSubAccountTagRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logz_sub_account_tag_rule", "test")
	r := LogzSubAccountTagRuleResource{}
	email := "a9c283cf-9045-40f0-aecf-d8738e5ed103@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LogzSubAccountTagRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tagrules.ParseAccountTagRuleID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Logz.TagRuleClient.SubAccountTagRulesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r LogzSubAccountTagRuleResource) template(data acceptance.TestData, email string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-logz-%d"
  location = "%s"
}

resource "azurerm_logz_monitor" "test" {
  name                = "acctest-lm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  plan {
    billing_cycle  = "MONTHLY"
    effective_date = "%s"
    usage_type     = "COMMITTED"
  }

  user {
    email        = "%s"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }
}

resource "azurerm_logz_sub_account" "test" {
  name            = "acctest-lsa-%d"
  logz_monitor_id = azurerm_logz_monitor.test.id
  user {
    email        = azurerm_logz_monitor.test.user[0].email
    first_name   = azurerm_logz_monitor.test.user[0].first_name
    last_name    = azurerm_logz_monitor.test.user[0].last_name
    phone_number = azurerm_logz_monitor.test.user[0].phone_number
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, getEffectiveDate(), email, data.RandomInteger)
}

func (r LogzSubAccountTagRuleResource) basic(data acceptance.TestData, email string) string {
	template := r.template(data, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_sub_account_tag_rule" "test" {
  logz_sub_account_id = azurerm_logz_sub_account.test.id
}
`, template)
}

func (r LogzSubAccountTagRuleResource) requiresImport(data acceptance.TestData, email string) string {
	config := r.basic(data, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_sub_account_tag_rule" "import" {
  logz_sub_account_id = azurerm_logz_sub_account_tag_rule.test.logz_sub_account_id
}
`, config)
}

func (r LogzSubAccountTagRuleResource) complete(data acceptance.TestData, email string) string {
	template := r.template(data, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_sub_account_tag_rule" "test" {
  logz_sub_account_id = azurerm_logz_sub_account.test.id
  tag_filter {
    name   = "ccc"
    action = "Include"
    value  = "ccc"
  }

  tag_filter {
    name   = "bbb"
    action = "Exclude"
    value  = ""
  }

  tag_filter {
    name   = "ccc"
    action = "Include"
    value  = "ccc"
  }
  send_aad_logs          = true
  send_activity_logs     = true
  send_subscription_logs = true
}
`, template)
}
