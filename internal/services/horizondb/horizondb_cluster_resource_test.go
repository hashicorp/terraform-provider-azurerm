// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package horizondb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/clusters"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HorizondbClusterResource struct{}

func TestAccHorizondbCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccHorizondbCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccHorizondbCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccHorizondbCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updated(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.updatedWithWriteOnlyPassword(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password", "administrator_password_wo_version"),
	})
}

func TestAccHorizondbCluster_geoRedundantBackup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoRedundantBackup(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccHorizondbCluster_restore(t *testing.T) {
	r := HorizondbClusterResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	restoreVariables := config.Variables{}
	// Azure can omit creation-time restore metadata from GET responses.
	restoreImport := data.ImportStepFor("azurerm_horizondb_cluster.restored", "administrator_password", "point_in_time_restore")
	restoreImport.ConfigVariables = restoreVariables
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
		{
			PreConfig: func() {
				restoreVariables["point_in_time_utc"] = config.StringVariable(time.Now().UTC().Format(time.RFC3339))
				time.Sleep(time.Minute)
			},
			Config:          r.restore(data),
			ConfigVariables: restoreVariables,
			Check:           check.That("azurerm_horizondb_cluster.restored").ExistsInAzure(r),
		},
		restoreImport,
	})
}

func TestAccHorizondbCluster_restoreCrossRegion(t *testing.T) {
	r := HorizondbClusterResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	if data.Locations.Primary == data.Locations.Secondary {
		t.Skip("cross-region restore requires distinct primary and secondary locations with HorizonDB support")
	}
	restoreVariables := config.Variables{}
	// Azure can omit creation-time restore metadata from GET responses.
	restoreImport := data.ImportStepFor("azurerm_horizondb_cluster.restored", "administrator_password", "point_in_time_restore")
	restoreImport.ConfigVariables = restoreVariables
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoRedundantBackup(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
		{
			PreConfig: func() {
				restoreVariables["point_in_time_utc"] = config.StringVariable(time.Now().UTC().Format(time.RFC3339))
				time.Sleep(time.Minute)
			},
			Config:          r.restoreCrossRegion(data),
			ConfigVariables: restoreVariables,
			Check:           check.That("azurerm_horizondb_cluster.restored").ExistsInAzure(r),
		},
		restoreImport,
	})
}

func TestAccHorizondbCluster_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.withoutIdentity(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.withTwoUserAssignedIdentities(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.withOneUserAssignedIdentity(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.withoutIdentity(data),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
	})
}

func (HorizondbClusterResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := clusters.ParseClusterID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.HorizonDB.ClustersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (HorizondbClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-horizondb-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r HorizondbClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_cluster" "test" {
  name                   = "acctest-horizondb-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  postgresql_version     = "18"
  administrator_login    = "acctestadmin"
  administrator_password = "P@ssword-123!"

  read_write_pool {
    name                   = "primary"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HorizondbClusterResource) updatedWithWriteOnlyPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

%s

resource "azurerm_horizondb_cluster" "test" {
  name                              = "acctest-horizondb-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  postgresql_version                = "18"
  administrator_login               = "acctestadmin"
  administrator_password_wo         = ephemeral.azurerm_key_vault_secret.test.value
  administrator_password_wo_version = 1
  backup_retention_in_days          = 14
  public_network_access             = "Enabled"

  read_write_pool {
    name                   = "primary"
    node_count             = 2
    high_availability_mode = "AnyZone"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 8
    }
  }

  tags = {
    environment = "acceptance"
    retention   = "14"
  }
}
`, r.template(data), string(hclwrite.Format([]byte(acceptance.WriteOnlyKeyVaultSecretTemplate(data, "P@ssword-789!")))), data.RandomInteger)
}

func (r HorizondbClusterResource) updated(data acceptance.TestData) string {
	// Updated configuration for backup retention, backup retention, high availability, multi-node pool and password
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_cluster" "test" {
  name                     = "acctest-horizondb-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  postgresql_version       = "18"
  administrator_login      = "acctestadmin"
  administrator_password   = "P@ssword-456!"
  backup_retention_in_days = 14
  public_network_access    = "Enabled"

  read_write_pool {
    name                   = "primary"
    node_count             = 2
    high_availability_mode = "AnyZone"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 8
    }
  }

  tags = {
    environment = "acceptance"
    retention   = "14"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HorizondbClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_cluster" "test" {
  name                     = "acctest-horizondb-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  postgresql_version       = "18"
  administrator_login      = "acctestadmin"
  administrator_password   = "P@ssword-456!"
  backup_retention_in_days = 14
  public_network_access    = "Enabled"

  read_write_pool {
    name                   = "primary"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }

  tags = {
    environment = "acceptance"
    retention   = "14"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HorizondbClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_horizondb_cluster" "import" {
  name                   = azurerm_horizondb_cluster.test.name
  resource_group_name    = azurerm_horizondb_cluster.test.resource_group_name
  location               = azurerm_horizondb_cluster.test.location
  postgresql_version     = azurerm_horizondb_cluster.test.postgresql_version
  administrator_login    = azurerm_horizondb_cluster.test.administrator_login
  administrator_password = azurerm_horizondb_cluster.test.administrator_password

  read_write_pool {
    name                   = "primary"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }
}
`, r.basic(data))
}

