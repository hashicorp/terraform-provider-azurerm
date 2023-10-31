# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "location" {
  description = "The Azure location where all resources in this example should be created."
}

variable "prefix" {
  description = "The prefix used for all resources used by this NetApp Volume"
}

variable "tenant_id" {
  description = "The Azure tenant ID used to create the user-assigned identity"
}
