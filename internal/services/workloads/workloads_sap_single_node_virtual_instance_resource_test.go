// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workloads_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPSingleNodeVirtualInstanceResource struct{}

func TestAccWorkloadsSAPSingleNodeVirtualInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_single_node_virtual_instance", "test")
	r := WorkloadsSAPSingleNodeVirtualInstanceResource{}
	sapVISNameSuffix := SAPSingleNodeVirtualInstanceNameSuffix()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("single_server_configuration.0.virtual_machine_configuration.0.os_profile.0.ssh_private_key"),
	})
}

func TestAccWorkloadsSAPSingleNodeVirtualInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_single_node_virtual_instance", "test")
	r := WorkloadsSAPSingleNodeVirtualInstanceResource{}
	sapVISNameSuffix := SAPSingleNodeVirtualInstanceNameSuffix()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data, sapVISNameSuffix),
			ExpectError: acceptance.RequiresImportError(data.ResourceType),
		},
	})
}

func TestAccWorkloadsSAPSingleNodeVirtualInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_single_node_virtual_instance", "test")
	r := WorkloadsSAPSingleNodeVirtualInstanceResource{}
	sapVISNameSuffix := SAPSingleNodeVirtualInstanceNameSuffix()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("single_server_configuration.0.virtual_machine_configuration.0.os_profile.0.ssh_private_key"),
	})
}

