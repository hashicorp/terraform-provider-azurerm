
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2025-09-01/containerinstance` Documentation

The `containerinstance` SDK allows for interaction with Azure Resource Manager `containerinstance` (API Version `2025-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2025-09-01/containerinstance"
```


### Client Initialization

```go
client := containerinstance.NewContainerInstanceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContainerInstanceClient.CGProfileCreateOrUpdate`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupProfileName")

payload := containerinstance.ContainerGroupProfile{
	// ...
}


read, err := client.CGProfileCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.CGProfileDelete`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupProfileName")

read, err := client.CGProfileDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.CGProfileGet`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupProfileName")

read, err := client.CGProfileGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.CGProfileGetByRevisionNumber`

```go
ctx := context.TODO()
id := containerinstance.NewRevisionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupProfileName", "revisionName")

read, err := client.CGProfileGetByRevisionNumber(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.CGProfileListAllRevisions`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupProfileName")

// alternatively `client.CGProfileListAllRevisions(ctx, id)` can be used to do batched pagination
items, err := client.CGProfileListAllRevisionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerInstanceClient.CGProfileUpdate`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupProfileID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupProfileName")

payload := containerinstance.ContainerGroupProfilePatch{
	// ...
}


read, err := client.CGProfileUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.CGProfilesListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.CGProfilesListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.CGProfilesListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerInstanceClient.CGProfilesListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.CGProfilesListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.CGProfilesListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsCreateOrUpdate`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName")

payload := containerinstance.ContainerGroup{
	// ...
}


if err := client.ContainerGroupsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsDelete`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName")

if err := client.ContainerGroupsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsGet`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName")

read, err := client.ContainerGroupsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsGetOutboundNetworkDependenciesEndpoints`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName")

read, err := client.ContainerGroupsGetOutboundNetworkDependenciesEndpoints(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ContainerGroupsList(ctx, id)` can be used to do batched pagination
items, err := client.ContainerGroupsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ContainerGroupsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ContainerGroupsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsRestart`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName")

if err := client.ContainerGroupsRestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsStart`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName")

if err := client.ContainerGroupsStartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsStop`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName")

read, err := client.ContainerGroupsStop(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsUpdate`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName")

payload := containerinstance.Resource{
	// ...
}


read, err := client.ContainerGroupsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.ContainersAttach`

```go
ctx := context.TODO()
id := containerinstance.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName", "containerName")

read, err := client.ContainersAttach(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.ContainersExecuteCommand`

```go
ctx := context.TODO()
id := containerinstance.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName", "containerName")

payload := containerinstance.ContainerExecRequest{
	// ...
}


read, err := client.ContainersExecuteCommand(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.ContainersListLogs`

```go
ctx := context.TODO()
id := containerinstance.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupName", "containerName")

read, err := client.ContainersListLogs(ctx, id, containerinstance.DefaultContainersListLogsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.LocationListCachedImages`

```go
ctx := context.TODO()
id := containerinstance.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.LocationListCachedImages(ctx, id)` can be used to do batched pagination
items, err := client.LocationListCachedImagesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerInstanceClient.LocationListCapabilities`

```go
ctx := context.TODO()
id := containerinstance.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.LocationListCapabilities(ctx, id)` can be used to do batched pagination
items, err := client.LocationListCapabilitiesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerInstanceClient.LocationListUsage`

```go
ctx := context.TODO()
id := containerinstance.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

read, err := client.LocationListUsage(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.NGroupsCreateOrUpdate`

```go
ctx := context.TODO()
id := containerinstance.NewNgroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ngroupName")

payload := containerinstance.NGroup{
	// ...
}


if err := client.NGroupsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.NGroupsDelete`

```go
ctx := context.TODO()
id := containerinstance.NewNgroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ngroupName")

if err := client.NGroupsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.NGroupsGet`

```go
ctx := context.TODO()
id := containerinstance.NewNgroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ngroupName")

read, err := client.NGroupsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.NGroupsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.NGroupsList(ctx, id)` can be used to do batched pagination
items, err := client.NGroupsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerInstanceClient.NGroupsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.NGroupsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.NGroupsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerInstanceClient.NGroupsRestart`

```go
ctx := context.TODO()
id := containerinstance.NewNgroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ngroupName")

if err := client.NGroupsRestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.NGroupsStart`

```go
ctx := context.TODO()
id := containerinstance.NewNgroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ngroupName")

if err := client.NGroupsStartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.NGroupsStop`

```go
ctx := context.TODO()
id := containerinstance.NewNgroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ngroupName")

read, err := client.NGroupsStop(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.NGroupsUpdate`

```go
ctx := context.TODO()
id := containerinstance.NewNgroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ngroupName")

payload := containerinstance.NGroupPatch{
	// ...
}


if err := client.NGroupsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.SubnetServiceAssociationLinkDelete`

```go
ctx := context.TODO()
id := commonids.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkName", "subnetName")

if err := client.SubnetServiceAssociationLinkDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```
