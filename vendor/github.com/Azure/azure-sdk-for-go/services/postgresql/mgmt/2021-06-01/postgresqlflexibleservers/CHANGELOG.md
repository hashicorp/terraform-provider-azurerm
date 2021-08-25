# Change History

## Breaking Changes

### Removed Constants

1. ResourceIdentityType.ResourceIdentityTypeSystemAssigned

### Removed Funcs

1. Identity.MarshalJSON() ([]byte, error)
1. PossibleResourceIdentityTypeValues() []ResourceIdentityType
1. ResourceModelWithAllowedPropertySet.MarshalJSON() ([]byte, error)
1. ResourceModelWithAllowedPropertySetIdentity.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. Identity
1. Plan
1. ResourceModelWithAllowedPropertySet
1. ResourceModelWithAllowedPropertySetIdentity
1. ResourceModelWithAllowedPropertySetPlan
1. ResourceModelWithAllowedPropertySetSku

#### Removed Struct Fields

1. Server.Identity
