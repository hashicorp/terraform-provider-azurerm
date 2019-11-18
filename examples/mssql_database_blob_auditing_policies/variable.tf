variable "prefix" {
  description = "The prefix which should be used for all resources in this example"
}

variable "location" {
  description = "The Azure Region in which all resources in this example should be created."
}

variable "login" {
  description = "The login of the sql server to be created"
}

variable "login_pwd" {
  description = "The login password of the sql server to be created"
}