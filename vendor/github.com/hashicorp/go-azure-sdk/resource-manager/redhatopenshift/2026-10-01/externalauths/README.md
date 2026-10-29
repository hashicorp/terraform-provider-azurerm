
## `github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/externalauths` Documentation

The `externalauths` SDK allows for interaction with Azure Resource Manager `redhatopenshift` (API Version `2026-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/externalauths"
```


### Client Initialization

```go
client := externalauths.NewExternalAuthsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ExternalAuthsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := externalauths.NewExternalAuthID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName", "externalAuthName")

payload := externalauths.ExternalAuth{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ExternalAuthsClient.Delete`

```go
ctx := context.TODO()
id := externalauths.NewExternalAuthID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName", "externalAuthName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ExternalAuthsClient.Get`

```go
ctx := context.TODO()
id := externalauths.NewExternalAuthID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName", "externalAuthName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ExternalAuthsClient.ListByParent`

```go
ctx := context.TODO()
id := externalauths.NewHcpOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName")

// alternatively `client.ListByParent(ctx, id)` can be used to do batched pagination
items, err := client.ListByParentComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ExternalAuthsClient.Update`

```go
ctx := context.TODO()
id := externalauths.NewExternalAuthID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName", "externalAuthName")

payload := externalauths.ExternalAuthUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
