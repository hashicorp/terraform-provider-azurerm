# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "location" {
  description = "The Azure location where all resources in this example should be created."
  default     = "East US"
}

variable "prefix" {
  description = "The prefix used for all resources used by this NetApp short-term clone example"
  default     = "netapp-shortclone"
}
