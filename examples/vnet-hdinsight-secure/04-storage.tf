# Create storage accounts

resource "random_id" "storage_account_name_unique" {
  count       = "${var.storage_account_count}"
  byte_length = 8
}

resource "azurerm_storage_account" "storage" {
  count                    = "${var.storage_account_count}"
  name                     = "${element(random_id.storage_account_name_unique.*.hex, count.index)}"
  resource_group_name      = "${azurerm_resource_group.rg.name}"
  location                 = "${azurerm_resource_group.rg.location}"
  account_kind             = "StorageV2"
  account_tier             = "Standard"
  access_tier              = "Hot"
  account_replication_type = "${var.account_replication_type}"

  network_rules {
    ip_rules                   = ["127.0.0.1"]
    virtual_network_subnet_ids = ["${azurerm_subnet.subnet.id}"]
  }

  tags = "${var.tags}"
}
