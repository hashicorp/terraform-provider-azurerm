# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  type        = string
  description = "Prefix for the name of a resource"
}

variable "parent_management_group_id" {
  type        = string
  description = "The ID of the parent management group"
}

variable "subscription_ids" {
  type        = list(string)
  description = "The list of subscription IDs to assign to the management group"
}
