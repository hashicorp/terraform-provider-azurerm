# Modify this file to customize the HDInsight Virtual Network configuration for your specific environment

# General 
resource_prefix = "hdi-dev" # Change this to a prefix that complies with your resource naming conventions

location = "West US 2" # Change this to the region you want to deploy 

tags = {
  environment = "Dev"
  costcenter  = "Unknown"
  project     = "HD Insight Secure VNet"
} # Change this to the tags you want to use for your resources

# VNet 
address_space = "10.0.0.0/16" # Change this to the address space you want to use, make sure it does not conflict with other VNets

dns_servers = [] # Change this if you require custom DNS settings for hybrid connectivity to your on-premises network

subnet_name = "hdinsight" # Change this to the name of the subnet you want to use for HDInsight

subnet_prefix = "10.0.1.0/24" # Change this to the address space you want to use for the subnet, must be contained within vnet address_space

service_endpoints = ["Microsoft.Sql","Microsoft.Storage"] # Change this to the specific service endpoints you want to enable.

# NSGs
source_address_prefixes_mgmt = ["168.61.49.99", "23.99.5.239", "168.61.48.131", "138.91.141.162"] # Required for all HDInsight Vnets

source_address_prefix_resolver = "168.63.129.16" # Required for all HDInsight VNets

source_address_prefixes_mgmt_region = ["52.175.211.210", "52.175.222.222"] # Limit to your specific region, West US 2 in this example

# Storage
storage_account_count = 1 # Change this to the number of storage accounts you want to create for HDInsight

account_replication_type = "LRS" # Change this to the type of replication you want to use for your storage accounts

# SQL
azuresqldb_databases = ["hivemetastoredb"]

sql_server_admin_user = "sqlserveradmin"
