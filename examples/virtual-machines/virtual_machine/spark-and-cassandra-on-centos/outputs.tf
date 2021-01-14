output "resource_group" {
  value = "${var.resource_group}"
}

output "primary_ip_address" {
  value = "${azurerm_public_ip.primary.ip_address}"
}

output "primary_ssh_command" {
  value = "ssh ${var.vm_admin_username}@${azurerm_public_ip.primary.ip_address}"
}

output "primary_web_ui_public_ip" {
  value = "${azurerm_public_ip.primary.ip_address}:8080"
}
