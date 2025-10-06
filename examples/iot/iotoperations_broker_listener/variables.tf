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

variable "broker_name" {
  description = "The name of the IoT Operations broker"
  type        = string
  default     = "broker-example"
}

variable "listener_name" {
  description = "The name of the broker listener"
  type        = string
  default     = "listener-example"
}

variable "listener_port" {
  description = "The port number for the broker listener"
  type        = number
  default     = 1883
}

variable "service_name" {
  description = "The service name for the broker listener"
  type        = string
  default     = "listener-service"
}

variable "service_type" {
  description = "The service type for the broker listener (LoadBalancer, ClusterIP, NodePort)"
  type        = string
  default     = "ClusterIP"
}

variable "enable_tls" {
  description = "Enable TLS for the broker listener"
  type        = bool
  default     = false
}

variable "tls_mode" {
  description = "TLS mode for the broker listener (Automatic, Manual)"
  type        = string
  default     = "Automatic"
}

variable "tls_cert_duration" {
  description = "TLS certificate duration"
  type        = string
  default     = "8760h"
}

variable "tls_issuer_name" {
  description = "TLS certificate issuer name"
  type        = string
  default     = "cluster-issuer"
}

variable "tls_issuer_kind" {
  description = "TLS certificate issuer kind"
  type        = string
  default     = "ClusterIssuer"
}

variable "tls_issuer_group" {
  description = "TLS certificate issuer group"
  type        = string
  default     = "cert-manager.io"
}

variable "tls_cert_renew_before" {
  description = "TLS certificate renew before duration"
  type        = string
  default     = "720h"
}

variable "tls_cert_secret_name" {
  description = "TLS certificate secret name"
  type        = string
  default     = "listener-tls-secret"
}

variable "tls_subject_organization" {
  description = "TLS certificate subject organization"
  type        = string
  default     = "Contoso"
}

variable "tls_subject_organizational_unit" {
  description = "TLS certificate subject organizational unit"
  type        = string
  default     = "IoT Operations"
}

variable "authentication_ref_name" {
  description = "Name of the authentication reference (optional)"
  type        = string
  default     = ""
}

variable "authorization_ref_name" {
  description = "Name of the authorization reference (optional)"
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