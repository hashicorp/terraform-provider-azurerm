package manageddevopspools_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-01-21/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedDevOpsPoolResource struct{}

func TestAccManagedDevOpsPoolSequential(t *testing.T) {
	if os.Getenv("ARM_MANAGED_DEVOPS_ORG_URL") == "" {
		t.Skip("Skipping as `ARM_MANAGED_DEVOPS_ORG_URL` is not specified")
	}

	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"managedDevOpsPool": {
			"basic":          TestAccManagedDevOpsPool_basic,
			"requiresImport": TestAccManagedDevOpsPool_requiresImport,
			"complete":       TestAccManagedDevOpsPool_complete,
			"update":         TestAccManagedDevOpsPool_update,
		},
	})
}

func TestAccManagedDevOpsPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
	r := ManagedDevOpsPoolResource{}

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

func TestAccManagedDevOpsPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
	r := ManagedDevOpsPoolResource{}

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

func TestAccManagedDevOpsPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
	r := ManagedDevOpsPoolResource{}

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

func TestAccManagedDevOpsPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_devops_pool", "test")
	r := ManagedDevOpsPoolResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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

func (ManagedDevOpsPoolResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := pools.ParsePoolID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ManagedDevOpsPools.PoolsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagedDevOpsPoolResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

resource "azurerm_managed_devops_pool" "test" {
  name                = "acctest-pool-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  maximum_concurrency            = 1
  dev_center_project_resource_id = azurerm_dev_center_project.test.id

  azure_devops_organization_profile {
    organization {
      url = "%s"
    }
  }

  stateless_agent_profile {}

  vmss_fabric_profile {
    image {
      well_known_image_name = "ubuntu-24.04"
      buffer                = "*"
    }
    sku_name = "Standard_D2ads_v5"
  }
}
`, r.template(data), data.RandomInteger, os.Getenv("ARM_MANAGED_DEVOPS_ORG_URL"))
}

func (r ManagedDevOpsPoolResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_devops_pool" "import" {
  name                = azurerm_managed_devops_pool.test.name
  location            = azurerm_managed_devops_pool.test.location
  resource_group_name = azurerm_managed_devops_pool.test.resource_group_name

  maximum_concurrency            = 1
  dev_center_project_resource_id = azurerm_dev_center_project.test.id

  azure_devops_organization_profile {
    organization {
      url = "%s"
    }
  }

  stateful_agent_profile {}

  vmss_fabric_profile {
    image {
      well_known_image_name = "ubuntu-24.04"
      buffer                = "*"
    }
    sku_name = "Standard_D2ads_v5"
  }
}
`, r.basic(data), os.Getenv("ARM_MANAGED_DEVOPS_ORG_URL"))
}

func (r ManagedDevOpsPoolResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkeyvault${var.random_string}"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    certificate_permissions = ["Create", "Delete", "Get", "Import", "Purge", "Recover", "Update", "List"]
    secret_permissions = [
      "Get",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Import",
      "Purge",
      "Recover",
      "Update",
      "List",
    ]

    key_permissions = [
      "Create",
    ]

    secret_permissions = [
      "Get",
      "Set",
    ]

    storage_permissions = [
      "Set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert${var.random_string}"
  key_vault_id = azurerm_key_vault.test.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyEncipherment",
        "keyCertSign",
      ]

      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}

provider "azuread" {}

data "azuread_service_principal" "devops_infrastructure" {
  display_name = "DevOpsInfrastructure"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "devops-infrastructure-delegation"
    service_delegation {
      name    = "Microsoft.DevOpsInfrastructure/pools"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_role_assignment" "devops_infrastructure_vnet_reader" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Reader"
  principal_id         = data.azuread_service_principal.devops_infrastructure.object_id
}

resource "azurerm_role_assignment" "devops_infrastructure_vnet_network_contributor" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Network Contributor"
  principal_id         = data.azuread_service_principal.devops_infrastructure.object_id
}

resource "azurerm_managed_devops_pool" "test" {
  name                = "acctest-pool-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  maximum_concurrency            = 1
  dev_center_project_resource_id = azurerm_dev_center_project.test.id

  azure_devops_organization_profile {
    organization {
      parallelism = 1
      url         = "%s"
    }
    permission_profile {
      kind = "SpecificAccounts"
      administrator_accounts {
        users = [data.azurerm_client_config.current.tenant_id]
      }
    }
  }

  stateful_agent_profile {
    grace_period_time_span = "00:10:00"
    max_agent_lifetime     = "08:00:00"
    manual_resource_predictions_profile {
      time_zone        = "UTC"
      monday_schedule  = { "09:00:00" = 1, "17:00:00" = 0 }
      tuesday_schedule = { "09:00:00" = 1, "17:00:00" = 0 }
    }
  }

  vmss_fabric_profile {
    image {
      resource_id = data.azurerm_platform_image.test.id
      aliases     = ["marketplace image"]
      buffer      = "0"
    }
    image {
      aliases               = ["well known image"]
      well_known_image_name = "ubuntu-24.04"
      buffer                = "100"
    }
    sku_name = "Standard_D2ads_v5"
    network_profile {
      subnet_id = azurerm_subnet.test.id
    }
    os_profile {
      logon_type = "Interactive"
      secrets_management {
        certificate_store_location = "/"
        certificate_store_name     = "My"

        key_export_enabled = false
        observed_certificates = [
          azurerm_key_vault_certificate.test.versionless_secret_id
        ]
      }
    }
    storage_profile {
      data_disk {
        caching              = "None"
        disk_size_gb         = 10
        drive_letter         = "F"
        storage_account_type = "Standard_LRS"
      }
      os_disk_storage_account_type = "Standard"
    }
  }

  tags = {
    Environment = "ppe"
    Project     = "Terraform"
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, os.Getenv("ARM_MANAGED_DEVOPS_ORG_URL"))
}

func (ManagedDevOpsPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
variable "primary_location" {
  default = %q
}
variable "random_integer" {
  default = %d
}
variable "random_string" {
  default = %q
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-${var.random_integer}"
  location = var.primary_location
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestuami-${var.random_string}"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_center" "test" {
  name                = "acctestdc-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_dev_center_project" "test" {
  name                = "acctestproj-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dev_center_id       = azurerm_dev_center.test.id
}

data "azurerm_platform_image" "test" {
  location  = azurerm_resource_group.test.location
  publisher = "Canonical"
  offer     = "0001-com-ubuntu-server-jammy"
  sku       = "22_04-lts"
}
`, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
