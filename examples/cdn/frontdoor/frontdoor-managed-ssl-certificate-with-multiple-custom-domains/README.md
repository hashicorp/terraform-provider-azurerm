## Example: CDN Frontdoor / Frontdoor Managed SSL Certificate and Multiple Custom Domains
TODO: Finish this

## Please Note

It is best practice for this example to create your Azure DNS Zone and redirect your domains NS before running this example configuration. Once the Domain's DNS has been redirected to the Azure DNS Zone you should import the Resource Group and the Azure DNS Zone into the Terraform state file with the below import commands. The examples Resource Group name is created in the format of `${var.prefix}-cdn-frontdoor-example`, so when you create your Resource Group make sure to follow this naming convention.

**How to Import the Existing Resources into Terraform:**
terraform import azurerm_resource_group.example /subscriptions/{subscription}/resourceGroups/`${var.prefix}-cdn-frontdoor-example`
terraform import azurerm_dns_zone.example /subscriptions/{subscription}/resourceGroups/`${var.prefix}-cdn-frontdoor-example`/providers/Microsoft.Network/dnszones/`dnsZoneName`

* `${var.prefix}` = The value you will enter when you run `Terraform Apply`(e.g. `Test`)
