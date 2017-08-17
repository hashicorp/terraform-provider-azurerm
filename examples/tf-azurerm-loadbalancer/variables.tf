variable "prefix" {
  description = "(Required) Default prefix to use with your resource names."
  default = "myapp"
}

variable "location" {
  description = "(Required) The location/region where the core network will be created. The full list of Azure regions can be found at https://azure.microsoft.com/regions"
  default = "West US 2"
}

variable "number_of_endpoints" {
  description = "Number of nics and associated inbound remote access rules to create which are load balanced on the vnet"
  default = 3
}

variable "remote_port"{ 
  description = "Protocols to be used for remote vm access. [protocol, backend_port].  Frontend port will be automatically generated starting at 50000 and in the output." 
  default = { 
      ssh = ["Tcp", "22"]
    }
}

variable "lb_port"{ 
  description = "Protocols to be used for lb health probes and rules. [frontend_port, protocol, backend_port]"
  default = { 
      http = ["80", "Tcp", "80"]
      https = ["443", "Tcp", "443"]
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