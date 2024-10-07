
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationgateways` Documentation

The `applicationgateways` SDK allows for interaction with Azure Resource Manager `network` (API Version `2023-11-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationgateways"
```


### Client Initialization

```go
client := applicationgateways.NewApplicationGatewaysClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ApplicationGatewaysClient.BackendHealth`

```go
ctx := context.TODO()
id := applicationgateways.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName")

if err := client.BackendHealthThenPoll(ctx, id, applicationgateways.DefaultBackendHealthOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationGatewaysClient.BackendHealthOnDemand`

```go
ctx := context.TODO()
id := applicationgateways.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName")

payload := applicationgateways.ApplicationGatewayOnDemandProbe{
	// ...
}


if err := client.BackendHealthOnDemandThenPoll(ctx, id, payload, applicationgateways.DefaultBackendHealthOnDemandOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationGatewaysClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := applicationgateways.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName")

payload := applicationgateways.ApplicationGateway{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationGatewaysClient.Delete`

```go
ctx := context.TODO()
id := applicationgateways.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationGatewaysClient.Get`

```go
ctx := context.TODO()
id := applicationgateways.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationGatewaysClient.GetSslPredefinedPolicy`

```go
ctx := context.TODO()
id := applicationgateways.NewPredefinedPolicyID("12345678-1234-9876-4563-123456789012", "predefinedPolicyName")

read, err := client.GetSslPredefinedPolicy(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationGatewaysClient.List`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationGatewaysClient.ListAll`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAll(ctx, id)` can be used to do batched pagination
items, err := client.ListAllComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationGatewaysClient.ListAvailableRequestHeaders`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListAvailableRequestHeaders(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationGatewaysClient.ListAvailableResponseHeaders`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListAvailableResponseHeaders(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationGatewaysClient.ListAvailableServerVariables`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListAvailableServerVariables(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationGatewaysClient.ListAvailableSslOptions`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListAvailableSslOptions(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationGatewaysClient.ListAvailableSslPredefinedPolicies`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListAvailableSslPredefinedPolicies(ctx, id)` can be used to do batched pagination
items, err := client.ListAvailableSslPredefinedPoliciesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ApplicationGatewaysClient.ListAvailableWafRuleSets`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListAvailableWafRuleSets(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ApplicationGatewaysClient.Start`

```go
ctx := context.TODO()
id := applicationgateways.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName")

if err := client.StartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationGatewaysClient.Stop`

```go
ctx := context.TODO()
id := applicationgateways.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName")

if err := client.StopThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ApplicationGatewaysClient.UpdateTags`

```go
ctx := context.TODO()
id := applicationgateways.NewApplicationGatewayID("12345678-1234-9876-4563-123456789012", "example-resource-group", "applicationGatewayName")

payload := applicationgateways.TagsObject{
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
