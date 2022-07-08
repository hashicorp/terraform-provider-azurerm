
## `github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/policyinsights` Documentation

The `policyinsights` SDK allows for interaction with the Azure Resource Manager Service `policyinsights` (API Version `2021-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/policyinsights"
```


### Client Initialization

```go
client := policyinsights.NewPolicyInsightsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PolicyInsightsClient.RemediationsCancelAtManagementGroup`

```go
ctx := context.TODO()
id := policyinsights.NewProviders2RemediationID("managementGroupIdValue", "remediationValue")

read, err := client.RemediationsCancelAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsCancelAtResource`

```go
ctx := context.TODO()
id := policyinsights.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationValue")

read, err := client.RemediationsCancelAtResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsCancelAtResourceGroup`

```go
ctx := context.TODO()
id := policyinsights.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationValue")

read, err := client.RemediationsCancelAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsCancelAtSubscription`

```go
ctx := context.TODO()
id := policyinsights.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationValue")

read, err := client.RemediationsCancelAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsCreateOrUpdateAtManagementGroup`

```go
ctx := context.TODO()
id := policyinsights.NewProviders2RemediationID("managementGroupIdValue", "remediationValue")

payload := policyinsights.Remediation{
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


### Example Usage: `PolicyInsightsClient.RemediationsCreateOrUpdateAtResource`

```go
ctx := context.TODO()
id := policyinsights.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationValue")

payload := policyinsights.Remediation{
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


### Example Usage: `PolicyInsightsClient.RemediationsCreateOrUpdateAtResourceGroup`

```go
ctx := context.TODO()
id := policyinsights.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationValue")

payload := policyinsights.Remediation{
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


### Example Usage: `PolicyInsightsClient.RemediationsCreateOrUpdateAtSubscription`

```go
ctx := context.TODO()
id := policyinsights.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationValue")

payload := policyinsights.Remediation{
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


### Example Usage: `PolicyInsightsClient.RemediationsDeleteAtManagementGroup`

```go
ctx := context.TODO()
id := policyinsights.NewProviders2RemediationID("managementGroupIdValue", "remediationValue")

read, err := client.RemediationsDeleteAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsDeleteAtResource`

```go
ctx := context.TODO()
id := policyinsights.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationValue")

read, err := client.RemediationsDeleteAtResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsDeleteAtResourceGroup`

```go
ctx := context.TODO()
id := policyinsights.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationValue")

read, err := client.RemediationsDeleteAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsDeleteAtSubscription`

```go
ctx := context.TODO()
id := policyinsights.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationValue")

read, err := client.RemediationsDeleteAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsGetAtManagementGroup`

```go
ctx := context.TODO()
id := policyinsights.NewProviders2RemediationID("managementGroupIdValue", "remediationValue")

read, err := client.RemediationsGetAtManagementGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsGetAtResource`

```go
ctx := context.TODO()
id := policyinsights.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationValue")

read, err := client.RemediationsGetAtResource(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsGetAtResourceGroup`

```go
ctx := context.TODO()
id := policyinsights.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationValue")

read, err := client.RemediationsGetAtResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsGetAtSubscription`

```go
ctx := context.TODO()
id := policyinsights.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationValue")

read, err := client.RemediationsGetAtSubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PolicyInsightsClient.RemediationsListDeploymentsAtManagementGroup`

```go
ctx := context.TODO()
id := policyinsights.NewProviders2RemediationID("managementGroupIdValue", "remediationValue")

// alternatively `client.RemediationsListDeploymentsAtManagementGroup(ctx, id, policyinsights.DefaultRemediationsListDeploymentsAtManagementGroupOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListDeploymentsAtManagementGroupComplete(ctx, id, policyinsights.DefaultRemediationsListDeploymentsAtManagementGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyInsightsClient.RemediationsListDeploymentsAtResource`

```go
ctx := context.TODO()
id := policyinsights.NewScopedRemediationID("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/some-resource-group", "remediationValue")

// alternatively `client.RemediationsListDeploymentsAtResource(ctx, id, policyinsights.DefaultRemediationsListDeploymentsAtResourceOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListDeploymentsAtResourceComplete(ctx, id, policyinsights.DefaultRemediationsListDeploymentsAtResourceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyInsightsClient.RemediationsListDeploymentsAtResourceGroup`

```go
ctx := context.TODO()
id := policyinsights.NewProviderRemediationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "remediationValue")

// alternatively `client.RemediationsListDeploymentsAtResourceGroup(ctx, id, policyinsights.DefaultRemediationsListDeploymentsAtResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListDeploymentsAtResourceGroupComplete(ctx, id, policyinsights.DefaultRemediationsListDeploymentsAtResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyInsightsClient.RemediationsListDeploymentsAtSubscription`

```go
ctx := context.TODO()
id := policyinsights.NewRemediationID("12345678-1234-9876-4563-123456789012", "remediationValue")

// alternatively `client.RemediationsListDeploymentsAtSubscription(ctx, id, policyinsights.DefaultRemediationsListDeploymentsAtSubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListDeploymentsAtSubscriptionComplete(ctx, id, policyinsights.DefaultRemediationsListDeploymentsAtSubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyInsightsClient.RemediationsListForManagementGroup`

```go
ctx := context.TODO()
id := policyinsights.NewManagementGroupID("managementGroupIdValue")

// alternatively `client.RemediationsListForManagementGroup(ctx, id, policyinsights.DefaultRemediationsListForManagementGroupOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListForManagementGroupComplete(ctx, id, policyinsights.DefaultRemediationsListForManagementGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyInsightsClient.RemediationsListForResource`

```go
ctx := context.TODO()
id := policyinsights.NewScopeID()

// alternatively `client.RemediationsListForResource(ctx, id, policyinsights.DefaultRemediationsListForResourceOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListForResourceComplete(ctx, id, policyinsights.DefaultRemediationsListForResourceOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyInsightsClient.RemediationsListForResourceGroup`

```go
ctx := context.TODO()
id := policyinsights.NewResourceGroupID()

// alternatively `client.RemediationsListForResourceGroup(ctx, id, policyinsights.DefaultRemediationsListForResourceGroupOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListForResourceGroupComplete(ctx, id, policyinsights.DefaultRemediationsListForResourceGroupOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PolicyInsightsClient.RemediationsListForSubscription`

```go
ctx := context.TODO()
id := policyinsights.NewSubscriptionID()

// alternatively `client.RemediationsListForSubscription(ctx, id, policyinsights.DefaultRemediationsListForSubscriptionOperationOptions())` can be used to do batched pagination
items, err := client.RemediationsListForSubscriptionComplete(ctx, id, policyinsights.DefaultRemediationsListForSubscriptionOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
