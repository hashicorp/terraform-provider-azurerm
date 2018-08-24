variable "resource_group" {
  description = "The name of the resource group in which to create the Service Bus"
  default     = "tfex-datalake"
}

variable "location" {
  description = "The location/region where the Service Bus is created. Changing this forces a new resource to be created."
  default     = "westeurope"
}
