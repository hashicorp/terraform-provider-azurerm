
## `github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/accesspolicies` Documentation

The `accesspolicies` SDK allows for interaction with the Azure Resource Manager Service `timeseriesinsights` (API Version `2020-05-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/accesspolicies"
```


### Client Initialization

```go
client := accesspolicies.NewAccessPoliciesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `AccessPoliciesClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := accesspolicies.NewAccessPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue", "accessPolicyValue")

payload := accesspolicies.AccessPolicyCreateOrUpdateParameters{
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


### Example Usage: `AccessPoliciesClient.Delete`

```go
ctx := context.TODO()
id := accesspolicies.NewAccessPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue", "accessPolicyValue")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccessPoliciesClient.Get`

```go
ctx := context.TODO()
id := accesspolicies.NewAccessPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue", "accessPolicyValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccessPoliciesClient.ListByEnvironment`

```go
ctx := context.TODO()
id := accesspolicies.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue")

read, err := client.ListByEnvironment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `AccessPoliciesClient.Update`

```go
ctx := context.TODO()
id := accesspolicies.NewAccessPolicyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "environmentValue", "accessPolicyValue")

payload := accesspolicies.AccessPolicyUpdateParameters{
	// ...
}


read, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
