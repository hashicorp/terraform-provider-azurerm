## Example CDN Frontdoor: Frontdoor Private Link Origin with a Linux Web Application

!>**IMPORTANT:** CDN Frontdoor is a **GLOBAL** resource. So make sure that your `${var.prefix}` value is sufficiently random else your resources may have naming collisions and cause unexpected errors when running this example and cause it to fail when applied(e.g. `foo`, `bar` or `example` would not be a good choice for the `${var.prefix}` argument as those values would have a high probability of collision).

->**INFORMATION:** The approval of the Private Link Endpoint is currently a manual step in this process. For more information please see the [product documentation](https://docs.microsoft.com/azure/frontdoor/private-link).

This example provisions a CDN Frontdoor with a Linux Web Application Private Link Origin within Azure.

### Variables

* `prefix` - (Required) The prefix used for all resources in this example.
