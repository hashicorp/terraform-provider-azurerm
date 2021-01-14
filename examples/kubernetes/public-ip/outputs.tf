output "cluster_egress_ip" {
  value = data.azurerm_public_ip.example.ip_address
}
