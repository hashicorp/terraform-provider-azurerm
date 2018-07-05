variable "prefix" {
  description = "The Prefix used for all resources in this HDInsight Cluster."
  default     = "tomdemo"
}

variable "location" {
  description = "The Azure Region in which the HDInsight Cluster should be provisioned,"
  default     = "West Europe"
}

variable "username" {
  description = "The Username used for the admin account on the Ambari server."
  default     = "sshuser"
}

variable "password" {
  description = "The Password associated with the admin account on the Ambari server."
  default     = "Password123!"
}

variable "tags" {
  type        = "map"
  description = "Any tags which should be applied to the HDInsight Cluster."
  default     = {}
}
