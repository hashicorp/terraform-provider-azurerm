# Need to add resource group for Terraform
resource "azurerm_resource_group" "resource_group" {
  name     = "${var.resourceGroupName}"
  location = "${var.resourceGroupLocation}"
}

# ARM Template reference
data "template_file" "azurerm_template" {
  template = "${file("azuredeploy.json")}"
}

# Azure template deployment (currently no HDInsight Terrsform provider)
resource "azurerm_template_deployment" "azurerm_template" {
  name                = "hdinsight-spark-linux-vnet"
  resource_group_name = "${azurerm_resource_group.resource_group.name}"

  parameters = {
    clusterName          = "${var.clusterName}"
    clusterLoginUserName = "${var.clusterLoginUserName}"
    clusterLoginPassword = "${var.clusterLoginPassword}"
    sshUserName          = "${var.sshUserName}"
    sshPassword          = "${var.sshPassword}"
  }

  template_body   = "${data.template_file.azurerm_template.rendered}"
  deployment_mode = "Incremental"
}