# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  type        = string
  description = "The prefix used for all resources in this example"
}

variable "location" {
  type        = string
  description = "The Azure location where all resources in this example should be created"
}

variable "role_definition_name" {
  type        = string
  default     = "Reader"
  description = "Desired role to assign your function (Reader, Contributor, Owner, etc.)"
}