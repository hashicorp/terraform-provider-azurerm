// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StorageSyncServerEndpointResource struct{}

func TestAccStorageSyncServerEndpointSequential(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"storageSyncServerEndpoint": {
			"basic":            testAccStorageSyncServerEndpoint_basic,
			"complete":         testAccStorageSyncServerEndpoint_complete,
			"listBasic":        testAccStorageSyncServerEndpoint_list_basic,
			"resourceIdentity": testAccStorageSyncServerEndpoint_resourceIdentity,
		},
	})
}

func testAccStorageSyncServerEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_server_endpoint", "test")
	r := StorageSyncServerEndpointResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccStorageSyncServerEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_sync_server_endpoint", "test")
	r := StorageSyncServerEndpointResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func (r StorageSyncServerEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := serverendpointresource.ParseServerEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Storage.SyncServerEndpointsClient.ServerEndpointsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StorageSyncServerEndpointResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_sync_server_endpoint" "test" {
  name                  = "acctestSE-%[2]s"
  storage_sync_group_id = azurerm_storage_sync_group.test.id
  registered_server_id  = trimspace(data.local_file.server_id.content)
  server_local_path     = "D:\\SyncFolder"

  depends_on = [terraform_data.afs_server_id, azurerm_storage_sync_cloud_endpoint.test]
}
`, r.template(data), data.RandomString)
}

func (r StorageSyncServerEndpointResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_storage_sync_server_endpoint" "test" {
  name                  = "acctestSE-%[2]s"
  storage_sync_group_id = azurerm_storage_sync_group.test.id
  registered_server_id  = trimspace(data.local_file.server_id.content)
  server_local_path     = "D:\\SyncFolder"

  cloud_tiering_enabled      = true
  volume_free_space_percent  = 30
  tier_files_older_than_days = 5
  local_cache_mode           = "UpdateLocallyCachedFiles"

  depends_on = [terraform_data.afs_server_id, azurerm_storage_sync_cloud_endpoint.test]
}
`, r.template(data), data.RandomString)
}

