TF_ACC=1 go test github.com/terraform-providers/terraform-provider-azurerm/azurerm -timeout 180m -v -run TestAccAzureRMVirtualMachineScaleSet_UserAssignedMSI &
TF_ACC=1 go test github.com/terraform-providers/terraform-provider-azurerm/azurerm -timeout 180m -v -run TestAccAzureRMVirtualMachineScaleSet_SystemAssignedMSI &
TF_ACC=1 go test github.com/terraform-providers/terraform-provider-azurerm/azurerm -timeout 180m -v -run TestAccAzureRMVirtualMachine_UserAssignedIdentity & 
TF_ACC=1 go test github.com/terraform-providers/terraform-provider-azurerm/azurerm -timeout 180m -v -run TestAccAzureRMVirtualMachine_SystemAssignedIdentity &
TF_ACC=1 go test github.com/terraform-providers/terraform-provider-azurerm/azurerm -timeout 180m -v -run TestAccAzureRMUserAssignedIdentity_create &
wait

