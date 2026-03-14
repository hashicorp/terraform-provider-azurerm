variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
  default     = "rg-iotoperations-endpoints"
}

variable "location" {
  description = "The Azure region where resources will be created"
  type        = string
  default     = "East US 2"
}

variable "instance_name" {
  description = "The name of the IoT Operations instance"
  type        = string
  default     = "iotops-instance-endpoints"
}

# MQTT Endpoint Variables
variable "mqtt_endpoint_name" {
  description = "The name of the MQTT dataflow endpoint"
  type        = string
  default     = "mqtt-endpoint"
}

variable "mqtt_host" {
  description = "The MQTT broker host"
  type        = string
  default     = "mqtt-broker.iotoperations.svc.cluster.local"
}

variable "mqtt_port" {
  description = "The MQTT broker port"
  type        = number
  default     = 1883
}

variable "mqtt_authentication_enabled" {
  description = "Enable MQTT authentication"
  type        = bool
  default     = false
}

variable "mqtt_auth_method" {
  description = "MQTT authentication method"
  type        = string
  default     = "UsernamePassword"
}

variable "mqtt_username" {
  description = "MQTT username"
  type        = string
  default     = ""
}

variable "mqtt_password_secret" {
  description = "Name of the secret containing MQTT password"
  type        = string
  default     = ""
}

variable "mqtt_tls_enabled" {
  description = "Enable TLS for MQTT connection"
  type        = bool
  default     = false
}

variable "mqtt_tls_mode" {
  description = "TLS mode for MQTT connection"
  type        = string
  default     = "Enabled"
}

variable "mqtt_tls_ca_cert_config_map" {
  description = "ConfigMap containing CA certificate for MQTT TLS"
  type        = string
  default     = ""
}

variable "mqtt_keep_alive_seconds" {
  description = "MQTT keep alive interval in seconds"
  type        = number
  default     = 60
}

variable "mqtt_retain" {
  description = "Enable MQTT message retention"
  type        = bool
  default     = false
}

variable "mqtt_session_expiry_seconds" {
  description = "MQTT session expiry in seconds"
  type        = number
  default     = 3600
}

variable "mqtt_qos" {
  description = "MQTT Quality of Service level"
  type        = number
  default     = 1
  validation {
    condition     = contains([0, 1, 2], var.mqtt_qos)
    error_message = "MQTT QoS must be 0, 1, or 2."
  }
}

# Azure Data Explorer Endpoint Variables
variable "adx_endpoint_name" {
  description = "The name of the Azure Data Explorer dataflow endpoint"
  type        = string
  default     = "adx-endpoint"
}

variable "adx_cluster_uri" {
  description = "Azure Data Explorer cluster URI"
  type        = string
  default     = "https://mycluster.eastus2.kusto.windows.net"
}

variable "adx_database_name" {
  description = "Azure Data Explorer database name"
  type        = string
  default     = "iottelemetry"
}

variable "adx_audience" {
  description = "Azure Data Explorer audience for authentication"
  type        = string
  default     = "https://kusto.windows.net"
}

variable "adx_batching_latency" {
  description = "ADX batching latency in seconds"
  type        = number
  default     = 5
}

variable "adx_batching_max_messages" {
  description = "ADX batching maximum messages"
  type        = number
  default     = 1000
}

# Azure Storage Endpoint Variables
variable "storage_endpoint_name" {
  description = "The name of the Azure Storage dataflow endpoint"
  type        = string
  default     = "storage-endpoint"
}

variable "storage_account_host" {
  description = "Azure Storage account host"
  type        = string
  default     = "https://mystorageaccount.blob.core.windows.net"
}

variable "storage_container_name" {
  description = "Azure Storage container name"
  type        = string
  default     = "iotdata"
}

variable "storage_audience" {
  description = "Azure Storage audience for authentication"
  type        = string
  default     = "https://storage.azure.com"
}

variable "storage_batching_latency" {
  description = "Storage batching latency in seconds"
  type        = number
  default     = 60
}

variable "storage_batching_max_messages" {
  description = "Storage batching maximum messages"
  type        = number
  default     = 10000
}

# Local Storage Endpoint Variables
variable "local_endpoint_name" {
  description = "The name of the local storage dataflow endpoint"
  type        = string
  default     = "local-endpoint"
}

variable "local_storage_pvc_name" {
  description = "Name of the Persistent Volume Claim for local storage"
  type        = string
  default     = "iot-local-storage"
}

# Fabric OneLake Endpoint Variables
variable "enable_fabric_endpoint" {
  description = "Enable Fabric OneLake endpoint"
  type        = bool
  default     = false
}

variable "fabric_endpoint_name" {
  description = "The name of the Fabric OneLake dataflow endpoint"
  type        = string
  default     = "fabric-endpoint"
}

variable "fabric_host" {
  description = "Fabric OneLake host"
  type        = string
  default     = "https://onelake.dfs.fabric.microsoft.com"
}

variable "fabric_workspace_id" {
  description = "Fabric workspace ID"
  type        = string
  default     = ""
}

variable "fabric_lakehouse_name" {
  description = "Fabric lakehouse name"
  type        = string
  default     = "iotlakehouse"
}

variable "fabric_audience" {
  description = "Fabric OneLake audience for authentication"
  type        = string
  default     = "https://onelake.dfs.fabric.microsoft.com"
}

variable "fabric_batching_latency" {
  description = "Fabric batching latency in seconds"
  type        = number
  default     = 30
}

variable "fabric_batching_max_messages" {
  description = "Fabric batching maximum messages"
  type        = number
  default     = 5000
}

variable "tags" {
  description = "A mapping of tags to assign to the resources"
  type        = map(string)
  default = {
    Environment = "Example"
    Purpose     = "IoT Operations Dataflow Endpoints Demo"
  }
}

variable "custom_location_id" {
  description = "The ID of the custom location to use for extended location."
  type        = string
}

variable "adx_host" {
  description = "Azure Data Explorer host"
  type        = string
}

variable "adx_database" {
  description = "Azure Data Explorer database"
  type        = string
}

variable "storage_host" {
  description = "Azure Storage host"
  type        = string
}

variable "fabric_workspace_name" {
  description = "Fabric workspace name"
  type        = string
}

variable "fabric_one_lake_path_type" {
  description = "Fabric OneLake path type"
  type        = string
}

variable "schema_registry_ref" {
  description = "The resource ID of the schema registry."
  type        = string
}