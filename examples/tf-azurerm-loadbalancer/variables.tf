variable "location" {
  description = "(Required) The location/region where the core network will be created. The full list of Azure regions can be found at https://azure.microsoft.com/regions"
}

variable "resource_group_name" {
  description = "(Required) The name of the resource group where the load balancer resources will be placed."
  default = "azure_lb-rg"
}

variable "prefix" {
  description = "(Required) Default prefix to use with your resource names."
  default = "azure_lb"
}

variable "number_of_endpoints" {
  description = "Number of inbound remote access rules to create which are load balanced on the vnet"
  default = 2
}

variable "remote_port"{ 
  description = "Protocols to be used for remote vm access. [protocol, backend_port].  Frontend port will be automatically generated starting at 50000 and in the output." 
  default = { 
    }
}

variable "lb_port"{ 
  description = "Protocols to be used for lb health probes and rules. [frontend_port, protocol, backend_port]"
  default = { 
    }
}

variable "lb_probe_unhealthy_threshold"{
  description = "Number of times the load balancer health probe has an unsuccessful attempt before considering the endpoint unhealthy."
  default = 2
}

variable "lb_probe_interval"{
  description = "Interval in seconds the load balancer health probe rule does a check"
  default = 5
}

variable "frontend_name" {
  description = "(Required) Specifies the name of the frontend ip configuration."
  default = "myPublicIP"
}

variable "public_ip_address_allocation" {
  description = "(Required) Defines how an IP address is assigned. Options are Static or Dynamic."
  default = "static"
}

variable "tags" {
  type = "map"
  default = {
    source = "terraform"
  }
}