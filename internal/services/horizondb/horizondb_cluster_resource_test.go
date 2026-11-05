// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package horizondb_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/horizondb/2026-10-01-preview/clusters"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HorizondbClusterResource struct{}

func TestAccHorizondbCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	importPlan := data.ImportBlockWithIDStep(true)
	importPlan.ImportPlanChecks.PreApply = []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionUpdate),
	}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep("administrator_password"),
		importPlan,
	})
}

func TestAccHorizondbCluster_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.complete(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep("administrator_password"),
	})
}

func TestAccHorizondbCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccHorizondbCluster_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{Config: r.settings(data, 7, false), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep("administrator_password"),
		{Config: r.complete(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep("administrator_password"),
		{Config: r.resetSettings(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		data.ImportStep("administrator_password"),
	})
}

func TestAccHorizondbCluster_authentication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withAdministratorPassword(data, "P@ssword-123!"),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.withAdministratorPassword(data, "P@ssword-456!"),
			Check:  check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
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

func TestAccHorizondbCluster_writeOnlyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.writeOnlyPassword(data, "P@ssword-123!", 1),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					resource.TestCheckNoResourceAttr(data.ResourceName, "administrator_password_wo"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_password", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_password_wo_version", "1"),
				),
			},
			data.ImportStep("administrator_password_wo_version"),
			{
				Config: r.writeOnlyPassword(data, "P@ssword-456!", 2),
				Check: resource.ComposeTestCheckFunc(
					check.That(data.ResourceName).ExistsInAzure(r),
					resource.TestCheckNoResourceAttr(data.ResourceName, "administrator_password_wo"),
					resource.TestCheckResourceAttr(data.ResourceName, "administrator_password_wo_version", "2"),
				),
			},
			data.ImportStep("administrator_password_wo_version"),
		},
	})
}

func TestAccHorizondbCluster_restore(t *testing.T) {
	r := HorizondbClusterResource{}
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	restoreVariables := config.Variables{}
	// Azure can omit creation-time restore metadata from GET responses.
	restoreImport := data.ImportStepFor("azurerm_horizondb_cluster.restored", "administrator_password", "restore")
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
	restoreImport := data.ImportStepFor("azurerm_horizondb_cluster.restored", "administrator_password", "restore")
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
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionUpdate),
				},
			},
			Check: check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.withOneUserAssignedIdentity(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionUpdate),
				},
			},
			Check: check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
		{
			Config: r.withoutIdentity(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionUpdate),
				},
			},
			Check: check.That(data.ResourceName).ExistsInAzure(r),
		},
		data.ImportStep("administrator_password"),
	})
}

func TestAccHorizondbCluster_poolRenameRejected(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{Config: r.basic(data), Check: check.That(data.ResourceName).ExistsInAzure(r)},
		{
			Config:      r.withPool(data, "renamed", 1, "Disabled", 4),
			ExpectError: regexp.MustCompile("changing `read_write_pool.name` in place is not supported"),
		},
	})
}

func TestAccHorizondbCluster_poolInvalidHighAvailability(t *testing.T) {
	testCases := map[string]struct {
		nodes int
		mode  string
	}{
		"singleNode":    {nodes: 1, mode: "ZoneRedundant"},
		"multipleNodes": {nodes: 2, mode: "Disabled"},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
			r := HorizondbClusterResource{}
			data.ResourceTest(t, r, []acceptance.TestStep{{
				Config:      r.withPool(data, "primary", testCase.nodes, testCase.mode, 4),
				ExpectError: regexp.MustCompile("`read_write_pool.high_availability_mode` must be `Disabled`"),
			}})
		})
	}
}

func TestAccHorizondbCluster_poolNameCollision(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_horizondb_cluster", "test")
	r := HorizondbClusterResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{{
		Config:      r.withPool(data, fmt.Sprintf("acctest-horizondb-%d", data.RandomInteger), 1, "Disabled", 4),
		ExpectError: regexp.MustCompile("`read_write_pool.name` must differ"),
	}})
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

func (r HorizondbClusterResource) template(data acceptance.TestData) string {
	return r.templateWithFeatures(data, "")
}

