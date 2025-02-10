
## `github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroups` Documentation

The `prometheusrulegroups` SDK allows for interaction with Azure Resource Manager `alertsmanagement` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroups"
```


### Client Initialization

```go
client := prometheusrulegroups.NewPrometheusRuleGroupsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrometheusRuleGroupsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := prometheusrulegroups.NewPrometheusRuleGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "prometheusRuleGroupName")

payload := prometheusrulegroups.PrometheusRuleGroupResource{
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


### Example Usage: `PrometheusRuleGroupsClient.Delete`

```go
ctx := context.TODO()
id := prometheusrulegroups.NewPrometheusRuleGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "prometheusRuleGroupName")

read, err := client.Delete(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrometheusRuleGroupsClient.Get`

```go
ctx := context.TODO()
id := prometheusrulegroups.NewPrometheusRuleGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "prometheusRuleGroupName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrometheusRuleGroupsClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

read, err := client.ListByResourceGroup(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrometheusRuleGroupsClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

read, err := client.ListBySubscription(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrometheusRuleGroupsClient.Update`

```go
ctx := context.TODO()
id := prometheusrulegroups.NewPrometheusRuleGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group", "prometheusRuleGroupName")

payload := prometheusrulegroups.PrometheusRuleGroupResourcePatchParameters{
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
