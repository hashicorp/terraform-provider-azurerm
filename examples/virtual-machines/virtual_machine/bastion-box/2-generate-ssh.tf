# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "tls_private_key" "example" {
  algorithm = "RSA"
  rsa_bits  = "2048"
}

locals {
  public_ssh_key = tls_private_key.example.public_key_openssh
}
