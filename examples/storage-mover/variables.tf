variable "primary_location" {
  type        = string
  default     = "East US"
  description = "Primary Azure region for resources"
}

variable "nfs_host" {
  type        = string
  default     = "192.168.0.1"
  description = "Host name or IP for NFS source endpoint (use a real NFS server for apply)"
}

# Uncomment when testing multi-cloud connector endpoint
# variable "multi_cloud_connector_id" {
#   type        = string
#   description = "Resource ID of Microsoft.HybridConnectivity/publicCloudConnectors resource"
# }
# variable "aws_s3_bucket_id" {
#   type        = string
#   description = "Resource ID of Microsoft.AwsConnector/s3Buckets resource"
# }
