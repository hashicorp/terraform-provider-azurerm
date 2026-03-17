variable "primary_location" {
  type        = string
  default     = "East US"
  description = "Primary Azure region for resources"
}

variable "nfs_host" {
  type        = string
  default     = "192.168.0.1"
  description = "Host for NFS source endpoint"
}

variable "smb_host" {
  type        = string
  default     = "192.168.0.2"
  description = "Host for SMB mount endpoint (use real SMB server for apply)"
}

variable "smb_share_name" {
  type        = string
  default     = "share1"
  description = "SMB share name for mount endpoint"
}
