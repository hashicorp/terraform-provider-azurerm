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

// Helper functions for common test requirements
func requiresBasicEnvVars(t *testing.T) {
	if os.Getenv("ARM_MANAGED_DEVOPS_ORG_URL") == "" {
		t.Skip("Skipping as `ARM_MANAGED_DEVOPS_ORG_URL` is not specified")
	}
}

func requiresCompleteEnvVars(t *testing.T) {
	envVars := []string{
		"ARM_MANAGED_DEVOPS_ORG_URL",
		"ARM_MANAGED_DEVOPS_ORG_URL_UPDATED",
		"ARM_MANAGED_DEVOPS_ADMIN_EMAIL",
		"ARM_MANAGED_DEVOPS_ORG_PROJECT",
	}

	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Skipf("Skipping as `%s` is not specified", envVar)
		}
	}
}

func TestAccManagedDevOpsPool_basic(t *testing.T) {
	requiresBasicEnvVars(t)

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
	requiresBasicEnvVars(t)

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
	requiresCompleteEnvVars(t)

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
	requiresCompleteEnvVars(t)

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
		// Intermediate step: remove network profile from pool but keep other resources
		// This is required to solve the dependency timing issue where Terraform tries to
		// delete the subnet before updating the pool. By removing the network_profile first,
		// the pool no longer references the subnet, allowing safe resource cleanup in the next step.
		// Without this intermediate step, we get "InUseSubnetCannotBeDeleted" errors because
		// the pool still has a service association link to the subnet when Terraform tries to delete it.
		{
			Config: r.completeWithoutNetwork(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		// Final step: go to basic (this should now work since pool no longer references subnet)
		{
			Config: r.basic(data),
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

func (ManagedDevOpsPoolResource) keyVaultConfig() string {
	return `
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
}`
}

func (ManagedDevOpsPoolResource) networkingConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
      name = "Microsoft.DevOpsInfrastructure/pools"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action"
      ]
    }
  }
}`, data.RandomInteger, data.RandomInteger)
}

func (ManagedDevOpsPoolResource) roleAssignmentsConfig() string {
	return `
provider "azuread" {}

data "azuread_service_principal" "devops_infrastructure" {
  display_name = "DevOpsInfrastructure"
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
}`
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
    sku_name = "Standard_B1s"
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

  stateless_agent_profile {}

  vmss_fabric_profile {
    image {
      well_known_image_name = "ubuntu-24.04"
      buffer                = "*"
    }
    sku_name = "Standard_B1s"
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

%s

resource "azurerm_dev_center_project" "test2" {
  name                = "acctestproj2-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dev_center_id       = azurerm_dev_center.test.id
}

%s

%s

resource "azurerm_managed_devops_pool" "test" {
  name                = "acctest-pool-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  maximum_concurrency            = 2
  dev_center_project_resource_id = azurerm_dev_center_project.test2.id

  azure_devops_organization_profile {
    organization {
      parallelism = 1
      url         = "%s"
      projects    = ["%s"]
    }
    permission_profile {
      kind = "SpecificAccounts"
      administrator_account {
        users = ["%s"]
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
      aliases               = ["well known image", "22.04 version"]
      well_known_image_name = "ubuntu-22.04"
      buffer                = "100"
    }
    sku_name = "Standard_B1ms"
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

  depends_on = [
    azurerm_role_assignment.devops_infrastructure_vnet_reader,
    azurerm_role_assignment.devops_infrastructure_vnet_network_contributor,
  ]
}
`, r.template(data), r.keyVaultConfig(), r.networkingConfig(data), r.roleAssignmentsConfig(), data.RandomInteger, os.Getenv("ARM_MANAGED_DEVOPS_ORG_URL_UPDATED"), os.Getenv("ARM_MANAGED_DEVOPS_ORG_PROJECT"), os.Getenv("ARM_MANAGED_DEVOPS_ADMIN_EMAIL"))
}

func (r ManagedDevOpsPoolResource) completeWithoutNetwork(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm" {
  features {}
}

%s

resource "azurerm_dev_center_project" "test2" {
  name                = "acctestproj2-${var.random_string}"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  dev_center_id       = azurerm_dev_center.test.id
}

%s

%s

resource "azurerm_managed_devops_pool" "test" {
  name                = "acctest-pool-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  maximum_concurrency            = 2
  dev_center_project_resource_id = azurerm_dev_center_project.test2.id

  azure_devops_organization_profile {
    organization {
      parallelism = 1
      url         = "%s"
      projects    = ["%s"]
    }
    permission_profile {
      kind = "SpecificAccounts"
      administrator_account {
        users = ["%s"]
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
      aliases               = ["well known image", "22.04 version"]
      well_known_image_name = "ubuntu-22.04"
      buffer                = "100"
    }
    sku_name = "Standard_B1ms"
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

  depends_on = [
    azurerm_role_assignment.devops_infrastructure_vnet_reader,
    azurerm_role_assignment.devops_infrastructure_vnet_network_contributor,
  ]
}
`, r.template(data), r.keyVaultConfig(), r.networkingConfig(data), r.roleAssignmentsConfig(), data.RandomInteger, os.Getenv("ARM_MANAGED_DEVOPS_ORG_URL_UPDATED"), os.Getenv("ARM_MANAGED_DEVOPS_ORG_PROJECT"), os.Getenv("ARM_MANAGED_DEVOPS_ADMIN_EMAIL"))
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
