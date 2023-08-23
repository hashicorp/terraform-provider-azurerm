# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "hostname" {
  value = "${var.hostname}"
}

output "vm_fqdn" {
  value = "${azurerm_public_ip.lbpip.fqdn}"
}

output "vms_rdp_access" {
  value = "${formatlist("RDP_URL=%v:%v", azurerm_public_ip.lbpip.fqdn, azurerm_lb_nat_rule.tcp.*.frontend_port)}"
}
