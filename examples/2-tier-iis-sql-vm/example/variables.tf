#################################################################
#   Variables
#################################################################

# Provider info
variable subscription_id {}

variable client_id {}
variable client_secret {}
variable tenant_id {}

# Generic info
variable location {}

variable resource_group_name {}
variable environment_name {}

# Network
variable address_space {}

variable dns_servers {
  type = "list"
}

variable wafsubnet_name {}
variable wafsubnet_prefix {}
variable rpsubnet_name {}
variable rpsubnet_prefix {}
variable issubnet_name {}
variable issubnet_prefix {}
variable dbsubnet_name {}
variable dbsubnet_prefix {}
variable dcsubnet_name {}
variable dcsubnet_prefix {}

# Active Directory & Domain Controller
variable prefix {}
variable private_ip_address {}
variable admin_username {}
variable admin_password {}

# IIS Servers
variable vmcount {}

# Domain Controller 2
variable "dc2private_ip_address" {}
variable "domainadmin_username" {}

# SQL LB
variable "lbprivate_ip_address" {}
# SQL DB Servers
variable sqlvmcount {}
