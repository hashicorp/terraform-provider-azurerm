resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}stor"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags {
    environment = "staging"
  }
}

resource "azurerm_template_deployment" "example" {
  name                = "${var.prefix}-deployment"
  resource_group_name = "${azurerm_resource_group.example.name}"
  deployment_mode     = "Incremental"
  template_body       = "${file("hdinsights.json")}"

  # these key-value pairs are passed into the ARM Template's `parameters` block
  parameters {
    clusterName             = "${var.prefix}-hdinsight"
    location                = "${azurerm_resource_group.example.location}"
    storageAccountEndpoint  = "${replace(replace(azurerm_storage_account.example.primary_blob_endpoint, "https://", ""), "/", "")}"
    storageAccountAccessKey = "${azurerm_storage_account.example.primary_access_key}"
    clusterLoginUserName    = "testadmin"
    clusterLoginPassword    = "Password1234!"
    sshUserName             = "testadmin2"
    sshPassword             = "Password12345!"
  }
}
