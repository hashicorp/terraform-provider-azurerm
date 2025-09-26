# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  description = "The prefix which should be used for all resources in this example"
  type        = string
  default     = "example"
}

variable "location" {
  description = "The Azure region in which all resources in this example should be created"
  type        = string
  default     = "West Europe"
}
