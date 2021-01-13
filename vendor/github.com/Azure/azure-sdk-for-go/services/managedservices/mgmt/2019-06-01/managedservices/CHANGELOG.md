Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewRegistrationAssignmentListPage` parameter(s) have been changed from `(func(context.Context, RegistrationAssignmentList) (RegistrationAssignmentList, error))` to `(RegistrationAssignmentList, func(context.Context, RegistrationAssignmentList) (RegistrationAssignmentList, error))`
- Function `NewRegistrationDefinitionListPage` parameter(s) have been changed from `(func(context.Context, RegistrationDefinitionList) (RegistrationDefinitionList, error))` to `(RegistrationDefinitionList, func(context.Context, RegistrationDefinitionList) (RegistrationDefinitionList, error))`
- Type of `ErrorResponse.Error` has been changed from `*ErrorResponseError` to `*ErrorDefinition`
- Struct `ErrorResponseError` has been removed

## New Content

- New struct `ErrorDefinition`
- New field `PrincipalIDDisplayName` in struct `Authorization`
- New field `DelegatedRoleDefinitionIds` in struct `Authorization`
