# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

variable "location" {
  description = "The Azure region in which to create the example resources."
  type        = string
  default     = "eastus2"
}

variable "prefix" {
  description = "The prefix used for the example resource names."
  type        = string
  default     = "example"

  validation {
    condition     = can(regex("^[a-z0-9]{4,12}$", var.prefix))
    error_message = "The prefix must contain between 4 and 12 lowercase letters or numbers."
  }
}