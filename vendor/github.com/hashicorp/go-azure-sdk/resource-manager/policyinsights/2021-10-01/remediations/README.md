
## `github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations` Documentation

The `remediations` SDK allows for interaction with Azure Resource Manager `policyinsights` (API Version `2021-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
```


### Client Initialization

```go
client := remediations.NewRemediationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RemediationsClient.CancelAtManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewProviders2RemediationID("managementGroupId", "remediationName")

read, err := client.CancelAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.CancelAtResource`

```go
ctx := context.TODO()
id := remediations.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationName")

read, err := client.CancelAtResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.CancelAtResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationName")

read, err := client.CancelAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.CancelAtSubscription`

```go
ctx := context.TODO()
id := remediations.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationName")

read, err := client.CancelAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.CreateOrUpdateAtManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewProviders2RemediationID("managementGroupId", "remediationName")

payload := remediations.Remediation{
	// ...
}


read, err := client.CreateOrUpdateAtManagementGroup(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.CreateOrUpdateAtResource`

```go
ctx := context.TODO()
id := remediations.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationName")

payload := remediations.Remediation{
	// ...
}


read, err := client.CreateOrUpdateAtResource(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.CreateOrUpdateAtResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationName")

payload := remediations.Remediation{
	// ...
}


read, err := client.CreateOrUpdateAtResourceGroup(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.CreateOrUpdateAtSubscription`

```go
ctx := context.TODO()
id := remediations.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationName")

payload := remediations.Remediation{
	// ...
}


read, err := client.CreateOrUpdateAtSubscription(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.DeleteAtManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewProviders2RemediationID("managementGroupId", "remediationName")

read, err := client.DeleteAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.DeleteAtResource`

```go
ctx := context.TODO()
id := remediations.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationName")

read, err := client.DeleteAtResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.DeleteAtResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationName")

read, err := client.DeleteAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.DeleteAtSubscription`

```go
ctx := context.TODO()
id := remediations.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationName")

read, err := client.DeleteAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.GetAtManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewProviders2RemediationID("managementGroupId", "remediationName")

read, err := client.GetAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.GetAtResource`

```go
ctx := context.TODO()
id := remediations.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationName")

read, err := client.GetAtResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.GetAtResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationName")

read, err := client.GetAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.GetAtSubscription`

```go
ctx := context.TODO()
id := remediations.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationName")

read, err := client.GetAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.ListDeploymentsAtManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewProviders2RemediationID("managementGroupId", "remediationName")

// alternatively `client.ListDeploymentsAtManagementGroup(ctx, id, remediations.DefaultListDeploymentsAtManagementGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListDeploymentsAtManagementGroupComplete(ctx, id, remediations.DefaultListDeploymentsAtManagementGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.ListDeploymentsAtResource`

```go
ctx := context.TODO()
id := remediations.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationName")

// alternatively `client.ListDeploymentsAtResource(ctx, id, remediations.DefaultListDeploymentsAtResourceOperationOptions())` can be used to do batched pagination
items, err := client.ListDeploymentsAtResourceComplete(ctx, id, remediations.DefaultListDeploymentsAtResourceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.ListDeploymentsAtResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationName")

// alternatively `client.ListDeploymentsAtResourceGroup(ctx, id, remediations.DefaultListDeploymentsAtResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListDeploymentsAtResourceGroupComplete(ctx, id, remediations.DefaultListDeploymentsAtResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.ListDeploymentsAtSubscription`

```go
ctx := context.TODO()
id := remediations.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationName")

// alternatively `client.ListDeploymentsAtSubscription(ctx, id, remediations.DefaultListDeploymentsAtSubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListDeploymentsAtSubscriptionComplete(ctx, id, remediations.DefaultListDeploymentsAtSubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.ListForManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewManagementGroupID("managementGroupId")

// alternatively `client.ListForManagementGroup(ctx, id, remediations.DefaultListForManagementGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListForManagementGroupComplete(ctx, id, remediations.DefaultListForManagementGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.ListForResource`

```go
ctx := context.TODO()
id := commonids.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.ListForResource(ctx, id, remediations.DefaultListForResourceOperationOptions())` can be used to do batched pagination
items, err := client.ListForResourceComplete(ctx, id, remediations.DefaultListForResourceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.ListForResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListForResourceGroup(ctx, id, remediations.DefaultListForResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.ListForResourceGroupComplete(ctx, id, remediations.DefaultListForResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.ListForSubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListForSubscription(ctx, id, remediations.DefaultListForSubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.ListForSubscriptionComplete(ctx, id, remediations.DefaultListForSubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
