# Storage Mover example - SMB mount endpoint branch.
# Run: make build && terraform init && terraform plan

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0"
    }
  }
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "acctest-rg-storage-mover-example"
  location = var.primary_location
}

resource "azurerm_storage_mover" "example" {
  name                = "acctest-ssm-example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  description         = "Example Storage Mover"
}

resource "azurerm_storage_account" "example" {
  name                     = "acctestsa${substr(md5(azurerm_resource_group.example.name), 0, 8)}"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  allow_nested_items_to_be_public = true
}

resource "azurerm_storage_container" "example" {
  name                  = "acctest-container"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "blob"
}

resource "azurerm_storage_mover_target_endpoint" "example" {
  name                   = "acctest-smte-example"
  storage_mover_id       = azurerm_storage_mover.example.id
  storage_account_id     = azurerm_storage_account.example.id
  storage_container_name = azurerm_storage_container.example.name
  description            = "Example target blob container endpoint"
}

resource "azurerm_storage_mover_source_endpoint" "example" {
  name             = "acctest-smse-example"
  storage_mover_id = azurerm_storage_mover.example.id
  host             = var.nfs_host
  export           = "/"
  nfs_version      = "NFSv4"
  description      = "Example NFS source endpoint"
}

resource "azurerm_storage_mover_smb_mount_endpoint" "example" {
  name             = "acctest-smsme-example"
  storage_mover_id = azurerm_storage_mover.example.id
  host             = var.smb_host
  share_name       = var.smb_share_name
  description      = "Example SMB mount endpoint"
}

resource "azurerm_storage_mover_project" "example" {
  name             = "acctest-smp-example"
  storage_mover_id = azurerm_storage_mover.example.id
  description      = "Example Storage Mover project"
}

resource "azurerm_storage_mover_job_definition" "example" {
  name                     = "acctest-smjd-example"
  storage_mover_project_id  = azurerm_storage_mover_project.example.id
  source_name               = azurerm_storage_mover_source_endpoint.example.name
  target_name               = azurerm_storage_mover_target_endpoint.example.name
  copy_mode                 = "Additive"
  description               = "Example job definition"
}
