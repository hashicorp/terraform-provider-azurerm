variable "resource_group" {
  description = "Resource group name into which your Spark and Cassandra deployment will go."
}

variable "location" {
  description = "The location/region where the virtual network is created. Changing this forces a new resource to be created."
  default     = "southcentralus"
}

variable "unique_prefix" {
  description = "This prefix is used for names which need to be globally unique."
}

variable "storage_primary_account_tier" {
  description = "Storage Tier that is used for primary Spark node. This storage account is used to store VM disks. Allowed values are Standard and Premium."
  default     = "Standard"
}

variable "storage_primary_replication_type" {
  description = "Storage Tier that is used for primary Spark node. This storage account is used to store VM disks. Possible values include LRS and GRS."
  default     = "LRS"
}

variable "storage_secondary_account_tier" {
  description = "Storage type that is used for each of the secondary Spark node. This storage account is used to store VM disks. Allowed values are Standard and Premium."
  default     = "Standard"
}

variable "storage_secondary_replication_type" {
  description = "Storage type that is used for each of the secondary Spark node. This storage account is used to store VM disks. Possible values include LRS and GRS."
  default     = "LRS"
}

variable "storage_cassandra_account_tier" {
  description = "Storage type that is used for Cassandra. This storage account is used to store VM disks. Allowed values are Standard and Premium."
  default     = "Standard"
}

variable "storage_cassandra_replication_type" {
  description = "Storage type that is used for Cassandra. This storage account is used to store VM disks. Possible values include LRS and GRS."
  default     = "LRS"
}

variable "vm_primary_vm_size" {
  description = "VM size for primary Spark node.  This VM can be sized smaller. Allowed values: Standard_D1_v2, Standard_D2_v2, Standard_D3_v2, Standard_D4_v2, Standard_D5_v2, Standard_D11_v2, Standard_D12_v2, Standard_D13_v2, Standard_D14_v2, Standard_A8, Standard_A9, Standard_A10, Standard_A11"
  default     = "Standard_D1_v2"
}

variable "vm_number_of_secondarys" {
  description = "Number of VMs to create to support the secondarys.  Each secondary is created on its own VM.  Minimum of 2 & Maximum of 200 VMs. min = 2, max = 200"
  default     = 2
}

variable "vm_secondary_vm_size" {
  description = "VM size for secondary Spark nodes.  This VM should be sized based on workloads. Allowed values: Standard_D1_v2, Standard_D2_v2, Standard_D3_v2, Standard_D4_v2, Standard_D5_v2, Standard_D11_v2, Standard_D12_v2, Standard_D13_v2, Standard_D14_v2, Standard_A8, Standard_A9, Standard_A10, Standard_A11"
  default     = "Standard_D3_v2"
}

variable "vm_cassandra_vm_size" {
  description = "VM size for Cassandra node.  This VM should be sized based on workloads. Allowed values: Standard_D1_v2, Standard_D2_v2, Standard_D3_v2, Standard_D4_v2, Standard_D5_v2, Standard_D11_v2, Standard_D12_v2, Standard_D13_v2, Standard_D14_v2, Standard_A8, Standard_A9, Standard_A10, Standard_A11"
  default     = "Standard_D3_v2"
}

variable "vm_admin_username" {
  description = "Specify an admin username that should be used to login to the VM. Min length: 1"
}

variable "vm_admin_password" {
  description = "Specify an admin password that should be used to login to the VM. Must be between 6-72 characters long and must satisfy at least 3 of password complexity requirements from the following: 1) Contains an uppercase character 2) Contains a lowercase character 3) Contains a numeric digit 4) Contains a special character"
}

variable "os_image_publisher" {
  description = "name of the publisher of the image (az vm image list)"
  default     = "OpenLogic"
}

variable "os_image_offer" {
  description = "the name of the offer (az vm image list)"
  default     = "CentOS"
}

variable "os_version" {
  description = "version of the image to apply (az vm image list)"
  default     = "7.3"
}

variable "api_version" {
  default = "2015-06-15"
}

variable "artifacts_location" {
  description = "The base URI where artifacts required by this template are located."
  default     = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/primary/spark-and-cassandra-on-centos/CustomScripts/"
}

variable "vnet_spark_prefix" {
  description = "The address space that is used by the virtual network. You can supply more than one address space. Changing this forces a new resource to be created."
  default     = "10.0.0.0/16"
}

