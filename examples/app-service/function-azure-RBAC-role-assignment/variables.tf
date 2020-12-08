variable "prefix" {
  description = "The prefix used for all resources in this example"
}

variable "location" {
  description = "The Azure location where all resources in this example should be created"
}

variable "role_definition_name" {
  description = "Desired role to assign your function (Reader, Contributor, Owner, etc.)"  
}