// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageAccountQueuePropertiesResource struct{}

func TestAccStorageAccountQueueProperties_queueProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	r := StorageAccountQueuePropertiesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.queuePropertiesLoggingOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.queuePropertiesMinuteMetricsOnly(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.queueProperties(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.queuePropertiesUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountQueueProperties_StorageV1_queuePropertiesLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	r := StorageAccountQueuePropertiesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageV1QueueProperties(data, "LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountQueueProperties_StorageV1_queuePropertiesGRS(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skipped as `queue_properties` have been moved to the `azurerm_storage_account_queue_properties` resource in v4.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	r := StorageAccountQueuePropertiesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageV1QueueProperties(data, "GRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccStorageAccountQueueProperties_StorageV1_queuePropertiesRAGRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_account", "test")
	r := StorageAccountQueuePropertiesResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageV1QueueProperties(data, "RAGRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r StorageAccountQueuePropertiesResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseStorageAccountID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Storage.ResourceManager.StorageAccounts.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r StorageAccountQueuePropertiesResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseStorageAccountID(state.ID)
	if err != nil {
		return nil, err
	}
	if _, err := client.Storage.ResourceManager.StorageAccounts.Delete(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r StorageAccountQueuePropertiesResource) queueProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

  queue_properties {
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
      enabled               = false
      retention_policy_days = 7
    }

    minute_metrics {
      version               = "1.0"
      enabled               = false
      retention_policy_days = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageAccountQueuePropertiesResource) queuePropertiesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

  queue_properties {
    cors_rule {
      allowed_origins    = ["http://www.example.com"]
      exposed_headers    = ["x-tempo-*", "x-method-*"]
      allowed_headers    = ["*"]
      allowed_methods    = ["GET"]
      max_age_in_seconds = "2000000000"
    }
    cors_rule {
      allowed_origins    = ["http://www.test.com"]
      exposed_headers    = ["x-tempo-*"]
      allowed_headers    = ["*"]
      allowed_methods    = ["PUT"]
      max_age_in_seconds = "1000"
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
      enabled               = true
      retention_policy_days = 7
      include_apis          = true
    }

    minute_metrics {
      version               = "1.0"
      enabled               = true
      include_apis          = false
      retention_policy_days = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

// TODO: Remove in 4.0
func (r StorageAccountQueuePropertiesResource) queuePropertiesLoggingOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

  queue_properties {
    logging {
      version               = "1.0"
      delete                = true
      read                  = true
      write                 = true
      retention_policy_days = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageAccountQueuePropertiesResource) queuePropertiesMinuteMetricsOnly(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

  queue_properties {
    minute_metrics {
      version               = "1.0"
      enabled               = false
      retention_policy_days = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r StorageAccountQueuePropertiesResource) storageV1QueueProperties(data acceptance.TestData, repl string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-storage-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_kind             = "Storage"
  account_tier             = "Standard"
  account_replication_type = "%s"

  queue_properties {
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
      enabled               = false
      retention_policy_days = 7
    }

    minute_metrics {
      version               = "1.0"
      enabled               = false
      retention_policy_days = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, repl)
}
