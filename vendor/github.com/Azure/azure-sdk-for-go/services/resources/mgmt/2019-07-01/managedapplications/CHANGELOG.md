Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewApplicationDefinitionListResultPage` parameter(s) have been changed from `(func(context.Context, ApplicationDefinitionListResult) (ApplicationDefinitionListResult, error))` to `(ApplicationDefinitionListResult, func(context.Context, ApplicationDefinitionListResult) (ApplicationDefinitionListResult, error))`
- Function `ApplicationDefinitionsClient.GetByIDPreparer` parameter(s) have been changed from `(context.Context, string)` to `(context.Context, string, string)`
- Function `ApplicationsClient.Update` parameter(s) have been changed from `(context.Context, string, string, *Application)` to `(context.Context, string, string, *ApplicationPatchable)`
- Function `ApplicationsClient.UpdatePreparer` parameter(s) have been changed from `(context.Context, string, string, *Application)` to `(context.Context, string, string, *ApplicationPatchable)`
- Function `ApplicationDefinitionsClient.CreateOrUpdateByIDPreparer` parameter(s) have been changed from `(context.Context, string, ApplicationDefinition)` to `(context.Context, string, string, ApplicationDefinition)`
- Function `NewApplicationListResultPage` parameter(s) have been changed from `(func(context.Context, ApplicationListResult) (ApplicationListResult, error))` to `(ApplicationListResult, func(context.Context, ApplicationListResult) (ApplicationListResult, error))`
- Function `ApplicationDefinitionsClient.GetByID` parameter(s) have been changed from `(context.Context, string)` to `(context.Context, string, string)`
- Function `ApplicationDefinitionsClient.DeleteByIDPreparer` parameter(s) have been changed from `(context.Context, string)` to `(context.Context, string, string)`
- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Function `ApplicationDefinitionsClient.DeleteByID` parameter(s) have been changed from `(context.Context, string)` to `(context.Context, string, string)`
- Function `ApplicationDefinitionsClient.CreateOrUpdateByID` parameter(s) have been changed from `(context.Context, string, ApplicationDefinition)` to `(context.Context, string, string, ApplicationDefinition)`
- Field `*ApplicationPropertiesPatchable` of struct `ApplicationPatchable` has been removed

## New Content

- New anonymous field `*ApplicationProperties` in struct `ApplicationPatchable`
