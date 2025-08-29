# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  description = "The prefix which should be used for all resources in this example"
  type        = string
  default     = "nfsv41-id-domain"
}

variable "location" {
  description = "The Azure Region in which all resources in this example should be created."
  type        = string
  default     = "WestUS3"
}

variable "nfsv4_id_domain" {
  description = "The NFSv4 ID domain for the NetApp account"
  type        = string
  default     = "example.com"
}
