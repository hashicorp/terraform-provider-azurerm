
## `github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns` Documentation

The `privatedns` SDK allows for interaction with Azure Resource Manager `privatedns` (API Version `2024-06-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatedns"
```


### Client Initialization

```go
client := privatedns.NewPrivateDNSClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `PrivateDNSClient.RecordSetsCreateOrUpdate`

```go
ctx := context.TODO()
id := privatedns.NewRecordTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateDnsZoneName", "A", "relativeRecordSetName")

payload := privatedns.RecordSet{
	// ...
}


read, err := client.RecordSetsCreateOrUpdate(ctx, id, payload, privatedns.DefaultRecordSetsCreateOrUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateDNSClient.RecordSetsDelete`

```go
ctx := context.TODO()
id := privatedns.NewRecordTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateDnsZoneName", "A", "relativeRecordSetName")

read, err := client.RecordSetsDelete(ctx, id, privatedns.DefaultRecordSetsDeleteOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateDNSClient.RecordSetsGet`

```go
ctx := context.TODO()
id := privatedns.NewRecordTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateDnsZoneName", "A", "relativeRecordSetName")

read, err := client.RecordSetsGet(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `PrivateDNSClient.RecordSetsListByType`

```go
ctx := context.TODO()
id := privatedns.NewPrivateZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateDnsZoneName", "A")

// alternatively `client.RecordSetsListByType(ctx, id, privatedns.DefaultRecordSetsListByTypeOperationOptions())` can be used to do batched pagination
items, err := client.RecordSetsListByTypeComplete(ctx, id, privatedns.DefaultRecordSetsListByTypeOperationOptions())
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `PrivateDNSClient.RecordSetsUpdate`

```go
ctx := context.TODO()
id := privatedns.NewRecordTypeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "privateDnsZoneName", "A", "relativeRecordSetName")

payload := privatedns.RecordSet{
	// ...
}


read, err := client.RecordSetsUpdate(ctx, id, payload, privatedns.DefaultRecordSetsUpdateOperationOptions())
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
