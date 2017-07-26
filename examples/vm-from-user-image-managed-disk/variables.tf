variable "customVmName" {
  description = "This is the name of the your VM"
}

variable "userImageStorageAccountName" {
  description = "This is the name of the your storage account" 
}

variable "userImageStorageAccountResourceGroupName" {
  description = "Resource group of the existing storage account"
}

variable "osDiskVhdUri" {
  description = "Uri of the your user image"
}

variable "dnsLabelPrefix" {
  description = "DNS Label for the Public IP. Must be lowercase. It should match with the following regular expression: ^[a-z][a-z0-9-]{1,61}[a-z0-9]$ or it will raise an error"
}

variable "adminUserName" {
  description = "User Name for the Virtual Machine"
}

variable "adminPassword" {
  description = "Password for the Virtual Machine"
}

variable "adminSSHKey" {
  description = "SSH Key for the Virtual Machine"  
}

variable "osType" {
  description = "This is the OS that your VM will be running"
}

variable "vmSize" {
  description = "This is the size of your VM"
}

variable "newOrExistingVnet" {
  description = "Select if this template needs a new VNet or will reference an existing VNet"
}

variable "newOrExistingVnetName" {
  description = "New or Existing VNet Name"
}

variable "newOrExistingSubnetName" {
  description = "New or Existing subnet Name"
}

variable "resourceGroupName" {
  description = "Name of the resource group container for all resources"
}

variable "resourceGroupLocation" {
 description = "Azure region used for resource deployment"
}