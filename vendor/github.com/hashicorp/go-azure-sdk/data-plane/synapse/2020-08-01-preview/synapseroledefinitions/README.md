
## `github.com/hashicorp/go-azure-sdk/data-plane/synapse/2020-08-01-preview/synapseroledefinitions` Documentation

The `synapseroledefinitions` SDK allows for interaction with <unknown source data type> `synapse` (API Version `2020-08-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/synapse/2020-08-01-preview/synapseroledefinitions"
```


### Client Initialization

```go
client := synapseroledefinitions.NewSynapseRoleDefinitionsClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `SynapseRoleDefinitionsClient.RoleDefinitionsGetRoleDefinitionById`

```go
ctx := context.TODO()
id := synapseroledefinitions.NewRoleDefinitionID("roleDefinitionId")

read, err := client.RoleDefinitionsGetRoleDefinitionById(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SynapseRoleDefinitionsClient.RoleDefinitionsListRoleDefinitions`

```go
ctx := context.TODO()


read, err := client.RoleDefinitionsListRoleDefinitions(ctx, synapseroledefinitions.DefaultRoleDefinitionsListRoleDefinitionsOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
