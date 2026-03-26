
## `github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroupresources` Documentation

The `prometheusrulegroupresources` SDK allows for interaction with Azure Resource Manager `alertsmanagement` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroupresources"
```


### Client Initialization

```go
client := prometheusrulegroupresources.NewPrometheusRuleGroupResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrometheusRuleGroupResourcesClient.PrometheusRuleGroupsCreateOrUpdate`

```go
ctx := context.TODO()
id := prometheusrulegroupresources.NewPrometheusRuleGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "prometheusRuleGroupName")

payload := prometheusrulegroupresources.PrometheusRuleGroupResource{
	// ...
}


read, err := client.PrometheusRuleGroupsCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrometheusRuleGroupResourcesClient.PrometheusRuleGroupsDelete`

```go
ctx := context.TODO()
id := prometheusrulegroupresources.NewPrometheusRuleGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "prometheusRuleGroupName")

read, err := client.PrometheusRuleGroupsDelete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrometheusRuleGroupResourcesClient.PrometheusRuleGroupsGet`

```go
ctx := context.TODO()
id := prometheusrulegroupresources.NewPrometheusRuleGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "prometheusRuleGroupName")

read, err := client.PrometheusRuleGroupsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrometheusRuleGroupResourcesClient.PrometheusRuleGroupsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.PrometheusRuleGroupsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.PrometheusRuleGroupsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrometheusRuleGroupResourcesClient.PrometheusRuleGroupsListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.PrometheusRuleGroupsListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.PrometheusRuleGroupsListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrometheusRuleGroupResourcesClient.PrometheusRuleGroupsUpdate`

```go
ctx := context.TODO()
id := prometheusrulegroupresources.NewPrometheusRuleGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "prometheusRuleGroupName")

payload := prometheusrulegroupresources.PrometheusRuleGroupResourcePatchParameters{
	// ...
}


read, err := client.PrometheusRuleGroupsUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
