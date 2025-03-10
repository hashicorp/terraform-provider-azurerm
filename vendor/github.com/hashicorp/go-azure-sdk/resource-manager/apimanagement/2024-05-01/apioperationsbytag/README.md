
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apioperationsbytag` Documentation

The `apioperationsbytag` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apioperationsbytag"
```


### Client Initialization

```go
client := apioperationsbytag.NewApiOperationsByTagClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApiOperationsByTagClient.OperationListByTags`

```go
ctx := context.TODO()
id := apioperationsbytag.NewApiID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId")

// alternatively `client.OperationListByTags(ctx, id, apioperationsbytag.DefaultOperationListByTagsOperationOptions())` can be used to do batched pagination
items, err := client.OperationListByTagsComplete(ctx, id, apioperationsbytag.DefaultOperationListByTagsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
