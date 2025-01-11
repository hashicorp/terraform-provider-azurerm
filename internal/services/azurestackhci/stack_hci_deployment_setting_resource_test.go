// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/deploymentsettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StackHCIDeploymentSettingResource struct{}

// The test must run in Windows environment with Administrator PowerShell, which will simulate HCI hardware by an Azure Windows VM based on a given image, and deploy HCI on it.
// based on https://github.com/Azure/Edge-infrastructure-quickstart-template
const (
	imageIdEnv                 = "ARM_TEST_HCI_IMAGE_ID"
	localAdminUserEnv          = "ARM_TEST_HCI_LOCAL_ADMIN_USER"
	localAdminUserPasswordEnv  = "ARM_TEST_HCI_LOCAL_ADMIN_USER_PASSWORD"
	domainAdminUserEnv         = "ARM_TEST_HCI_DOMAIN_ADMIN_USER"
	domainAdminUserPasswordEnv = "ARM_TEST_HCI_DOMAIN_ADMIN_USER_PASSWORD"
	deploymentUserEnv          = "ARM_TEST_HCI_DEPLOYMENT_USER"
	deploymentUserPasswordEnv  = "ARM_TEST_HCI_DEPLOYMENT_USER_PASSWORD"
)

func TestAccStackHCIDeploymentSetting_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_deployment_setting", "test")
	r := StackHCIDeploymentSettingResource{}

	if os.Getenv(imageIdEnv) == "" || os.Getenv(localAdminUserEnv) == "" || os.Getenv(domainAdminUserEnv) == "" || os.Getenv(domainAdminUserEnv) == "" || os.Getenv(domainAdminUserPasswordEnv) == "" || os.Getenv(deploymentUserEnv) == "" || os.Getenv(deploymentUserPasswordEnv) == "" {
		t.Skipf("skip the test as one or more of below environment variables are not specified: %q, %q, %q, %q, %q, %q, %q", imageIdEnv, localAdminUserEnv, localAdminUserPasswordEnv, domainAdminUserEnv, domainAdminUserPasswordEnv, deploymentUserEnv, deploymentUserPasswordEnv)
	}

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

func TestAccStackHCIDeploymentSetting_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_deployment_setting", "test")
	r := StackHCIDeploymentSettingResource{}

	if os.Getenv(imageIdEnv) == "" || os.Getenv(localAdminUserEnv) == "" || os.Getenv(domainAdminUserEnv) == "" || os.Getenv(domainAdminUserEnv) == "" || os.Getenv(domainAdminUserPasswordEnv) == "" || os.Getenv(deploymentUserEnv) == "" || os.Getenv(deploymentUserPasswordEnv) == "" {
		t.Skipf("skip the test as one or more of below environment variables are not specified: %q, %q, %q, %q, %q, %q, %q", imageIdEnv, localAdminUserEnv, localAdminUserPasswordEnv, domainAdminUserEnv, domainAdminUserPasswordEnv, deploymentUserEnv, deploymentUserPasswordEnv)
	}

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

func TestAccStackHCIDeploymentSetting_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_stack_hci_deployment_setting", "test")
	r := StackHCIDeploymentSettingResource{}

	if os.Getenv(imageIdEnv) == "" || os.Getenv(localAdminUserEnv) == "" || os.Getenv(domainAdminUserEnv) == "" || os.Getenv(domainAdminUserEnv) == "" || os.Getenv(domainAdminUserPasswordEnv) == "" || os.Getenv(deploymentUserEnv) == "" || os.Getenv(deploymentUserPasswordEnv) == "" {
		t.Skipf("skip the test as one or more of below environment variables are not specified: %q, %q, %q, %q, %q, %q, %q", imageIdEnv, localAdminUserEnv, localAdminUserPasswordEnv, domainAdminUserEnv, domainAdminUserPasswordEnv, deploymentUserEnv, deploymentUserPasswordEnv)
	}

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

func (r StackHCIDeploymentSettingResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	clusterClient := client.AzureStackHCI.DeploymentSettings
	id, err := deploymentsettings.ParseDeploymentSettingID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r StackHCIDeploymentSettingResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    resource_group {
      // need to delete arc vm
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_stack_hci_deployment_setting" "test" {
  stack_hci_cluster_id = azurerm_stack_hci_cluster.test.id
  arc_resource_ids     = [for server in data.azurerm_arc_machine.server : server.id]
  version              = "10.0.0.0"

  scale_unit {
    active_directory_organizational_unit_path = "OU=hci${var.random_string},DC=jumpstart,DC=local"
    domain_fqdn                               = "jumpstart.local"
    secrets_location                          = azurerm_key_vault.DeploymentKeyVault.vault_uri
    name_prefix                               = "hci${var.random_string}"
    streaming_data_client_enabled             = true
    eu_location_enabled                       = false
    episodic_data_upload_enabled              = true

    bitlocker_boot_volume_enabled   = true
    bitlocker_data_volume_enabled   = true
    credential_guard_enabled        = true
    drift_control_enabled           = true
    drtm_protection_enabled         = true
    hvci_protection_enabled         = true
    side_channel_mitigation_enabled = true
    smb_cluster_encryption_enabled  = false
    smb_signing_enabled             = true
    wdac_enabled                    = true

    cluster {
      azure_service_endpoint = "core.windows.net"
      cloud_account_name     = azurerm_storage_account.witness.name
      name                   = azurerm_stack_hci_cluster.test.name
      witness_type           = "Cloud"
      witness_path           = "Cloud"
    }

    host_network {
      intent {
        name                              = "ManagementCompute"
        adapter_property_override_enabled = true
        adapter_property_override {
          network_direct = "Disabled"
        }
        adapter = [
          "FABRIC",
          "FABRIC2",
        ]
        traffic_type = [
          "Management",
          "Compute",
        ]
      }

      intent {
        name                              = "Storage"
        adapter_property_override_enabled = true
        adapter_property_override {
          network_direct = "Disabled"
        }
        adapter = [
          "StorageA",
          "StorageB",
        ]
        traffic_type = [
          "Storage",
        ]
      }

      storage_network {
        name                 = "Storage1Network"
        network_adapter_name = "StorageA"
        vlan_id              = "711"
      }

      storage_network {
        name                 = "Storage2Network"
        network_adapter_name = "StorageB"
        vlan_id              = "712"
      }
    }

    infrastructure_network {
      gateway     = "192.168.1.1"
      subnet_mask = "255.255.255.0"
      dns_server = [
        "192.168.1.254"
      ]
      ip_pool {
        ending_address   = "192.168.1.65"
        starting_address = "192.168.1.55"
      }
    }

    optional_service {
      custom_location = "customlocation${var.random_string}"
    }

    physical_node {
      ipv4_address = "192.168.1.12"
      name         = "AzSHOST1"
    }

    physical_node {
      ipv4_address = "192.168.1.13"
      name         = "AzSHOST2"
    }

    storage {
      configuration_mode = "Express"
    }
  }
}
`, template)
}

func (r StackHCIDeploymentSettingResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)

	return fmt.Sprintf(`
%s

resource "azurerm_stack_hci_deployment_setting" "import" {
  stack_hci_cluster_id = azurerm_stack_hci_cluster.test.id
  arc_resource_ids     = azurerm_stack_hci_deployment_setting.test.arc_resource_ids
  version              = azurerm_stack_hci_deployment_setting.test.version

  scale_unit {
    active_directory_organizational_unit_path = azurerm_stack_hci_deployment_setting.test.scale_unit.0.active_directory_organizational_unit_path
    domain_fqdn                               = azurerm_stack_hci_deployment_setting.test.scale_unit.0.domain_fqdn
    secrets_location                          = azurerm_stack_hci_deployment_setting.test.scale_unit.0.secrets_location
    name_prefix                               = azurerm_stack_hci_deployment_setting.test.scale_unit.0.name_prefix
    streaming_data_client_enabled             = azurerm_stack_hci_deployment_setting.test.scale_unit.0.streaming_data_client_enabled
    eu_location_enabled                       = azurerm_stack_hci_deployment_setting.test.scale_unit.0.eu_location_enabled
    episodic_data_upload_enabled              = azurerm_stack_hci_deployment_setting.test.scale_unit.0.episodic_data_upload_enabled

    bitlocker_boot_volume_enabled   = azurerm_stack_hci_deployment_setting.test.scale_unit.0.bitlocker_boot_volume_enabled
    bitlocker_data_volume_enabled   = azurerm_stack_hci_deployment_setting.test.scale_unit.0.bitlocker_data_volume_enabled
    credential_guard_enabled        = azurerm_stack_hci_deployment_setting.test.scale_unit.0.credential_guard_enabled
    drift_control_enabled           = azurerm_stack_hci_deployment_setting.test.scale_unit.0.drift_control_enabled
    drtm_protection_enabled         = azurerm_stack_hci_deployment_setting.test.scale_unit.0.drtm_protection_enabled
    hvci_protection_enabled         = azurerm_stack_hci_deployment_setting.test.scale_unit.0.hvci_protection_enabled
    side_channel_mitigation_enabled = azurerm_stack_hci_deployment_setting.test.scale_unit.0.side_channel_mitigation_enabled
    smb_cluster_encryption_enabled  = azurerm_stack_hci_deployment_setting.test.scale_unit.0.smb_cluster_encryption_enabled
    smb_signing_enabled             = azurerm_stack_hci_deployment_setting.test.scale_unit.0.smb_signing_enabled
    wdac_enabled                    = azurerm_stack_hci_deployment_setting.test.scale_unit.0.wdac_enabled


    cluster {
      azure_service_endpoint = azurerm_stack_hci_deployment_setting.test.scale_unit.0.cluster.0.azure_service_endpoint
      cloud_account_name     = azurerm_stack_hci_deployment_setting.test.scale_unit.0.cluster.0.cloud_account_name
      name                   = azurerm_stack_hci_deployment_setting.test.scale_unit.0.cluster.0.name
      witness_type           = azurerm_stack_hci_deployment_setting.test.scale_unit.0.cluster.0.witness_type
      witness_path           = azurerm_stack_hci_deployment_setting.test.scale_unit.0.cluster.0.witness_path
    }

    host_network {
      intent {
        name                              = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.intent.0.name
        adapter                           = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.intent.0.adapter
        adapter_property_override_enabled = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.intent.0.adapter_property_override_enabled
        traffic_type                      = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.intent.0.traffic_type
        adapter_property_override {
          network_direct = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.intent.0.adapter_property_override.0.network_direct
        }
      }

      intent {
        name                              = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.intent.1.name
        adapter                           = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.intent.1.adapter
        adapter_property_override_enabled = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.intent.0.adapter_property_override_enabled
        traffic_type                      = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.intent.1.traffic_type
        adapter_property_override {
          network_direct = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.intent.0.adapter_property_override.0.network_direct
        }
      }

      storage_network {
        name                 = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.storage_network.0.name
        network_adapter_name = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.storage_network.0.network_adapter_name
        vlan_id              = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.storage_network.0.vlan_id
      }

      storage_network {
        name                 = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.storage_network.1.name
        network_adapter_name = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.storage_network.1.network_adapter_name
        vlan_id              = azurerm_stack_hci_deployment_setting.test.scale_unit.0.host_network.0.storage_network.1.vlan_id
      }
    }

    infrastructure_network {
      gateway     = azurerm_stack_hci_deployment_setting.test.scale_unit.0.infrastructure_network.0.gateway
      subnet_mask = azurerm_stack_hci_deployment_setting.test.scale_unit.0.infrastructure_network.0.subnet_mask
      dns_server  = azurerm_stack_hci_deployment_setting.test.scale_unit.0.infrastructure_network.0.dns_server
      ip_pool {
        ending_address   = azurerm_stack_hci_deployment_setting.test.scale_unit.0.infrastructure_network.0.ip_pool.0.ending_address
        starting_address = azurerm_stack_hci_deployment_setting.test.scale_unit.0.infrastructure_network.0.ip_pool.0.starting_address
      }
    }

    optional_service {
      custom_location = azurerm_stack_hci_deployment_setting.test.scale_unit.0.optional_service.0.custom_location
    }

    physical_node {
      ipv4_address = azurerm_stack_hci_deployment_setting.test.scale_unit.0.physical_node.0.ipv4_address
      name         = azurerm_stack_hci_deployment_setting.test.scale_unit.0.physical_node.0.name
    }

    physical_node {
      ipv4_address = azurerm_stack_hci_deployment_setting.test.scale_unit.0.physical_node.1.ipv4_address
      name         = azurerm_stack_hci_deployment_setting.test.scale_unit.0.physical_node.1.name
    }

    storage {
      configuration_mode = azurerm_stack_hci_deployment_setting.test.scale_unit.0.storage.0.configuration_mode
    }
  }
}
`, config)
}

func (r StackHCIDeploymentSettingResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    resource_group {
      // need to delete arc vm
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_stack_hci_deployment_setting" "test" {
  stack_hci_cluster_id = azurerm_stack_hci_cluster.test.id
  arc_resource_ids     = [for server in data.azurerm_arc_machine.server : server.id]
  version              = "10.0.0.0"

  scale_unit {
    active_directory_organizational_unit_path = "OU=hci${var.random_string},DC=jumpstart,DC=local"
    domain_fqdn                               = "jumpstart.local"
    secrets_location                          = azurerm_key_vault.DeploymentKeyVault.vault_uri
    name_prefix                               = "hci${var.random_string}"
    streaming_data_client_enabled             = true
    eu_location_enabled                       = false
    episodic_data_upload_enabled              = true
    bitlocker_boot_volume_enabled             = true
    bitlocker_data_volume_enabled             = true
    credential_guard_enabled                  = true
    drift_control_enabled                     = true
    drtm_protection_enabled                   = true
    hvci_protection_enabled                   = true
    side_channel_mitigation_enabled           = true
    smb_cluster_encryption_enabled            = false
    smb_signing_enabled                       = true
    wdac_enabled                              = true


    cluster {
      azure_service_endpoint = "core.windows.net"
      cloud_account_name     = azurerm_storage_account.witness.name
      name                   = azurerm_stack_hci_cluster.test.name
      witness_type           = "Cloud"
      witness_path           = "Cloud"
    }

    host_network {
      storage_auto_ip_enabled                 = true
      storage_connectivity_switchless_enabled = false

      intent {
        name                                          = "ManagementCompute"
        adapter_property_override_enabled             = true
        qos_policy_override_enabled                   = false
        virtual_switch_configuration_override_enabled = false
        adapter = [
          "FABRIC",
          "FABRIC2",
        ]
        traffic_type = [
          "Management",
          "Compute",
        ]
        qos_policy_override {
          priority_value8021_action_cluster = "7"
          priority_value8021_action_smb     = "3"
          bandwidth_percentage_smb          = "50"
        }
        adapter_property_override {
          network_direct = "Disabled"
        }
      }

      intent {
        name                                          = "Storage"
        adapter_property_override_enabled             = true
        qos_policy_override_enabled                   = false
        virtual_switch_configuration_override_enabled = false
        adapter = [
          "StorageA",
          "StorageB",
        ]
        traffic_type = [
          "Storage",
        ]
        qos_policy_override {
          priority_value8021_action_cluster = "7"
          priority_value8021_action_smb     = "3"
          bandwidth_percentage_smb          = "50"
        }
        adapter_property_override {
          network_direct = "Disabled"
        }
      }

      storage_network {
        name                 = "Storage1Network"
        network_adapter_name = "StorageA"
        vlan_id              = "711"
      }

      storage_network {
        name                 = "Storage2Network"
        network_adapter_name = "StorageB"
        vlan_id              = "712"
      }
    }

    infrastructure_network {
      gateway      = "192.168.1.1"
      subnet_mask  = "255.255.255.0"
      dhcp_enabled = false
      dns_server = [
        "192.168.1.254"
      ]
      ip_pool {
        ending_address   = "192.168.1.65"
        starting_address = "192.168.1.55"
      }
    }

    optional_service {
      custom_location = "customlocation${var.random_string}"
    }

    physical_node {
      ipv4_address = "192.168.1.12"
      name         = "AzSHOST1"
    }

    physical_node {
      ipv4_address = "192.168.1.13"
      name         = "AzSHOST2"
    }

    storage {
      configuration_mode = "Express"
    }
  }
}
`, template)
}

func (r StackHCIDeploymentSettingResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}

variable "random_string" {
  default = %q
}

variable "image_id" {
  description = "The ID of the image to use for the HCI virtual host machine"
  sensitive   = true
  type        = string
  default     = %q
}

variable "local_admin_user" {
  description = "The username of the local administrator account."
  sensitive   = true
  type        = string
  default     = %q
}

variable "local_admin_password" {
  description = "The password of the local administrator account."
  sensitive   = true
  type        = string
  default     = %q
}

variable "domain_admin_user" {
  description = "The username of the domain account."
  sensitive   = true
  type        = string
  default     = %q
}

variable "domain_admin_password" {
  description = "The password of the domain account."
  sensitive   = true
  type        = string
  default     = %q
}

variable "deployment_user" {
  sensitive   = true
  type        = string
  description = "The username for deployment user."
  default     = %q
}

variable "deployment_user_password" {
  sensitive   = true
  type        = string
  description = "The password for deployment user."
  default     = %q
}

provider "azuread" {}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-hci-vm-${var.random_string}"
  location = var.primary_location
}

resource "azurerm_public_ip" "test" {
  name                = "ip-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network" "test" {
  name                = "vnet-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["172.17.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "default"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["172.17.0.0/24"]
}

resource "azurerm_network_interface" "test" {
  name                = "ni-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  ip_configuration {
    name                          = "ipconfig1"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.test.id
    subnet_id                     = azurerm_subnet.test.id
  }
}

