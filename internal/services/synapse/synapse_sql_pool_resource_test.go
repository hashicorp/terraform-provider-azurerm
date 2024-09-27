// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SynapseSqlPoolResource struct{}

func TestAccSynapseSqlPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoBackupDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_policy_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("GRS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseSqlPool_basicThreePointOh(t *testing.T) {
	// NOTE: Validate that the original resources default values during create are preserved...
	if features.FourPointOhBeta() {
		t.Skipf("Skippped as 'storage_account_type' is now a Required field in 4.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoBackupThreePointOhDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_policy_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("GRS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseSqlPool_threePointOhUpdate(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("Skippped as 'storage_account_type' is now a Required field in 4.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoBackupThreePointOhDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_policy_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("GRS"),
			),
		},
		data.ImportStep(),
		{
			Config: r.geoBackup(data, false, "GRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_policy_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("GRS"),
			),
		},
		data.ImportStep(),
		{
			Config: r.geoBackupThreePointOhDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_policy_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("GRS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseSqlPool_utf8(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.utf8(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseSqlPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoBackupDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSynapseSqlPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

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

func TestAccSynapseSqlPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoBackupDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_policy_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("GRS"),
			),
		},
		data.ImportStep(),
		{
			Config: r.geoBackup(data, false, "GRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_policy_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("GRS"),
			),
		},
		data.ImportStep(),
		{
			Config: r.geoBackupDefault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_policy_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("GRS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseSqlPool_geoBackupDisabledGRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoBackup(data, false, "GRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_policy_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("GRS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseSqlPool_geoBackupDisabledLRS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoBackup(data, false, "LRS"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_policy_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("LRS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSynapseSqlPool_geoBackupInvalid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_synapse_sql_pool", "test")
	r := SynapseSqlPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.geoBackup(data, true, "LRS"),
			ExpectError: regexp.MustCompile("`geo_backup_policy_enabled` cannot be `true` if the `storage_account_type` is `LRS`"),
		},
	})
}

func (r SynapseSqlPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SqlPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Synapse.SqlPoolClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Synapse Sql Pool %q (Workspace %q / Resource Group %q): %+v", id.Name, id.WorkspaceName, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r SynapseSqlPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  name               = "acctest-%d"
  storage_account_id = azurerm_storage_account.test.id
}

resource "azurerm_synapse_workspace" "test" {
  name                                 = "acctestsw%d"
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func (r SynapseSqlPoolResource) utf8(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_synapse_sql_pool" "test" {
  name                 = "販売管理"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  sku_name             = "DW100c"
  create_mode          = "Default"
  storage_account_type = "GRS"
}
`, template)
}

func (r SynapseSqlPoolResource) requiresImport(data acceptance.TestData) string {
	config := r.geoBackupDefault(data)
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_sql_pool" "import" {
  name                 = azurerm_synapse_sql_pool.test.name
  synapse_workspace_id = azurerm_synapse_sql_pool.test.synapse_workspace_id
  sku_name             = azurerm_synapse_sql_pool.test.sku_name
  create_mode          = azurerm_synapse_sql_pool.test.create_mode
  storage_account_type = "GRS"
}
`, config)
}

func (r SynapseSqlPoolResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_synapse_sql_pool" "test" {
  name                      = "acctestSP%s"
  synapse_workspace_id      = azurerm_synapse_workspace.test.id
  sku_name                  = "DW500c"
  create_mode               = "Default"
  collation                 = "SQL_Latin1_General_CP1_CI_AS"
  data_encrypted            = true
  geo_backup_policy_enabled = true
  storage_account_type      = "GRS"

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomString)
}

func (r SynapseSqlPoolResource) geoBackupThreePointOhDefault(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_synapse_sql_pool" "test" {
  name                 = "acctestSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  sku_name             = "DW100c"
  create_mode          = "Default"
}
`, template, data.RandomString)
}

func (r SynapseSqlPoolResource) geoBackupDefault(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_synapse_sql_pool" "test" {
  name                 = "acctestSP%s"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  sku_name             = "DW100c"
  create_mode          = "Default"
  storage_account_type = "GRS"
}
`, template, data.RandomString)
}

func (r SynapseSqlPoolResource) geoBackup(data acceptance.TestData, geoBackupPolicy bool, accountType string) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_synapse_sql_pool" "test" {
  name                      = "acctestSP%s"
  synapse_workspace_id      = azurerm_synapse_workspace.test.id
  sku_name                  = "DW100c"
  create_mode               = "Default"
  geo_backup_policy_enabled = %t
  storage_account_type      = "%s"
}
`, template, data.RandomString, geoBackupPolicy, accountType)
}
