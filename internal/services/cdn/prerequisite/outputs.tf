output "dns_zone_id" {
  description = "The resource ID of the created Azure DNS zone."
  value       = azurerm_dns_zone.parent.id
}

output "name_servers" {
  description = "The Azure DNS name servers you must configure as NS records at your registrar for the subdomain."
  value       = azurerm_dns_zone.parent.name_servers
}
