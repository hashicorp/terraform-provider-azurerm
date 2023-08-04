# Property Naming

As with naming variables, property naming can also be a laborious task. Given the nature of the provider careful consideration should be given to property names, since changing it is a non-negligible amount of effort.

Whilst there are many cases where the property name can be taken over 1 to 1 from the Azure API, there are many instances where this is not the case.

Here are some general guidelines you can turn to when naming properties:

* The name should describe what the property is for succinctly, but as with many things a balance should be struck between too short or too long.

* Abbreviations should not be used and the full words should be used instead e.g. `resource_group_name` instead of `rg_name` or `virtual_machine` instead of `vm`. 

* For blocks avoid redundant words in the name that don't add informational value e.g.`firewall_properties` can be shortened to `firewall`, the same can apply to individual properties e.g. `email_address` to `email`.

* As a general rule, booleans should be appended with `_enabled`, e.g. `public_network_access_enabled`.

* Choose the officially marketed name for new properties over the ones used in the API if they differ.