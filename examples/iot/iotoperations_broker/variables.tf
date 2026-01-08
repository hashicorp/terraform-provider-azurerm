# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "prefix" {
  description = "The prefix used for all resources in this example"
  type        = string
}

variable "resource_group_name" {
  description = "The name of an existing resource group where resources will be created"
  type        = string
}

variable "instance_name" {
  description = "The name of the existing IoT Operations instance"
  type        = string
}

variable "custom_location_id" {
  description = "The ARM resource ID of the Custom Location (Arc-enabled Kubernetes cluster)"
  type        = string
}

variable "broker_name" {
  description = "The name of the IoT Operations broker"
  type        = string
}