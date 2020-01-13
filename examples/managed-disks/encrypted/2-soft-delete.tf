resource "null_resource" "enable-soft-delete-and-purge-protection" {
  # Soft Delete and Purge Protection aren't available in the Azure Provider at this time
  # as such we'll use the Azure CLI to enable them for the moment
  # TODO: fix in 2.0 once these become available
  provisioner "local-exec" {
    command = "az keyvault update --name ${azurerm_key_vault.test.name} --resource-group ${azurerm_key_vault.test.resource_group_name} --enable-soft-delete true --enable-purge-protection true"
  }
}