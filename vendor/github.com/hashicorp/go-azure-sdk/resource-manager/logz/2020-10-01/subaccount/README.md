
## `github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/subaccount` Documentation

The `subaccount` SDK allows for interaction with the Azure Resource Manager Service `logz` (API Version `2020-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/logz/2020-10-01/subaccount"
```


### Client Initialization

```go
client := subaccount.NewSubAccountClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SubAccountClient.Create`

```go
ctx := context.TODO()
id := subaccount.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "accountValue")

payload := subaccount.LogzMonitorResource{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SubAccountClient.Delete`

```go
ctx := context.TODO()
id := subaccount.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "accountValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SubAccountClient.Get`

```go
ctx := context.TODO()
id := subaccount.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "accountValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SubAccountClient.List`

```go
ctx := context.TODO()
id := subaccount.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SubAccountClient.ListMonitoredResources`

```go
ctx := context.TODO()
id := subaccount.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "accountValue")

// alternatively `client.ListMonitoredResources(ctx, id)` can be used to do batched pagination
items, err := client.ListMonitoredResourcesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SubAccountClient.Update`

```go
ctx := context.TODO()
id := subaccount.NewAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue", "accountValue")

payload := subaccount.LogzMonitorResourceUpdateParameters{
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