func (HorizondbClusterResource) templateWithFeatures(data acceptance.TestData, features string) string {
	featuresConfig := "features {}"
	if features != "" {
		featuresConfig = fmt.Sprintf("features {\n    %s\n  }", features)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  %s
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-horizondb-%d"
  location = "%s"
}
`, featuresConfig, data.RandomInteger, data.Locations.Primary)
}

func (r HorizondbClusterResource) basic(data acceptance.TestData) string {
	return r.withPool(data, "primary", 1, "Disabled", 4)
}

func (HorizondbClusterResource) pool() string {
	lines := strings.Split(strings.TrimSpace(`
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
`), "\n")
	for i := 1; i < len(lines); i++ {
		if lines[i] != "" {
			lines[i] = "  " + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}

func (r HorizondbClusterResource) createModeValidation(data acceptance.TestData, version, mode string, restore bool) string {
	versionConfig := ""
	if version != "" {
		versionConfig = fmt.Sprintf("  postgresql_version     = %q", version)
	}
	restoreConfig := ""
	if restore {
		restoreConfig = `
  restore {
    source_horizondb_cluster_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.HorizonDb/clusters/source"
    point_in_time_utc           = "2026-10-01T00:00:00Z"
  }
`
	}
	return fmt.Sprintf(`
%s
resource "azurerm_horizondb_cluster" "test" {
  name                   = "acctest-horizondb-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "acctestadmin"
  administrator_password = "P@ssword-123!"
  create_mode            = %q
%s
%s
  %s
}
`, r.template(data), data.RandomInteger, mode, versionConfig, restoreConfig, r.pool())
}

func (r HorizondbClusterResource) withAdministratorPassword(data acceptance.TestData, password string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_horizondb_cluster" "test" {
  name                   = "acctest-horizondb-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  postgresql_version     = "18"
  administrator_login    = "acctestadmin"
  administrator_password = %q

  %s
}
`, r.template(data), data.RandomInteger, password, r.pool())
}

func (r HorizondbClusterResource) withPool(data acceptance.TestData, name string, nodes int, mode string, cores int) string {
	return r.withPoolAndFeatures(data, name, nodes, mode, cores, "")
}

func (r HorizondbClusterResource) writeOnlyPassword(data acceptance.TestData, secret string, version int) string {
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
  administrator_password_wo_version = %d

  %s
}
`, r.template(data), string(hclwrite.Format([]byte(acceptance.WriteOnlyKeyVaultSecretTemplate(data, secret)))), data.RandomInteger, version, r.pool())
}

func (r HorizondbClusterResource) withPoolAndFeatures(data acceptance.TestData, name string, nodes int, mode string, cores int, features string) string {
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
    name                   = %q
    node_count             = %d
    high_availability_mode = %q

    compute {
      processor  = "AMD"
      generation = 1
      vcores     = %d
    }
  }

}
`, r.templateWithFeatures(data, features), data.RandomInteger, name, nodes, mode, cores)
}

func (r HorizondbClusterResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_cluster" "test" {
  name                          = "acctest-horizondb-%d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  postgresql_version            = "18"
  administrator_login           = "acctestadmin"
  administrator_password        = "P@ssword-456!"
  backup_retention_days         = 14
  public_network_access_enabled = true

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

func (r HorizondbClusterResource) resetSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_horizondb_cluster" "test" {
  name                          = "acctest-horizondb-%d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  postgresql_version            = "18"
  administrator_login           = "acctestadmin"
  administrator_password        = "P@ssword-123!"
  backup_retention_days         = 7
  public_network_access_enabled = false

  %s
}
`, r.template(data), data.RandomInteger, r.pool())
}

func (r HorizondbClusterResource) settings(data acceptance.TestData, retention int, publicNetworkAccessEnabled bool) string {
	return fmt.Sprintf(`
%s
resource "azurerm_horizondb_cluster" "test" {
  name                          = "acctest-horizondb-%d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  postgresql_version            = "18"
  administrator_login           = "acctestadmin"
  administrator_password        = "P@ssword-123!"
  backup_retention_days         = %d
  public_network_access_enabled = %t

  %s

  tags = {
    environment = "acceptance"
    retention   = "%d"
  }
}
`, r.template(data), data.RandomInteger, retention, publicNetworkAccessEnabled, r.pool(), retention)
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

  %s
}
`, r.basic(data), r.pool())
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

  %s
}
`, r.template(data), data.RandomInteger, r.pool())
}

func (r HorizondbClusterResource) restore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

variable "point_in_time_utc" {
  type = string
}

resource "azurerm_horizondb_cluster" "restored" {
  name                   = "acctest-horizondb-restored-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  create_mode            = "PointInTimeRestore"
  administrator_login    = "acctestadmin"
  administrator_password = "P@ssword-123!"

  restore {
    source_horizondb_cluster_id = azurerm_horizondb_cluster.test.id
    point_in_time_utc           = var.point_in_time_utc
  }

  %s
}
`, r.basic(data), data.RandomInteger, r.pool())
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
  create_mode            = "PointInTimeRestore"
  administrator_login    = "acctestadmin"
  administrator_password = "P@ssword-123!"

  restore {
    source_horizondb_cluster_id = azurerm_horizondb_cluster.test.id
    point_in_time_utc           = var.point_in_time_utc
    source_location             = azurerm_horizondb_cluster.test.location
  }

  %s
}
`, r.geoRedundantBackup(data), data.RandomInteger, data.Locations.Secondary, r.pool())
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

  %s

  tags = {
    identity_count = "0"
  }
}
`, r.identityTemplate(data), data.RandomInteger, r.pool())
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

  %s

  tags = {
    identity_count = "1"
  }
}
`, r.identityTemplate(data), data.RandomInteger, r.pool())
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

  %s

  tags = {
    identity_count = "2"
  }
}
`, r.identityTemplate(data), data.RandomInteger, r.pool())
}
