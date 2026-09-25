
## `github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2026-03-15/openapis` Documentation

The `openapis` SDK allows for interaction with Azure Resource Manager `cosmosdb` (API Version `2026-03-15`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2026-03-15/openapis"
```


### Client Initialization

```go
client := openapis.NewOpenapisClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `OpenapisClient.CassandraClustersCreateUpdate`

```go
ctx := context.TODO()
id := openapis.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

payload := openapis.ClusterResource{
	// ...
}


if err := client.CassandraClustersCreateUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraClustersDeallocate`

```go
ctx := context.TODO()
id := openapis.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

if err := client.CassandraClustersDeallocateThenPoll(ctx, id, openapis.DefaultCassandraClustersDeallocateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraClustersDelete`

```go
ctx := context.TODO()
id := openapis.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

if err := client.CassandraClustersDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraClustersGet`

```go
ctx := context.TODO()
id := openapis.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

read, err := client.CassandraClustersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CassandraClustersInvokeCommand`

```go
ctx := context.TODO()
id := openapis.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

payload := openapis.CommandPostBody{
	// ...
}


if err := client.CassandraClustersInvokeCommandThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraClustersListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.CassandraClustersListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.CassandraClustersListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CassandraClustersListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.CassandraClustersListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.CassandraClustersListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CassandraClustersStart`

```go
ctx := context.TODO()
id := openapis.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

if err := client.CassandraClustersStartThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraClustersStatus`

```go
ctx := context.TODO()
id := openapis.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

read, err := client.CassandraClustersStatus(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CassandraClustersUpdate`

```go
ctx := context.TODO()
id := openapis.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

payload := openapis.ClusterResource{
	// ...
}


if err := client.CassandraClustersUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraDataCentersCreateUpdate`

```go
ctx := context.TODO()
id := openapis.NewDataCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName", "dataCenterName")

payload := openapis.DataCenterResource{
	// ...
}


if err := client.CassandraDataCentersCreateUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraDataCentersDelete`

```go
ctx := context.TODO()
id := openapis.NewDataCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName", "dataCenterName")

if err := client.CassandraDataCentersDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraDataCentersGet`

```go
ctx := context.TODO()
id := openapis.NewDataCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName", "dataCenterName")

read, err := client.CassandraDataCentersGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CassandraDataCentersList`

```go
ctx := context.TODO()
id := openapis.NewCassandraClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName")

// alternatively `client.CassandraDataCentersList(ctx, id)` can be used to do batched pagination
items, err := client.CassandraDataCentersListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CassandraDataCentersUpdate`

```go
ctx := context.TODO()
id := openapis.NewDataCenterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "cassandraClusterName", "dataCenterName")

payload := openapis.DataCenterResource{
	// ...
}


if err := client.CassandraDataCentersUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesCreateUpdateCassandraKeyspace`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName")

payload := openapis.CassandraKeyspaceCreateUpdateParameters{
	// ...
}


if err := client.CassandraResourcesCreateUpdateCassandraKeyspaceThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesCreateUpdateCassandraRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewCassandraRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

payload := openapis.CassandraRoleAssignmentResource{
	// ...
}


if err := client.CassandraResourcesCreateUpdateCassandraRoleAssignmentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesCreateUpdateCassandraRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewCassandraRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

payload := openapis.CassandraRoleDefinitionResource{
	// ...
}


if err := client.CassandraResourcesCreateUpdateCassandraRoleDefinitionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesCreateUpdateCassandraTable`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName", "tableName")

payload := openapis.CassandraTableCreateUpdateParameters{
	// ...
}


if err := client.CassandraResourcesCreateUpdateCassandraTableThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesDeleteCassandraKeyspace`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName")

if err := client.CassandraResourcesDeleteCassandraKeyspaceThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesDeleteCassandraRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewCassandraRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

if err := client.CassandraResourcesDeleteCassandraRoleAssignmentThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesDeleteCassandraRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewCassandraRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

if err := client.CassandraResourcesDeleteCassandraRoleDefinitionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesDeleteCassandraTable`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName", "tableName")

if err := client.CassandraResourcesDeleteCassandraTableThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesGetCassandraKeyspace`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName")

read, err := client.CassandraResourcesGetCassandraKeyspace(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CassandraResourcesGetCassandraKeyspaceThroughput`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName")

read, err := client.CassandraResourcesGetCassandraKeyspaceThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CassandraResourcesGetCassandraRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewCassandraRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

read, err := client.CassandraResourcesGetCassandraRoleAssignment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CassandraResourcesGetCassandraRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewCassandraRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

read, err := client.CassandraResourcesGetCassandraRoleDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CassandraResourcesGetCassandraTable`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName", "tableName")

read, err := client.CassandraResourcesGetCassandraTable(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CassandraResourcesGetCassandraTableThroughput`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName", "tableName")

read, err := client.CassandraResourcesGetCassandraTableThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.CassandraResourcesListCassandraKeyspaces`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.CassandraResourcesListCassandraKeyspaces(ctx, id)` can be used to do batched pagination
items, err := client.CassandraResourcesListCassandraKeyspacesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CassandraResourcesListCassandraRoleAssignments`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.CassandraResourcesListCassandraRoleAssignments(ctx, id)` can be used to do batched pagination
items, err := client.CassandraResourcesListCassandraRoleAssignmentsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CassandraResourcesListCassandraRoleDefinitions`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.CassandraResourcesListCassandraRoleDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.CassandraResourcesListCassandraRoleDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CassandraResourcesListCassandraTables`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName")

// alternatively `client.CassandraResourcesListCassandraTables(ctx, id)` can be used to do batched pagination
items, err := client.CassandraResourcesListCassandraTablesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CassandraResourcesMigrateCassandraKeyspaceToAutoscale`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName")

if err := client.CassandraResourcesMigrateCassandraKeyspaceToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesMigrateCassandraKeyspaceToManualThroughput`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName")

if err := client.CassandraResourcesMigrateCassandraKeyspaceToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesMigrateCassandraTableToAutoscale`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName", "tableName")

if err := client.CassandraResourcesMigrateCassandraTableToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesMigrateCassandraTableToManualThroughput`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName", "tableName")

if err := client.CassandraResourcesMigrateCassandraTableToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesUpdateCassandraKeyspaceThroughput`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName")

payload := openapis.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.CassandraResourcesUpdateCassandraKeyspaceThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CassandraResourcesUpdateCassandraTableThroughput`

```go
ctx := context.TODO()
id := openapis.NewCassandraKeyspaceTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "cassandraKeyspaceName", "tableName")

payload := openapis.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.CassandraResourcesUpdateCassandraTableThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.CollectionListMetricDefinitions`

```go
ctx := context.TODO()
id := openapis.NewCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "databaseName", "collectionName")

// alternatively `client.CollectionListMetricDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.CollectionListMetricDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CollectionListMetrics`

```go
ctx := context.TODO()
id := openapis.NewCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "databaseName", "collectionName")

// alternatively `client.CollectionListMetrics(ctx, id, openapis.DefaultCollectionListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.CollectionListMetricsComplete(ctx, id, openapis.DefaultCollectionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CollectionListUsages`

```go
ctx := context.TODO()
id := openapis.NewCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "databaseName", "collectionName")

// alternatively `client.CollectionListUsages(ctx, id, openapis.DefaultCollectionListUsagesOperationOptions())` can be used to do batched pagination
items, err := client.CollectionListUsagesComplete(ctx, id, openapis.DefaultCollectionListUsagesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CollectionPartitionListMetrics`

```go
ctx := context.TODO()
id := openapis.NewCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "databaseName", "collectionName")

// alternatively `client.CollectionPartitionListMetrics(ctx, id, openapis.DefaultCollectionPartitionListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.CollectionPartitionListMetricsComplete(ctx, id, openapis.DefaultCollectionPartitionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CollectionPartitionListUsages`

```go
ctx := context.TODO()
id := openapis.NewCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "databaseName", "collectionName")

// alternatively `client.CollectionPartitionListUsages(ctx, id, openapis.DefaultCollectionPartitionListUsagesOperationOptions())` can be used to do batched pagination
items, err := client.CollectionPartitionListUsagesComplete(ctx, id, openapis.DefaultCollectionPartitionListUsagesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CollectionPartitionRegionListMetrics`

```go
ctx := context.TODO()
id := openapis.NewDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "regionName", "databaseName", "collectionName")

// alternatively `client.CollectionPartitionRegionListMetrics(ctx, id, openapis.DefaultCollectionPartitionRegionListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.CollectionPartitionRegionListMetricsComplete(ctx, id, openapis.DefaultCollectionPartitionRegionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.CollectionRegionListMetrics`

```go
ctx := context.TODO()
id := openapis.NewDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "regionName", "databaseName", "collectionName")

// alternatively `client.CollectionRegionListMetrics(ctx, id, openapis.DefaultCollectionRegionListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.CollectionRegionListMetricsComplete(ctx, id, openapis.DefaultCollectionRegionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.DatabaseAccountRegionListMetrics`

```go
ctx := context.TODO()
id := openapis.NewRegionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "regionName")

// alternatively `client.DatabaseAccountRegionListMetrics(ctx, id, openapis.DefaultDatabaseAccountRegionListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.DatabaseAccountRegionListMetricsComplete(ctx, id, openapis.DefaultDatabaseAccountRegionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsCheckNameExists`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountNameID("accountName")

read, err := client.DatabaseAccountsCheckNameExists(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsCreateOrUpdate`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

payload := openapis.DatabaseAccountCreateUpdateParameters{
	// ...
}


if err := client.DatabaseAccountsCreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsDelete`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

if err := client.DatabaseAccountsDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsFailoverPriorityChange`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

payload := openapis.FailoverPolicies{
	// ...
}


if err := client.DatabaseAccountsFailoverPriorityChangeThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsGet`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

read, err := client.DatabaseAccountsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsGetReadOnlyKeys`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

read, err := client.DatabaseAccountsGetReadOnlyKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.DatabaseAccountsList(ctx, id)` can be used to do batched pagination
items, err := client.DatabaseAccountsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.DatabaseAccountsListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.DatabaseAccountsListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsListConnectionStrings`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

read, err := client.DatabaseAccountsListConnectionStrings(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsListKeys`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

read, err := client.DatabaseAccountsListKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsListMetricDefinitions`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.DatabaseAccountsListMetricDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.DatabaseAccountsListMetricDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsListMetrics`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.DatabaseAccountsListMetrics(ctx, id, openapis.DefaultDatabaseAccountsListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.DatabaseAccountsListMetricsComplete(ctx, id, openapis.DefaultDatabaseAccountsListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsListReadOnlyKeys`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

read, err := client.DatabaseAccountsListReadOnlyKeys(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsListUsages`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.DatabaseAccountsListUsages(ctx, id, openapis.DefaultDatabaseAccountsListUsagesOperationOptions())` can be used to do batched pagination
items, err := client.DatabaseAccountsListUsagesComplete(ctx, id, openapis.DefaultDatabaseAccountsListUsagesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsOfflineRegion`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

payload := openapis.RegionForOnlineOffline{
	// ...
}


if err := client.DatabaseAccountsOfflineRegionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsOnlineRegion`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

payload := openapis.RegionForOnlineOffline{
	// ...
}


if err := client.DatabaseAccountsOnlineRegionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsRegenerateKey`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

payload := openapis.DatabaseAccountRegenerateKeyParameters{
	// ...
}


if err := client.DatabaseAccountsRegenerateKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.DatabaseAccountsUpdate`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

payload := openapis.DatabaseAccountUpdateParameters{
	// ...
}


if err := client.DatabaseAccountsUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.DatabaseListMetricDefinitions`

```go
ctx := context.TODO()
id := openapis.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "databaseName")

// alternatively `client.DatabaseListMetricDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.DatabaseListMetricDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.DatabaseListMetrics`

```go
ctx := context.TODO()
id := openapis.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "databaseName")

// alternatively `client.DatabaseListMetrics(ctx, id, openapis.DefaultDatabaseListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.DatabaseListMetricsComplete(ctx, id, openapis.DefaultDatabaseListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.DatabaseListUsages`

```go
ctx := context.TODO()
id := openapis.NewDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "databaseName")

// alternatively `client.DatabaseListUsages(ctx, id, openapis.DefaultDatabaseListUsagesOperationOptions())` can be used to do batched pagination
items, err := client.DatabaseListUsagesComplete(ctx, id, openapis.DefaultDatabaseListUsagesOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.FleetCreate`

```go
ctx := context.TODO()
id := openapis.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

payload := openapis.FleetResource{
	// ...
}


read, err := client.FleetCreate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.FleetDelete`

```go
ctx := context.TODO()
id := openapis.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

if err := client.FleetDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.FleetGet`

```go
ctx := context.TODO()
id := openapis.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

read, err := client.FleetGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.FleetList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.FleetList(ctx, id)` can be used to do batched pagination
items, err := client.FleetListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.FleetUpdate`

```go
ctx := context.TODO()
id := openapis.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

payload := openapis.FleetResourceUpdate{
	// ...
}


read, err := client.FleetUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.FleetspaceAccountCreate`

```go
ctx := context.TODO()
id := openapis.NewFleetspaceAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName", "fleetspaceAccountName")

payload := openapis.FleetspaceAccountResource{
	// ...
}


if err := client.FleetspaceAccountCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.FleetspaceAccountDelete`

```go
ctx := context.TODO()
id := openapis.NewFleetspaceAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName", "fleetspaceAccountName")

if err := client.FleetspaceAccountDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.FleetspaceAccountGet`

```go
ctx := context.TODO()
id := openapis.NewFleetspaceAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName", "fleetspaceAccountName")

read, err := client.FleetspaceAccountGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.FleetspaceAccountList`

```go
ctx := context.TODO()
id := openapis.NewFleetspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName")

// alternatively `client.FleetspaceAccountList(ctx, id)` can be used to do batched pagination
items, err := client.FleetspaceAccountListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.FleetspaceCreate`

```go
ctx := context.TODO()
id := openapis.NewFleetspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName")

payload := openapis.FleetspaceResource{
	// ...
}


if err := client.FleetspaceCreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.FleetspaceDelete`

```go
ctx := context.TODO()
id := openapis.NewFleetspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName")

if err := client.FleetspaceDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.FleetspaceGet`

```go
ctx := context.TODO()
id := openapis.NewFleetspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName")

read, err := client.FleetspaceGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.FleetspaceList`

```go
ctx := context.TODO()
id := openapis.NewFleetID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName")

// alternatively `client.FleetspaceList(ctx, id)` can be used to do batched pagination
items, err := client.FleetspaceListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.FleetspaceUpdate`

```go
ctx := context.TODO()
id := openapis.NewFleetspaceID("12345678-1234-9876-4563-123456789012", "example-resource-group", "fleetName", "fleetspaceName")

payload := openapis.FleetspaceUpdate{
	// ...
}


if err := client.FleetspaceUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesCreateUpdateGremlinDatabase`

```go
ctx := context.TODO()
id := openapis.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName")

payload := openapis.GremlinDatabaseCreateUpdateParameters{
	// ...
}


if err := client.GremlinResourcesCreateUpdateGremlinDatabaseThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesCreateUpdateGremlinGraph`

```go
ctx := context.TODO()
id := openapis.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName", "graphName")

payload := openapis.GremlinGraphCreateUpdateParameters{
	// ...
}


if err := client.GremlinResourcesCreateUpdateGremlinGraphThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesCreateUpdateGremlinRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewGremlinRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

payload := openapis.GremlinRoleAssignmentResource{
	// ...
}


if err := client.GremlinResourcesCreateUpdateGremlinRoleAssignmentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesCreateUpdateGremlinRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewGremlinRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

payload := openapis.GremlinRoleDefinitionResource{
	// ...
}


if err := client.GremlinResourcesCreateUpdateGremlinRoleDefinitionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesDeleteGremlinDatabase`

```go
ctx := context.TODO()
id := openapis.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName")

if err := client.GremlinResourcesDeleteGremlinDatabaseThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesDeleteGremlinGraph`

```go
ctx := context.TODO()
id := openapis.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName", "graphName")

if err := client.GremlinResourcesDeleteGremlinGraphThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesDeleteGremlinRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewGremlinRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

if err := client.GremlinResourcesDeleteGremlinRoleAssignmentThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesDeleteGremlinRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewGremlinRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

if err := client.GremlinResourcesDeleteGremlinRoleDefinitionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesGetGremlinDatabase`

```go
ctx := context.TODO()
id := openapis.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName")

read, err := client.GremlinResourcesGetGremlinDatabase(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.GremlinResourcesGetGremlinDatabaseThroughput`

```go
ctx := context.TODO()
id := openapis.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName")

read, err := client.GremlinResourcesGetGremlinDatabaseThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.GremlinResourcesGetGremlinGraph`

```go
ctx := context.TODO()
id := openapis.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName", "graphName")

read, err := client.GremlinResourcesGetGremlinGraph(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.GremlinResourcesGetGremlinGraphThroughput`

```go
ctx := context.TODO()
id := openapis.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName", "graphName")

read, err := client.GremlinResourcesGetGremlinGraphThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.GremlinResourcesGetGremlinRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewGremlinRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

read, err := client.GremlinResourcesGetGremlinRoleAssignment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.GremlinResourcesGetGremlinRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewGremlinRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

read, err := client.GremlinResourcesGetGremlinRoleDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.GremlinResourcesListGremlinDatabases`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.GremlinResourcesListGremlinDatabases(ctx, id)` can be used to do batched pagination
items, err := client.GremlinResourcesListGremlinDatabasesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.GremlinResourcesListGremlinGraphs`

```go
ctx := context.TODO()
id := openapis.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName")

// alternatively `client.GremlinResourcesListGremlinGraphs(ctx, id)` can be used to do batched pagination
items, err := client.GremlinResourcesListGremlinGraphsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.GremlinResourcesListGremlinRoleAssignments`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.GremlinResourcesListGremlinRoleAssignments(ctx, id)` can be used to do batched pagination
items, err := client.GremlinResourcesListGremlinRoleAssignmentsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.GremlinResourcesListGremlinRoleDefinitions`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.GremlinResourcesListGremlinRoleDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.GremlinResourcesListGremlinRoleDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.GremlinResourcesMigrateGremlinDatabaseToAutoscale`

```go
ctx := context.TODO()
id := openapis.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName")

if err := client.GremlinResourcesMigrateGremlinDatabaseToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesMigrateGremlinDatabaseToManualThroughput`

```go
ctx := context.TODO()
id := openapis.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName")

if err := client.GremlinResourcesMigrateGremlinDatabaseToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesMigrateGremlinGraphToAutoscale`

```go
ctx := context.TODO()
id := openapis.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName", "graphName")

if err := client.GremlinResourcesMigrateGremlinGraphToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesMigrateGremlinGraphToManualThroughput`

```go
ctx := context.TODO()
id := openapis.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName", "graphName")

if err := client.GremlinResourcesMigrateGremlinGraphToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesRetrieveContinuousBackupInformation`

```go
ctx := context.TODO()
id := openapis.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName", "graphName")

payload := openapis.ContinuousBackupRestoreLocation{
	// ...
}


if err := client.GremlinResourcesRetrieveContinuousBackupInformationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesUpdateGremlinDatabaseThroughput`

```go
ctx := context.TODO()
id := openapis.NewGremlinDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName")

payload := openapis.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.GremlinResourcesUpdateGremlinDatabaseThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.GremlinResourcesUpdateGremlinGraphThroughput`

```go
ctx := context.TODO()
id := openapis.NewGraphID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "gremlinDatabaseName", "graphName")

payload := openapis.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.GremlinResourcesUpdateGremlinGraphThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.LocationsGet`

```go
ctx := context.TODO()
id := openapis.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

read, err := client.LocationsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.LocationsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.LocationsList(ctx, id)` can be used to do batched pagination
items, err := client.LocationsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesCreateUpdateMongoDBCollection`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName", "collectionName")

payload := openapis.MongoDBCollectionCreateUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesCreateUpdateMongoDBCollectionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesCreateUpdateMongoDBDatabase`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName")

payload := openapis.MongoDBDatabaseCreateUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesCreateUpdateMongoDBDatabaseThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesCreateUpdateMongoRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewMongodbRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoRoleDefinitionId")

payload := openapis.MongoRoleDefinitionCreateUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesCreateUpdateMongoRoleDefinitionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesCreateUpdateMongoUserDefinition`

```go
ctx := context.TODO()
id := openapis.NewMongodbUserDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoUserDefinitionId")

payload := openapis.MongoUserDefinitionCreateUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesCreateUpdateMongoUserDefinitionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesDeleteMongoDBCollection`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName", "collectionName")

if err := client.MongoDBResourcesDeleteMongoDBCollectionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesDeleteMongoDBDatabase`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName")

if err := client.MongoDBResourcesDeleteMongoDBDatabaseThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesDeleteMongoRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewMongodbRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoRoleDefinitionId")

if err := client.MongoDBResourcesDeleteMongoRoleDefinitionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesDeleteMongoUserDefinition`

```go
ctx := context.TODO()
id := openapis.NewMongodbUserDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoUserDefinitionId")

if err := client.MongoDBResourcesDeleteMongoUserDefinitionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesGetMongoDBCollection`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName", "collectionName")

read, err := client.MongoDBResourcesGetMongoDBCollection(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesGetMongoDBCollectionThroughput`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName", "collectionName")

read, err := client.MongoDBResourcesGetMongoDBCollectionThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesGetMongoDBDatabase`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName")

read, err := client.MongoDBResourcesGetMongoDBDatabase(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesGetMongoDBDatabaseThroughput`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName")

read, err := client.MongoDBResourcesGetMongoDBDatabaseThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesGetMongoRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewMongodbRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoRoleDefinitionId")

read, err := client.MongoDBResourcesGetMongoRoleDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesGetMongoUserDefinition`

```go
ctx := context.TODO()
id := openapis.NewMongodbUserDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongoUserDefinitionId")

read, err := client.MongoDBResourcesGetMongoUserDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesListMongoDBCollections`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName")

// alternatively `client.MongoDBResourcesListMongoDBCollections(ctx, id)` can be used to do batched pagination
items, err := client.MongoDBResourcesListMongoDBCollectionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesListMongoDBDatabases`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.MongoDBResourcesListMongoDBDatabases(ctx, id)` can be used to do batched pagination
items, err := client.MongoDBResourcesListMongoDBDatabasesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesListMongoRoleDefinitions`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.MongoDBResourcesListMongoRoleDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.MongoDBResourcesListMongoRoleDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesListMongoUserDefinitions`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.MongoDBResourcesListMongoUserDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.MongoDBResourcesListMongoUserDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesMigrateMongoDBCollectionToAutoscale`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName", "collectionName")

if err := client.MongoDBResourcesMigrateMongoDBCollectionToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesMigrateMongoDBCollectionToManualThroughput`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName", "collectionName")

if err := client.MongoDBResourcesMigrateMongoDBCollectionToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesMigrateMongoDBDatabaseToAutoscale`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName")

if err := client.MongoDBResourcesMigrateMongoDBDatabaseToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesMigrateMongoDBDatabaseToManualThroughput`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName")

if err := client.MongoDBResourcesMigrateMongoDBDatabaseToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesRetrieveContinuousBackupInformation`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName", "collectionName")

payload := openapis.ContinuousBackupRestoreLocation{
	// ...
}


if err := client.MongoDBResourcesRetrieveContinuousBackupInformationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesUpdateMongoDBCollectionThroughput`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseCollectionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName", "collectionName")

payload := openapis.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesUpdateMongoDBCollectionThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoDBResourcesUpdateMongoDBDatabaseThroughput`

```go
ctx := context.TODO()
id := openapis.NewMongodbDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "mongodbDatabaseName")

payload := openapis.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.MongoDBResourcesUpdateMongoDBDatabaseThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoMIResourcesCreateUpdateMongoMIRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewMongoMIRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

payload := openapis.MongoMIRoleAssignmentResource{
	// ...
}


if err := client.MongoMIResourcesCreateUpdateMongoMIRoleAssignmentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoMIResourcesCreateUpdateMongoMIRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewMongoMIRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

payload := openapis.MongoMIRoleDefinitionResource{
	// ...
}


if err := client.MongoMIResourcesCreateUpdateMongoMIRoleDefinitionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoMIResourcesDeleteMongoMIRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewMongoMIRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

if err := client.MongoMIResourcesDeleteMongoMIRoleAssignmentThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoMIResourcesDeleteMongoMIRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewMongoMIRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

if err := client.MongoMIResourcesDeleteMongoMIRoleDefinitionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.MongoMIResourcesGetMongoMIRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewMongoMIRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

read, err := client.MongoMIResourcesGetMongoMIRoleAssignment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.MongoMIResourcesGetMongoMIRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewMongoMIRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

read, err := client.MongoMIResourcesGetMongoMIRoleDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.MongoMIResourcesListMongoMIRoleAssignments`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.MongoMIResourcesListMongoMIRoleAssignments(ctx, id)` can be used to do batched pagination
items, err := client.MongoMIResourcesListMongoMIRoleAssignmentsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.MongoMIResourcesListMongoMIRoleDefinitions`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.MongoMIResourcesListMongoMIRoleDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.MongoMIResourcesListMongoMIRoleDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.PartitionKeyRangeIdListMetrics`

```go
ctx := context.TODO()
id := openapis.NewPartitionKeyRangeIdID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "databaseName", "collectionName", "partitionKeyRangeId")

// alternatively `client.PartitionKeyRangeIdListMetrics(ctx, id, openapis.DefaultPartitionKeyRangeIdListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.PartitionKeyRangeIdListMetricsComplete(ctx, id, openapis.DefaultPartitionKeyRangeIdListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.PartitionKeyRangeIdRegionListMetrics`

```go
ctx := context.TODO()
id := openapis.NewCollectionPartitionKeyRangeIdID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "regionName", "databaseName", "collectionName", "partitionKeyRangeId")

// alternatively `client.PartitionKeyRangeIdRegionListMetrics(ctx, id, openapis.DefaultPartitionKeyRangeIdRegionListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.PartitionKeyRangeIdRegionListMetricsComplete(ctx, id, openapis.DefaultPartitionKeyRangeIdRegionListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.PercentileListMetrics`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.PercentileListMetrics(ctx, id, openapis.DefaultPercentileListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.PercentileListMetricsComplete(ctx, id, openapis.DefaultPercentileListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.PercentileSourceTargetListMetrics`

```go
ctx := context.TODO()
id := openapis.NewSourceRegionTargetRegionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sourceRegionName", "targetRegionName")

// alternatively `client.PercentileSourceTargetListMetrics(ctx, id, openapis.DefaultPercentileSourceTargetListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.PercentileSourceTargetListMetricsComplete(ctx, id, openapis.DefaultPercentileSourceTargetListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.PercentileTargetListMetrics`

```go
ctx := context.TODO()
id := openapis.NewTargetRegionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "targetRegionName")

// alternatively `client.PercentileTargetListMetrics(ctx, id, openapis.DefaultPercentileTargetListMetricsOperationOptions())` can be used to do batched pagination
items, err := client.PercentileTargetListMetricsComplete(ctx, id, openapis.DefaultPercentileTargetListMetricsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableDatabaseAccountsGetByLocation`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

read, err := client.RestorableDatabaseAccountsGetByLocation(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.RestorableDatabaseAccountsList`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.RestorableDatabaseAccountsList(ctx, id)` can be used to do batched pagination
items, err := client.RestorableDatabaseAccountsListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableDatabaseAccountsListByLocation`

```go
ctx := context.TODO()
id := openapis.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

// alternatively `client.RestorableDatabaseAccountsListByLocation(ctx, id)` can be used to do batched pagination
items, err := client.RestorableDatabaseAccountsListByLocationComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableGremlinDatabasesList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableGremlinDatabasesList(ctx, id)` can be used to do batched pagination
items, err := client.RestorableGremlinDatabasesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableGremlinGraphsList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableGremlinGraphsList(ctx, id, openapis.DefaultRestorableGremlinGraphsListOperationOptions())` can be used to do batched pagination
items, err := client.RestorableGremlinGraphsListComplete(ctx, id, openapis.DefaultRestorableGremlinGraphsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableGremlinResourcesList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableGremlinResourcesList(ctx, id, openapis.DefaultRestorableGremlinResourcesListOperationOptions())` can be used to do batched pagination
items, err := client.RestorableGremlinResourcesListComplete(ctx, id, openapis.DefaultRestorableGremlinResourcesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableMongodbCollectionsList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableMongodbCollectionsList(ctx, id, openapis.DefaultRestorableMongodbCollectionsListOperationOptions())` can be used to do batched pagination
items, err := client.RestorableMongodbCollectionsListComplete(ctx, id, openapis.DefaultRestorableMongodbCollectionsListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableMongodbDatabasesList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableMongodbDatabasesList(ctx, id)` can be used to do batched pagination
items, err := client.RestorableMongodbDatabasesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableMongodbResourcesList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableMongodbResourcesList(ctx, id, openapis.DefaultRestorableMongodbResourcesListOperationOptions())` can be used to do batched pagination
items, err := client.RestorableMongodbResourcesListComplete(ctx, id, openapis.DefaultRestorableMongodbResourcesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableSqlContainersList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableSqlContainersList(ctx, id, openapis.DefaultRestorableSqlContainersListOperationOptions())` can be used to do batched pagination
items, err := client.RestorableSqlContainersListComplete(ctx, id, openapis.DefaultRestorableSqlContainersListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableSqlDatabasesList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableSqlDatabasesList(ctx, id)` can be used to do batched pagination
items, err := client.RestorableSqlDatabasesListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableSqlResourcesList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableSqlResourcesList(ctx, id, openapis.DefaultRestorableSqlResourcesListOperationOptions())` can be used to do batched pagination
items, err := client.RestorableSqlResourcesListComplete(ctx, id, openapis.DefaultRestorableSqlResourcesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableTableResourcesList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableTableResourcesList(ctx, id, openapis.DefaultRestorableTableResourcesListOperationOptions())` can be used to do batched pagination
items, err := client.RestorableTableResourcesListComplete(ctx, id, openapis.DefaultRestorableTableResourcesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.RestorableTablesList`

```go
ctx := context.TODO()
id := openapis.NewRestorableDatabaseAccountID("12345678-1234-9876-4563-123456789012", "locationName", "instanceId")

// alternatively `client.RestorableTablesList(ctx, id, openapis.DefaultRestorableTablesListOperationOptions())` can be used to do batched pagination
items, err := client.RestorableTablesListComplete(ctx, id, openapis.DefaultRestorableTablesListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.ServiceList`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.ServiceList(ctx, id)` can be used to do batched pagination
items, err := client.ServiceListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.SqlResourcesCreateUpdateClientEncryptionKey`

```go
ctx := context.TODO()
id := openapis.NewClientEncryptionKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "clientEncryptionKeyName")

payload := openapis.ClientEncryptionKeyCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateClientEncryptionKeyThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesCreateUpdateSqlContainer`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

payload := openapis.SqlContainerCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlContainerThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesCreateUpdateSqlDatabase`

```go
ctx := context.TODO()
id := openapis.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName")

payload := openapis.SqlDatabaseCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlDatabaseThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesCreateUpdateSqlRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewSqlRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

payload := openapis.SqlRoleAssignmentCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlRoleAssignmentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesCreateUpdateSqlRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewSqlRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

payload := openapis.SqlRoleDefinitionCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlRoleDefinitionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesCreateUpdateSqlStoredProcedure`

```go
ctx := context.TODO()
id := openapis.NewStoredProcedureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName", "storedProcedureName")

payload := openapis.SqlStoredProcedureCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlStoredProcedureThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesCreateUpdateSqlTrigger`

```go
ctx := context.TODO()
id := openapis.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName", "triggerName")

payload := openapis.SqlTriggerCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlTriggerThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesCreateUpdateSqlUserDefinedFunction`

```go
ctx := context.TODO()
id := openapis.NewUserDefinedFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName", "userDefinedFunctionName")

payload := openapis.SqlUserDefinedFunctionCreateUpdateParameters{
	// ...
}


if err := client.SqlResourcesCreateUpdateSqlUserDefinedFunctionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesDeleteSqlContainer`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

if err := client.SqlResourcesDeleteSqlContainerThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesDeleteSqlDatabase`

```go
ctx := context.TODO()
id := openapis.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName")

if err := client.SqlResourcesDeleteSqlDatabaseThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesDeleteSqlRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewSqlRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

if err := client.SqlResourcesDeleteSqlRoleAssignmentThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesDeleteSqlRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewSqlRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

if err := client.SqlResourcesDeleteSqlRoleDefinitionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesDeleteSqlStoredProcedure`

```go
ctx := context.TODO()
id := openapis.NewStoredProcedureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName", "storedProcedureName")

if err := client.SqlResourcesDeleteSqlStoredProcedureThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesDeleteSqlTrigger`

```go
ctx := context.TODO()
id := openapis.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName", "triggerName")

if err := client.SqlResourcesDeleteSqlTriggerThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesDeleteSqlUserDefinedFunction`

```go
ctx := context.TODO()
id := openapis.NewUserDefinedFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName", "userDefinedFunctionName")

if err := client.SqlResourcesDeleteSqlUserDefinedFunctionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesGetClientEncryptionKey`

```go
ctx := context.TODO()
id := openapis.NewClientEncryptionKeyID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "clientEncryptionKeyName")

read, err := client.SqlResourcesGetClientEncryptionKey(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.SqlResourcesGetSqlContainer`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

read, err := client.SqlResourcesGetSqlContainer(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.SqlResourcesGetSqlContainerThroughput`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

read, err := client.SqlResourcesGetSqlContainerThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.SqlResourcesGetSqlDatabase`

```go
ctx := context.TODO()
id := openapis.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName")

read, err := client.SqlResourcesGetSqlDatabase(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.SqlResourcesGetSqlDatabaseThroughput`

```go
ctx := context.TODO()
id := openapis.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName")

read, err := client.SqlResourcesGetSqlDatabaseThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.SqlResourcesGetSqlRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewSqlRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

read, err := client.SqlResourcesGetSqlRoleAssignment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.SqlResourcesGetSqlRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewSqlRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

read, err := client.SqlResourcesGetSqlRoleDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.SqlResourcesGetSqlStoredProcedure`

```go
ctx := context.TODO()
id := openapis.NewStoredProcedureID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName", "storedProcedureName")

read, err := client.SqlResourcesGetSqlStoredProcedure(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.SqlResourcesGetSqlTrigger`

```go
ctx := context.TODO()
id := openapis.NewTriggerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName", "triggerName")

read, err := client.SqlResourcesGetSqlTrigger(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.SqlResourcesGetSqlUserDefinedFunction`

```go
ctx := context.TODO()
id := openapis.NewUserDefinedFunctionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName", "userDefinedFunctionName")

read, err := client.SqlResourcesGetSqlUserDefinedFunction(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.SqlResourcesListClientEncryptionKeys`

```go
ctx := context.TODO()
id := openapis.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName")

// alternatively `client.SqlResourcesListClientEncryptionKeys(ctx, id)` can be used to do batched pagination
items, err := client.SqlResourcesListClientEncryptionKeysComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.SqlResourcesListSqlContainers`

```go
ctx := context.TODO()
id := openapis.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName")

// alternatively `client.SqlResourcesListSqlContainers(ctx, id)` can be used to do batched pagination
items, err := client.SqlResourcesListSqlContainersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.SqlResourcesListSqlDatabases`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.SqlResourcesListSqlDatabases(ctx, id)` can be used to do batched pagination
items, err := client.SqlResourcesListSqlDatabasesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.SqlResourcesListSqlRoleAssignments`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.SqlResourcesListSqlRoleAssignments(ctx, id)` can be used to do batched pagination
items, err := client.SqlResourcesListSqlRoleAssignmentsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.SqlResourcesListSqlRoleDefinitions`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.SqlResourcesListSqlRoleDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.SqlResourcesListSqlRoleDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.SqlResourcesListSqlStoredProcedures`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

// alternatively `client.SqlResourcesListSqlStoredProcedures(ctx, id)` can be used to do batched pagination
items, err := client.SqlResourcesListSqlStoredProceduresComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.SqlResourcesListSqlTriggers`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

// alternatively `client.SqlResourcesListSqlTriggers(ctx, id)` can be used to do batched pagination
items, err := client.SqlResourcesListSqlTriggersComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.SqlResourcesListSqlUserDefinedFunctions`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

// alternatively `client.SqlResourcesListSqlUserDefinedFunctions(ctx, id)` can be used to do batched pagination
items, err := client.SqlResourcesListSqlUserDefinedFunctionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.SqlResourcesMigrateSqlContainerToAutoscale`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

if err := client.SqlResourcesMigrateSqlContainerToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesMigrateSqlContainerToManualThroughput`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

if err := client.SqlResourcesMigrateSqlContainerToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesMigrateSqlDatabaseToAutoscale`

```go
ctx := context.TODO()
id := openapis.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName")

if err := client.SqlResourcesMigrateSqlDatabaseToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesMigrateSqlDatabaseToManualThroughput`

```go
ctx := context.TODO()
id := openapis.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName")

if err := client.SqlResourcesMigrateSqlDatabaseToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesRetrieveContinuousBackupInformation`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

payload := openapis.ContinuousBackupRestoreLocation{
	// ...
}


if err := client.SqlResourcesRetrieveContinuousBackupInformationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesUpdateSqlContainerThroughput`

```go
ctx := context.TODO()
id := openapis.NewContainerID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName", "containerName")

payload := openapis.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.SqlResourcesUpdateSqlContainerThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.SqlResourcesUpdateSqlDatabaseThroughput`

```go
ctx := context.TODO()
id := openapis.NewSqlDatabaseID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "sqlDatabaseName")

payload := openapis.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.SqlResourcesUpdateSqlDatabaseThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.TableResourcesCreateUpdateTable`

```go
ctx := context.TODO()
id := openapis.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "tableName")

payload := openapis.TableCreateUpdateParameters{
	// ...
}


if err := client.TableResourcesCreateUpdateTableThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.TableResourcesCreateUpdateTableRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewTableRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

payload := openapis.TableRoleAssignmentResource{
	// ...
}


if err := client.TableResourcesCreateUpdateTableRoleAssignmentThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.TableResourcesCreateUpdateTableRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewTableRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

payload := openapis.TableRoleDefinitionResource{
	// ...
}


if err := client.TableResourcesCreateUpdateTableRoleDefinitionThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.TableResourcesDeleteTable`

```go
ctx := context.TODO()
id := openapis.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "tableName")

if err := client.TableResourcesDeleteTableThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.TableResourcesDeleteTableRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewTableRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

if err := client.TableResourcesDeleteTableRoleAssignmentThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.TableResourcesDeleteTableRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewTableRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

if err := client.TableResourcesDeleteTableRoleDefinitionThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.TableResourcesGetTable`

```go
ctx := context.TODO()
id := openapis.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "tableName")

read, err := client.TableResourcesGetTable(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.TableResourcesGetTableRoleAssignment`

```go
ctx := context.TODO()
id := openapis.NewTableRoleAssignmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleAssignmentId")

read, err := client.TableResourcesGetTableRoleAssignment(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.TableResourcesGetTableRoleDefinition`

```go
ctx := context.TODO()
id := openapis.NewTableRoleDefinitionID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "roleDefinitionId")

read, err := client.TableResourcesGetTableRoleDefinition(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.TableResourcesGetTableThroughput`

```go
ctx := context.TODO()
id := openapis.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "tableName")

read, err := client.TableResourcesGetTableThroughput(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `OpenapisClient.TableResourcesListTableRoleAssignments`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.TableResourcesListTableRoleAssignments(ctx, id)` can be used to do batched pagination
items, err := client.TableResourcesListTableRoleAssignmentsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.TableResourcesListTableRoleDefinitions`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.TableResourcesListTableRoleDefinitions(ctx, id)` can be used to do batched pagination
items, err := client.TableResourcesListTableRoleDefinitionsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.TableResourcesListTables`

```go
ctx := context.TODO()
id := openapis.NewDatabaseAccountID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName")

// alternatively `client.TableResourcesListTables(ctx, id)` can be used to do batched pagination
items, err := client.TableResourcesListTablesComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `OpenapisClient.TableResourcesMigrateTableToAutoscale`

```go
ctx := context.TODO()
id := openapis.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "tableName")

if err := client.TableResourcesMigrateTableToAutoscaleThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.TableResourcesMigrateTableToManualThroughput`

```go
ctx := context.TODO()
id := openapis.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "tableName")

if err := client.TableResourcesMigrateTableToManualThroughputThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.TableResourcesRetrieveContinuousBackupInformation`

```go
ctx := context.TODO()
id := openapis.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "tableName")

payload := openapis.ContinuousBackupRestoreLocation{
	// ...
}


if err := client.TableResourcesRetrieveContinuousBackupInformationThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `OpenapisClient.TableResourcesUpdateTableThroughput`

```go
ctx := context.TODO()
id := openapis.NewTableID("12345678-1234-9876-4563-123456789012", "example-resource-group", "databaseAccountName", "tableName")

payload := openapis.ThroughputSettingsUpdateParameters{
	// ...
}


if err := client.TableResourcesUpdateTableThroughputThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
