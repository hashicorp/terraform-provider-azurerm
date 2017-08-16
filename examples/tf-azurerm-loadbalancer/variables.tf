variable "prefix" {
  description = "Default prefix to use with your resource names."
  default = "myapp"
}

variable "location" {
  description = "The location/region where the core network will be created. The full list of Azure regions can be found at https://azure.microsoft.com/regions"
  default = "West US 2"
}

variable "frontend_name" {
  description = "(Required) Specifies the name of the frontend ip configuration."
  default = "myPublicIP"
}

variable "public_ip_address_allocation" {
  description = "Defines how a private IP address is assigned. Options are Static or Dynamic."
  default = "static"
}

variable "frontend_subnet_id" {
  description = "(Optional) Reference to subnet associated with the IP Configuration."
  default = ""
}

variable "frontend_private_ip" {
  description = "(Optional) Private IP Address to assign to the Load Balancer. The last one and first four IPs in any range are reserved and cannot be manually assigned."
  default = ""
}

variable "tags" {
  type = "map"
  default = {
    tag1 = ""
    tag2 = ""
  }
}
/* TODO: Move to provision via module not within LB module itself
variable "frontend_public_ip_address_id" {
  description = "(Optional) Reference to Public IP address to be associated with the Load Balancer."
  default = ""
}*/