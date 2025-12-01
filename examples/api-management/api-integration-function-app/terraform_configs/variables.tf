

# these are sandbox credentials so dear hackers, don't bother ;)
variable "subscription_id" {
  default = "YOUR-SUBSCRIPTION-ID"
  type    = string
}

variable "resource_group_name" {
  default = "YOUR-RESOURCEGROUP-NAME"
  type    = string
}


variable "sku" {
  description = "The pricing tier of this API Management service"
  default     = "Developer"
  type        = string
  validation {
    condition     = contains(["Developer", "Standard", "Premium"], var.sku)
    error_message = "The sku must be one of the following: Developer, Standard, Premium."
  }
}

