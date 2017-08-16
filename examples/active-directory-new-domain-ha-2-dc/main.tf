# Unique name for the storage account
resource "random_id" "storage_account_name" {
  keepers = {
    # Generate a new id each time a new resource group is defined
    resource_group = "${azurerm_resource_group.quickstartad.name}"
  }

  byte_length = 8
}

resource "azurerm_resource_group" "quickstartad" {
  name     = "${var.resource_group_name}"
  location = "${var.resource_group_location}"

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

# Using a storage account until managed disks are supported by terraform provider
resource "azurerm_storage_account" "adha_storage_account" {
  name                = "haad${random_id.storage_account_name.hex}"
  resource_group_name = "${azurerm_resource_group.quickstartad.name}"
  location            = "${azurerm_resource_group.quickstartad.location}"

  account_type = "${var.config["storage_account_type"]}"

  tags {
    Source = "Azure Quickstarts for Terraform"
  }
}

# Need a storage containers until managed disks supported by terraform provider
resource "azurerm_storage_container" "pdc_storage_container" {
  name                  = "${var.config["pdc_storage_container_name"]}"
  resource_group_name   = "${azurerm_resource_group.quickstartad.name}"
  storage_account_name  = "${azurerm_storage_account.adha_storage_account.name}"
  container_access_type = "private"
}

resource "azurerm_storage_container" "bdc_storage_container" {
  name                  = "${var.config["bdc_storage_container_name"]}"
  resource_group_name   = "${azurerm_resource_group.quickstartad.name}"
  storage_account_name  = "${azurerm_storage_account.adha_storage_account.name}"
  container_access_type = "private"
}

resource "azurerm_availability_set" "adha_availability_set" {
  name                = "${var.config["availability_set_name"]}"
  resource_group_name = "${azurerm_resource_group.quickstartad.name}"
  location            = "${azurerm_resource_group.quickstartad.location}"
}