func (r StorageSyncServerEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-StorageSync-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestnw-%[1]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "acctestnic-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = "acctest-%[1]d"
  computer_name       = "afs%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  license_type        = "Windows_Server"
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  identity {
    type = "SystemAssigned"
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctest-disk-%[1]d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = 32
}

resource "azurerm_virtual_machine_data_disk_attachment" "test" {
  managed_disk_id    = azurerm_managed_disk.test.id
  virtual_machine_id = azurerm_windows_virtual_machine.test.id
  lun                = 0
  caching            = "None"
}

resource "azurerm_storage_sync" "test" {
  name                = "acctest-StorageSync-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_role_assignment" "vm_filesync_admin" {
  scope                = azurerm_storage_sync.test.id
  role_definition_name = "Azure File Sync Administrator"
  principal_id         = azurerm_windows_virtual_machine.test.identity[0].principal_id
}

resource "azurerm_storage_sync_group" "test" {
  name            = "acctest-StorageSyncGroup-%[1]d"
  storage_sync_id = azurerm_storage_sync.test.id
}

resource "azurerm_storage_account" "test" {
  name                     = "accstr%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name               = "acctest-share-%[1]d"
  storage_account_id = azurerm_storage_account.test.id
  quota              = 1

  acl {
    id = "GhostedRecall"
    access_policy {
      permissions = "r"
    }
  }
}

resource "azurerm_storage_sync_cloud_endpoint" "test" {
  name                  = "acctest-CEP-%[1]d"
  storage_sync_group_id = azurerm_storage_sync_group.test.id
  storage_account_id    = azurerm_storage_account.test.id
  file_share_name       = azurerm_storage_share.test.name
}

resource "azurerm_virtual_machine_run_command" "afs_register" {
  name               = "afs-register-%[1]d"
  location           = azurerm_resource_group.test.location
  virtual_machine_id = azurerm_windows_virtual_machine.test.id

  source {
    script = <<-SCRIPT
      $ErrorActionPreference = "Stop"
      $disk = Get-Disk | Where-Object PartitionStyle -eq RAW | Select-Object -First 1
      if ($disk) { $p = $disk | Initialize-Disk -PartitionStyle GPT -PassThru | New-Partition -AssignDriveLetter -UseMaximumSize; $p | Format-Volume -FileSystem NTFS | Out-Null }
      New-Item -Path D:\SyncFolder -ItemType Directory -Force | Out-Null
      New-Item -Path C:\Temp\AFS -ItemType Directory -Force | Out-Null
      [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
      Invoke-WebRequest -Uri https://aka.ms/afs/agent/Server2016 -OutFile C:\Temp\AFS\StorageSyncAgent.msi -UseBasicParsing
      Start-Process msiexec.exe -ArgumentList "/i","C:\Temp\AFS\StorageSyncAgent.msi","/qn","/norestart" -Wait -NoNewWindow
      $waited = 0; while ((Get-Service FileSyncSvc -ErrorAction SilentlyContinue).Status -ne "Running" -and $waited -lt 180) { Start-Sleep 10; $waited += 10 }
      if ((Get-Service FileSyncSvc -ErrorAction SilentlyContinue).Status -ne "Running") { throw "FileSyncSvc did not start within 3 minutes" }
      Install-PackageProvider -Name NuGet -MinimumVersion 2.8.5.201 -Force -Scope AllUsers
      Install-Module -Name Az.Accounts -Force -AllowClobber -Scope AllUsers
      Install-Module -Name Az.StorageSync -Force -AllowClobber -Scope AllUsers
      Import-Module Az.Accounts -Force
      Import-Module Az.StorageSync -Force
      $subId = (Invoke-RestMethod -Uri "http://169.254.169.254/metadata/instance?api-version=2021-02-01" -Headers @{Metadata="true"} -UseBasicParsing).compute.subscriptionId
      Connect-AzAccount -Identity
      Set-AzContext -Subscription $subId
      Register-AzStorageSyncServer -ResourceGroupName "${azurerm_resource_group.test.name}" -StorageSyncServiceName "${azurerm_storage_sync.test.name}"
    SCRIPT
  }

  depends_on = [
    azurerm_virtual_machine_data_disk_attachment.test,
    azurerm_role_assignment.vm_filesync_admin,
    azurerm_storage_sync.test,
  ]
}

resource "terraform_data" "afs_server_id" {
  depends_on = [azurerm_virtual_machine_run_command.afs_register]

  triggers_replace = {
    vm_id   = azurerm_windows_virtual_machine.test.id
    rg_name = azurerm_resource_group.test.name
    ss_name = azurerm_storage_sync.test.name
  }

  provisioner "local-exec" {
    interpreter = ["bash", "-c"]
    command     = <<-EOT
      az login --service-principal -u "$ARM_CLIENT_ID" -p "$ARM_CLIENT_SECRET" --tenant "$ARM_TENANT_ID" > /dev/null
      az account set --subscription "$ARM_SUBSCRIPTION_ID"
      sub=$(az account show --query id -o tsv)
      server_id=""
      for i in $(seq 1 12); do
        server_id=$(az rest \
          --method GET \
          --url "https://management.azure.com/subscriptions/$sub/resourceGroups/${azurerm_resource_group.test.name}/providers/Microsoft.StorageSync/storageSyncServices/${azurerm_storage_sync.test.name}/registeredServers?api-version=2020-03-01" \
          --query 'value[0].id' -o tsv 2>/dev/null)
        if [ -n "$server_id" ]; then
          break
        fi
        echo "Attempt $i/12: registered server not yet visible in API, retrying in 30s..."
        sleep 30
      done
      if [ -z "$server_id" ]; then
        echo "ERROR: registered server ID not found after waiting"
        exit 1
      fi
      echo "$server_id" > "/tmp/afs_server_id_${azurerm_resource_group.test.name}.txt"
    EOT
  }

  provisioner "local-exec" {
    when        = destroy
    interpreter = ["bash", "-c"]
    command     = <<-EOT
      az login --service-principal -u "$ARM_CLIENT_ID" -p "$ARM_CLIENT_SECRET" --tenant "$ARM_TENANT_ID" > /dev/null
      az account set --subscription "$ARM_SUBSCRIPTION_ID"
      sub=$(az account show --query id -o tsv)
      server_ids=$(az rest \
        --method GET \
        --url "https://management.azure.com/subscriptions/$sub/resourceGroups/${self.triggers_replace.rg_name}/providers/Microsoft.StorageSync/storageSyncServices/${self.triggers_replace.ss_name}/registeredServers?api-version=2020-03-01" \
        --query 'value[].id' -o tsv)
      echo "Servers to deregister: $server_ids"
      for sid in $server_ids; do
        echo "Deregistering: $sid"
        az rest --method DELETE \
          --url "https://management.azure.com$sid?api-version=2020-03-01" || true
      done
      if [ -n "$server_ids" ]; then
        for i in 1 2 3 4 5 6; do
          remaining=$(az rest --method GET \
            --url "https://management.azure.com/subscriptions/$sub/resourceGroups/${self.triggers_replace.rg_name}/providers/Microsoft.StorageSync/storageSyncServices/${self.triggers_replace.ss_name}/registeredServers?api-version=2020-03-01" \
            --query 'length(value)' -o tsv 2>/dev/null || echo "0")
          [ "$remaining" = "0" ] && break
          echo "Waiting for deregistration... ($i/6)"
          sleep 20
        done
      fi
      rm -f "/tmp/afs_server_id_${self.triggers_replace.rg_name}.txt"
    EOT
  }
}

data "local_file" "server_id" {
  filename   = "/tmp/afs_server_id_${azurerm_resource_group.test.name}.txt"
  depends_on = [terraform_data.afs_server_id]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
