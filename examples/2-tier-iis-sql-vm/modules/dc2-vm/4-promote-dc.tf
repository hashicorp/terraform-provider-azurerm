// the `exit_code_hack` is to keep the VM Extension resource happy
locals {
  dc2import_command       = "Import-Module ADDSDeployment"
  dc2user_command         = "$dc2user = ${var.domainadmin_username}" 
  dc2password_command     = "$password = ConvertTo-SecureString ${var.admin_password} -AsPlainText -Force"
  dc2creds_command        = "$mycreds = New-Object System.Management.Automation.PSCredential -ArgumentList $dc2user, $password"
  dc2install_ad_command   = "Add-WindowsFeature -name ad-domain-services -IncludeManagementTools"
  dc2configure_ad_command = "Install-ADDSDomainController -Credential $mycreds -CreateDnsDelegation:$false -DomainName ${var.active_directory_domain} -InstallDns:$true -SafeModeAdministratorPassword $password -Force:$true"
  dc2shutdown_command     = "shutdown -r -t 10"
  dc2exit_code_hack       = "exit 0"
  dc2powershell_command   = "${local.dc2import_command}; ${local.dc2user_command}; ${local.dc2password_command}; ${local.dc2creds_command}; ${local.dc2install_ad_command}; ${local.dc2configure_ad_command}; ${local.dc2shutdown_command}; ${local.dc2exit_code_hack}"
}

resource "azurerm_virtual_machine_extension" "promote-dc" {
  name                 = "promote-dc"
  location             = "${azurerm_virtual_machine_extension.join-domain.location}"
  resource_group_name  = "${var.resource_group_name}"
  virtual_machine_name = "${azurerm_virtual_machine.domain-controller2.name}"
  publisher            = "Microsoft.Compute"
  type                 = "CustomScriptExtension"
  type_handler_version = "1.9"

  settings = <<SETTINGS
    {
        "commandToExecute": "powershell.exe -Command \"${local.dc2powershell_command}\""
    }
SETTINGS
}
