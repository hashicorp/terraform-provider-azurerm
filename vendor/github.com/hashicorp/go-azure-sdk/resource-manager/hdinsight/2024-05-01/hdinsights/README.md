
## `github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2024-05-01/hdinsights` Documentation

The `hdinsights` SDK allows for interaction with the Azure Resource Manager Service `hdinsight` (API Version `2024-05-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2024-05-01/hdinsights"
```


### Client Initialization

```go
client := hdinsights.NewHdinsightsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `HdinsightsClient.AvailableClusterPoolVersionsListByLocation`

```go
ctx := context.TODO()
id := hdinsights.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

// alternatively `client.AvailableClusterPoolVersionsListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.AvailableClusterPoolVersionsListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.AvailableClusterVersionsListByLocation`

```go
ctx := context.TODO()
id := hdinsights.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

// alternatively `client.AvailableClusterVersionsListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.AvailableClusterVersionsListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClusterAvailableUpgradesList`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

// alternatively `client.ClusterAvailableUpgradesList(ctx, id)` can be used to do batched pagination
items, err := client.ClusterAvailableUpgradesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClusterJobsList`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

// alternatively `client.ClusterJobsList(ctx, id, hdinsights.DefaultClusterJobsListOperationOptions())` can be used to do batched pagination
items, err := client.ClusterJobsListComplete(ctx, id, hdinsights.DefaultClusterJobsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClusterJobsRunJob`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

payload := hdinsights.ClusterJob{
	// ...
}


if err := client.ClusterJobsRunJobThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClusterLibrariesList`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

// alternatively `client.ClusterLibrariesList(ctx, id, hdinsights.DefaultClusterLibrariesListOperationOptions())` can be used to do batched pagination
items, err := client.ClusterLibrariesListComplete(ctx, id, hdinsights.DefaultClusterLibrariesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClusterLibrariesManageLibraries`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

payload := hdinsights.ClusterLibraryManagementOperation{
	// ...
}


if err := client.ClusterLibrariesManageLibrariesThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClusterPoolAvailableUpgradesList`

```go
ctx := context.TODO()
id := hdinsights.NewClusterPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue")

// alternatively `client.ClusterPoolAvailableUpgradesList(ctx, id)` can be used to do batched pagination
items, err := client.ClusterPoolAvailableUpgradesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClusterPoolUpgradeHistoriesList`

```go
ctx := context.TODO()
id := hdinsights.NewClusterPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue")

// alternatively `client.ClusterPoolUpgradeHistoriesList(ctx, id)` can be used to do batched pagination
items, err := client.ClusterPoolUpgradeHistoriesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClusterPoolsCreateOrUpdate`

```go
ctx := context.TODO()
id := hdinsights.NewClusterPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue")

payload := hdinsights.ClusterPool{
	// ...
}


if err := client.ClusterPoolsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClusterPoolsDelete`

```go
ctx := context.TODO()
id := hdinsights.NewClusterPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue")

if err := client.ClusterPoolsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClusterPoolsGet`

```go
ctx := context.TODO()
id := hdinsights.NewClusterPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue")

read, err := client.ClusterPoolsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HdinsightsClient.ClusterPoolsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ClusterPoolsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ClusterPoolsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClusterPoolsListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ClusterPoolsListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ClusterPoolsListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClusterPoolsUpdateTags`

```go
ctx := context.TODO()
id := hdinsights.NewClusterPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue")

payload := hdinsights.TagsObject{
	// ...
}


if err := client.ClusterPoolsUpdateTagsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClusterPoolsUpgrade`

```go
ctx := context.TODO()
id := hdinsights.NewClusterPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue")

payload := hdinsights.ClusterPoolUpgrade{
	// ...
}


if err := client.ClusterPoolsUpgradeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClusterUpgradeHistoriesList`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

// alternatively `client.ClusterUpgradeHistoriesList(ctx, id)` can be used to do batched pagination
items, err := client.ClusterUpgradeHistoriesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClustersCreate`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

payload := hdinsights.Cluster{
	// ...
}


if err := client.ClustersCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClustersDelete`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

if err := client.ClustersDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClustersGet`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

read, err := client.ClustersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HdinsightsClient.ClustersGetInstanceView`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

read, err := client.ClustersGetInstanceView(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HdinsightsClient.ClustersListByClusterPoolName`

```go
ctx := context.TODO()
id := hdinsights.NewClusterPoolID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue")

// alternatively `client.ClustersListByClusterPoolName(ctx, id)` can be used to do batched pagination
items, err := client.ClustersListByClusterPoolNameComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClustersListInstanceViews`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

// alternatively `client.ClustersListInstanceViews(ctx, id)` can be used to do batched pagination
items, err := client.ClustersListInstanceViewsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClustersListServiceConfigs`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

// alternatively `client.ClustersListServiceConfigs(ctx, id)` can be used to do batched pagination
items, err := client.ClustersListServiceConfigsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HdinsightsClient.ClustersResize`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

payload := hdinsights.ClusterResizeData{
	// ...
}


if err := client.ClustersResizeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClustersUpdate`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

payload := hdinsights.ClusterPatch{
	// ...
}


if err := client.ClustersUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClustersUpgrade`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

payload := hdinsights.ClusterUpgrade{
	// ...
}


if err := client.ClustersUpgradeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.ClustersUpgradeManualRollback`

```go
ctx := context.TODO()
id := hdinsights.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterPoolValue", "clusterValue")

payload := hdinsights.ClusterUpgradeRollback{
	// ...
}


if err := client.ClustersUpgradeManualRollbackThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HdinsightsClient.LocationsCheckNameAvailability`

```go
ctx := context.TODO()
id := hdinsights.NewLocationID("12345678-1234-9876-4563-123456789012", "locationValue")

payload := hdinsights.NameAvailabilityParameters{
	// ...
}


read, err := client.LocationsCheckNameAvailability(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
