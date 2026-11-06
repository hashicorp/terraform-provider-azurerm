change "resource-fix" {
  body = "`azurerm_vpn_gateway_connection` - fix a perpetual diff on `vpn_link.shared_key` by retrieving the shared key from the VPN link connection shared key API"
}