variable "vnet_spark_subnet1_name" {
  description = "The name used for the Primary subnet."
  default     = "Subnet-Primary"
}

variable "vnet_spark_subnet1_prefix" {
  description = "The address prefix to use for the Primary subnet."
  default     = "10.0.0.0/24"
}

variable "vnet_spark_subnet2_name" {
  description = "The name used for the secondary/agent subnet."
  default     = "Subnet-Secondary"
}

variable "vnet_spark_subnet2_prefix" {
  description = "The address prefix to use for the secondary/agent subnet."
  default     = "10.0.1.0/24"
}

variable "vnet_spark_subnet3_name" {
  description = "The name used for the subnet used by Cassandra."
  default     = "Subnet-Cassandra"
}

variable "vnet_spark_subnet3_prefix" {
  description = "The address prefix to use for the subnet used by Cassandra."
  default     = "10.0.2.0/24"
}

variable "nsg_spark_primary_name" {
  description = "The name of the network security group for Spark's Primary"
  default     = "nsg-spark-primary"
}

variable "nsg_spark_secondary_name" {
  description = "The name of the network security group for Spark's secondary/agent nodes"
  default     = "nsg-spark-secondary"
}

variable "nsg_cassandra_name" {
  description = "The name of the network security group for Cassandra"
  default     = "nsg-cassandra"
}

variable "nic_primary_name" {
  description = "The name of the network interface card for Primary"
  default     = "nic-primary"
}

variable "nic_primary_node_ip" {
  description = "The private IP address used by the Primary's network interface card"
  default     = "10.0.0.5"
}

variable "nic_cassandra_name" {
  description = "The name of the network interface card used by Cassandra"
  default     = "nic-cassandra"
}

variable "nic_cassandra_node_ip" {
  description = "The private IP address of Cassandra's network interface card"
  default     = "10.0.2.5"
}

variable "nic_secondary_name_prefix" {
  description = "The prefix used to constitute the secondary/agents' names"
  default     = "nic-secondary-"
}

variable "nic_secondary_node_ip_prefix" {
  description = "The prefix of the private IP address used by the network interface card of the secondary/agent nodes"
  default     = "10.0.1."
}

variable "public_ip_primary_name" {
  description = "The name of the primary node's public IP address"
  default     = "public-ip-primary"
}

variable "public_ip_secondary_name_prefix" {
  description = "The prefix to the secondary/agent nodes' IP address names"
  default     = "public-ip-secondary-"
}

variable "public_ip_cassandra_name" {
  description = "The name of Cassandra's node's public IP address"
  default     = "public-ip-cassandra"
}

variable "vm_primary_name" {
  description = "The name of Spark's Primary virtual machine"
  default     = "spark-primary"
}

variable "vm_primary_os_disk_name" {
  description = "The name of the os disk used by Spark's Primary virtual machine"
  default     = "vmPrimaryOSDisk"
}

variable "vm_primary_storage_account_container_name" {
  description = "The name of the storage account container used by Spark's primary"
  default     = "vhds"
}

variable "vm_secondary_name_prefix" {
  description = "The name prefix used by Spark's secondary/agent nodes"
  default     = "spark-secondary-"
}

variable "vm_secondary_os_disk_name_prefix" {
  description = "The prefix used to constitute the names of the os disks used by the secondary/agent nodes"
  default     = "vmSecondaryOSDisk-"
}

variable "vm_secondary_storage_account_container_name" {
  description = "The name of the storage account container used by the secondary/agent nodes"
  default     = "vhds"
}

variable "vm_cassandra_name" {
  description = "The name of the virtual machine used by Cassandra"
  default     = "cassandra"
}

variable "vm_cassandra_os_disk_name" {
  description = "The name of the os disk used by the Cassandra virtual machine"
  default     = "vmCassandraOSDisk"
}

variable "vm_cassandra_storage_account_container_name" {
  description = "The name of the storage account container used by the Cassandra node"
  default     = "vhds"
}

variable "availability_secondary_name" {
  description = "The name of the availability set for the secondary/agent machines"
  default     = "availability-secondary"
}

variable "script_spark_provisioner_script_file_name" {
  description = "The name of the script kept in version control which will provision Spark"
  default     = "scriptSparkProvisioner.sh"
}

variable "script_cassandra_provisioner_script_file_name" {
  description = "The name of the script kept in version control which will provision Cassandra"
  default     = "scriptCassandraProvisioner.sh"
}
