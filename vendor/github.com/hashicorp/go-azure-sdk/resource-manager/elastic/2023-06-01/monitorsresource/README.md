
## `github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/monitorsresource` Documentation

The `monitorsresource` SDK allows for interaction with the Azure Resource Manager Service `elastic` (API Version `2023-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/monitorsresource"
```


### Client Initialization

```go
client := monitorsresource.NewMonitorsResourceClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `MonitorsResourceClient.MonitorsCreate`

```go
ctx := context.TODO()
id := monitorsresource.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue")

payload := monitorsresource.ElasticMonitorResource{
	// ...
}


if err := client.MonitorsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `MonitorsResourceClient.MonitorsDelete`

```go
ctx := context.TODO()
id := monitorsresource.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue")

if err := client.MonitorsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `MonitorsResourceClient.MonitorsGet`

```go
ctx := context.TODO()
id := monitorsresource.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue")

read, err := client.MonitorsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `MonitorsResourceClient.MonitorsList`

```go
ctx := context.TODO()
id := monitorsresource.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.MonitorsList(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MonitorsResourceClient.MonitorsListByResourceGroup`

```go
ctx := context.TODO()
id := monitorsresource.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.MonitorsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `MonitorsResourceClient.MonitorsUpdate`

```go
ctx := context.TODO()
id := monitorsresource.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorValue")

payload := monitorsresource.ElasticMonitorResourceUpdateParameters{
	// ...
}


read, err := client.MonitorsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
