output "public_ip_address" {
  value = "${azurerm_public_ip.static.ip_address}"
}
