
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/ddoscustompolicies` Documentation

The `ddoscustompolicies` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/ddoscustompolicies"
```


### Client Initialization

```go
client := ddoscustompolicies.NewDdosCustomPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DdosCustomPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := ddoscustompolicies.NewDdosCustomPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ddosCustomPolicyValue")

payload := ddoscustompolicies.DdosCustomPolicy{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `DdosCustomPoliciesClient.Delete`

```go
ctx := context.TODO()
id := ddoscustompolicies.NewDdosCustomPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ddosCustomPolicyValue")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `DdosCustomPoliciesClient.Get`

```go
ctx := context.TODO()
id := ddoscustompolicies.NewDdosCustomPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ddosCustomPolicyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DdosCustomPoliciesClient.UpdateTags`

```go
ctx := context.TODO()
id := ddoscustompolicies.NewDdosCustomPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "ddosCustomPolicyValue")

payload := ddoscustompolicies.TagsObject{
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
