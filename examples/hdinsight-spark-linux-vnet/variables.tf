variable "resourceGroupName" {
  description = "Name of the resource group container for all resources"
}

variable "resourceGroupLocation" {
  description = "Azure region used for resource deployment"
}

variable "clusterName" {
  description = "The name of the HDInsight cluster to create."
}

variable "clusterLoginUserName" {
  description = "These credentials can be used to submit jobs to the cluster and to log into cluster dashboards."
  default     = "admin"
}

variable "clusterLoginPassword" {
  description = "The password must be at least 10 characters in length and must contain at least one digit, one non-alphanumeric character, and one upper or lower case letter."
}

variable "sshUserName" {
  description = "These credentials can be used to remotely access the cluster."
  default     = "sshuser"
}

variable "sshPassword" {
  description = "The password must be at least 10 characters in length and must contain at least one digit, one non-alphanumeric character, and one upper or lower case letter."
}
