
## `github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-04-03/testlines` Documentation

The `testlines` SDK allows for interaction with Azure Resource Manager `voiceservices` (API Version `2023-04-03`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/voiceservices/2023-04-03/testlines"
```


### Client Initialization

```go
client := testlines.NewTestLinesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `TestLinesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := testlines.NewTestLineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationsGatewayName", "testLineName")

payload := testlines.TestLine{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `TestLinesClient.Delete`

```go
ctx := context.TODO()
id := testlines.NewTestLineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationsGatewayName", "testLineName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `TestLinesClient.Get`

```go
ctx := context.TODO()
id := testlines.NewTestLineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationsGatewayName", "testLineName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `TestLinesClient.ListByCommunicationsGateway`

```go
ctx := context.TODO()
id := testlines.NewCommunicationsGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationsGatewayName")

// alternatively `client.ListByCommunicationsGateway(ctx, id)` can be used to do batched pagination
items, err := client.ListByCommunicationsGatewayComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `TestLinesClient.Update`

```go
ctx := context.TODO()
id := testlines.NewTestLineID("12345678-1234-9876-4563-123456789012", "example-resource-group", "communicationsGatewayName", "testLineName")

payload := testlines.TestLineUpdate{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
