output "service_id" {
  description = "The ID of the API Management Service created"
  value       = azurerm_api_management.apim_service.id
}

output "gateway_url" {
  description = "The URL of the Gateway for the API Management Service"
  value       = azurerm_api_management.apim_service.gateway_url
}

output "service_public_ip_addresses" {
  description = "The Public IP addresses of the API Management Service"
  value       = azurerm_api_management.apim_service.public_ip_addresses
}

output "api_outputs" {
  description = "The IDs, state, and version outputs of the APIs created"
  value = {
      id             = azurerm_api_management_api.api.id
      is_current     = azurerm_api_management_api.api.is_current
      is_online      = azurerm_api_management_api.api.is_online
      version        = azurerm_api_management_api.api.version
      version_set_id = azurerm_api_management_api.api.version_set_id
    }
}

output "group_id" {
  description = "The ID of the API Management Group created"
  value = azurerm_api_management_group.group.id
}

output "product_ids" {
  description = "The ID of the Product created"
  value = azurerm_api_management_product.product.id
}

output "product_api_ids" {
  description = "The ID of the Product/API association created"
  value       = azurerm_api_management_product_api.product_api.id
}

output "product_group_ids" {
  description = "The ID of the Product/Group association created"
  value       = azurerm_api_management_product_group.product_group.id
}
