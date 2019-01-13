# Deploy a virtual appliance VM from Azure Marketplace

In order to enable deployment in a programatic way (such as with Terrafrom), you first need to enable it in Azure. This is a one time action per Subscription per Solution.

Two ways to acheive that:

A. using PowerShell:

https://docs.microsoft.com/en-us/powershell/module/azurerm.marketplaceordering/set-azurermmarketplaceterms?view=azurermps-6.13.0

A validation example:

![getazurermmarketplaceterms](https://user-images.githubusercontent.com/18166141/50479182-ecfdf700-09dd-11e9-9afa-d47f0a77fcb9.JPG)


B. using GUI:

search and choose for the solution in Azure Marketplace. scroll down and click on "want to deploy programmatically?"

![wanttodeployprogrammatically](https://user-images.githubusercontent.com/18166141/50479221-25053a00-09de-11e9-82cf-c779acaa272a.jpg)

scroll down and click on "Enable" and then "Save"

![configureprogrammaticdeployment](https://user-images.githubusercontent.com/18166141/50479242-377f7380-09de-11e9-94b3-3e07492533bd.JPG)

further information on the GUI way can be found here:

https://azure.microsoft.com/en-us/blog/working-with-marketplace-images-on-azure-resource-manager/

In my example I used a Radware solution but it could be anything with Azure MArketplace
