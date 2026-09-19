## Example: Container App deployed into Virtual Network

A container app deployed into a virtual network and only accessible from other resources in the virtual network, not from the public internet.

An example of when this would be required is if API-Management (APIM) is used to manage access to the container app. APIM would be publicly accessible, but the container app would be internal-only.

**Note:** An internal-only Azure Container Apps environment receives an IP address directly from the target subnet and, by default, there is no domain registration created with the internal load balancer. This means we need to configure DNS resolution for the environment endpoint. In Azure, this can be accomplished by creating and configuring an Azure Private DNS Zone.

## Creates

1. A Resource Group
1. A [Virtual Network](https://azure.microsoft.com/en-us/products/virtual-network/) and subnet
1. A Container App Environment deployed into the subnet
1. A [Container App](https://azure.microsoft.com/en-us/products/container-apps/) only accessible from other resources within the virtual network
1. A Private DNS Zone that maps the FQDN of the container app to the static IP address of the container app environment
