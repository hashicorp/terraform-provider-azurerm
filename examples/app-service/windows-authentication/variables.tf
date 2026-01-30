# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  type        = string
  description = "The prefix used for all resources in this example"
}

variable "location" {
  type        = string
  description = "The Azure location where all resources in this example should be created"
}
