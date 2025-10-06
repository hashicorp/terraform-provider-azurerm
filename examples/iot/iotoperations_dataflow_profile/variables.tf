variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
  default     = "rg-iotoperations-dataflow-profiles"
}

variable "location" {
  description = "The Azure region where resources will be created"
  type        = string
  default     = "East US 2"
}

variable "instance_name" {
  description = "The name of the IoT Operations instance"
  type        = string
  default     = "iotops-instance-profiles"
}

# High Performance Profile Variables
variable "high_performance_profile_name" {
  description = "The name of the high-performance dataflow profile"
  type        = string
  default     = "high-performance-profile"
}

variable "high_performance_instance_count" {
  description = "Number of instances for high-performance profile"
  type        = number
  default     = 4
  validation {
    condition     = var.high_performance_instance_count >= 1 && var.high_performance_instance_count <= 10
    error_message = "Instance count must be between 1 and 10."
  }
}

variable "high_performance_log_level" {
  description = "Log level for high-performance profile"
  type        = string
  default     = "warn"
  validation {
    condition     = contains(["trace", "debug", "info", "warn", "error"], var.high_performance_log_level)
    error_message = "Log level must be one of: trace, debug, info, warn, error."
  }
}

variable "high_performance_prometheus_port" {
  description = "Prometheus metrics port for high-performance profile"
  type        = number
  default     = 9090
}

variable "high_performance_self_check_mode" {
  description = "Self-check mode for high-performance profile"
  type        = string
  default     = "Enabled"
  validation {
    condition     = contains(["Enabled", "Disabled"], var.high_performance_self_check_mode)
    error_message = "Self-check mode must be either 'Enabled' or 'Disabled'."
  }
}

variable "high_performance_self_check_interval" {
  description = "Self-check interval in seconds for high-performance profile"
  type        = number
  default     = 30
}

variable "high_performance_self_check_timeout" {
  description = "Self-check timeout in seconds for high-performance profile"
  type        = number
  default     = 15
}

# Standard Profile Variables
variable "standard_profile_name" {
  description = "The name of the standard dataflow profile"
  type        = string
  default     = "standard-profile"
}

variable "standard_instance_count" {
  description = "Number of instances for standard profile"
  type        = number
  default     = 2
  validation {
    condition     = var.standard_instance_count >= 1 && var.standard_instance_count <= 10
    error_message = "Instance count must be between 1 and 10."
  }
}

variable "standard_log_level" {
  description = "Log level for standard profile"
  type        = string
  default     = "info"
  validation {
    condition     = contains(["trace", "debug", "info", "warn", "error"], var.standard_log_level)
    error_message = "Log level must be one of: trace, debug, info, warn, error."
  }
}

variable "standard_prometheus_port" {
  description = "Prometheus metrics port for standard profile"
  type        = number
  default     = 9091
}

variable "standard_self_check_mode" {
  description = "Self-check mode for standard profile"
  type        = string
  default     = "Enabled"
  validation {
    condition     = contains(["Enabled", "Disabled"], var.standard_self_check_mode)
    error_message = "Self-check mode must be either 'Enabled' or 'Disabled'."
  }
}

variable "standard_self_check_interval" {
  description = "Self-check interval in seconds for standard profile"
  type        = number
  default     = 60
}

variable "standard_self_check_timeout" {
  description = "Self-check timeout in seconds for standard profile"
  type        = number
  default     = 30
}

# Edge Profile Variables
variable "edge_profile_name" {
  description = "The name of the edge dataflow profile"
  type        = string
  default     = "edge-profile"
}

variable "edge_instance_count" {
  description = "Number of instances for edge profile"
  type        = number
  default     = 1
  validation {
    condition     = var.edge_instance_count >= 1 && var.edge_instance_count <= 10
    error_message = "Instance count must be between 1 and 10."
  }
}

variable "edge_log_level" {
  description = "Log level for edge profile"
  type        = string
  default     = "error"
  validation {
    condition     = contains(["trace", "debug", "info", "warn", "error"], var.edge_log_level)
    error_message = "Log level must be one of: trace, debug, info, warn, error."
  }
}

variable "edge_prometheus_port" {
  description = "Prometheus metrics port for edge profile"
  type        = number
  default     = 9092
}

variable "edge_self_check_mode" {
  description = "Self-check mode for edge profile"
  type        = string
  default     = "Enabled"
  validation {
    condition     = contains(["Enabled", "Disabled"], var.edge_self_check_mode)
    error_message = "Self-check mode must be either 'Enabled' or 'Disabled'."
  }
}

variable "edge_self_check_interval" {
  description = "Self-check interval in seconds for edge profile"
  type        = number
  default     = 120
}

variable "edge_self_check_timeout" {
  description = "Self-check timeout in seconds for edge profile"
  type        = number
  default     = 60
}

# Development Profile Variables
variable "create_development_profile" {
  description = "Whether to create a development profile"
  type        = bool
  default     = false
}

variable "development_profile_name" {
  description = "The name of the development dataflow profile"
  type        = string
  default     = "development-profile"
}

variable "development_instance_count" {
  description = "Number of instances for development profile"
  type        = number
  default     = 1
  validation {
    condition     = var.development_instance_count >= 1 && var.development_instance_count <= 10
    error_message = "Instance count must be between 1 and 10."
  }
}

variable "development_log_level" {
  description = "Log level for development profile"
  type        = string
  default     = "debug"
  validation {
    condition     = contains(["trace", "debug", "info", "warn", "error"], var.development_log_level)
    error_message = "Log level must be one of: trace, debug, info, warn, error."
  }
}

variable "development_prometheus_port" {
  description = "Prometheus metrics port for development profile"
  type        = number
  default     = 9093
}

variable "development_self_check_mode" {
  description = "Self-check mode for development profile"
  type        = string
  default     = "Enabled"
  validation {
    condition     = contains(["Enabled", "Disabled"], var.development_self_check_mode)
    error_message = "Self-check mode must be either 'Enabled' or 'Disabled'."
  }
}

variable "development_self_check_interval" {
  description = "Self-check interval in seconds for development profile"
  type        = number
  default     = 30
}

variable "development_self_check_timeout" {
  description = "Self-check timeout in seconds for development profile"
  type        = number
  default     = 15
}

variable "tags" {
  description = "A mapping of tags to assign to the resources"
  type        = map(string)
  default = {
    Environment = "Example"
    Purpose     = "IoT Operations Dataflow Profiles Demo"
  }
}