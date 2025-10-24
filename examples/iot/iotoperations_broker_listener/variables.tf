variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
  default     = "rg-iotoperations-example"
}

variable "location" {
  description = "The Azure region where resources will be created"
  type        = string
  default     = "East US 2"
}

variable "instance_name" {
  description = "The name of the IoT Operations instance"
  type        = string
  default     = "iotops-instance-example"
}

variable "custom_location_id" {
  description = "The resource ID of the custom location (Arc-enabled Kubernetes cluster)"
  type        = string
}

variable "broker_name" {
  description = "The name of the IoT Operations broker"
  type        = string
  default     = "broker-example"
}

# Simple Listener Variables
variable "simple_listener_name" {
  description = "The name of the simple broker listener"
  type        = string
  default     = "simple-listener"
}

variable "simple_listener_port" {
  description = "The port number for the simple broker listener"
  type        = number
  default     = 1883
}

variable "simple_service_name" {
  description = "The service name for the simple broker listener"
  type        = string
  default     = "simple-listener-service"
}

variable "simple_service_type" {
  description = "The service type for the simple broker listener"
  type        = string
  default     = "ClusterIP"
}

# Complex Listener Variables
variable "complex_listener_name" {
  description = "The name of the complex broker listener"
  type        = string
  default     = "complex-listener"
}

variable "complex_listener_port" {
  description = "The port number for the complex broker listener"
  type        = number
  default     = 8883
}

variable "complex_node_port" {
  description = "The node port for the complex broker listener"
  type        = number
  default     = 30883
}

variable "complex_protocol" {
  description = "The protocol for the complex broker listener"
  type        = string
  default     = "Mqtt"
}

variable "complex_service_name" {
  description = "The service name for the complex broker listener"
  type        = string
  default     = "complex-listener-service"
}

variable "complex_service_type" {
  description = "The service type for the complex broker listener"
  type        = string
  default     = "LoadBalancer"
}

variable "complex_authentication_ref" {
  description = "Authentication reference for the complex broker listener"
  type        = string
  default     = ""
}

variable "complex_authorization_ref" {
  description = "Authorization reference for the complex broker listener"
  type        = string
  default     = ""
}

# Complex TLS Variables
variable "complex_tls_mode" {
  description = "TLS mode for the complex broker listener"
  type        = string
  default     = "Automatic"
}

variable "complex_tls_cert_duration" {
  description = "TLS certificate duration for complex listener"
  type        = string
  default     = "8760h"
}

variable "complex_tls_cert_secret_name" {
  description = "TLS certificate secret name for complex listener"
  type        = string
  default     = "complex-tls-secret"
}

variable "complex_tls_cert_renew_before" {
  description = "TLS certificate renew before duration for complex listener"
  type        = string
  default     = "720h"
}

variable "complex_tls_issuer_name" {
  description = "TLS certificate issuer name for complex listener"
  type        = string
  default     = "cluster-issuer"
}

variable "complex_tls_issuer_kind" {
  description = "TLS certificate issuer kind for complex listener"
  type        = string
  default     = "ClusterIssuer"
}

variable "complex_tls_issuer_group" {
  description = "TLS certificate issuer group for complex listener"
  type        = string
  default     = "cert-manager.io"
}

variable "complex_tls_private_key_algorithm" {
  description = "TLS private key algorithm for complex listener"
  type        = string
  default     = "RSA"
}

variable "complex_tls_private_key_rotation_policy" {
  description = "TLS private key rotation policy for complex listener"
  type        = string
  default     = "Always"
}

variable "complex_tls_san_dns" {
  description = "TLS SAN DNS names for complex listener"
  type        = list(string)
  default     = ["broker.example.com"]
}

variable "complex_tls_san_ip" {
  description = "TLS SAN IP addresses for complex listener"
  type        = list(string)
  default     = ["10.0.0.1"]
}

# Full Listener Variables
variable "enable_full_listener" {
  description = "Enable the full/advanced broker listener"
  type        = bool
  default     = false
}

variable "full_listener_name" {
  description = "The name of the full broker listener"
  type        = string
  default     = "full-listener"
}

variable "full_service_name" {
  description = "The service name for the full broker listener"
  type        = string
  default     = "full-listener-service"
}

variable "full_service_type" {
  description = "The service type for the full broker listener"
  type        = string
  default     = "LoadBalancer"
}

variable "full_authentication_ref" {
  description = "Authentication reference for the full broker listener"
  type        = string
  default     = ""
}

variable "full_authorization_ref" {
  description = "Authorization reference for the full broker listener"
  type        = string
  default     = ""
}

variable "tags" {
  description = "A mapping of tags to assign to the resources"
  type        = map(string)
  default = {
    Environment = "Example"
    Purpose     = "IoT Operations Broker Listener Demo"
  }
}