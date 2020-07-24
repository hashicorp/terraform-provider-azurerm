---
subcategory: "attestation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_attestation_attestation_provider"
description: |-
  Manages a attestation AttestationProvider.
---

# azurerm_attestation_attestation_provider

Manages a attestation AttestationProvider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_attestation_attestation_provider" "example" {
  name = "example-attestationprovider"
  resource_group_name = azurerm_resource_group.example.name
  location = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this attestation AttestationProvider. Changing this forces a new attestation AttestationProvider to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the attestation AttestationProvider should exist. Changing this forces a new attestation AttestationProvider to be created.

* `location` - (Required) The Azure Region where the attestation AttestationProvider should exist. Changing this forces a new attestation AttestationProvider to be created.

---

* `attestation_policy` - (Optional) Name of attestation policy. Changing this forces a new attestation AttestationProvider to be created.

* `policy_signing_certificate` - (Optional)  A `policy_signing_certificate` block as defined below. Changing this forces a new attestation AttestationProvider to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the attestation AttestationProvider.

---

An `policy_signing_certificate` block exports the following:

* `key` - (Optional)  A `key` block as defined below. Changing this forces a new attestation AttestationProvider to be created.

---

An `key` block exports the following:

* `alg` - (Required) The "alg" (algorithm) parameter identifies the algorithm intended for
use with the key.  The values used should either be registered in the
IANA "JSON Web Signature and Encryption Algorithms" registry
established by [JWA] or be a value that contains a Collision-
Resistant Name. Changing this forces a new attestation AttestationProvider to be created.

* `kid` - (Required) The "kid" (key ID) parameter is used to match a specific key.  This
is used, for instance, to choose among a set of keys within a JWK Set
during key rollover.  The structure of the "kid" value is
unspecified.  When "kid" values are used within a JWK Set, different
keys within the JWK Set SHOULD use distinct "kid" values.  (One
example in which different keys might use the same "kid" value is if
they have different "kty" (key type) values but are considered to be
equivalent alternatives by the application using them.)  The "kid"
value is a case-sensitive string. Changing this forces a new attestation AttestationProvider to be created.

* `kty` - (Required) The "kty" (key type) parameter identifies the cryptographic algorithm
family used with the key, such as "RSA" or "EC". "kty" values should
either be registered in the IANA "JSON Web Key Types" registry
established by [JWA] or be a value that contains a Collision-
Resistant Name.  The "kty" value is a case-sensitive string. Changing this forces a new attestation AttestationProvider to be created.

* `use` - (Required) Use ("public key use") identifies the intended use of
the public key. The "use" parameter is employed to indicate whether
a public key is used for encrypting data or verifying the signature
on data. Values are commonly "sig" (signature) or "enc" (encryption). Changing this forces a new attestation AttestationProvider to be created.

---

* `crv` - (Optional) The "crv" (curve) parameter identifies the curve type. Changing this forces a new attestation AttestationProvider to be created.

* `d` - (Optional) RSA private exponent or ECC private key. Changing this forces a new attestation AttestationProvider to be created.

* `dp` - (Optional) RSA Private Key Parameter. Changing this forces a new attestation AttestationProvider to be created.

* `dq` - (Optional) RSA Private Key Parameter. Changing this forces a new attestation AttestationProvider to be created.

* `e` - (Optional) RSA public exponent, in Base64. Changing this forces a new attestation AttestationProvider to be created.

* `k` - (Optional) Symmetric key. Changing this forces a new attestation AttestationProvider to be created.

* `n` - (Optional) RSA modulus, in Base64. Changing this forces a new attestation AttestationProvider to be created.

* `p` - (Optional) RSA secret prime. Changing this forces a new attestation AttestationProvider to be created.

* `q` - (Optional) RSA secret prime, with p < q. Changing this forces a new attestation AttestationProvider to be created.

* `qi` - (Optional) RSA Private Key Parameter. Changing this forces a new attestation AttestationProvider to be created.

* `x` - (Optional) X coordinate for the Elliptic Curve point. Changing this forces a new attestation AttestationProvider to be created.

* `x5cs` - (Optional) The "x5c" (X.509 certificate chain) parameter contains a chain of one or more PKIX certificates [RFC5280].  The certificate chain is represented as a JSON array of certificate value strings.  Each string in the array is a base64-encoded (Section 4 of [RFC4648] not base64url-encoded) DER [ITU.X690.1994] PKIX certificate value. The PKIX certificate containing the key value MUST be the first certificate. Changing this forces a new attestation AttestationProvider to be created.

* `y` - (Optional) Y coordinate for the Elliptic Curve point. Changing this forces a new attestation AttestationProvider to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the attestation AttestationProvider.

* `attest_uri` - Gets the uri of attestation service.

* `trust_model` - Trust model for the attestation service instance.

* `type` - The type of the resource. Ex- Microsoft.Compute/virtualMachines or Microsoft.Storage/storageAccounts.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the attestation AttestationProvider.
* `read` - (Defaults to 5 minutes) Used when retrieving the attestation AttestationProvider.
* `update` - (Defaults to 30 minutes) Used when updating the attestation AttestationProvider.
* `delete` - (Defaults to 30 minutes) Used when deleting the attestation AttestationProvider.

## Import

attestation AttestationProviders can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_attestation_attestation_provider.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Attestation/attestationProviders/provider1
```
