variable "resource_group" {
  description = "The name of the resource group in which to create the Event Hub"
  default     = "tfex-eventhub"
}

variable "location" {
  description = "The location/region where the Event Hub is created. Changing this forces a new resource to be created."
  default     = "southcentralus"
}
