# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

variable "location" {
  description = "The Azure location where all resources in this example should be created."
  type        = string
}

variable "prefix" {
  description = "The prefix used for all resources used by this NetApp Flexible Service Level example"
  type        = string
}
