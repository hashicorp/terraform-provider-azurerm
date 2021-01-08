Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewSpatialAnchorsAccountListPage` parameter(s) have been changed from `(func(context.Context, SpatialAnchorsAccountList) (SpatialAnchorsAccountList, error))` to `(SpatialAnchorsAccountList, func(context.Context, SpatialAnchorsAccountList) (SpatialAnchorsAccountList, error))`
- Function `NewOperationListPage` parameter(s) have been changed from `(func(context.Context, OperationList) (OperationList, error))` to `(OperationList, func(context.Context, OperationList) (OperationList, error))`
- Type of `CheckNameAvailabilityResponse.NameAvailable` has been changed from `NameAvailability` to `*bool`
- Const `False` has been removed
- Const `True` has been removed
- Function `PossibleNameAvailabilityValues` has been removed

## New Content

- New const `Premium`
- New const `Standard`
- New const `SystemAssigned`
- New const `Basic`
- New const `Free`
- New function `ResourceModelWithAllowedPropertySet.MarshalJSON() ([]byte, error)`
- New function `PossibleResourceIdentityTypeValues() []ResourceIdentityType`
- New function `Identity.MarshalJSON() ([]byte, error)`
- New function `ResourceModelWithAllowedPropertySetIdentity.MarshalJSON() ([]byte, error)`
- New function `PossibleSkuTierValues() []SkuTier`
- New struct `Identity`
- New struct `Plan`
- New struct `ResourceModelWithAllowedPropertySet`
- New struct `ResourceModelWithAllowedPropertySetIdentity`
- New struct `ResourceModelWithAllowedPropertySetPlan`
- New struct `ResourceModelWithAllowedPropertySetSku`
- New struct `Sku`
- New field `IsDataAction` in struct `Operation`
- New field `Identity` in struct `SpatialAnchorsAccount`
