
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/scripts` Documentation

The `scripts` SDK allows for interaction with the Azure Resource Manager Service `kusto` (API Version `2023-05-02`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-05-02/scripts"
```


### Client Initialization

```go
client := scripts.NewScriptsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScriptsClient.CheckNameAvailability`

```go
ctx := context.TODO()
id := scripts.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

payload := scripts.ScriptCheckNameRequest{
	// ...
}


read, err := client.CheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScriptsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := scripts.NewScriptID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "scriptValue")

payload := scripts.Script{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ScriptsClient.Delete`

```go
ctx := context.TODO()
id := scripts.NewScriptID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "scriptValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ScriptsClient.Get`

```go
ctx := context.TODO()
id := scripts.NewScriptID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "scriptValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScriptsClient.ListByDatabase`

```go
ctx := context.TODO()
id := scripts.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

read, err := client.ListByDatabase(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ScriptsClient.Update`

```go
ctx := context.TODO()
id := scripts.NewScriptID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue", "scriptValue")

payload := scripts.Script{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
