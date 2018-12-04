output "public_ip_address" {
  value = "${azurerm_public_ip.main.ip_address}"
}

output "ssh_command" {
  value = "ssh ${local.admin_username}@${azurerm_public_ip.main.ip_address}"
}
