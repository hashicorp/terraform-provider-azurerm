variable "location" {
  type        = string
  description = "Azure region where the DNS zone resource group will be created."
  default = "westeurope"
}

variable "resource_group_name" {
  type        = string
  description = "Resource group name to create/use for the Azure DNS zone."
  default = "acctest-afdx-dns-prereq"
}

variable "dns_zone_name" {
  type        = string
  description = "The public DNS zone name to create in Azure DNS, e.g. 'acctest.example.com'."

  validation {
    condition = (
      length(trimspace(var.dns_zone_name)) > 0 &&
      length(regexall("\\.", trimspace(var.dns_zone_name))) >= 1 &&
      !endswith(trimspace(var.dns_zone_name), ".")
    )
    error_message = "dns_zone_name must be a fully-qualified DNS zone name with at least two labels, e.g. 'acctest.example.com'."
  }
}

variable "tags" {
  type        = map(string)
  description = "Optional tags to apply to the Azure DNS zone."
  default     = {}
}
