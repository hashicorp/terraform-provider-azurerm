
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2023-05-01/containerinstance` Documentation

The `containerinstance` SDK allows for interaction with the Azure Resource Manager Service `containerinstance` (API Version `2023-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2023-05-01/containerinstance"
```


### Client Initialization

```go
client := containerinstance.NewContainerInstanceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsCreateOrUpdate`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue")

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
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue")

if err := client.ContainerGroupsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsGet`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue")

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
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue")

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
id := containerinstance.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

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
id := containerinstance.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

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
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue")

if err := client.ContainerGroupsRestartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsStart`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue")

if err := client.ContainerGroupsStartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerInstanceClient.ContainerGroupsStop`

```go
ctx := context.TODO()
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue")

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
id := containerinstance.NewContainerGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue")

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
id := containerinstance.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue", "containerValue")

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
id := containerinstance.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue", "containerValue")

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
id := containerinstance.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerGroupValue", "containerValue")

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
id := containerinstance.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

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
id := containerinstance.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

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
id := containerinstance.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

read, err := client.LocationListUsage(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerInstanceClient.SubnetServiceAssociationLinkDelete`

```go
ctx := context.TODO()
id := containerinstance.NewSubnetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "virtualNetworkValue", "subnetValue")

if err := client.SubnetServiceAssociationLinkDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```
