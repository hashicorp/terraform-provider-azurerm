# Property Naming

As with naming variables, property naming can also be a laborious task. Given the nature of the provider careful consideration should be given to property names, since changing it is a non-negligible amount of effort.

Whilst there are many cases where the property name can be taken over 1 to 1 from the Azure API, there are many instances where this is not the case.

Here are some general guidelines you can turn to when naming properties:

* The name should describe what the property is for succinctly, but as with many things a balance should be struck between too short or too long.

* Choose the officially marketed name for new properties over the ones used in the API if they differ.

* Abbreviations should not be used and the full words should be used instead e.g. 
>`resource_group_name` instead of `rg_name` or `virtual_machine` instead of `vm`.

* For blocks avoid redundant words in the name that don't add informational value e.g.
>`firewall_properties` can be shortened to `firewall`, the same can apply to individual properties e.g. `email_address` to `email`.

* Properties for certificates or artifacts that must be in a certain format should be appended with the format e.g.
> A certificate that must be base64 encoded should be named `certificate_base64`

* Similarly, properties that pertain to sizes or durations/windows/occurences should be appended with the appropriate unit of measure e.g.
> `duration_in_seconds` or `size_in_gb`

* Time properties that are not in the format of RFC3339 or are specified as UTC in the documentation should have that appended e.g.
 > `timestamp_in_utc`

* For booleans these guidelines apply:

  * As a general rule, booleans should be appended with `_enabled` e.g.
  >`public_network_access_enabled`

  * Booleans named `disableSomething` in the API should be flipped and exposed as `something_enabled` in the provider.
  
  * Avoid redundant verbs like `is` at the beginning of the property e.g.
  >`is_storage_enabled` must be renamed to `storage_enabled`.

  * Avoid double negatives which obfuscate the purpose of the property these should be removed and flipped e.g.
  >`no_storage_enabled` becomes `storage_enabled` or `block_user_upload_enabled` becomes `user_upload_enabled`.
