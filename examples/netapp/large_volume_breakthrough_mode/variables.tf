# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

variable "location" {
  description = "The Azure location where all resources in this example should be created."
}

variable "prefix" {
  description = "The prefix used for all resources used by this NetApp Volume"
}

variable "pool_size_in_tb" {
  description = "The size of the capacity pool in TiB, it must be large enough to hold the Breakthrough Mode volume."
  default     = 4
}

variable "storage_quota_in_gb" {
  description = "The size of the Breakthrough Mode volume in GiB, valid values are between 2400 (2,400 GiB) and 2457600 (2,400 TiB)."
  default     = 2400
}

variable "throughput_in_mibps" {
  description = "The throughput assigned to the volume in MiB/s, it cannot exceed the throughput available in the manual QoS capacity pool."
  default     = 38.4
}
