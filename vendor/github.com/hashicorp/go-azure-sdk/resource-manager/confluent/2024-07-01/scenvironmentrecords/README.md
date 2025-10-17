
## `github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/scenvironmentrecords` Documentation

The `scenvironmentrecords` SDK allows for interaction with Azure Resource Manager `confluent` (API Version `2024-07-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/scenvironmentrecords"
```


### Client Initialization

```go
client := scenvironmentrecords.NewSCEnvironmentRecordsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SCEnvironmentRecordsClient.EnvironmentCreateOrUpdate`

```go
ctx := context.TODO()
id := scenvironmentrecords.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId")

payload := scenvironmentrecords.SCEnvironmentRecord{
	// ...
}


read, err := client.EnvironmentCreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SCEnvironmentRecordsClient.EnvironmentDelete`

```go
ctx := context.TODO()
id := scenvironmentrecords.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId")

if err := client.EnvironmentDeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `SCEnvironmentRecordsClient.OrganizationGetEnvironmentById`

```go
ctx := context.TODO()
id := scenvironmentrecords.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId")

read, err := client.OrganizationGetEnvironmentById(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `SCEnvironmentRecordsClient.OrganizationListEnvironments`

```go
ctx := context.TODO()
id := scenvironmentrecords.NewOrganizationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName")

// alternatively `client.OrganizationListEnvironments(ctx, id, scenvironmentrecords.DefaultOrganizationListEnvironmentsOperationOptions())` can be used to do batched pagination
items, err := client.OrganizationListEnvironmentsComplete(ctx, id, scenvironmentrecords.DefaultOrganizationListEnvironmentsOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `SCEnvironmentRecordsClient.OrganizationListSchemaRegistryClusters`

```go
ctx := context.TODO()
id := scenvironmentrecords.NewEnvironmentID("12345678-1234-9876-4563-123456789012", "example-resource-group", "organizationName", "environmentId")

// alternatively `client.OrganizationListSchemaRegistryClusters(ctx, id, scenvironmentrecords.DefaultOrganizationListSchemaRegistryClustersOperationOptions())` can be used to do batched pagination
items, err := client.OrganizationListSchemaRegistryClustersComplete(ctx, id, scenvironmentrecords.DefaultOrganizationListSchemaRegistryClustersOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
