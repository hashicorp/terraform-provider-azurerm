variable "resource_group_name" {
  type        = "string"
  description = "Name of the azure resource group."
  default     = "tfex-batch"
}

variable "location" {
  type        = "string"
  description = "Location of the azure resource group."
  default     = "westeurope"
}

variable "storage_account_tier" {
  description = "Defines the Tier of storage account to be created. Valid options are Standard and Premium."
  default     = "Standard"
}

variable "storage_replication_type" {
  description = "Defines the Replication Type to use for this storage account. Valid options include LRS, GRS etc."
  default     = "LRS"
}

variable "batch_pool_nodes_vm_size" {
  description = "Size of the virtual machines in the nodes pool"
  default     = "Standard_A1_v2"
}

variable "batch_pool_nodes_count" {
  description = "Number of virtual machines in the nodes pool"
  default     = 2
}

variable "batch_pool_nodes_vm_image" {
  description = "Virtual machine image to use for the nodes in the pool"
  default     = "canonical:ubuntuserver:16.04-LTS"
}

variable "batch_pool_nodes_agent_sku_id" {
  description = "Agent Sku to use for the nodes in the pool"
  default     = "batch.node.ubuntu 16.04"
}



