# This is the name of the your VM - Needs to be unique in the resource group
customVmName = "gen-vmname"

# This is the name of the your storage account - Needs to be the same storage account as the source user image
userImageStorageAccountName = "image-storage-account"

# Resource group of the existing storage account
userImageStorageAccountResourceGroupName = "image-storage-account-resource-group"

# "Uri of the your user image
osDiskVhdUri = ""

# DNS Label for the Public IP. Must be lowercase. It should match with the following regular expression: ^[a-z][a-z0-9-]{1,61}[a-z0-9]$ or it will raise an error.
# Needs to be unique
dnsLabelPrefix = "gen-unique"

# User Name for the Virtual Machine
adminUserName = "gen-unique"

# Password for the Virtual Machine
adminPassword = "gen-unique"

# This is the OS that your VM will be running. 
# Needs to be "linux" or "windows"
osType = "linux"

# This is the size of your VM
vmSize = "Standard_A1"

# Select if this template needs a new VNet or will reference an existing VNet
# Needs to be "new" or "existing"
newOrExistingVnet = "new"

# New or existing VNet name
newOrExistingVnetName = "uivnet"

# New or existing subnet name
newOrExistingSubnetName = "uisubnet"

# Name of the resource group container for all new resources
resourceGroupName = "vmuserimagelinux"

# Resource group location
resourceGroupLocation = "West US"

# SSH Key
admin_sshkey=""