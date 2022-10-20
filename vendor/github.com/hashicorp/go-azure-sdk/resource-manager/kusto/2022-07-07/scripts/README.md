
## `github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-07-07/scripts` Documentation

The `scripts` SDK allows for interaction with the Azure Resource Manager Service `kusto` (API Version `2022-07-07`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2022-07-07/scripts"
```


### Client Initialization

```go
client := scripts.NewScriptsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ScriptsClient.ScriptsListByDatabase`

```go
ctx := context.TODO()
id := scripts.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue", "databaseValue")

read, err := client.ScriptsListByDatabase(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
