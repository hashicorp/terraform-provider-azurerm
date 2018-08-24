resource "azurerm_virtual_machine_extension" "wsfc" {
  count                = "${var.sqlvmcount}"
  name                 = "create-cluster"
  resource_group_name  = "${var.resource_group_name}"
  location             = "${var.location}"
  virtual_machine_name = "${element(azurerm_virtual_machine.sql.*.name, count.index)}"
  publisher            = "Microsoft.Compute"
  type                 = "CustomScriptExtension"
  type_handler_version = "1.9"

  settings = <<SETTINGS
    { 
      "commandToExecute": "powershell Install-WindowsFeature -Name Failover-Clustering -IncludeManagementTools"
    } 
SETTINGS

  depends_on = ["azurerm_virtual_machine_extension.join-domain"]
}
