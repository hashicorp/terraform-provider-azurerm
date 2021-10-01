variable "ase_resource_group_name" {
  type    = string
}

variable "use_existing_vnet_and_subnet" {
  type    = bool
  default = false 
}

variable "vnet_resource_group_name" {
  type    = string
}

variable "virtual_network_name" {
  type    = string
}

variable "location" {
  type    = string
}

variable "vnet_address_prefixes" {
  type    = list(string)
  default = ["172.16.0.0/16"] 
}

variable "subnet_name" {
  type    = string
}

variable "subnet_address_prefixes" {
  type    = list(string)
  default = ["172.16.0.0/24"]
}

variable "ase_name" {
  type    = string
}

variable "dedicated_host_count" {
  type    = number 
  default = 0
}

variable "zone_redundant" {
  type    = bool
  default = false
}

variable "create_private_dns" {
  type    = bool
  default = true
}

variable "internal_load_balancing_mode" {
  type    = string
  default = "Web, Publishing"
}

variable "network_security_group_name" {
  type    = string
}

variable "network_security_group_security_rules" {
  type    = any
  default = []
}