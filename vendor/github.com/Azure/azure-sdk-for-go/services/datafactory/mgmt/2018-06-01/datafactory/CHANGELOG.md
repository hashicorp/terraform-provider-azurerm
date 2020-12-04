Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewOperationListResponsePage` parameter(s) have been changed from `(func(context.Context, OperationListResponse) (OperationListResponse, error))` to `(OperationListResponse, func(context.Context, OperationListResponse) (OperationListResponse, error))`
- Function `NewFactoryListResponsePage` parameter(s) have been changed from `(func(context.Context, FactoryListResponse) (FactoryListResponse, error))` to `(FactoryListResponse, func(context.Context, FactoryListResponse) (FactoryListResponse, error))`
- Function `NewLinkedServiceListResponsePage` parameter(s) have been changed from `(func(context.Context, LinkedServiceListResponse) (LinkedServiceListResponse, error))` to `(LinkedServiceListResponse, func(context.Context, LinkedServiceListResponse) (LinkedServiceListResponse, error))`
- Function `NewTriggerListResponsePage` parameter(s) have been changed from `(func(context.Context, TriggerListResponse) (TriggerListResponse, error))` to `(TriggerListResponse, func(context.Context, TriggerListResponse) (TriggerListResponse, error))`
- Function `NewManagedVirtualNetworkListResponsePage` parameter(s) have been changed from `(func(context.Context, ManagedVirtualNetworkListResponse) (ManagedVirtualNetworkListResponse, error))` to `(ManagedVirtualNetworkListResponse, func(context.Context, ManagedVirtualNetworkListResponse) (ManagedVirtualNetworkListResponse, error))`
- Function `NewIntegrationRuntimeListResponsePage` parameter(s) have been changed from `(func(context.Context, IntegrationRuntimeListResponse) (IntegrationRuntimeListResponse, error))` to `(IntegrationRuntimeListResponse, func(context.Context, IntegrationRuntimeListResponse) (IntegrationRuntimeListResponse, error))`
- Function `NewDatasetListResponsePage` parameter(s) have been changed from `(func(context.Context, DatasetListResponse) (DatasetListResponse, error))` to `(DatasetListResponse, func(context.Context, DatasetListResponse) (DatasetListResponse, error))`
- Function `NewPipelineListResponsePage` parameter(s) have been changed from `(func(context.Context, PipelineListResponse) (PipelineListResponse, error))` to `(PipelineListResponse, func(context.Context, PipelineListResponse) (PipelineListResponse, error))`
- Function `NewQueryDataFlowDebugSessionsResponsePage` parameter(s) have been changed from `(func(context.Context, QueryDataFlowDebugSessionsResponse) (QueryDataFlowDebugSessionsResponse, error))` to `(QueryDataFlowDebugSessionsResponse, func(context.Context, QueryDataFlowDebugSessionsResponse) (QueryDataFlowDebugSessionsResponse, error))`
- Function `NewManagedPrivateEndpointListResponsePage` parameter(s) have been changed from `(func(context.Context, ManagedPrivateEndpointListResponse) (ManagedPrivateEndpointListResponse, error))` to `(ManagedPrivateEndpointListResponse, func(context.Context, ManagedPrivateEndpointListResponse) (ManagedPrivateEndpointListResponse, error))`
- Function `NewDataFlowListResponsePage` parameter(s) have been changed from `(func(context.Context, DataFlowListResponse) (DataFlowListResponse, error))` to `(DataFlowListResponse, func(context.Context, DataFlowListResponse) (DataFlowListResponse, error))`
- Type of `NetezzaSource.PartitionOption` has been changed from `NetezzaPartitionOption` to `interface{}`
- Type of `SapHanaSource.PartitionOption` has been changed from `SapHanaPartitionOption` to `interface{}`
- Type of `SQLSource.PartitionOption` has been changed from `SQLPartitionOption` to `interface{}`
- Type of `AzureSQLSource.PartitionOption` has been changed from `SQLPartitionOption` to `interface{}`
- Type of `OracleSource.PartitionOption` has been changed from `OraclePartitionOption` to `interface{}`
- Type of `TeradataSource.PartitionOption` has been changed from `TeradataPartitionOption` to `interface{}`
- Type of `SapTableSource.PartitionOption` has been changed from `SapTablePartitionOption` to `interface{}`
- Type of `SQLMISource.PartitionOption` has been changed from `SQLPartitionOption` to `interface{}`
- Type of `SQLDWSource.PartitionOption` has been changed from `SQLPartitionOption` to `interface{}`
- Type of `SQLServerSource.PartitionOption` has been changed from `SQLPartitionOption` to `interface{}`
- Type of `ExecuteDataFlowActivityTypePropertiesCompute.ComputeType` has been changed from `DataFlowComputeType` to `interface{}`
- Type of `ExecuteDataFlowActivityTypePropertiesCompute.CoreCount` has been changed from `*int32` to `interface{}`

## New Content

- New const `TumblingWindowFrequencyMonth`
- New struct `CopyActivityLogSettings`
- New struct `LogLocationSettings`
- New struct `LogSettings`
- New field `ContinueOnError` in struct `ExecuteDataFlowActivityTypeProperties`
- New field `RunConcurrently` in struct `ExecuteDataFlowActivityTypeProperties`
- New field `TraceLevel` in struct `ExecuteDataFlowActivityTypeProperties`
- New field `ConnectionProperties` in struct `ConcurLinkedServiceTypeProperties`
- New field `LogSettings` in struct `CopyActivityTypeProperties`
- New field `AuthenticationType` in struct `AmazonS3LinkedServiceTypeProperties`
- New field `SessionToken` in struct `AmazonS3LinkedServiceTypeProperties`
