variable "prefix" {
  description = "A prefix used for all resources in this example"
}

variable "location" {
  description = "The Azure Region in which all resources in this example should be provisioned"
}

variable "fwprivate_ip" {
  description = "The IP address packets should be forwarded to when using the VirtualAppliance hop type."
}

variable "service_principal_pw" {
  description = "The password for the service principal."
}