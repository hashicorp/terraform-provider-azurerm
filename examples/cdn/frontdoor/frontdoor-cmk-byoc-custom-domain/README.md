## Example CDN Frontdoor: Frontdoor CMK(Customer Managed Key)/BYOC(Bring Your Own Certificate) TLS/SSL Certificate with Custom Domain

!>**IMPORTANT:** CDN Frontdoor is a **GLOBAL** resource. So make sure that your `${var.prefix}` value is sufficiently random else your resources may have naming collisions and cause unexpected errors when running this example and cause it to fail when applied(e.g. `foo`, `bar` or `example` would not be a good choice for the `${var.prefix}` argument as those values would have a high probability of collision).

!>**IMPORTANT:** If you don't already have a custom domain, you must first purchase one with a domain provider. For example, see [buy a custom domain name](https://docs.microsoft.com/azure/app-service/manage-custom-dns-buy-domain).

This example provisions a CDN Frontdoor with a CMK/BYOC TLS/SSL Custom Domain within Azure.

---

To successfully complete this example you will need to create an `Azure DNS Zone` and delegate your domain provider's domain name system (DNS) to the `Azure DNS Zone`. For more information on how to delegate your domain provider's DNS to the `Azure DNS Zone` please see the [delegate a domain to Azure DNS](https://docs.microsoft.com/azure/dns/dns-delegate-domain-azure-dns) product documentation. You may create the `Azure DNS Zone` via Portal or with Terraform by using the below `Example Azure DNS Zone` HCL. However, if you use Portal to create your `Azure DNS Zone` pay close attention to follow the naming convention of this example for your `Resource Group` name (e.g. `${var.prefix}-cdn-frontdoor-byoc-example`).

## Example Azure DNS Zone:

```hcl
resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-cdn-frontdoor-byoc-example"
  location = "westeurope"
} 

resource "azurerm_dns_zone" "example" {
  name                = "example.com" # change this to be your domain name
  resource_group_name = azurerm_resource_group.example.name
}
```

## Looking up Object IDs in Azure Portal:

For this example you will need to look up the object IDs of the Frontdoor service, your personal user account, if you want to view your secrets in the Portal UI, and your service principal that Terraform is running as. To look up the required object IDs you will need to open [Portal](https://portal.azure.com/) and follow the below steps:

* From the left hand menu select `Azure Active Directory`.

* In the search filter box, near the top of the page, type `Microsoft.Azure.Cdn`.

* Click on the `Microsoft.Azure.Cdn` entry in the `Enterprise Applications` results view.

* This will open the `Enterprise Applications Properties`, copy the `Object ID` and paste it into the examples `main.tf` file where is says `<- Object Id for the Microsoft.Azure.Cdn Enterprise Application.`.

Repeat the above steps for all of the object IDs needed for this example.

## Key Vault Permissions:

The following Key Vault permission are granted by this example:

| Object ID                                | Key Permissions | Secret Permissions   | Certificate Permissions                       |
|:-----------------------------------------|:---------------:|:--------------------:|:---------------------------------------------:|
| `Microsoft.Azure.Cdn` Object ID | -               | **Get**              | -                                             |
| Your Personal AAD Object ID              | -               | **Get** and **List** | **Get**, **List**, **Purge** and **Recover**  |
| Terraform Service Principal              | -               | **Get**              | **Get**, **Import**, **Delete** and **Purge** |

Once you have created your `Azure DNS Zone` and delegated your domain provider's DNS to the `Azure DNS Zone` you will need to import the `Resource Group` and the `Azure DNS Zone` into the Terraform state file by running the following import commands:

* terraform import azurerm_resource_group.example /subscriptions/{subscription}/resourceGroups/`${var.prefix}-cdn-frontdoor-byoc-example`
* terraform import azurerm_dns_zone.example /subscriptions/{subscription}/resourceGroups/`${var.prefix}-cdn-frontdoor-byoc-example`/providers/Microsoft.Network/dnszones/`dnsZoneName`

Now that the Prerequisites have been completed and your state file contains your `Resource Group` and `Azure DNS Zone` you can simply run `terraform apply` to create a `CDN Frontdoor` with a BYOC TLS/SSL custom domain.

### Variables

* `prefix` - (Required) The prefix used for all resources in this example.
