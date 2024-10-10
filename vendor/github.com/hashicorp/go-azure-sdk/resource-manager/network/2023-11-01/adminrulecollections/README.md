
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/adminrulecollections` Documentation

The `adminrulecollections` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/adminrulecollections"
```


### Client Initialization

```go
client := adminrulecollections.NewAdminRuleCollectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AdminRuleCollectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := adminrulecollections.NewRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityAdminConfigurationName", "ruleCollectionName")

payload := adminrulecollections.AdminRuleCollection{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AdminRuleCollectionsClient.Delete`

```go
ctx := context.TODO()
id := adminrulecollections.NewRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityAdminConfigurationName", "ruleCollectionName")

if err := client.DeleteThenPoll(ctx, id, adminrulecollections.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `AdminRuleCollectionsClient.Get`

```go
ctx := context.TODO()
id := adminrulecollections.NewRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityAdminConfigurationName", "ruleCollectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AdminRuleCollectionsClient.List`

```go
ctx := context.TODO()
id := adminrulecollections.NewSecurityAdminConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityAdminConfigurationName")

// alternatively `client.List(ctx, id, adminrulecollections.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, adminrulecollections.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
