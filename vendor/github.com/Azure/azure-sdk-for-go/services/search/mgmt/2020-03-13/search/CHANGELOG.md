Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewListQueryKeysResultPage` parameter(s) have been changed from `(func(context.Context, ListQueryKeysResult) (ListQueryKeysResult, error))` to `(ListQueryKeysResult, func(context.Context, ListQueryKeysResult) (ListQueryKeysResult, error))`
- Function `NewPrivateEndpointConnectionListResultPage` parameter(s) have been changed from `(func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error))` to `(PrivateEndpointConnectionListResult, func(context.Context, PrivateEndpointConnectionListResult) (PrivateEndpointConnectionListResult, error))`
- Function `NewServiceListResultPage` parameter(s) have been changed from `(func(context.Context, ServiceListResult) (ServiceListResult, error))` to `(ServiceListResult, func(context.Context, ServiceListResult) (ServiceListResult, error))`
- Const `SharedPrivateLinkResourceStatusTimeout` has been removed
- Const `SharedPrivateLinkResourceStatusDisconnected` has been removed
- Const `SharedPrivateLinkResourceStatusApproved` has been removed
- Const `SharedPrivateLinkResourceStatusRejected` has been removed
- Const `SharedPrivateLinkResourceStatusPending` has been removed
- Function `*SharedPrivateLinkResourceListResultPage.NextWithContext` has been removed
- Function `SharedPrivateLinkResourceListResultIterator.Response` has been removed
- Function `SharedPrivateLinkResourcesClient.ListByServiceResponder` has been removed
- Function `SharedPrivateLinkResourcesClient.GetResponder` has been removed
- Function `*SharedPrivateLinkResourceListResultPage.Next` has been removed
- Function `SharedPrivateLinkResourcesClient.DeleteSender` has been removed
- Function `SharedPrivateLinkResourcesClient.GetSender` has been removed
- Function `NewSharedPrivateLinkResourcesClient` has been removed
- Function `SharedPrivateLinkResourceListResult.MarshalJSON` has been removed
- Function `SharedPrivateLinkResourcesClient.ListByServiceSender` has been removed
- Function `SharedPrivateLinkResourceListResultIterator.Value` has been removed
- Function `SharedPrivateLinkResourceListResult.IsEmpty` has been removed
- Function `SharedPrivateLinkResource.MarshalJSON` has been removed
- Function `NewSharedPrivateLinkResourceListResultPage` has been removed
- Function `*SharedPrivateLinkResourceListResultIterator.Next` has been removed
- Function `SharedPrivateLinkResourcesClient.Delete` has been removed
- Function `NewSharedPrivateLinkResourcesClientWithBaseURI` has been removed
- Function `SharedPrivateLinkResourcesClient.ListByService` has been removed
- Function `SharedPrivateLinkResourcesClient.DeleteResponder` has been removed
- Function `SharedPrivateLinkResourcesClient.DeletePreparer` has been removed
- Function `SharedPrivateLinkResourcesClient.CreateOrUpdatePreparer` has been removed
- Function `NewSharedPrivateLinkResourceListResultIterator` has been removed
- Function `SharedPrivateLinkResourcesClient.ListByServiceComplete` has been removed
- Function `SharedPrivateLinkResourcesClient.CreateOrUpdateSender` has been removed
- Function `SharedPrivateLinkResourceListResultPage.Response` has been removed
- Function `PossibleSharedPrivateLinkResourceStatusValues` has been removed
- Function `SharedPrivateLinkResourcesClient.ListByServicePreparer` has been removed
- Function `SharedPrivateLinkResourceListResultPage.Values` has been removed
- Function `*SharedPrivateLinkResourceListResultIterator.NextWithContext` has been removed
- Function `SharedPrivateLinkResourcesClient.Get` has been removed
- Function `SharedPrivateLinkResourcesClient.CreateOrUpdate` has been removed
- Function `SharedPrivateLinkResourceListResultIterator.NotDone` has been removed
- Function `SharedPrivateLinkResourcesClient.CreateOrUpdateResponder` has been removed
- Function `SharedPrivateLinkResourcesClient.GetPreparer` has been removed
- Function `SharedPrivateLinkResourceListResultPage.NotDone` has been removed
- Struct `ShareablePrivateLinkResourceProperties` has been removed
- Struct `ShareablePrivateLinkResourceType` has been removed
- Struct `SharedPrivateLinkResource` has been removed
- Struct `SharedPrivateLinkResourceListResult` has been removed
- Struct `SharedPrivateLinkResourceListResultIterator` has been removed
- Struct `SharedPrivateLinkResourceListResultPage` has been removed
- Struct `SharedPrivateLinkResourceProperties` has been removed
- Struct `SharedPrivateLinkResourcesClient` has been removed
- Field `SharedPrivateLinkResources` of struct `ServiceProperties` has been removed
- Field `ShareablePrivateLinkResourceTypes` of struct `PrivateLinkResourceProperties` has been removed
