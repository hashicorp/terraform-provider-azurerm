
## `github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountcapabilityhost` Documentation

The `accountcapabilityhost` SDK allows for interaction with Azure Resource Manager `cognitive` (API Version `2026-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountcapabilityhost"
```


### Client Initialization

```go
client := accountcapabilityhost.NewAccountCapabilityHostClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccountCapabilityHostClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := accountcapabilityhost.NewCapabilityHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "capabilityHostName")

payload := accountcapabilityhost.CapabilityHost{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `AccountCapabilityHostClient.Delete`

```go
ctx := context.TODO()
id := accountcapabilityhost.NewCapabilityHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "capabilityHostName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `AccountCapabilityHostClient.Get`

```go
ctx := context.TODO()
id := accountcapabilityhost.NewCapabilityHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName", "capabilityHostName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccountCapabilityHostClient.List`

```go
ctx := context.TODO()
id := accountcapabilityhost.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "accountName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
