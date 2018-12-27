# Deploy a VM from Azure Marketplace

deploy a VM from Azure Marketplace. In order to enable deployment in a programatic way (such as with Terrafrom), you first need to enable it in Azure. This is a one time action per Subscription per Solution.

Two ways to acheive that:

A. using PowerShell:

https://docs.microsoft.com/en-us/powershell/module/azurerm.marketplaceordering/set-azurermmarketplaceterms?view=azurermps-6.13.0

a validation example:


B. using GUI:
https://azure.microsoft.com/en-us/blog/working-with-marketplace-images-on-azure-resource-manager/

In my example I used a Radware solution but it could be anything with Azure MArketplace