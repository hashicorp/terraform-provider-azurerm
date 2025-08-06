
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/securityuserrulecollections` Documentation

The `securityuserrulecollections` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/securityuserrulecollections"
```


### Client Initialization

```go
client := securityuserrulecollections.NewSecurityUserRuleCollectionsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SecurityUserRuleCollectionsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := securityuserrulecollections.NewSecurityUserConfigurationRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName", "ruleCollectionName")

payload := securityuserrulecollections.SecurityUserRuleCollection{
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


### Example Usage: `SecurityUserRuleCollectionsClient.Delete`

```go
ctx := context.TODO()
id := securityuserrulecollections.NewSecurityUserConfigurationRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName", "ruleCollectionName")

if err := client.DeleteThenPoll(ctx, id, securityuserrulecollections.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `SecurityUserRuleCollectionsClient.Get`

```go
ctx := context.TODO()
id := securityuserrulecollections.NewSecurityUserConfigurationRuleCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName", "ruleCollectionName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SecurityUserRuleCollectionsClient.List`

```go
ctx := context.TODO()
id := securityuserrulecollections.NewSecurityUserConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "networkManagerName", "securityUserConfigurationName")

// alternatively `client.List(ctx, id, securityuserrulecollections.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, securityuserrulecollections.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
