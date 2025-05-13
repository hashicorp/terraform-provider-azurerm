// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AccountQueuePropertiesResource struct{}

func TestAccStorageAccountQueueProperties_corsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_queue_properties", "test")
	r := AccountQueuePropertiesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.corsOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountQueueProperties_loggingOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_queue_properties", "test")
	r := AccountQueuePropertiesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.loggingOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountQueueProperties_hourMetricsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_queue_properties", "test")
	r := AccountQueuePropertiesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hourMetricsOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountQueueProperties_minuteMetricsOnly(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_queue_properties", "test")
	r := AccountQueuePropertiesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.minuteMetricsOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountQueueProperties_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_queue_properties", "test")
	r := AccountQueuePropertiesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.corsOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.loggingOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.hourMetricsOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.minuteMetricsOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.corsOnly(data),
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

func TestAccStorageAccountQueueProperties_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account_queue_properties", "test")
	r := AccountQueuePropertiesResource{}

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

func (r AccountQueuePropertiesResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseStorageAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	account, err := client.Storage.GetAccount(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if account == nil {
		return nil, fmt.Errorf("unable to locate %s", *id)
	}

	queuesClient, err := client.Storage.QueuesDataPlaneClient(ctx, *account, client.Storage.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return nil, fmt.Errorf("building Queues Client for %s: %v", *id, err)
	}

	props, err := queuesClient.GetServiceProperties(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrieving Queue Properties for %s: %+v", *id, err)
	}

	present := !reflect.DeepEqual(storage.DefaultValueForAccountQueueProperties(), props)
	return pointer.To(present), nil
}

func (r AccountQueuePropertiesResource) corsOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_account_queue_properties" "test" {
  storage_account_id = azurerm_storage_account.test.id
  cors_rule {
    allowed_origins    = ["http://www.example.com"]
    exposed_headers    = ["x-tempo-*"]
    allowed_headers    = ["x-tempo-*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = "500"
  }

  cors_rule {
    allowed_origins    = ["http://www.contoso.com"]
    exposed_headers    = ["x-example-*"]
    allowed_headers    = ["x-example-*"]
    allowed_methods    = ["GET"]
    max_age_in_seconds = "60"
  }
}
`, r.template(data))
}

func (r AccountQueuePropertiesResource) hourMetricsOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_account_queue_properties" "test" {
  storage_account_id = azurerm_storage_account.test.id
  hour_metrics {
    version               = "1.0"
    retention_policy_days = 7
  }
}
`, r.template(data))
}

func (r AccountQueuePropertiesResource) minuteMetricsOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_account_queue_properties" "test" {
  storage_account_id = azurerm_storage_account.test.id
  minute_metrics {
    version               = "1.0"
    retention_policy_days = 7
  }
}
`, r.template(data))
}

func (r AccountQueuePropertiesResource) loggingOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_account_queue_properties" "test" {
  storage_account_id = azurerm_storage_account.test.id
  logging {
    version               = "1.0"
    delete                = true
    read                  = true
    write                 = true
    retention_policy_days = 7
  }
}
`, r.template(data))
}

func (r AccountQueuePropertiesResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_account_queue_properties" "test" {
  storage_account_id = azurerm_storage_account.test.id
  cors_rule {
    allowed_origins    = ["http://www.example.com"]
    exposed_headers    = ["x-tempo-*"]
    allowed_headers    = ["x-tempo-*"]
    allowed_methods    = ["GET", "PUT"]
    max_age_in_seconds = "500"
  }

  logging {
    version               = "1.0"
    delete                = true
    read                  = true
    write                 = true
    retention_policy_days = 7
  }

  hour_metrics {
    version               = "1.0"
    retention_policy_days = 7
  }

  minute_metrics {
    version               = "1.0"
    retention_policy_days = 7
  }
}
`, r.template(data))
}

func (r AccountQueuePropertiesResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
