
## `github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources` Documentation

The `elasticmonitorresources` SDK allows for interaction with Azure Resource Manager `elastic` (API Version `2025-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
```


### Client Initialization

```go
client := elasticmonitorresources.NewElasticMonitorResourcesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ElasticMonitorResourcesClient.AllTrafficFilterslist`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.AllTrafficFilterslist(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticMonitorResourcesClient.AssociateTrafficFilterAssociate`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

if err := client.AssociateTrafficFilterAssociateThenPoll(ctx, id, elasticmonitorresources.DefaultAssociateTrafficFilterAssociateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticMonitorResourcesClient.BillingInfoGet`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.BillingInfoGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticMonitorResourcesClient.ConnectedPartnerResourcesList`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.ConnectedPartnerResourcesList(ctx, id)` can be used to do batched pagination
items, err := client.ConnectedPartnerResourcesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ElasticMonitorResourcesClient.CreateAndAssociateIPFilterCreate`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

if err := client.CreateAndAssociateIPFilterCreateThenPoll(ctx, id, elasticmonitorresources.DefaultCreateAndAssociateIPFilterCreateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticMonitorResourcesClient.CreateAndAssociatePLFilterCreate`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

if err := client.CreateAndAssociatePLFilterCreateThenPoll(ctx, id, elasticmonitorresources.DefaultCreateAndAssociatePLFilterCreateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticMonitorResourcesClient.DeploymentInfoList`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.DeploymentInfoList(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticMonitorResourcesClient.DetachAndDeleteTrafficFilterDelete`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.DetachAndDeleteTrafficFilterDelete(ctx, id, elasticmonitorresources.DefaultDetachAndDeleteTrafficFilterDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticMonitorResourcesClient.DetachTrafficFilterUpdate`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

if err := client.DetachTrafficFilterUpdateThenPoll(ctx, id, elasticmonitorresources.DefaultDetachTrafficFilterUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticMonitorResourcesClient.ExternalUserCreateOrUpdate`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := elasticmonitorresources.ExternalUserInfo{
	// ...
}


read, err := client.ExternalUserCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticMonitorResourcesClient.ListAssociatedTrafficFilterslist`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.ListAssociatedTrafficFilterslist(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticMonitorResourcesClient.MonitorUpgrade`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := elasticmonitorresources.ElasticMonitorUpgrade{
	// ...
}


if err := client.MonitorUpgradeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticMonitorResourcesClient.MonitoredResourcesList`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.MonitoredResourcesList(ctx, id)` can be used to do batched pagination
items, err := client.MonitoredResourcesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ElasticMonitorResourcesClient.MonitorsCreate`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := elasticmonitorresources.ElasticMonitorResource{
	// ...
}


if err := client.MonitorsCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticMonitorResourcesClient.MonitorsDelete`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

if err := client.MonitorsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticMonitorResourcesClient.MonitorsGet`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.MonitorsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticMonitorResourcesClient.MonitorsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.MonitorsList(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ElasticMonitorResourcesClient.MonitorsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.MonitorsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.MonitorsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ElasticMonitorResourcesClient.MonitorsUpdate`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := elasticmonitorresources.ElasticMonitorResourceUpdateParameters{
	// ...
}


if err := client.MonitorsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticMonitorResourcesClient.OrganizationsResubscribe`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := elasticmonitorresources.ResubscribeProperties{
	// ...
}


if err := client.OrganizationsResubscribeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ElasticMonitorResourcesClient.TrafficFiltersDelete`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.TrafficFiltersDelete(ctx, id, elasticmonitorresources.DefaultTrafficFiltersDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticMonitorResourcesClient.UpgradableVersionsDetails`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.UpgradableVersionsDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticMonitorResourcesClient.VMCollectionUpdate`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

payload := elasticmonitorresources.VMCollectionUpdate{
	// ...
}


read, err := client.VMCollectionUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `ElasticMonitorResourcesClient.VMHostList`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

// alternatively `client.VMHostList(ctx, id)` can be used to do batched pagination
items, err := client.VMHostListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `ElasticMonitorResourcesClient.VMIngestionDetails`

```go
ctx := context.TODO()
id := elasticmonitorresources.NewMonitorID("12345678-1234-9876-4563-123456789012", "example-resource-group", "monitorName")

read, err := client.VMIngestionDetails(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
