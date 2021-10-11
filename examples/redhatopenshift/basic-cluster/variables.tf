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
