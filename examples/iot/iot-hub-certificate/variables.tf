# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  type        = string
  description = "Value to be prepended to all resources in this example"
}

variable "location" {
  type        = string
  description = "Region where resources will be deployed"
  default     = "westeurope"
}

variable "iot_hub_dps_name" {
  type        = string
  description = "Name of the IoT Hub DPS instance"
}

variable "dps_cert_content" {
  type        = string
  description = "The Base-64 representation of the X509 leaf certificate .cer file or just a .pem file content"
}
