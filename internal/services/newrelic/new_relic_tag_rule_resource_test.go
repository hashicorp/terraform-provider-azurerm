// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package newrelic_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NewRelicTagRuleResource struct{}

func TestAccNewRelicTagRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_tag_rule", "test")
	r := NewRelicTagRuleResource{}
	email := "27362230-e2d8-4c73-9ee3-fdef83459ca3@example.com"
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

func TestAccNewRelicTagRule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_tag_rule", "test")
	r := NewRelicTagRuleResource{}
	email := "85b5febd-127d-4633-9c25-bcfea555af46@example.com"
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

func TestAccNewRelicTagRule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_tag_rule", "test")
	r := NewRelicTagRuleResource{}
	email := "672d9312-65a7-484c-870d-94584850a423@example.com"
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

func TestAccNewRelicTagRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_tag_rule", "test")
	r := NewRelicTagRuleResource{}
	email := "f0ff47c3-3aed-45b0-b239-260d9625045a@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data, email),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r NewRelicTagRuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tagrules.ParseTagRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.NewRelic.TagRulesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r NewRelicTagRuleResource) template(data acceptance.TestData, email string) string {
	year, month, day := time.Now().Add(time.Hour * 72).Date()
	effectiveDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_new_relic_monitor" "test" {
  name                = "acctest-nrm-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  plan {
    effective_date = "%[3]s"
  }

  user {
    email        = "%[4]s"
    first_name   = "first"
    last_name    = "last"
    phone_number = "123456"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, effectiveDate, email)
}

func (r NewRelicTagRuleResource) basic(data acceptance.TestData, email string) string {
	template := r.template(data, email)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
				%s

resource "azurerm_new_relic_tag_rule" "test" {
  monitor_id = azurerm_new_relic_monitor.test.id
}
`, template)
}

func (r NewRelicTagRuleResource) requiresImport(data acceptance.TestData, email string) string {
	config := r.basic(data, email)
	return fmt.Sprintf(`
			%s

resource "azurerm_new_relic_tag_rule" "import" {
  monitor_id = azurerm_new_relic_tag_rule.test.monitor_id
}
`, config)
}

func (r NewRelicTagRuleResource) complete(data acceptance.TestData, email string) string {
	template := r.template(data, email)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
			%s

resource "azurerm_new_relic_tag_rule" "test" {
  monitor_id                         = azurerm_new_relic_monitor.test.id
  azure_active_directory_log_enabled = true
  activity_log_enabled               = true
  metric_enabled                     = true
  subscription_log_enabled           = true

  log_tag_filter {
    name   = "log1"
    action = "Include"
    value  = "log1"
  }

  log_tag_filter {
    name   = "log2"
    action = "Exclude"
    value  = ""
  }

  metric_tag_filter {
    name   = "metric1"
    action = "Include"
    value  = "metric1"
  }

  metric_tag_filter {
    name   = "metric2"
    action = "Exclude"
    value  = ""
  }
}
`, template)
}

func (r NewRelicTagRuleResource) update(data acceptance.TestData, email string) string {
	template := r.template(data, email)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
			%s

resource "azurerm_new_relic_tag_rule" "test" {
  monitor_id                         = azurerm_new_relic_monitor.test.id
  azure_active_directory_log_enabled = false
  activity_log_enabled               = false
  metric_enabled                     = false
  subscription_log_enabled           = false

  log_tag_filter {
    name   = "log2"
    action = "Exclude"
    value  = ""
  }

  log_tag_filter {
    name   = "log1"
    action = "Include"
    value  = "log1"
  }

  metric_tag_filter {
    name   = "metric1"
    action = "Exclude"
    value  = ""
  }

  metric_tag_filter {
    name   = "metric2"
    action = "Include"
    value  = "metric1"
  }
}
`, template)
}
