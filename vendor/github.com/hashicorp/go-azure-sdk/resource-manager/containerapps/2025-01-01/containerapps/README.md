
## `github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/containerapps` Documentation

The `containerapps` SDK allows for interaction with Azure Resource Manager `containerapps` (API Version `2025-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/containerapps"
```


### Client Initialization

```go
client := containerapps.NewContainerAppsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ContainerAppsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

payload := containerapps.ContainerApp{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerAppsClient.Delete`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerAppsClient.DiagnosticsGetDetector`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppDetectorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName", "detectorName")

read, err := client.DiagnosticsGetDetector(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsClient.DiagnosticsGetRevision`

```go
ctx := context.TODO()
id := containerapps.NewRevisionsApiRevisionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName", "revisionName")

read, err := client.DiagnosticsGetRevision(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsClient.DiagnosticsGetRoot`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

read, err := client.DiagnosticsGetRoot(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsClient.DiagnosticsListDetectors`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

// alternatively `client.DiagnosticsListDetectors(ctx, id)` can be used to do batched pagination
items, err := client.DiagnosticsListDetectorsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerAppsClient.DiagnosticsListRevisions`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

// alternatively `client.DiagnosticsListRevisions(ctx, id, containerapps.DefaultDiagnosticsListRevisionsOperationOptions())` can be used to do batched pagination
items, err := client.DiagnosticsListRevisionsComplete(ctx, id, containerapps.DefaultDiagnosticsListRevisionsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerAppsClient.Get`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsClient.GetAuthToken`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

read, err := client.GetAuthToken(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerAppsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ContainerAppsClient.ListCustomHostNameAnalysis`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

read, err := client.ListCustomHostNameAnalysis(ctx, id, containerapps.DefaultListCustomHostNameAnalysisOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsClient.ListSecrets`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

read, err := client.ListSecrets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ContainerAppsClient.Start`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerAppsClient.Stop`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ContainerAppsClient.Update`

```go
ctx := context.TODO()
id := containerapps.NewContainerAppID("12345678-1234-9876-4563-123456789012", "example-resource-group", "containerAppName")

payload := containerapps.ContainerApp{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
