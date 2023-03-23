
## `github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations` Documentation

The `remediations` SDK allows for interaction with the Azure Resource Manager Service `policyinsights` (API Version `2021-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
```


### Client Initialization

```go
client := remediations.NewRemediationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RemediationsClient.RemediationsCancelAtManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewProviders2RemediationID("managementGroupIdValue", "remediationValue")

read, err := client.RemediationsCancelAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsCancelAtResource`

```go
ctx := context.TODO()
id := remediations.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationValue")

read, err := client.RemediationsCancelAtResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsCancelAtResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationValue")

read, err := client.RemediationsCancelAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsCancelAtSubscription`

```go
ctx := context.TODO()
id := remediations.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationValue")

read, err := client.RemediationsCancelAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsCreateOrUpdateAtManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewProviders2RemediationID("managementGroupIdValue", "remediationValue")

payload := remediations.Remediation{
	// ...
}


read, err := client.RemediationsCreateOrUpdateAtManagementGroup(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsCreateOrUpdateAtResource`

```go
ctx := context.TODO()
id := remediations.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationValue")

payload := remediations.Remediation{
	// ...
}


read, err := client.RemediationsCreateOrUpdateAtResource(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsCreateOrUpdateAtResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationValue")

payload := remediations.Remediation{
	// ...
}


read, err := client.RemediationsCreateOrUpdateAtResourceGroup(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsCreateOrUpdateAtSubscription`

```go
ctx := context.TODO()
id := remediations.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationValue")

payload := remediations.Remediation{
	// ...
}


read, err := client.RemediationsCreateOrUpdateAtSubscription(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsDeleteAtManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewProviders2RemediationID("managementGroupIdValue", "remediationValue")

read, err := client.RemediationsDeleteAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsDeleteAtResource`

```go
ctx := context.TODO()
id := remediations.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationValue")

read, err := client.RemediationsDeleteAtResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsDeleteAtResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationValue")

read, err := client.RemediationsDeleteAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsDeleteAtSubscription`

```go
ctx := context.TODO()
id := remediations.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationValue")

read, err := client.RemediationsDeleteAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsGetAtManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewProviders2RemediationID("managementGroupIdValue", "remediationValue")

read, err := client.RemediationsGetAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsGetAtResource`

```go
ctx := context.TODO()
id := remediations.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationValue")

read, err := client.RemediationsGetAtResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsGetAtResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationValue")

read, err := client.RemediationsGetAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsGetAtSubscription`

```go
ctx := context.TODO()
id := remediations.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationValue")

read, err := client.RemediationsGetAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RemediationsClient.RemediationsListDeploymentsAtManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewProviders2RemediationID("managementGroupIdValue", "remediationValue")

// alternatively `client.RemediationsListDeploymentsAtManagementGroup(ctx, id, remediations.DefaultRemediationsListDeploymentsAtManagementGroupOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListDeploymentsAtManagementGroupComplete(ctx, id, remediations.DefaultRemediationsListDeploymentsAtManagementGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.RemediationsListDeploymentsAtResource`

```go
ctx := context.TODO()
id := remediations.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationValue")

// alternatively `client.RemediationsListDeploymentsAtResource(ctx, id, remediations.DefaultRemediationsListDeploymentsAtResourceOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListDeploymentsAtResourceComplete(ctx, id, remediations.DefaultRemediationsListDeploymentsAtResourceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.RemediationsListDeploymentsAtResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationValue")

// alternatively `client.RemediationsListDeploymentsAtResourceGroup(ctx, id, remediations.DefaultRemediationsListDeploymentsAtResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListDeploymentsAtResourceGroupComplete(ctx, id, remediations.DefaultRemediationsListDeploymentsAtResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.RemediationsListDeploymentsAtSubscription`

```go
ctx := context.TODO()
id := remediations.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationValue")

// alternatively `client.RemediationsListDeploymentsAtSubscription(ctx, id, remediations.DefaultRemediationsListDeploymentsAtSubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListDeploymentsAtSubscriptionComplete(ctx, id, remediations.DefaultRemediationsListDeploymentsAtSubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.RemediationsListForManagementGroup`

```go
ctx := context.TODO()
id := remediations.NewManagementGroupID("managementGroupIdValue")

// alternatively `client.RemediationsListForManagementGroup(ctx, id, remediations.DefaultRemediationsListForManagementGroupOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListForManagementGroupComplete(ctx, id, remediations.DefaultRemediationsListForManagementGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.RemediationsListForResource`

```go
ctx := context.TODO()
id := remediations.NewScopeID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group")

// alternatively `client.RemediationsListForResource(ctx, id, remediations.DefaultRemediationsListForResourceOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListForResourceComplete(ctx, id, remediations.DefaultRemediationsListForResourceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.RemediationsListForResourceGroup`

```go
ctx := context.TODO()
id := remediations.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.RemediationsListForResourceGroup(ctx, id, remediations.DefaultRemediationsListForResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListForResourceGroupComplete(ctx, id, remediations.DefaultRemediationsListForResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RemediationsClient.RemediationsListForSubscription`

```go
ctx := context.TODO()
id := remediations.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.RemediationsListForSubscription(ctx, id, remediations.DefaultRemediationsListForSubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListForSubscriptionComplete(ctx, id, remediations.DefaultRemediationsListForSubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