func TestAccWorkloadsSAPSingleNodeVirtualInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_single_node_virtual_instance", "test")
	r := WorkloadsSAPSingleNodeVirtualInstanceResource{}
	sapVISNameSuffix := SAPSingleNodeVirtualInstanceNameSuffix()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("single_server_configuration.0.virtual_machine_configuration.0.os_profile.0.ssh_private_key"),
		{
			Config: r.update(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("single_server_configuration.0.virtual_machine_configuration.0.os_profile.0.ssh_private_key"),
		{
			Config: r.complete(data, sapVISNameSuffix),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("single_server_configuration.0.virtual_machine_configuration.0.os_profile.0.ssh_private_key"),
	})
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sapvirtualinstances.ParseSapVirtualInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Workloads.SAPVirtualInstances
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "tls_private_key" "test" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

data "tls_public_key" "test" {
  private_key_pem = tls_private_key.test.private_key_pem
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sapvis-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-uai-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_role_assignment" "test" {
  scope                = data.azurerm_subscription.current.id
  role_definition_name = "Azure Center for SAP solutions service role"
  principal_id         = azurerm_user_assigned_identity.test.principal_id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_resource_group" "app" {
  name     = "acctestRG-sapapp-%d"
  location = "%s"

  depends_on = [
    azurerm_subnet.test
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.Locations.Primary)
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) basic(data acceptance.TestData, sapVISNameSuffix int) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_workloads_sap_single_node_virtual_instance" "test" {
  name                        = "X%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  app_location                = azurerm_resource_group.app.location
  managed_resource_group_name = "acctestManagedRG%d"
  sap_fqdn                    = "sap.bpaas.com"

  single_server_configuration {
    app_resource_group_name = azurerm_resource_group.app.name
    subnet_id               = azurerm_subnet.test.id

    virtual_machine_configuration {
      virtual_machine_size = "Standard_E32ds_v4"

      image {
        offer     = "RHEL-SAP-HA"
        publisher = "RedHat"
        sku       = "82sapha-gen2"
        version   = "latest"
      }

      os_profile {
        admin_username  = "testAdmin"
        ssh_private_key = tls_private_key.test.private_key_pem
        ssh_public_key  = data.tls_public_key.test.public_key_openssh
      }
    }
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, r.template(data), sapVISNameSuffix, data.RandomInteger)
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) requiresImport(data acceptance.TestData, sapVISNameSuffix int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_workloads_sap_single_node_virtual_instance" "import" {
  name                        = azurerm_workloads_sap_single_node_virtual_instance.test.name
  resource_group_name         = azurerm_workloads_sap_single_node_virtual_instance.test.resource_group_name
  location                    = azurerm_workloads_sap_single_node_virtual_instance.test.location
  environment                 = azurerm_workloads_sap_single_node_virtual_instance.test.environment
  sap_product                 = azurerm_workloads_sap_single_node_virtual_instance.test.sap_product
  app_location                = azurerm_workloads_sap_single_node_virtual_instance.test.app_location
  managed_resource_group_name = azurerm_workloads_sap_single_node_virtual_instance.test.managed_resource_group_name
  sap_fqdn                    = "sap.bpaas.com"

  single_server_configuration {
    app_resource_group_name = azurerm_resource_group.app.name
    subnet_id               = azurerm_subnet.test.id

    virtual_machine_configuration {
      virtual_machine_size = "Standard_E32ds_v4"

      image {
        offer     = "RHEL-SAP-HA"
        publisher = "RedHat"
        sku       = "82sapha-gen2"
        version   = "latest"
      }

      os_profile {
        admin_username  = "testAdmin"
        ssh_private_key = tls_private_key.test.private_key_pem
        ssh_public_key  = data.tls_public_key.test.public_key_openssh
      }
    }
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}
`, r.basic(data, sapVISNameSuffix))
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) complete(data acceptance.TestData, sapVISNameSuffix int) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_workloads_sap_single_node_virtual_instance" "test" {
  name                        = "X%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "acctestManagedRG%d"
  app_location                = azurerm_resource_group.app.location
  sap_fqdn                    = "sap.bpaas.com"

  single_server_configuration {
    app_resource_group_name = azurerm_resource_group.app.name
    subnet_id               = azurerm_subnet.test.id
    database_type           = "HANA"
    secondary_ip_enabled    = true

    virtual_machine_configuration {
      virtual_machine_size = "Standard_E32ds_v4"

      image {
        offer     = "RHEL-SAP-HA"
        publisher = "RedHat"
        sku       = "82sapha-gen2"
        version   = "latest"
      }

      os_profile {
        admin_username  = "testAdmin"
        ssh_private_key = tls_private_key.test.private_key_pem
        ssh_public_key  = data.tls_public_key.test.public_key_openssh
      }
    }

    disk_volume_configuration {
      volume_name     = "hana/data"
      number_of_disks = 3
      size_in_gb      = 128
      sku_name        = "Premium_LRS"
    }

    disk_volume_configuration {
      volume_name     = "hana/log"
      number_of_disks = 3
      size_in_gb      = 128
      sku_name        = "Premium_LRS"
    }

    disk_volume_configuration {
      volume_name     = "hana/shared"
      number_of_disks = 1
      size_in_gb      = 256
      sku_name        = "Premium_LRS"
    }

    disk_volume_configuration {
      volume_name     = "usr/sap"
      number_of_disks = 1
      size_in_gb      = 128
      sku_name        = "Premium_LRS"
    }

    disk_volume_configuration {
      volume_name     = "backup"
      number_of_disks = 2
      size_in_gb      = 256
      sku_name        = "StandardSSD_LRS"
    }

    disk_volume_configuration {
      volume_name     = "os"
      number_of_disks = 1
      size_in_gb      = 64
      sku_name        = "StandardSSD_LRS"
    }

    virtual_machine_resource_names {
      host_name               = "apphostName0"
      os_disk_name            = "app0osdisk"
      virtual_machine_name    = "appvm0"
      network_interface_names = ["appnic0"]

      data_disk {
        volume_name = "default"
        names       = ["app0disk0"]
      }
    }
  }

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  tags = {
    Env = "Test"
  }

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, r.template(data), data.RandomString, sapVISNameSuffix, data.RandomInteger)
}

func (r WorkloadsSAPSingleNodeVirtualInstanceResource) update(data acceptance.TestData, sapVISNameSuffix int) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_workloads_sap_single_node_virtual_instance" "test" {
  name                        = "X%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  environment                 = "NonProd"
  sap_product                 = "S4HANA"
  managed_resource_group_name = "acctestManagedRG%d"
  app_location                = azurerm_resource_group.app.location
  sap_fqdn                    = "sap.bpaas.com"

  single_server_configuration {
    app_resource_group_name = azurerm_resource_group.app.name
    subnet_id               = azurerm_subnet.test.id
    database_type           = "HANA"
    secondary_ip_enabled    = true

    virtual_machine_configuration {
      virtual_machine_size = "Standard_E32ds_v4"

      image {
        offer     = "RHEL-SAP-HA"
        publisher = "RedHat"
        sku       = "82sapha-gen2"
        version   = "latest"
      }

      os_profile {
        admin_username  = "testAdmin"
        ssh_private_key = tls_private_key.test.private_key_pem
        ssh_public_key  = data.tls_public_key.test.public_key_openssh
      }
    }

    disk_volume_configuration {
      volume_name     = "hana/data"
      number_of_disks = 3
      size_in_gb      = 128
      sku_name        = "Premium_LRS"
    }

    disk_volume_configuration {
      volume_name     = "hana/log"
      number_of_disks = 3
      size_in_gb      = 128
      sku_name        = "Premium_LRS"
    }

    disk_volume_configuration {
      volume_name     = "hana/shared"
      number_of_disks = 1
      size_in_gb      = 256
      sku_name        = "Premium_LRS"
    }

    disk_volume_configuration {
      volume_name     = "usr/sap"
      number_of_disks = 1
      size_in_gb      = 128
      sku_name        = "Premium_LRS"
    }

    disk_volume_configuration {
      volume_name     = "backup"
      number_of_disks = 2
      size_in_gb      = 256
      sku_name        = "StandardSSD_LRS"
    }

    disk_volume_configuration {
      volume_name     = "os"
      number_of_disks = 1
      size_in_gb      = 64
      sku_name        = "StandardSSD_LRS"
    }

    virtual_machine_resource_names {
      host_name               = "apphostName0"
      os_disk_name            = "app0osdisk"
      virtual_machine_name    = "appvm0"
      network_interface_names = ["appnic0"]

      data_disk {
        volume_name = "default"
        names       = ["app0disk0"]
      }
    }
  }

  depends_on = [
    azurerm_role_assignment.test
  ]
}
`, r.template(data), data.RandomString, sapVISNameSuffix, data.RandomInteger)
}

func SAPSingleNodeVirtualInstanceNameSuffix() int {
	rand.NewSource(time.Now().UnixNano())
	num := rand.Intn(90) + 10

	return num
}
