package synapse_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkspacePackageResource struct {
}

func TestAccSynapseWorkspacePackage_basic(t *testing.T) {
	sourceBlob, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}

	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace_package", "test")
	r := WorkspacePackageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, sourceBlob.Name()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("source", "source_md5"),
	})
}

func TestAccSynapseWorkspacePackage_requiresImport(t *testing.T) {
	sourceBlob, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create local source blob file")
	}

	if err := populateTempFile(sourceBlob); err != nil {
		t.Fatalf("Error populating temp file: %s", err)
	}
	data := acceptance.BuildTestData(t, "azurerm_synapse_workspace_package", "test")
	r := WorkspacePackageResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, sourceBlob.Name()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, sourceBlob.Name()),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func (r WorkspacePackageResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WorkspacePackageID(state.ID)
	if err != nil {
		return nil, err
	}

	environment := clients.Account.Environment
	client, err := clients.Synapse.LibraryClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, id.LibraryName)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r WorkspacePackageResource) basic(data acceptance.TestData, filename string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_workspace_package" "test" {
  name                 = "acctest%d.jar"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  source               = "%s"
  source_md5           = "${filemd5("%s")}"
  depends_on = [
    azurerm_synapse_firewall_rule.test,
  ]
}
`, r.template(data), data.RandomInteger, filename, filename)
}

func (r WorkspacePackageResource) requiresImport(data acceptance.TestData, filename string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_synapse_workspace_package" "import" {
  name                 = azurerm_synapse_workspace_package.test.name
  synapse_workspace_id = azurerm_synapse_workspace_package.test.synapse_workspace_id
  source               = azurerm_synapse_workspace_package.test.source
  source_md5           = azurerm_synapse_workspace_package.test.source_md5
}
`, r.basic(data, filename))
}

func (WorkspacePackageResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-synapse-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
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
  managed_virtual_network_enabled      = true
}

resource "azurerm_synapse_firewall_rule" "test" {
  name                 = "allowAll"
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func populateTempFile(input *os.File) error {
	if err := input.Truncate(25*1024*1024 + 512); err != nil {
		return fmt.Errorf("Failed to truncate file to 25M")
	}

	for i := int64(0); i < 20; i += 2 {
		randomBytes := make([]byte, 1*1024*1024)
		if _, err := rand.Read(randomBytes); err != nil {
			return fmt.Errorf("Failed to read random bytes")
		}

		if _, err := input.WriteAt(randomBytes, i*1024*1024); err != nil {
			return fmt.Errorf("Failed to write random bytes to file")
		}
	}

	randomBytes := make([]byte, 5*1024*1024)
	if _, err := rand.Read(randomBytes); err != nil {
		return fmt.Errorf("Failed to read random bytes")
	}

	if _, err := input.WriteAt(randomBytes, 20*1024*1024); err != nil {
		return fmt.Errorf("Failed to write random bytes to file")
	}

	if err := input.Close(); err != nil {
		return fmt.Errorf("Failed to close source blob")
	}

	return nil
}
