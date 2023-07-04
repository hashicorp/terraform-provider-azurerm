# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  description = "The prefix used for all resources in this example"
}

variable "location" {
  description = "The Azure location where all resources in this example should be created"
}

variable "user_name" {
  description = "The user name of virtual machine"
}

variable "password" {
  description = "The password of virtual machine"
  sensitive   = true
}

# Refer to https://github.com/Azure/azure-cli-extensions/blob/ed3f463e9ef7980eff196504a8bb29800c123eba/src/connectedk8s/azext_connectedk8s/custom.py#L365 to generate the private key
variable "private_pem" {
  description = "The private certificate used by the agent"
}

# Refer to https://github.com/Azure/azure-cli-extensions/blob/ed3f463e9ef7980eff196504a8bb29800c123eba/src/connectedk8s/azext_connectedk8s/custom.py#L359 to generate the public key
variable "public_key" {
  description = "The base64-encoded public certificate used by the agent"
}