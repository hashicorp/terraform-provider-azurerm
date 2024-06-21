// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logz_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogzTagRuleResource struct{}

func TestAccLogzTagRule_basic(t *testing.T) {
	t.Skipf("This resource is being deprecated by Azure and will be removed by end of 2025.")

	data := acceptance.BuildTestData(t, "azurerm_logz_tag_rule", "test")
	r := LogzTagRuleResource{}
	email := "71212724-7a73-48c3-9399-de59313d4905@example.com"
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

func TestAccLogzTagRule_requiresImport(t *testing.T) {
	t.Skipf("This resource is being deprecated by Azure and will be removed by end of 2025.")

	data := acceptance.BuildTestData(t, "azurerm_logz_tag_rule", "test")
	r := LogzTagRuleResource{}
	email := "b993f18e-9094-4a38-9e80-a0530ebbc6e2@example.com"
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

func TestAccLogzTagRule_complete(t *testing.T) {
	t.Skipf("This resource is being deprecated by Azure and will be removed by end of 2025.")

	data := acceptance.BuildTestData(t, "azurerm_logz_tag_rule", "test")
	r := LogzTagRuleResource{}
	email := "9121e724-355c-4b74-9d73-ff118ce7241e@example.com"
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

func TestAccLogzTagRule_update(t *testing.T) {
	t.Skipf("This resource is being deprecated by Azure and will be removed by end of 2025.")

	data := acceptance.BuildTestData(t, "azurerm_logz_tag_rule", "test")
	r := LogzTagRuleResource{}
	email := "41a35aed-12d8-46f3-a2a7-9f89404d7989@example.com"
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

func (r LogzTagRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tagrules.ParseTagRuleID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.Logz.TagRuleClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r LogzTagRuleResource) template(data acceptance.TestData, email string) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, getEffectiveDate(), email)
}

func (r LogzTagRuleResource) basic(data acceptance.TestData, email string) string {
	template := r.template(data, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_tag_rule" "test" {
  logz_monitor_id = azurerm_logz_monitor.test.id
}
`, template)
}

func (r LogzTagRuleResource) requiresImport(data acceptance.TestData, email string) string {
	config := r.basic(data, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_tag_rule" "import" {
  logz_monitor_id = azurerm_logz_tag_rule.test.logz_monitor_id
}
`, config)
}

func (r LogzTagRuleResource) complete(data acceptance.TestData, email string) string {
	template := r.template(data, email)
	return fmt.Sprintf(`
%s

resource "azurerm_logz_tag_rule" "test" {
  logz_monitor_id = azurerm_logz_monitor.test.id
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

func getEffectiveDate() string {
	effectiveDate := time.Now().Add(time.Hour * 72)
	year, month, day := effectiveDate.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, effectiveDate.Location()).Format(time.RFC3339)
}
