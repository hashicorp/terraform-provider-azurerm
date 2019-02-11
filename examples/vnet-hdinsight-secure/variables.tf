# General 
variable "resource_prefix" {
  type = "string"
  description = "The prefix used for all resources in this example"
  default     = "hdi-secure-vnet"
}

variable "location" {
  type = "string"
  description = "The Azure Region in which the resources in this example should exist"
  default     = "West US 2"
}

variable "tags" {
  type        = "map"
  description = "Any tags which should be assigned to the resources in this plan"

  default = {
    environment = "Unknown"
    costcenter  = "Unknown"
    project     = "Unknown"
  }
}

# VNet 
variable "address_space" {
  type = "string"
  description = "The address space that is used by the virtual network."
  default     = "10.0.0.0/16"
}

variable "dns_servers" {
  type = "list"
  description = "The DNS servers to be used by the virtual network"
  default     = []
}

variable "subnet_prefix" {
  type = "string"
  description = "The address prefix to use for the subnet."
  default     = "10.0.1.0/24"
}

variable "subnet_name" {
  type = "string"
  description = "The name to use for the subnet."
  default     = "hdi-secure"
}

# Note: New Service Endpoints may be available in the future
# Docs: https://docs.microsoft.com/en-us/azure/virtual-network/virtual-network-service-endpoints-overview
variable "service_endpoints" {
  type = "list"
  description = "The service endoints to enable on the subnet."
  default     = ["Microsoft.AzureCosmosDb","Microsoft.KeyVault","Microsoft.Sql","Microsoft.Storage","Microsoft.ServiceBus","Microsoft.EventHub"]
}

# NSGs
# Note: HDInsight Static IPs for management traffic
# Docs: https://docs.microsoft.com/en-us/azure/hdinsight/hdinsight-extend-hadoop-virtual-network#hdinsight-ip-1
# Warning: Be sure to check for updates, this is a static list of IPs that may change over time.

variable "source_address_prefixes_mgmt" {
  type        = "list"
  description = "Used in NSG rule to enable management traffic for HD Insight"
  default     = ["168.61.49.99", "23.99.5.239", "168.61.48.131", "138.91.141.162"]
}

variable "source_address_prefix_resolver" {
  type = "string"
  description = "Used in NSG rule to enable Azure Resolver traffic to HD Insight"
  default     = "168.63.129.16"
}

variable "source_address_prefixes_mgmt_region" {
  type        = "list"
  description = "Used in NSG rule to enable management traffic for HD Insight"
  default = [
    "23.102.235.122",
    "52.175.38.134",
    "13.76.245.160",
    "13.76.136.249",
    "104.210.84.115",
    "13.75.152.195",
    "13.77.2.56",
    "13.77.2.94",
    "191.235.84.104",
    "191.235.87.113",
    "52.229.127.96",
    "52.229.123.172",
    "52.228.37.66",
    "52.228.45.222",
    "42.159.96.170",
    "139.217.2.219",
    "42.159.198.178",
    "42.159.234.157",
    "40.73.37.141",
    "40.73.38.172",
    "52.164.210.96",
    "13.74.153.132",
    "52.166.243.90",
    "52.174.36.244",
    "51.4.146.68",
    "51.4.146.80",
    "51.5.150.132",
    "51.5.144.101",
    "52.172.153.209",
    "52.172.152.49",
    "104.211.223.67",
    "104.211.216.210",
    "13.78.125.90",
    "13.78.89.60",
    "40.74.125.69",
    "138.91.29.150",
    "52.231.39.142",
    "52.231.36.209",
    "52.231.203.16",
    "52.231.205.214",
    "51.141.13.110",
    "51.141.7.20",
    "51.140.47.39",
    "51.140.52.16",
    "13.67.223.215",
    "40.86.83.253",
    "13.82.225.233",
    "40.71.175.99",
    "157.56.8.38",
    "157.55.213.99",
    "52.161.23.15",
    "52.161.10.167",
    "13.64.254.98",
    "23.101.196.19",
    "52.175.211.210",
    "52.175.222.222"
  ]
}

# Azure Storage Accounts  for HDInsight
variable "storage_account_count" {
  description = "The number of storage acccounts you want to deploy"
  default = 1
}
variable "account_replication_type" {
  type = "string"
  description = "Account type for storage accounts ('LRS', 'GRS' or 'RAGRS'.  'ZRS' probably does not apply to HD Insight)"
  default = "LRS"
}

# Azure SQL Databases for HDInsight Hive Metastore or Oozie
variable "azuresqldb_databases" {
  type = "list"
  description = "The names of the Azure SQL Databases to create. Leave this empty if you don't want to create any."
  default = ["hivemetastoredb","ooziedb"]
}

variable "sql_server_admin_user" {
  type = "string"
  description = "The admin user for the Azure SQL Database server."
  default = "sqlserveradmin"
}

variable "sql_server_admin_password" {
  type = "string"
  description = "The password for the Azure SQL Database admin user."
}
