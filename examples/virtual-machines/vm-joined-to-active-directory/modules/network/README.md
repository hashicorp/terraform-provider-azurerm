## Module: Network

This module creates a Network with 2 x Subnets:

 - Domain Controllers
 - Domain Clients

This module shouldn't be used as-is in a Production Environment - where you'd probably have [Network Security Rules](https://www.terraform.io/docs/providers/azurerm/r/network_security_rule.html) configured - it's designed to be a simplified configuration for the purposes of this example.
