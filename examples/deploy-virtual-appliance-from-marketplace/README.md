# Deploy a Virtual Appliance VM from Azure Marketplace

In order to enable deployment in a programatic way (such as with Terraform), you first need to enable it in Azure. This is a one time action per Subscription per Solution.

Two ways to acheive that:

A. using PowerShell:

https://docs.microsoft.com/en-us/powershell/module/azurerm.marketplaceordering/set-azurermmarketplaceterms?view=azurermps-6.13.0

A validation example:

https://raw.githubusercontent.com/shayshahak/terraform-provider-azurerm/master/examples/deploy-virtual-appliance-from-marketplace/images/GetAzureRmMarketPlaceTerms.JPG


B. using GUI:

Search and choose for the solution in Azure Marketplace. scroll down and click on "want to deploy programmatically?"

https://github.com/shayshahak/terraform-provider-azurerm/blob/master/examples/deploy-virtual-appliance-from-marketplace/images/wanttodeployprogrammatically.jpg

scroll down and click on "Enable" and then "Save"

https://github.com/shayshahak/terraform-provider-azurerm/blob/master/examples/deploy-virtual-appliance-from-marketplace/images/configureprogrammaticdeployment.JPG

further information on the GUI way can be found here:

https://azure.microsoft.com/en-us/blog/working-with-marketplace-images-on-azure-resource-manager/

In my example I used a Radware solution but it could be anything with Azure MArketplace
