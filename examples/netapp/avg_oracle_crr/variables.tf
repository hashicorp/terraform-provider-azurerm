# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  description = "The prefix which should be used for all resources in this example"
  type        = string
  default     = "testnaocrr"
}

variable "location" {
  description = "The Azure Region in which all resources in this example should be created"
  type        = string
  default     = "West Europe"
}

variable "alt_location" {
  description = "The alternative Azure Region for the secondary volume group"
  type        = string
  default     = "East US"
}
