
## `github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/scclusterrecords` Documentation

The `scclusterrecords` SDK allows for interaction with Azure Resource Manager `confluent` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/scclusterrecords"
```


### Client Initialization

```go
client := scclusterrecords.NewSCClusterRecordsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SCClusterRecordsClient.ClusterCreateOrUpdate`

```go
ctx := context.TODO()
id := scclusterrecords.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId")

payload := scclusterrecords.SCClusterRecord{
	// ...
}


read, err := client.ClusterCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SCClusterRecordsClient.ClusterDelete`

```go
ctx := context.TODO()
id := scclusterrecords.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId")

if err := client.ClusterDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SCClusterRecordsClient.OrganizationCreateAPIKey`

```go
ctx := context.TODO()
id := scclusterrecords.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId")

payload := scclusterrecords.CreateAPIKeyModel{
	// ...
}


read, err := client.OrganizationCreateAPIKey(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SCClusterRecordsClient.OrganizationGetClusterById`

```go
ctx := context.TODO()
id := scclusterrecords.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId", "clusterId")

read, err := client.OrganizationGetClusterById(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SCClusterRecordsClient.OrganizationListClusters`

```go
ctx := context.TODO()
id := scclusterrecords.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId")

// alternatively `client.OrganizationListClusters(ctx, id, scclusterrecords.DefaultOrganizationListClustersOperationOptions())` can be used to do batched pagination
items, err := client.OrganizationListClustersComplete(ctx, id, scclusterrecords.DefaultOrganizationListClustersOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
