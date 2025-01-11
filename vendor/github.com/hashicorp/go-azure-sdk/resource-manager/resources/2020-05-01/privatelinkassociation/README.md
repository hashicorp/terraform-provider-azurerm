
## `github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/privatelinkassociation` Documentation

The `privatelinkassociation` SDK allows for interaction with Azure Resource Manager `resources` (API Version `2020-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/privatelinkassociation"
```


### Client Initialization

```go
client := privatelinkassociation.NewPrivateLinkAssociationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateLinkAssociationClient.Delete`

```go
ctx := context.TODO()
id := privatelinkassociation.NewPrivateLinkAssociationID("groupId", "plaId")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkAssociationClient.Get`

```go
ctx := context.TODO()
id := privatelinkassociation.NewPrivateLinkAssociationID("groupId", "plaId")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkAssociationClient.List`

```go
ctx := context.TODO()
id := commonids.NewManagementGroupID("groupId")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateLinkAssociationClient.Put`

```go
ctx := context.TODO()
id := privatelinkassociation.NewPrivateLinkAssociationID("groupId", "plaId")

payload := privatelinkassociation.PrivateLinkAssociationObject{
	// ...
}


read, err := client.Put(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
