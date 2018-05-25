variable "policy_definition_name" {
  description = "Policy definition name must only contain lowercase letters, digits or dashes, cannot use dash as the first two or last one characters, cannot contain consecutive dashes, and is limited between 2 and 60 characters in length."
  default     = "demoPolicy"
}

variable "policy_type" {
  description = "Valid values are 'BuiltIn', 'Custom' and 'NotSpecified'."
  default     = "Custom"
}

variable "mode" {
  description = "Valid values are 'All', 'Indexed' and 'NotSpecified'."
  default     = "All"
}

variable "display_name" {
  description = "Policy display name must only contain lowercase letters, digits or dashes, cannot use dash as the first two or last one characters, cannot contain consecutive dashes, and is limited between 2 and 60 characters in length."
  default     = "demoPolicy"
}
