# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "cdn_endpoint_id" {
  value = "${azurerm_cdn_endpoint.cdnendpt.name}.azureedge.net"
}
