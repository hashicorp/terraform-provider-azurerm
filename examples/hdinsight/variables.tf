variable "azure_resource_group_name" {
  type        = "string"
  description = "Name of the azure resource group."
}

variable "azure_resource_group_location" {
  type        = "string"
  description = "Location of the azure resource group."
}

variable "azure_storage_account_name" {
  type = "string"
  description = "Name of the storage account"
}

variable "azure_hdinsight_cluster_name" {
  type        = "string"
  description = "HDInsight cluster name"
}

variable "azure_hdinsight_cluster_version" {
  type        = "string"
  description = "HDInsight cluster version"
}

variable "azure_hdinsight_cluster_rest_username" {
  type        = "string"
  description = "HDInsight gateway REST username"
}

variable "azure_hdinsight_cluster_rest_password" {
  type        = "string"
  description = "HDInsight gateway REST password"
}

variable "azure_hdinsight_cluster_headnode_instance_count" {
  type        = "string"
  description = "HDInsight cluster headnode instance count"
}

variable "azure_hdinsight_cluster_headnode_vmsize" {
  type        = "string"
  description = "HDInsight cluster headnode vm size"
}

variable "azure_hdinsight_cluster_headnode_username" {
  type        = "string"
  description = "HDInsight cluster headnode username"
}

variable "azure_hdinsight_cluster_headnode_password" {
  type        = "string"
  description = "HDInsight cluster headnode password"
}

variable "azure_hdinsight_cluster_workernode_instance_count" {
  type        = "string"
  description = "HDInsight cluster workernode instance count"
}

variable "azure_hdinsight_cluster_workernode_vmsize" {
  type        = "string"
  description = "HDInsight cluster workernode instance count"
}

variable "azure_hdinsight_cluster_workernode_username" {
  type        = "string"
  description = "HDInsight cluster workernode username"
}

variable "azure_hdinsight_cluster_workernode_sshkey" {
  type        = "string"
  description = "HDInsight cluster workernode SSH key"
}