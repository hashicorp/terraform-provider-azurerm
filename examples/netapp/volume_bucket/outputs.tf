# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# The IP address that backs the bucket S3 endpoint. Combine this with the
# server certificate FQDN (or your DNS record) to point S3 clients at the
# buckets.
output "bucket_server_ip_address" {
  value = azurerm_netapp_volume_bucket_with_server.first.server_ip_address
}

output "first_bucket_id" {
  value = azurerm_netapp_volume_bucket_with_server.first.id
}

output "second_bucket_id" {
  value = azurerm_netapp_volume_bucket.second.id
}
