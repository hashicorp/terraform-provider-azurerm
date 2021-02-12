Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewOperationListResultPage` parameter(s) have been changed from `(func(context.Context, OperationListResult) (OperationListResult, error))` to `(OperationListResult, func(context.Context, OperationListResult) (OperationListResult, error))`
- Struct `ErrorResponseUpdatedFormat` has been removed

## New Content

- New const `CustomRole`
- New const `BuiltInRole`
- New function `SQLResourcesClient.DeleteSQLRoleAssignmentSender(*http.Request) (SQLResourcesDeleteSQLRoleAssignmentFuture, error)`
- New function `SQLResourcesClient.GetSQLRoleDefinitionPreparer(context.Context, string, string, string) (*http.Request, error)`
- New function `SQLResourcesClient.ListSQLRoleAssignments(context.Context, string, string) (SQLRoleAssignmentListResult, error)`
- New function `PossibleRoleDefinitionTypeValues() []RoleDefinitionType`
- New function `*SQLResourcesDeleteSQLRoleAssignmentFuture.Result(SQLResourcesClient) (autorest.Response, error)`
- New function `*SQLRoleAssignmentCreateUpdateParameters.UnmarshalJSON([]byte) error`
- New function `SQLResourcesClient.CreateUpdateSQLRoleAssignmentPreparer(context.Context, string, string, string, SQLRoleAssignmentCreateUpdateParameters) (*http.Request, error)`
- New function `SQLResourcesClient.DeleteSQLRoleDefinitionSender(*http.Request) (SQLResourcesDeleteSQLRoleDefinitionFuture, error)`
- New function `SQLRoleDefinitionCreateUpdateParameters.MarshalJSON() ([]byte, error)`
- New function `SQLResourcesClient.GetSQLRoleAssignment(context.Context, string, string, string) (SQLRoleAssignmentGetResults, error)`
- New function `SQLResourcesClient.GetSQLRoleAssignmentResponder(*http.Response) (SQLRoleAssignmentGetResults, error)`
- New function `SQLResourcesClient.DeleteSQLRoleDefinition(context.Context, string, string, string) (SQLResourcesDeleteSQLRoleDefinitionFuture, error)`
- New function `SQLResourcesClient.CreateUpdateSQLRoleAssignmentSender(*http.Request) (SQLResourcesCreateUpdateSQLRoleAssignmentFuture, error)`
- New function `SQLResourcesClient.ListSQLRoleDefinitionsSender(*http.Request) (*http.Response, error)`
- New function `SQLResourcesClient.GetSQLRoleAssignmentSender(*http.Request) (*http.Response, error)`
- New function `SQLResourcesClient.CreateUpdateSQLRoleDefinitionPreparer(context.Context, string, string, string, SQLRoleDefinitionCreateUpdateParameters) (*http.Request, error)`
- New function `SQLRoleAssignmentGetResults.MarshalJSON() ([]byte, error)`
- New function `SQLResourcesClient.DeleteSQLRoleAssignmentResponder(*http.Response) (autorest.Response, error)`
- New function `SQLResourcesClient.DeleteSQLRoleDefinitionResponder(*http.Response) (autorest.Response, error)`
- New function `SQLRoleAssignmentCreateUpdateParameters.MarshalJSON() ([]byte, error)`
- New function `SQLResourcesClient.CreateUpdateSQLRoleDefinitionSender(*http.Request) (SQLResourcesCreateUpdateSQLRoleDefinitionFuture, error)`
- New function `*SQLResourcesCreateUpdateSQLRoleAssignmentFuture.Result(SQLResourcesClient) (SQLRoleAssignmentGetResults, error)`
- New function `SQLRoleDefinitionGetResults.MarshalJSON() ([]byte, error)`
- New function `SQLResourcesClient.ListSQLRoleAssignmentsSender(*http.Request) (*http.Response, error)`
- New function `SQLResourcesClient.CreateUpdateSQLRoleAssignmentResponder(*http.Response) (SQLRoleAssignmentGetResults, error)`
- New function `SQLResourcesClient.ListSQLRoleAssignmentsPreparer(context.Context, string, string) (*http.Request, error)`
- New function `SQLResourcesClient.GetSQLRoleDefinitionResponder(*http.Response) (SQLRoleDefinitionGetResults, error)`
- New function `SQLResourcesClient.CreateUpdateSQLRoleAssignment(context.Context, string, string, string, SQLRoleAssignmentCreateUpdateParameters) (SQLResourcesCreateUpdateSQLRoleAssignmentFuture, error)`
- New function `SQLResourcesClient.DeleteSQLRoleAssignmentPreparer(context.Context, string, string, string) (*http.Request, error)`
- New function `SQLResourcesClient.DeleteSQLRoleDefinitionPreparer(context.Context, string, string, string) (*http.Request, error)`
- New function `*SQLRoleAssignmentGetResults.UnmarshalJSON([]byte) error`
- New function `*SQLRoleDefinitionCreateUpdateParameters.UnmarshalJSON([]byte) error`
- New function `SQLResourcesClient.ListSQLRoleDefinitionsPreparer(context.Context, string, string) (*http.Request, error)`
- New function `SQLResourcesClient.ListSQLRoleDefinitionsResponder(*http.Response) (SQLRoleDefinitionListResult, error)`
- New function `SQLResourcesClient.CreateUpdateSQLRoleDefinitionResponder(*http.Response) (SQLRoleDefinitionGetResults, error)`
- New function `SQLResourcesClient.GetSQLRoleDefinition(context.Context, string, string, string) (SQLRoleDefinitionGetResults, error)`
- New function `SQLResourcesClient.ListSQLRoleAssignmentsResponder(*http.Response) (SQLRoleAssignmentListResult, error)`
- New function `*SQLRoleDefinitionGetResults.UnmarshalJSON([]byte) error`
- New function `SQLResourcesClient.DeleteSQLRoleAssignment(context.Context, string, string, string) (SQLResourcesDeleteSQLRoleAssignmentFuture, error)`
- New function `SQLResourcesClient.GetSQLRoleDefinitionSender(*http.Request) (*http.Response, error)`
- New function `SQLResourcesClient.ListSQLRoleDefinitions(context.Context, string, string) (SQLRoleDefinitionListResult, error)`
- New function `*SQLResourcesCreateUpdateSQLRoleDefinitionFuture.Result(SQLResourcesClient) (SQLRoleDefinitionGetResults, error)`
- New function `SQLResourcesClient.GetSQLRoleAssignmentPreparer(context.Context, string, string, string) (*http.Request, error)`
- New function `*SQLResourcesDeleteSQLRoleDefinitionFuture.Result(SQLResourcesClient) (autorest.Response, error)`
- New function `SQLResourcesClient.CreateUpdateSQLRoleDefinition(context.Context, string, string, string, SQLRoleDefinitionCreateUpdateParameters) (SQLResourcesCreateUpdateSQLRoleDefinitionFuture, error)`
- New struct `CorsPolicy`
- New struct `DefaultErrorResponse`
- New struct `ManagedServiceIdentityUserAssignedIdentitiesValue`
- New struct `Permission`
- New struct `SQLResourcesCreateUpdateSQLRoleAssignmentFuture`
- New struct `SQLResourcesCreateUpdateSQLRoleDefinitionFuture`
- New struct `SQLResourcesDeleteSQLRoleAssignmentFuture`
- New struct `SQLResourcesDeleteSQLRoleDefinitionFuture`
- New struct `SQLRoleAssignmentCreateUpdateParameters`
- New struct `SQLRoleAssignmentGetResults`
- New struct `SQLRoleAssignmentListResult`
- New struct `SQLRoleAssignmentResource`
- New struct `SQLRoleDefinitionCreateUpdateParameters`
- New struct `SQLRoleDefinitionGetResults`
- New struct `SQLRoleDefinitionListResult`
- New struct `SQLRoleDefinitionResource`
- New field `Identity` in struct `DatabaseAccountUpdateParameters`
- New field `Cors` in struct `DefaultRequestDatabaseAccountCreateUpdateProperties`
- New field `Cors` in struct `DatabaseAccountCreateUpdateProperties`
- New field `Cors` in struct `DatabaseAccountUpdateProperties`
- New field `Cors` in struct `RestoreReqeustDatabaseAccountCreateUpdateProperties`
- New field `UserAssignedIdentities` in struct `ManagedServiceIdentity`
- New field `Cors` in struct `DatabaseAccountGetProperties`
