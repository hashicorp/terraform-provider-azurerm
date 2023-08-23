
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationsecuritygroups` Documentation

The `applicationsecuritygroups` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationsecuritygroups"
```


### Client Initialization

```go
client := applicationsecuritygroups.NewApplicationSecurityGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationSecurityGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := applicationsecuritygroups.NewApplicationSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationSecurityGroupValue")

payload := applicationsecuritygroups.ApplicationSecurityGroup{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationSecurityGroupsClient.Delete`

```go
ctx := context.TODO()
id := applicationsecuritygroups.NewApplicationSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationSecurityGroupValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationSecurityGroupsClient.Get`

```go
ctx := context.TODO()
id := applicationsecuritygroups.NewApplicationSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationSecurityGroupValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationSecurityGroupsClient.List`

```go
ctx := context.TODO()
id := applicationsecuritygroups.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationSecurityGroupsClient.ListAll`

```go
ctx := context.TODO()
id := applicationsecuritygroups.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationSecurityGroupsClient.UpdateTags`

```go
ctx := context.TODO()
id := applicationsecuritygroups.NewApplicationSecurityGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationSecurityGroupValue")

payload := applicationsecuritygroups.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
