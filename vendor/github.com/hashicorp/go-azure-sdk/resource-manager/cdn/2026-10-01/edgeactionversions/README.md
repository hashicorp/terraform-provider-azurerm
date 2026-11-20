
## `github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions` Documentation

The `edgeactionversions` SDK allows for interaction with Azure Resource Manager `cdn` (API Version `2026-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
```


### Client Initialization

```go
client := edgeactionversions.NewEdgeActionVersionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `EdgeActionVersionsClient.Create`

```go
ctx := context.TODO()
id := edgeactionversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "versionName")

payload := edgeactionversions.EdgeActionVersion{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeActionVersionsClient.Delete`

```go
ctx := context.TODO()
id := edgeactionversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "versionName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeActionVersionsClient.DeployVersionCode`

```go
ctx := context.TODO()
id := edgeactionversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "versionName")

payload := edgeactionversions.VersionCode{
	// ...
}


if err := client.DeployVersionCodeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeActionVersionsClient.Get`

```go
ctx := context.TODO()
id := edgeactionversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "versionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `EdgeActionVersionsClient.GetVersionCode`

```go
ctx := context.TODO()
id := edgeactionversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "versionName")

if err := client.GetVersionCodeThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeActionVersionsClient.ListByEdgeAction`

```go
ctx := context.TODO()
id := edgeactionversions.NewEdgeActionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName")

// alternatively `client.ListByEdgeAction(ctx, id)` can be used to do batched pagination
items, err := client.ListByEdgeActionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `EdgeActionVersionsClient.SwapDefault`

```go
ctx := context.TODO()
id := edgeactionversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "versionName")

if err := client.SwapDefaultThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `EdgeActionVersionsClient.Update`

```go
ctx := context.TODO()
id := edgeactionversions.NewVersionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "edgeActionName", "versionName")

payload := edgeactionversions.EdgeActionVersionUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
