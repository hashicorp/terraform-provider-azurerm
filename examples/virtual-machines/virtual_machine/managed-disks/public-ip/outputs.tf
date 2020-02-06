output "public_ip_address" {
  value = "${azurerm_public_ip.example.ip_address}"
}

output "ssh_command" {
  value = "ssh testadmin@${azurerm_public_ip.example.ip_address}"
}
