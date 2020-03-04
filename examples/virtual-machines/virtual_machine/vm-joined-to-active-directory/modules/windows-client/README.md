## Module: Windows Client

This module provisions a Windows Client which will be bound to the Active Directory Domain created in the other module.

There's a few hacks in here as we have to wait for Active Directory to become available, but this takes advantage of the `azurerm_virtual_machine_extension` resource. It's worth noting that the keys in this resource are case sensitive.
