# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

variable "location" {
  description = "The Azure location where all resources in this example should be created."
}

variable "prefix" {
  description = "The prefix used for all resources used by this NetApp Volume Bucket (Key Vault) example."
}

variable "server_fqdn" {
  description = "The DNS name that resolves to the bucket endpoint IP address. Must match the Subject Alternative Name of the bucket server certificate."
}
