
## `github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/replicas` Documentation

The `replicas` SDK allows for interaction with the Azure Resource Manager Service `postgresql` (API Version `2017-12-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/replicas"
```


### Client Initialization

```go
client := replicas.NewReplicasClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ReplicasClient.ListByServer`

```go
ctx := context.TODO()
id := replicas.NewServerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serverValue")

read, err := client.ListByServer(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
