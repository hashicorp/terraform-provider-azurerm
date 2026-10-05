
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/namedvalue` Documentation

The `namedvalue` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2022-08-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/namedvalue"
```


### Client Initialization

```go
client := namedvalue.NewNamedValueClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `NamedValueClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "namedValueId")

payload := namedvalue.NamedValueCreateContract{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload, namedvalue.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `NamedValueClient.Delete`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "namedValueId")

read, err := client.Delete(ctx, id, namedvalue.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.Get`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "namedValueId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.GetEntityTag`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "namedValueId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.ListByService`

```go
ctx := context.TODO()
id := namedvalue.NewServiceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName")

// alternatively `client.ListByService(ctx, id, namedvalue.DefaultListByServiceOperationOptions())` can be used to do batched pagination
items, err := client.ListByServiceComplete(ctx, id, namedvalue.DefaultListByServiceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `NamedValueClient.ListValue`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "namedValueId")

read, err := client.ListValue(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `NamedValueClient.RefreshSecret`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "namedValueId")

if err := client.RefreshSecretThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `NamedValueClient.Update`

```go
ctx := context.TODO()
id := namedvalue.NewNamedValueID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "namedValueId")

payload := namedvalue.NamedValueUpdateParameters{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload, namedvalue.DefaultUpdateOperationOptions()); err != nil {
	// handle the error
}
```
