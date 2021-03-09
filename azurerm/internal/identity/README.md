# Identity

This package contains helpers for working with Managed Identities.

Azure supports up to 4 different combinations of Managed Identities:

* `None` - where no Managed Identity is available/configured for this Azure Resource.
* `SystemAssigned` - where Azure will generate a Managed Identity (Service Principal) for this Azure Resource.
* `SystemAssigned, UserAssigned` - where Azure will generate a Managed Identity (Service Principal) for this Azure Resource, but they can also be assigned.
* `UserAssigned` - where specific Managed Identities can be assigned to this Azure Resource.

Since Managed Identities are an optional feature - within Terarform we're exposing this in 3 manners, exposed in this package as 3 types:

* `SystemAssigned`
* `SystemAssignedUserAssigned` (coming soon)
* `UserAssigned`

Where the block is Optional within Terraform - for consistency across the Provider we've opted to treat the absence of the `identity` block to represent "None" - and the presence of the block to indicate one of the Managed Identity types above.

## Usage

Within the resource itself, assign a type reference via:

```go
type resourceNameIdentity = identity.SystemAssigned
```

which can then be instantiated and used to call the Expand, Flatten and Schema functions:

```go
resourceNameIdentity{}.Schema()
resourceNameIdentity{}.Expand(d.Get("identity").([]interface{}))
resourceNameIdentity{}.Flatten(input)
```

Due to the Azure SDK using a different Type for each Service Package, at this time an Expand and Flatten function are needed to cast from the intermediate type `*identity.ExpandedConfig` to the type used within the Azure SDK for the specified Service Package, for example:

```go
func expandResourceNameIdentity(input []interface{}) (*somepackage.PackageTypeForManagedIdentity, error) {
	config, err := resourceNameIdentity{}.Expand(input)
	if err != nil {
		return nil, err
	}

	return &somepackage.ManagedIdentityProperties{
		Type:        somepackage.ManagedIdentityType(config.Type),
		TenantID:    config.TenantId,
		PrincipalID: config.PrincipalId,
	}, nil
}

func flattenResourceNameIdentity(input *somepackage.ManagedIdentityProperties) []interface{} {
	var config *identity.ExpandedConfig
	if input != nil {
		config = &identity.ExpandedConfig{
			Type:        string(input.Type),
			PrincipalId: input.PrincipalID,
			TenantId:    input.TenantID,
		}
	}
	return resourceNameIdentity{}.Flatten(config)
}
```