output "domain_controllers_subnet_id" {
  value = azurerm_subnet.domain-controllers.id
}

output "domain_members_subnet_id" {
  value = azurerm_subnet.domain-members.id
}