## Example CDN Frontdoor: Frontdoor CMK(Customer Managed Key)/BYOC(Bring Your Own Certificate) TLS/SSL Certificate with Custom Domain

!>**IMPORTANT:** CDN Frontdoor is a **GLOBAL** resource. So make sure that your `${var.prefix}` value is sufficiently random else your resources may have naming collisions and cause unexpected errors when running this example and cause it to fail when applied(e.g. `foo`, `bar` or `example` would not be a good choice for the `${var.prefix}` argument as those values would have a high probability of collision).

!>**IMPORTANT:** If you don't already have a custom domain, you must first purchase one with a domain provider. For example, see [buy a custom domain name](https://docs.microsoft.com/azure/app-service/manage-custom-dns-buy-domain).

To successfully complete this example you will need to create an `Azure DNS Zone` and delegate your domain provider's domain name system (DNS) to the `Azure DNS Zone`. For more information on how to delegate your domain provider's DNS to the `Azure DNS Zone` please see the [delegate a domain to Azure DNS](https://docs.microsoft.com/azure/dns/dns-delegate-domain-azure-dns) product documentation. You may create the `Azure DNS Zone` via Portal or with Terraform by using the below `Example Azure DNS Zone` HCL. However, if you use Portal to create your `Azure DNS Zone` pay close attention to follow the naming convention of this example for your `Resource Group` name (e.g. `${var.prefix}-cdn-frontdoor-managed-ssl-example`).

## Example Azure DNS Zone:

```hcl
resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-cdn-frontdoor-managed-ssl-example"
  location = "westeurope"
} 

resource "azurerm_dns_zone" "example" {
  name                = "example.com" # change this to be your domain name
  resource_group_name = azurerm_resource_group.example.name
}
```

## Key Vault Permissions:

**TODO:** More explanation here around the Object IDs, Key Vault and Key Vault Certificate aspects of this example.

The following Key Vault permission are granted by this example, but you will need to update the `object_id` fields in the `main.tf` file of this example for your tenant to run this example:

| Object ID                                | Key Permissions | Secret Permissions   | Certificate Permissions                       |
|:-----------------------------------------|:---------------:|:--------------------:|:---------------------------------------------:|
| `Microsoft.AzureFrontDoor-Cdn` Object ID | -               | **Get**              | -                                             |
| Your Personal AAD Object ID              | -               | **Get** and **List** | **Get**, **List**, **Purge** and **Recover**  |
| Terraform Service Principal              | -               | **Get**              | **Get**, **Import**, **Delete** and **Purge** |

Once you have created your `Azure DNS Zone`, delegated your domain provider's DNS to the `Azure DNS Zone` you will need to import the `Resource Group` and the `Azure DNS Zone` into the Terraform state file by running the following import commands:

* terraform import azurerm_resource_group.example /subscriptions/{subscription}/resourceGroups/`${var.prefix}-cdn-frontdoor-managed-ssl-example`
* terraform import azurerm_dns_zone.example /subscriptions/{subscription}/resourceGroups/`${var.prefix}-cdn-frontdoor-managed-ssl-example`/providers/Microsoft.Network/dnszones/`dnsZoneName`

Now that the Prerequisites have been completed and your state file contains your `Resource Group` and `Azure DNS Zone` you can simply run `terraform apply` to create a `CDN Frontdoor` with two Frontdoor managed TLS/SSL certificate custom domains.
