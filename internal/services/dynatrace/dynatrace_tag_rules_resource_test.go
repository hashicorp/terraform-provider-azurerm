// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynatrace_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type TagRulesResource struct {
	dynatraceInfo dynatraceInfo
}

func NewTagRulesResource() TagRulesResource {
	return TagRulesResource{
		dynatraceInfo: dynatraceInfo{
			UserCountry:     os.Getenv("DYNATRACE_USER_COUNTRY"),
			UserEmail:       os.Getenv("DYNATRACE_USER_EMAIL"),
			UserFirstName:   os.Getenv("DYNATRACE_USER_FIRST_NAME"),
			UserLastName:    os.Getenv("DYNATRACE_USER_LAST_NAME"),
			UserPhoneNumber: os.Getenv("DYNATRACE_USER_PHONE_NUMBER"),
		},
	}
}

func TestAccDynatraceTagRules_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dynatrace_tag_rules", "test")
	r := NewTagRulesResource()
	r.preCheck(t)

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

func TestAccDynatraceTagRules_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dynatrace_tag_rules", "test")
	r := NewTagRulesResource()
	r.preCheck(t)

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

func TestAccDynatraceTagRules_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dynatrace_tag_rules", "test")
	r := NewTagRulesResource()
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDynatraceTagRules_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dynatrace_tag_rules", "test")
	r := NewTagRulesResource()
	r.preCheck(t)

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

func (r TagRulesResource) preCheck(t *testing.T) {
	if r.dynatraceInfo.UserCountry == "" {
		t.Skipf("DYNATRACE_USER_COUNTRY must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserEmail == "" {
		t.Skipf("DYNATRACE_USER_EMAIL must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserFirstName == "" {
		t.Skipf("DYNATRACE_USER_FIRST_NAME must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserLastName == "" {
		t.Skipf("DYNATRACE_USER_LAST_NAME must be set for acceptance tests")
	}
	if r.dynatraceInfo.UserPhoneNumber == "" {
		t.Skipf("DYNATRACE_USER_PHONE_NUMBER must be set for acceptance tests")
	}
}

func (r TagRulesResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tagrules.ParseTagRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Dynatrace.TagRulesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(true), nil
}

func (r TagRulesResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_dynatrace_tag_rules" "test" {
  name       = "default"
  monitor_id = azurerm_dynatrace_monitor.test.id
}
`, MonitorsResource{}.basic(data))
}

func (r TagRulesResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_dynatrace_tag_rules" "test" {
  name       = "default"
  monitor_id = azurerm_dynatrace_monitor.test.id

  log_rule {
    filtering_tag {
      name   = "Environment"
      value  = "Prod"
      action = "Include"
    }
    send_azure_active_directory_logs_enabled = true
    send_activity_logs_enabled               = true
    send_subscription_logs_enabled           = true
  }

  metric_rule {
    filtering_tag {
      name   = "Environment"
      value  = "Prod"
      action = "Include"
    }
  }
}
`, MonitorsResource{}.basic(data))
}

func (r TagRulesResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dynatrace_tag_rules" "test" {
  name       = "default"
  monitor_id = azurerm_dynatrace_monitor.test.id

  log_rule {
    filtering_tag {
      name   = "Foo"
      value  = "Bar"
      action = "Exclude"
    }
    send_azure_active_directory_logs_enabled = false
    send_activity_logs_enabled               = false
    send_subscription_logs_enabled           = false
  }

  metric_rule {
    filtering_tag {
      name   = "Foo"
      value  = "Bar"
      action = "Exclude"
    }
  }
}
`, MonitorsResource{}.basic(data))
}

func (r TagRulesResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dynatrace_tag_rules" "import" {
  name       = azurerm_dynatrace_tag_rules.test.name
  monitor_id = azurerm_dynatrace_tag_rules.test.monitor_id
}
`, template)
}
