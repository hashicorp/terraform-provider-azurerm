# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# The S3 access key returned by Azure NetApp Files when the credentials are
# generated inline (i.e. `store_in_key_vault = false`).
output "bucket_access_key" {
  value     = azurerm_netapp_volume_bucket_credentials.example.access_key
  sensitive = true
}

# The S3 secret key returned by Azure NetApp Files when the credentials are
# generated inline.
output "bucket_secret_key" {
  value     = azurerm_netapp_volume_bucket_credentials.example.secret_key
  sensitive = true
}

# The expiry date of the credential pair. After this date the keys stop
# working and the bucket reports `status = "CredentialsExpired"`. Rotate by
# tainting `azurerm_netapp_volume_bucket_credentials.example`.
output "bucket_key_pair_expiry" {
  value = azurerm_netapp_volume_bucket_credentials.example.key_pair_expiry
}

# The IP address that backs the bucket S3 endpoint. Combine this with the
# server certificate FQDN (or your DNS record) to point S3 clients at the
# bucket.
output "bucket_server_ip_address" {
  value = azurerm_netapp_volume_bucket.example.server_ip_address
}

# Read the sensitive values from the CLI with:
#   terraform output -raw bucket_access_key
#   terraform output -raw bucket_secret_key
