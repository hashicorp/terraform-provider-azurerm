package newrelic_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/newrelic/2024-03-01/monitoredsubscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/newrelic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NewRelicMonitoredSubscriptionResource struct{}

func TestAccNewRelicMonitoredSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_monitored_subscription", "test")
	r := NewRelicMonitoredSubscriptionResource{}
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

func TestAccNewRelicMonitoredSubscription_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_monitored_subscription", "test")
	r := NewRelicMonitoredSubscriptionResource{}
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

func TestAccNewRelicMonitoredSubscription_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_monitored_subscription", "test")
	r := NewRelicMonitoredSubscriptionResource{}
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

func TestAccNewRelicMonitoredSubscription_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_new_relic_monitored_subscription", "test")
	r := NewRelicMonitoredSubscriptionResource{}
	email := "f0ff47c3-3aed-45b0-b239-260d9625045a@example.com"
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, email),
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

func (r NewRelicMonitoredSubscriptionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NewRelicMonitoredSubscriptionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.NewRelic.MonitoredSubscriptionsClient
	monitorId := monitoredsubscriptions.NewMonitorID(id.SubscriptionId, id.ResourceGroup, id.MonitorName)
	resp, err := client.Get(ctx, monitorId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r NewRelicMonitoredSubscriptionResource) basic(data acceptance.TestData, email string) string {
	template := r.template(data, email)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_new_relic_monitored_subscription" "test" {
  monitor_id = azurerm_new_relic_monitor.test.id
}
`, template)
}

func (r NewRelicMonitoredSubscriptionResource) requiresImport(data acceptance.TestData, email string) string {
	config := r.basic(data, email)
	return fmt.Sprintf(`
%s

resource "azurerm_new_relic_monitored_subscription" "import" {
  monitor_id = azurerm_new_relic_monitored_subscription.test.monitor_id
}
`, config)
}

func (r NewRelicMonitoredSubscriptionResource) update(data acceptance.TestData, email string) string {
	template := r.template(data, email)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_subscription" "test" {
  subscription_id = "%s"
}

resource "azurerm_new_relic_monitored_subscription" "test" {
  monitor_id = azurerm_new_relic_monitor.test.id

  monitored_subscription {
    subscription_id = data.azurerm_subscription.test.subscription_id

    log_tag_filter {
      name   = "log2"
      action = "Include"
      value  = ""
    }

    metric_tag_filter {
      name   = "metric3"
      action = "Include"
      value  = ""
    }
  }
}
`, template, data.Subscriptions.Secondary)
}

func (r NewRelicMonitoredSubscriptionResource) complete(data acceptance.TestData, email string) string {
	template := r.template(data, email)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_subscription" "test" {
  subscription_id = "%s"
}

resource "azurerm_new_relic_monitored_subscription" "test" {
  monitor_id = azurerm_new_relic_monitor.test.id
  monitored_subscription {
    subscription_id                    = data.azurerm_subscription.test.subscription_id
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
}
`, template, data.Subscriptions.Secondary)
}

func (r NewRelicMonitoredSubscriptionResource) template(data acceptance.TestData, email string) string {
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
