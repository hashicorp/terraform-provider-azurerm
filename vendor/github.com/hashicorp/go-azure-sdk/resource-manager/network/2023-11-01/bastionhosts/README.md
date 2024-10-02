
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/bastionhosts` Documentation

The `bastionhosts` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/bastionhosts"
```


### Client Initialization

```go
client := bastionhosts.NewBastionHostsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BastionHostsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := bastionhosts.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostName")

payload := bastionhosts.BastionHost{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BastionHostsClient.Delete`

```go
ctx := context.TODO()
id := bastionhosts.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BastionHostsClient.DeleteBastionShareableLink`

```go
ctx := context.TODO()
id := bastionhosts.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostName")

payload := bastionhosts.BastionShareableLinkListRequest{
	// ...
}


if err := client.DeleteBastionShareableLinkThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BastionHostsClient.DeleteBastionShareableLinkByToken`

```go
ctx := context.TODO()
id := bastionhosts.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostName")

payload := bastionhosts.BastionShareableLinkTokenListRequest{
	// ...
}


if err := client.DeleteBastionShareableLinkByTokenThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BastionHostsClient.DisconnectActiveSessions`

```go
ctx := context.TODO()
id := bastionhosts.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostName")

payload := bastionhosts.SessionIds{
	// ...
}


// alternatively `client.DisconnectActiveSessions(ctx, id, payload)` can be used to do batched pagination
items, err := client.DisconnectActiveSessionsComplete(ctx, id, payload)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BastionHostsClient.Get`

```go
ctx := context.TODO()
id := bastionhosts.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BastionHostsClient.GetActiveSessions`

```go
ctx := context.TODO()
id := bastionhosts.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostName")

// alternatively `client.GetActiveSessions(ctx, id)` can be used to do batched pagination
items, err := client.GetActiveSessionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BastionHostsClient.GetBastionShareableLink`

```go
ctx := context.TODO()
id := bastionhosts.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostName")

payload := bastionhosts.BastionShareableLinkListRequest{
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


### Example Usage: `BastionHostsClient.List`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BastionHostsClient.ListByResourceGroup`

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


### Example Usage: `BastionHostsClient.PutBastionShareableLink`

```go
ctx := context.TODO()
id := bastionhosts.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostName")

payload := bastionhosts.BastionShareableLinkListRequest{
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


### Example Usage: `BastionHostsClient.UpdateTags`

```go
ctx := context.TODO()
id := bastionhosts.NewBastionHostID("12345678-1234-9876-4563-123456789012", "example-resource-group", "bastionHostName")

payload := bastionhosts.TagsObject{
	// ...
}


if err := client.UpdateTagsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
