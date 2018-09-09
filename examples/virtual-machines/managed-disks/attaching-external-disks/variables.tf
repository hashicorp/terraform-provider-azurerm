variable "prefix" {
  description = "The Prefix used for all resources in this example"
}

variable "location" {
  description = "The Azure Region in which the resources in this example should exist"
}

variable "tags" {
  type        = "map"
  default     = {}
  description = "Any tags which should be assigned to the resources in this example"
}

variable "number_of_disks" {
  description = "The number of Data Disks which should be attached"
  default     = 2
}
