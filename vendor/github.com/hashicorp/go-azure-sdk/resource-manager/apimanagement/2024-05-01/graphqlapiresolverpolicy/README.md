
## `github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/graphqlapiresolverpolicy` Documentation

The `graphqlapiresolverpolicy` SDK allows for interaction with Azure Resource Manager `apimanagement` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/graphqlapiresolverpolicy"
```


### Client Initialization

```go
client := graphqlapiresolverpolicy.NewGraphQLApiResolverPolicyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GraphQLApiResolverPolicyClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := graphqlapiresolverpolicy.NewResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "resolverId")

payload := graphqlapiresolverpolicy.PolicyContract{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, graphqlapiresolverpolicy.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GraphQLApiResolverPolicyClient.Delete`

```go
ctx := context.TODO()
id := graphqlapiresolverpolicy.NewResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "resolverId")

read, err := client.Delete(ctx, id, graphqlapiresolverpolicy.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GraphQLApiResolverPolicyClient.Get`

```go
ctx := context.TODO()
id := graphqlapiresolverpolicy.NewResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "resolverId")

read, err := client.Get(ctx, id, graphqlapiresolverpolicy.DefaultGetOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GraphQLApiResolverPolicyClient.GetEntityTag`

```go
ctx := context.TODO()
id := graphqlapiresolverpolicy.NewResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "resolverId")

read, err := client.GetEntityTag(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `GraphQLApiResolverPolicyClient.ListByResolver`

```go
ctx := context.TODO()
id := graphqlapiresolverpolicy.NewResolverID("12345678-1234-9876-4563-123456789012", "example-resource-group", "serviceName", "apiId", "resolverId")

// alternatively `client.ListByResolver(ctx, id)` can be used to do batched pagination
items, err := client.ListByResolverComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