resource "azurerm_network_interface_security_group_association" "test" {
  network_interface_id      = azurerm_network_interface.test.id
  network_security_group_id = azurerm_network_security_group.test.id
  depends_on = [
    azurerm_network_interface.test,
    azurerm_network_security_group.test,
  ]
}

resource "azurerm_network_security_group" "test" {
  name                = "nsg-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_network_security_rule" "test1" {
  access                      = "Allow"
  destination_address_prefix  = "*"
  destination_port_ranges     = ["6985", "15985", "25985"]
  direction                   = "Inbound"
  name                        = "rule1-${var.random_string}"
  network_security_group_name = azurerm_network_security_group.test.name
  priority                    = 103
  protocol                    = "*"
  resource_group_name         = azurerm_resource_group.test.name
  source_address_prefix       = "*"
  source_port_range           = "*"
}

resource "azurerm_virtual_machine" "test" {
  name                             = "vm-${var.random_string}"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  network_interface_ids            = [azurerm_network_interface.test.id]
  vm_size                          = "Standard_E32s_v5"
  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true
  storage_image_reference {
    id = var.image_id
  }
  boot_diagnostics {
    enabled     = true
    storage_uri = ""
  }

  storage_os_disk {
    create_option = "FromImage"
    name          = "vm-${var.random_string}_OsDisk"
  }

  depends_on = [
    azurerm_network_interface_security_group_association.test,
  ]
}

locals {
  servers = [
    {
      name         = "AzSHOST1"
      ipv4_address = "192.168.1.12"
      port         = 15985
    },
    {
      name         = "AzSHOST2"
      ipv4_address = "192.168.1.13"
      port         = 25985
    }
  ]
  connection_roles = [
    "Azure Connected Machine Onboarding",
    "Azure Connected Machine Resource Administrator",
    "Azure Resource Bridge Deployment Role"
  ]
  machine_roles = [
    "Key Vault Secrets User",
    "Azure Connected Machine Resource Manager",
    "Azure Stack HCI Device Management Role",
    "Reader"
  ]
}

resource "azurerm_resource_group" "test2" {
  name     = "acctest-hci-${var.random_string}"
  location = var.primary_location
}

resource "azuread_application" "test" {
  display_name = "acctest-hci-onboard-${var.random_string}"
}

resource "azuread_service_principal" "test" {
  client_id = azuread_application.test.client_id
}

resource "azuread_service_principal_password" "test" {
  service_principal_id = azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "connect" {
  count                = length(local.connection_roles)
  scope                = azurerm_resource_group.test2.id
  role_definition_name = local.connection_roles[count.index]
  principal_id         = azuread_service_principal.test.object_id
}

resource "azurerm_role_assignment" "connect2" {
  scope                = data.azurerm_subscription.current.id
  role_definition_name = "Contributor"
  principal_id         = azuread_service_principal.test.object_id
}

// this is following https://learn.microsoft.com/en-us/azure-stack/hci/deploy/deployment-tool-active-directory
resource "terraform_data" "ad_creation_provisioner" {
  depends_on = [azurerm_virtual_machine.test]

  provisioner "local-exec" {
    command     = "powershell.exe -ExecutionPolicy Bypass -NoProfile -File ./testdata/ad.ps1 -userName ${var.domain_admin_user} -password \"${var.domain_admin_password}\" -authType Credssp -ip ${azurerm_public_ip.test.ip_address} -port 6985 -adouPath OU=hci${var.random_string},DC=jumpstart,DC=local -domainFqdn jumpstart.local -ifdeleteadou false -deploymentUserName ${var.deployment_user} -deploymentUserPassword \"${var.deployment_user_password}\""
    interpreter = ["PowerShell", "-Command"]
  }

  lifecycle {
    replace_triggered_by = [azurerm_resource_group.test2.name]
  }
}

resource "terraform_data" "provisioner" {
  count = length(local.servers)

  depends_on = [
    terraform_data.ad_creation_provisioner,
    azurerm_role_assignment.connect,
    azurerm_role_assignment.connect2,
  ]

  provisioner "local-exec" {
    command = "echo Connect ${local.servers[count.index].name} to Azure Arc..."
  }

  provisioner "local-exec" {
    command     = "powershell.exe -ExecutionPolicy Bypass -NoProfile -File ./testdata/connect.ps1 -userName ${var.local_admin_user} -password \"${var.local_admin_password}\" -authType Credssp -ip ${azurerm_public_ip.test.ip_address} -port ${local.servers[count.index].port} -subscriptionId ${data.azurerm_client_config.current.subscription_id} -resourceGroupName ${azurerm_resource_group.test2.name} -region ${azurerm_resource_group.test2.location} -tenant ${data.azurerm_client_config.current.tenant_id} -servicePrincipalId ${azuread_service_principal.test.client_id} -servicePrincipalSecret ${azuread_service_principal_password.test.value} -expandC true"
    interpreter = ["PowerShell", "-Command"]
  }

  provisioner "local-exec" {
    command = "echo connected ${local.servers[count.index].name}"
  }

  lifecycle {
    replace_triggered_by = [azurerm_resource_group.test2.name]
  }
}

data "azurerm_arc_machine" "server" {
  count               = length(local.servers)
  name                = local.servers[count.index].name
  resource_group_name = azurerm_resource_group.test2.name
  depends_on          = [terraform_data.provisioner]
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "DeploymentKeyVault" {
  name                            = "hci${var.random_string}-testkv"
  location                        = azurerm_resource_group.test2.location
  resource_group_name             = azurerm_resource_group.test2.name
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  enabled_for_disk_encryption     = true
  tenant_id                       = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days      = 30
  enable_rbac_authorization       = true
  public_network_access_enabled   = true
  sku_name                        = "standard"
}

resource "azurerm_role_assignment" "KeyVault" {
  scope                = azurerm_key_vault.DeploymentKeyVault.id
  role_definition_name = "Key Vault Secrets Officer"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_secret" "AzureStackLCMUserCredential" {
  name         = "AzureStackLCMUserCredential"
  content_type = "Secret"
  value        = base64encode("${var.deployment_user}:${var.deployment_user_password}")
  key_vault_id = azurerm_key_vault.DeploymentKeyVault.id
  depends_on   = [azurerm_role_assignment.KeyVault]
}

resource "azurerm_key_vault_secret" "LocalAdminCredential" {
  name         = "LocalAdminCredential"
  content_type = "Secret"
  value        = base64encode("${var.local_admin_user}:${var.local_admin_password}")
  key_vault_id = azurerm_key_vault.DeploymentKeyVault.id
  depends_on   = [azurerm_role_assignment.KeyVault]
}

resource "azurerm_key_vault_secret" "DefaultARBApplication" {
  name         = "DefaultARBApplication"
  content_type = "Secret"
  value        = base64encode("${azuread_service_principal.test.client_id}:${azuread_service_principal_password.test.value}")
  key_vault_id = azurerm_key_vault.DeploymentKeyVault.id
  depends_on   = [azurerm_role_assignment.KeyVault]
}

resource "azurerm_key_vault_secret" "WitnessStorageKey" {
  name         = "WitnessStorageKey"
  content_type = "Secret"
  value        = base64encode(azurerm_storage_account.witness.primary_access_key)
  key_vault_id = azurerm_key_vault.DeploymentKeyVault.id
  depends_on   = [azurerm_role_assignment.KeyVault]
}

resource "azurerm_storage_account" "witness" {
  name                     = "hci${var.random_string}teststa"
  location                 = azurerm_resource_group.test2.location
  resource_group_name      = azurerm_resource_group.test2.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

// service principal of 'Microsoft.AzureStackHCI Resource Provider'
data "azuread_service_principal" "hciRp" {
  client_id = "1412d89f-b8a8-4111-b4fd-e82905cbd85d"
}

resource "azurerm_role_assignment" "MachineRoleAssign1" {
  count                = length(local.machine_roles)
  scope                = azurerm_resource_group.test2.id
  role_definition_name = local.machine_roles[count.index]
  principal_id         = data.azurerm_arc_machine.server[0].identity[0].principal_id
}

resource "azurerm_role_assignment" "MachineRoleAssign2" {
  count                = length(local.machine_roles)
  scope                = azurerm_resource_group.test2.id
  role_definition_name = local.machine_roles[count.index]
  principal_id         = data.azurerm_arc_machine.server[1].identity[0].principal_id
}

resource "azurerm_role_assignment" "ServicePrincipalRoleAssign" {
  scope                = azurerm_resource_group.test2.id
  role_definition_name = "Azure Connected Machine Resource Manager"
  principal_id         = data.azuread_service_principal.hciRp.object_id
}

resource "azurerm_stack_hci_cluster" "test" {
  depends_on = [
    azurerm_key_vault_secret.DefaultARBApplication,
    azurerm_key_vault_secret.AzureStackLCMUserCredential,
    azurerm_key_vault_secret.LocalAdminCredential,
    azurerm_key_vault_secret.WitnessStorageKey,
    azurerm_role_assignment.ServicePrincipalRoleAssign,
    azurerm_role_assignment.MachineRoleAssign1,
    azurerm_role_assignment.MachineRoleAssign2,
  ]
  name                = "hci${var.random_string}-cl"
  resource_group_name = azurerm_resource_group.test2.name
  location            = azurerm_resource_group.test2.location
  identity {
    type = "SystemAssigned"
  }

  // the client_id will be populated after deployment
  lifecycle {
    ignore_changes = [
      client_id
    ]
  }
}
`, data.Locations.Primary, data.RandomString, os.Getenv(imageIdEnv), os.Getenv(localAdminUserEnv), os.Getenv(localAdminUserPasswordEnv), os.Getenv(domainAdminUserEnv), os.Getenv(domainAdminUserPasswordEnv), os.Getenv(deploymentUserEnv), os.Getenv(deploymentUserPasswordEnv))
}
