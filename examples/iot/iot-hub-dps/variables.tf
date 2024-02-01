# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  type        = string
  description = "Value to be prepended to all resources in this example"
}

variable "iot_hub_connection_string" {
  type        = string
  description = "Connection string for the IoT Hub"
}
