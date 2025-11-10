# Property Naming

As with naming variables, property naming can also be a laborious task. Given the nature of the provider careful consideration should be given to property names, since changing it is a non-negligible amount of effort.

Whilst there are many cases where the property name can be taken over 1 to 1 from the Azure API, there are many instances where this is not the case.

Here are some general guidelines you can turn to when naming properties:

## General Property Naming Conventions

* The name should describe what the property is for succinctly, but as with many things a balance should be struck between too short or too long.

* Choose the officially marketed name (i.e. what is used in Azure Portal or the documentation) for new properties over the ones used in the API if they differ.

* Abbreviations should not be used and the full words should be used instead e.g. 
> `resource_group_name` instead of `rg_name` or `virtual_machine` instead of `vm`.

* For blocks avoid redundant words in the name that don't add informational value e.g.
> `firewall_properties` can be shortened to `firewall`, the same can apply to individual properties e.g. `email_address` to `email`.

* Properties for certificates or artifacts that must be in a certain format should be appended with the format e.g.
> A certificate that must be base64 encoded should be named `certificate_base64`

* Similarly, properties that pertain to sizes or durations/windows/occurrences should be appended with the appropriate unit of measure. However, refrain from appending unit of measure where the possible values are in [ISO8601 Duration format](https://docs.digi.com/resources/documentation/digidocs/90001488-13/reference/r_iso_8601_duration_format.htm) to avoid complex two-way integer-string mapping logic e.g.
> `duration_in_seconds` or `size_in_gb`

> `duration = "12h"`

* Time properties that are not in the format of RFC3339 or are specified as UTC in the documentation should have that appended e.g.
 > `timestamp_in_utc`

## Boolean Property Naming Conventions

* As a general rule, booleans should be appended with `_enabled` e.g.
>`compression_enabled`

* Booleans named `disableSomething` in the API should be flipped and exposed as `something_enabled` in the provider.

* Do note that the above is not a hard requirement, there are cases where appending `_enabled` leads to odd and/or inaccurate naming, e.g.
> `requireMtls` being named `require_mtls_enabled`, in this scenario it would be preferable to name the property `mtls_required`.

* Similarly, properties that indicate a state of something likely make more sense without appending `_enabled`, e.g.
> `acceptedTermsOfService` becomes `terms_of_service_accepted`.

* Avoid redundant verbs like `is` at the beginning of the property e.g.
>`is_storage_enabled` must be renamed to `storage_enabled`.

* Avoid double negatives which obfuscate the purpose of the property these should be removed and flipped e.g.
>`no_storage_enabled` becomes `storage_enabled` or `block_user_upload_enabled` becomes `user_upload_enabled`.
