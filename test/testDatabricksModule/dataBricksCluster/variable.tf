variable "databrick_workspace_id" {
  description = "workspace for the cluster to be built"
  type = string
}

variable "databrick_workspace_URL" {
  description = "workspace URL of the workspace"
  type = string
}

variable "databrick_cluster_depends_on" {
  type = any
  default = []
}