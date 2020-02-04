output "domain_controllers_subnet_id" {
  value = "${azurerm_subnet.domain-controllers.id}"
}

output "domain_clients_subnet_id" {
  value = "${azurerm_subnet.domain-clients.id}"
}