func (r HorizondbClusterResource) geoRedundantBackup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_cluster" "test" {
  name                         = "acctest-horizondb-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  postgresql_version           = "18"
  administrator_login          = "acctestadmin"
  administrator_password       = "P@ssword-123!"
  geo_redundant_backup_enabled = true
  backup_retention_in_days     = 14

  read_write_pool {
    name                   = "primary"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r HorizondbClusterResource) restore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

variable "point_in_time_utc" {
  type = string
}

resource "azurerm_horizondb_cluster" "restored" {
  name                     = "acctest-horizondb-restored-%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  administrator_login      = "acctestadmin"
  administrator_password   = "P@ssword-123!"
  backup_retention_in_days = 14

  point_in_time_restore {
    source_horizondb_cluster_id = azurerm_horizondb_cluster.test.id
    point_in_time_utc           = var.point_in_time_utc
  }

  read_write_pool {
    name                   = "primary"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }
}
`, r.basic(data), data.RandomInteger)
}

func (r HorizondbClusterResource) restoreCrossRegion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

variable "point_in_time_utc" {
  type = string
}

resource "azurerm_horizondb_cluster" "restored" {
  name                   = "acctest-horizondb-restored-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = %q
  administrator_login    = "acctestadmin"
  administrator_password = "P@ssword-123!"

  point_in_time_restore {
    source_horizondb_cluster_id = azurerm_horizondb_cluster.test.id
    point_in_time_utc           = var.point_in_time_utc
    source_location             = azurerm_horizondb_cluster.test.location
  }

  read_write_pool {
    name                   = "primary"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }
}
`, r.geoRedundantBackup(data), data.RandomInteger, data.Locations.Secondary)
}

func (r HorizondbClusterResource) identityTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  count               = 2
  name                = "acctest-horizondb-%d-${count.index}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r HorizondbClusterResource) withoutIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_cluster" "test" {
  name                   = "acctest-horizondb-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  postgresql_version     = "18"
  administrator_login    = "acctestadmin"
  administrator_password = "P@ssword-123!"

  read_write_pool {
    name                   = "primary"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }

  tags = {
    identity_count = "0"
  }
}
`, r.identityTemplate(data), data.RandomInteger)
}

func (r HorizondbClusterResource) withOneUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_cluster" "test" {
  name                   = "acctest-horizondb-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  postgresql_version     = "18"
  administrator_login    = "acctestadmin"
  administrator_password = "P@ssword-123!"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test[0].id]
  }

  read_write_pool {
    name                   = "primary"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }

  tags = {
    identity_count = "1"
  }
}
`, r.identityTemplate(data), data.RandomInteger)
}

func (r HorizondbClusterResource) withTwoUserAssignedIdentities(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_cluster" "test" {
  name                   = "acctest-horizondb-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  postgresql_version     = "18"
  administrator_login    = "acctestadmin"
  administrator_password = "P@ssword-123!"

  identity {
    type         = "UserAssigned"
    identity_ids = azurerm_user_assigned_identity.test[*].id
  }

  read_write_pool {
    name                   = "primary"
    node_count             = 1
    high_availability_mode = "Disabled"

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = 4
    }
  }

  tags = {
    identity_count = "2"
  }
}
`, r.identityTemplate(data), data.RandomInteger)
}
