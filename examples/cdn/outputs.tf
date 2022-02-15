output "cdn_endpoint_id" {
  value = "${azurerm_cdn_endpoint.cdnendpt.name}.azureedge.net"
}
