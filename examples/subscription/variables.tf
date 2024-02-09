# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  type        = string
  description = "Prefix for the name of a resource"
}

variable "subscription_billing_scope_id" {
  type        = string
  description = "Subscription billing scope id"
}
