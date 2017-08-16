variable "prefix" {
  description = "(Required) Default prefix to use with your resource names."
  default = "myapp"
}

variable "location" {
  description = "(Required) The location/region where the core network will be created. The full list of Azure regions can be found at https://azure.microsoft.com/regions"
  default = "West US 2"
}

variable "frontend_name" {
  description = "(Required) Specifies the name of the frontend ip configuration."
  default = "myPublicIP"
}

variable "public_ip_address_allocation" {
  description = "(Required) Defines how a private IP address is assigned. Options are Static or Dynamic."
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

variable "address_space" {
  description = "The address space that is used by the virtual network."
  default     = "10.0.0.0/16"
}

# If no values specified, this defaults to Azure DNS 
variable "dns_servers" {
  description = "The DNS servers to be used with vNet"
  default     = []
}

variable "subnet_prefixes" {
  description = "The address prefix to use for the subnet."
  default     = ["10.0.1.0/24"]
}

variable "subnet_names" {
  description = "A list of public subnets inside the vNet."
  default     = ["subnet1"]
}

variable "tags" {
  type = "map"
  default = {
    tag1 = ""
    tag2 = ""
  }
}