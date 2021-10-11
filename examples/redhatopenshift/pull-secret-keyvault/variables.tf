variable "prefix" {
  description = "A prefix used for all resources in this example"
}

variable "location" {
  description = "The Azure Region in which all resources in this example should be provisioned"
}

variable "client_id" {
  description = "Service principal client ID"
}

variable "client_secret" {
  description = "Service principal client secret"
  sensitive   = true
}

variable "key_vault_name" {
  description 		= "Name of Azure Key Vault that hosts the Red Hat pull secret"
}

variable "key_vault_resource_group_name" {
  description 		= "Name of the resource group where the Azure Key Vault is hosted"
}
