
## `github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/recordsets` Documentation

The `recordsets` SDK allows for interaction with the Azure Resource Manager Service `privatedns` (API Version `2018-09-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/recordsets"
```


### Client Initialization

```go
client := recordsets.NewRecordSetsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `RecordSetsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := recordsets.NewRecordTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue", "A", "relativeRecordSetValue")

payload := recordsets.RecordSet{
	// ...
}


read, err := client.CreateOrUpdate(ctx, id, payload, recordsets.DefaultCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RecordSetsClient.Delete`

```go
ctx := context.TODO()
id := recordsets.NewRecordTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue", "A", "relativeRecordSetValue")

read, err := client.Delete(ctx, id, recordsets.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RecordSetsClient.Get`

```go
ctx := context.TODO()
id := recordsets.NewRecordTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue", "A", "relativeRecordSetValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `RecordSetsClient.List`

```go
ctx := context.TODO()
id := recordsets.NewPrivateDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue")

// alternatively `client.List(ctx, id, recordsets.DefaultListOperationOptions())` can be used to do batched pagination
items, err := client.ListComplete(ctx, id, recordsets.DefaultListOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RecordSetsClient.ListByType`

```go
ctx := context.TODO()
id := recordsets.NewPrivateZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue", "A")

// alternatively `client.ListByType(ctx, id, recordsets.DefaultListByTypeOperationOptions())` can be used to do batched pagination
items, err := client.ListByTypeComplete(ctx, id, recordsets.DefaultListByTypeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `RecordSetsClient.Update`

```go
ctx := context.TODO()
id := recordsets.NewRecordTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateZoneValue", "A", "relativeRecordSetValue")

payload := recordsets.RecordSet{
	// ...
}


read, err := client.Update(ctx, id, payload, recordsets.DefaultUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
