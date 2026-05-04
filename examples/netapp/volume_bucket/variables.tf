# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

variable "location" {
  description = "The Azure location where all resources in this example should be created."
}

variable "prefix" {
  description = "The prefix used for all resources used by this NetApp Volume Bucket example."
}

variable "server_fqdn" {
  description = "The DNS name to embed in the bucket server certificate. The first bucket on a given volume backing IP must include this value plus a matching certificate."
}
