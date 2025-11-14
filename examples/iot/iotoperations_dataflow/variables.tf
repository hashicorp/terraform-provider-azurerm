variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
  default     = "rg-iotoperations-dataflow"
}

variable "location" {
  description = "The Azure region where resources will be created"
  type        = string
  default     = "East US 2"
}

variable "instance_name" {
  description = "The name of the IoT Operations instance"
  type        = string
  default     = "iotops-instance-dataflow"
}

variable "custom_location_id" {
  description = "The resource ID of the custom location (Arc-enabled Kubernetes cluster)"
  type        = string
}

variable "tags" {
  description = "A mapping of tags to assign to the resources"
  type        = map(string)
  default     = {}
}

variable "dataflow_profile_name" {
  description = "The name of the dataflow profile"
  type        = string
  default     = "dataflow-profile-example"
}

variable "dataflow_profile_instance_count" {
  description = "The number of dataflow profile instances"
  type        = number
  default     = 1
}

variable "dataflow_profile_log_level" {
  description = "The log level for the dataflow profile"
  type        = string
  default     = "info"
  validation {
    condition     = contains(["trace", "debug", "info", "warn", "error"], var.dataflow_profile_log_level)
    error_message = "Log level must be one of: trace, debug, info, warn, error."
  }
}

variable "dataflow_profile_prometheus_port" {
  description = "The Prometheus metrics port for the dataflow profile"
  type        = number
  default     = 9090
}

variable "dataflow_name" {
  description = "The name of the dataflow"
  type        = string
  default     = "dataflow-example"
}

variable "dataflow_mode" {
  description = "The mode of the dataflow (Enabled, Disabled)"
  type        = string
  default     = "Enabled"
  validation {
    condition     = contains(["Enabled", "Disabled"], var.dataflow_mode)
    error_message = "Dataflow mode must be either 'Enabled' or 'Disabled'."
  }
}

variable "dataflow_sources" {
  description = "List of dataflow sources"
  type = list(object({
    name         = string
    endpoint_ref = string
    asset_ref    = optional(string)
    schema_ref   = optional(string)
    serialization = optional(object({
      format = string
    }))
  }))
  default = [
    {
      name         = "source-mqtt"
      endpoint_ref = "mqtt-endpoint"
      asset_ref    = "temperature-asset"
      schema_ref   = "temperature-schema"
      serialization = {
        format = "Json"
      }
    }
  ]
}

variable "dataflow_destinations" {
  description = "List of dataflow destinations"
  type = list(object({
    name         = string
    endpoint_ref = string
    schema_ref   = optional(string)
    serialization = optional(object({
      format = string
    }))
  }))
  default = [
    {
      name         = "destination-adx"
      endpoint_ref = "adx-endpoint"
      schema_ref   = "processed-schema"
      serialization = {
        format = "Json"
      }
    }
  ]
}

variable "dataflow_transformations" {
  description = "List of dataflow transformations"
  type = list(object({
    type = string
    filter = optional(object({
      expression = string
      type      = string
    }))
    map = optional(object({
      expression = string
      type      = string
    }))
  }))
  default = [
    {
      type = "filter"
      filter = {
        expression = "temperature > 20"
        type      = "condition"
      }
    },
    {
      type = "map"
      map = {
        expression = "{ temp_celsius: temperature, timestamp: $metadata.timestamp }"
        type      = "newProperties"
      }
    }
  ]
}

variable "dataflow_operations" {
  description = "List of dataflow operations"
  type = list(object({
    name                = string
    operation_type      = string
    source_name         = string
    destination_name    = string
    built_in_transformation = optional(object({
      serialize_type = string
      schema_ref     = string
    }))
  }))
  default = [
    {
      name             = "operation-transform"
      operation_type   = "source"
      source_name      = "source-mqtt"
      destination_name = "destination-adx"
      built_in_transformation = {
        serialize_type = "Json"
        schema_ref     = "processed-schema"
      }
    }
  ]
}

variable "schema_registry_ref" {
  description = "The resource ID of the schema registry to associate with the instance."
  type        = string
}

