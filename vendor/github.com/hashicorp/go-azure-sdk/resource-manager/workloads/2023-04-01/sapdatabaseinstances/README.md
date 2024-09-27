
## `github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapdatabaseinstances` Documentation

The `sapdatabaseinstances` SDK allows for interaction with Azure Resource Manager `workloads` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapdatabaseinstances"
```


### Client Initialization

```go
client := sapdatabaseinstances.NewSAPDatabaseInstancesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SAPDatabaseInstancesClient.Create`

```go
ctx := context.TODO()
id := sapdatabaseinstances.NewDatabaseInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "databaseInstanceName")

payload := sapdatabaseinstances.SAPDatabaseInstance{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPDatabaseInstancesClient.Delete`

```go
ctx := context.TODO()
id := sapdatabaseinstances.NewDatabaseInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "databaseInstanceName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SAPDatabaseInstancesClient.Get`

```go
ctx := context.TODO()
id := sapdatabaseinstances.NewDatabaseInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "databaseInstanceName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SAPDatabaseInstancesClient.List`

```go
ctx := context.TODO()
id := sapdatabaseinstances.NewSapVirtualInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SAPDatabaseInstancesClient.StartInstance`

```go
ctx := context.TODO()
id := sapdatabaseinstances.NewDatabaseInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "databaseInstanceName")

if err := client.StartInstanceThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SAPDatabaseInstancesClient.StopInstance`

```go
ctx := context.TODO()
id := sapdatabaseinstances.NewDatabaseInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "databaseInstanceName")

payload := sapdatabaseinstances.StopRequest{
	// ...
}


if err := client.StopInstanceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SAPDatabaseInstancesClient.Update`

```go
ctx := context.TODO()
id := sapdatabaseinstances.NewDatabaseInstanceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "sapVirtualInstanceName", "databaseInstanceName")

payload := sapdatabaseinstances.UpdateSAPDatabaseInstanceRequest{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
