
## `github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simpolicy` Documentation

The `simpolicy` SDK allows for interaction with Azure Resource Manager `mobilenetwork` (API Version `2022-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simpolicy"
```


### Client Initialization

```go
client := simpolicy.NewSIMPolicyClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SIMPolicyClient.SimPoliciesCreateOrUpdate`

```go
ctx := context.TODO()
id := simpolicy.NewSimPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName", "simPolicyName")

payload := simpolicy.SimPolicy{
	// ...
}


if err := client.SimPoliciesCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `SIMPolicyClient.SimPoliciesDelete`

```go
ctx := context.TODO()
id := simpolicy.NewSimPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName", "simPolicyName")

if err := client.SimPoliciesDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SIMPolicyClient.SimPoliciesGet`

```go
ctx := context.TODO()
id := simpolicy.NewSimPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName", "simPolicyName")

read, err := client.SimPoliciesGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SIMPolicyClient.SimPoliciesUpdateTags`

```go
ctx := context.TODO()
id := simpolicy.NewSimPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "mobileNetworkName", "simPolicyName")

payload := simpolicy.TagsObject{
	// ...
}


read, err := client.SimPoliciesUpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
