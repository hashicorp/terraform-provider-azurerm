provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

data "azurerm_platform_image" "test" {
  location  = var.location
  publisher = "Canonical"
  offer     = "UbuntuServer"
  sku       = "18.04-LTS"
}

resource "azurerm_resource_group" "test" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_key_vault" "test" {
  name                        = "${var.prefix}kv"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "premium"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
}


# grant the service principal/user access to the key vault to be able to create the key
resource "azurerm_key_vault_access_policy" "service-principal" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "create",
    "delete",
    "get",
    "update",
  ]

  secret_permissions = [
    "get",
    "delete",
    "set",
  ]
}

# Disk Encryption Key
resource "azurerm_key_vault_secret" "test" {
  name         = "dek"
  value        = "<Disk Encryption Key value>"
  key_vault_id = azurerm_key_vault.test.id

  tags = {
    # DiskEncryptionKeyFileName is required when using the encrypted disk to create a Linux Virtual Machine
    DiskEncryptionKeyFileName = "LinuxPassPhraseFileName"
  }

  depends_on = [
    azurerm_key_vault_access_policy.service-principal
  ]
}

resource "azurerm_managed_disk" "test" {
  name                 = "${var.prefix}-disk"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "FromImage"
  image_reference_id   = data.azurerm_platform_image.test.id
  os_type              = "Linux"
  hyper_v_generation   = "V1"

  encryption_settings {
    disk_encryption_key {
      secret_url      = azurerm_key_vault_secret.test.id
      source_vault_id = azurerm_key_vault.test.id
    }
  }
}
