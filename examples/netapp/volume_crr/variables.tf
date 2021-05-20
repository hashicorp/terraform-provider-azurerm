variable "location" {
  description = "The Azure location where all resources in this example should be created."
}

variable "alt_location" {
  description = "The Azure location where the secondary volume will be created."
}

variable "prefix" {
  description = "The prefix used for all resources used by this NetApp Volume"
}
