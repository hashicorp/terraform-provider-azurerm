variable "prefix" {
  description = "The prefix which should be used for all resources in this example"
}

variable "location" {
  description = "The Azure Region in which all resources in this example should be created. Use the `Name` value from `az account list-locations -otable` for your chosen region"
  default = "westeurope"
}
