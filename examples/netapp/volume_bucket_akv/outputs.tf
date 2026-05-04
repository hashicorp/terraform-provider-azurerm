# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

# The bucket access keys live in the credentials Key Vault, not in Terraform
# state. Read them through `azurerm_key_vault_secret` so consumers don't have
# to call the Azure CLI separately.
data "azurerm_key_vault_secret" "bucket_credentials" {
  name         = "${var.prefix}-bucket-creds"
  key_vault_id = azurerm_key_vault.credentials.id

  depends_on = [
    azurerm_netapp_volume_bucket_credentials.example,
  ]
}

# Raw JSON secret value, e.g. {"access_key_id":"...","secret_access_key":"..."}.
output "bucket_credentials_secret_json" {
  value     = data.azurerm_key_vault_secret.bucket_credentials.value
  sensitive = true
}

# Convenience outputs that pre-decode the JSON.
output "bucket_access_key" {
  value     = jsondecode(data.azurerm_key_vault_secret.bucket_credentials.value).access_key_id
  sensitive = true
}

output "bucket_secret_key" {
  value     = jsondecode(data.azurerm_key_vault_secret.bucket_credentials.value).secret_access_key
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
