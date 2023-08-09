
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/bastionshareablelink` Documentation

The `bastionshareablelink` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/bastionshareablelink"
```


### Client Initialization

```go
client := bastionshareablelink.NewBastionShareableLinkClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BastionShareableLinkClient.DeleteBastionShareableLink`

```go
ctx := context.TODO()
id := bastionshareablelink.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostValue")

payload := bastionshareablelink.BastionShareableLinkListRequest{
	// ...
}


if err := client.DeleteBastionShareableLinkThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BastionShareableLinkClient.GetBastionShareableLink`

```go
ctx := context.TODO()
id := bastionshareablelink.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostValue")

payload := bastionshareablelink.BastionShareableLinkListRequest{
	// ...
}


// alternatively `client.GetBastionShareableLink(ctx, id, payload)` can be used to do batched pagination
items, err := client.GetBastionShareableLinkComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BastionShareableLinkClient.PutBastionShareableLink`

```go
ctx := context.TODO()
id := bastionshareablelink.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostValue")

payload := bastionshareablelink.BastionShareableLinkListRequest{
	// ...
}


// alternatively `client.PutBastionShareableLink(ctx, id, payload)` can be used to do batched pagination
items, err := client.PutBastionShareableLinkComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
