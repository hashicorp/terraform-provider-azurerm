// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataProtectionBackupPolicyKubernatesClusterTestResource struct{}

func TestAccDataProtectionBackupPolicyKubernatesCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_kubernetes_cluster", "test")
	r := DataProtectionBackupPolicyKubernatesClusterTestResource{}
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

func TestAccDataProtectionBackupPolicyKubernatesCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_kubernetes_cluster", "test")
	r := DataProtectionBackupPolicyKubernatesClusterTestResource{}
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

func TestAccDataProtectionBackupPolicyKubernatesCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_protection_backup_policy_kubernetes_cluster", "test")
	r := DataProtectionBackupPolicyKubernatesClusterTestResource{}
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

func (r DataProtectionBackupPolicyKubernatesClusterTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := backuppolicies.ParseBackupPolicyID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DataProtection.BackupPolicyClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving DataProtection BackupPolicy (%q): %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r DataProtectionBackupPolicyKubernatesClusterTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-dataprotection-%d"
  location = "%s"
}

resource "azurerm_data_protection_backup_vault" "test" {
  name                = "acctest-dbv-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DataProtectionBackupPolicyKubernatesClusterTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_kubernetes_cluster" "test" {
  name                = "acctest-aks-%d"
  resource_group_name = azurerm_resource_group.test.name
  vault_name          = azurerm_data_protection_backup_vault.test.name

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  retention_rule {
    name     = "Daily"
    priority = 25

    life_cycle {
      duration        = "P84D"
      data_store_type = "OperationalStore"
    }

    criteria {
      days_of_week           = ["Thursday"]
      months_of_year         = ["November"]
      weeks_of_month         = ["First"]
      scheduled_backup_times = ["2021-05-23T02:30:00Z"]
    }
  }

  default_retention_rule {
    life_cycle {
      duration        = "P7D"
      data_store_type = "OperationalStore"
    }
  }
}
`, template, data.RandomInteger)
}

func (r DataProtectionBackupPolicyKubernatesClusterTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_kubernetes_cluster" "import" {
  name                = azurerm_data_protection_backup_policy_kubernetes_cluster.test.name
  resource_group_name = azurerm_data_protection_backup_policy_kubernetes_cluster.test.resource_group_name
  vault_name          = azurerm_data_protection_backup_policy_kubernetes_cluster.test.vault_name

  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]

  retention_rule {
    name     = "Daily"
    priority = 25

    life_cycle {
      duration        = "P84D"
      data_store_type = "OperationalStore"
    }

    criteria {
      days_of_week           = ["Thursday"]
      months_of_year         = ["November"]
      scheduled_backup_times = ["2021-05-23T02:30:00Z"]
    }
  }

  default_retention_rule {
    life_cycle {
      duration        = "P7D"
      data_store_type = "OperationalStore"
    }
  }
}
`, config)
}

func (r DataProtectionBackupPolicyKubernatesClusterTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_data_protection_backup_policy_kubernetes_cluster" "test" {
  name                            = "acctest-aks-%d"
  resource_group_name             = azurerm_resource_group.test.name
  vault_name                      = azurerm_data_protection_backup_vault.test.name
  backup_repeating_time_intervals = ["R/2021-05-23T02:30:00+00:00/P1W"]
  time_zone                       = "India Standard Time"

  retention_rule {
    name     = "Daily"
    priority = 25

    life_cycle {
      duration        = "P7D"
      data_store_type = "OperationalStore"
    }

    life_cycle {
      duration        = "P84D"
      data_store_type = "OperationalStore"
    }

    criteria {
      absolute_criteria = "FirstOfDay"
    }
  }

  default_retention_rule {
    life_cycle {
      duration        = "P7D"
      data_store_type = "OperationalStore"
    }
  }
}
`, template, data.RandomInteger)
}
