
## `github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/dnssecconfigs` Documentation

The `dnssecconfigs` SDK allows for interaction with Azure Resource Manager `dns` (API Version `2023-07-01-preview`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/dnssecconfigs"
```


### Client Initialization

```go
client := dnssecconfigs.NewDnssecConfigsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `DnssecConfigsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := dnssecconfigs.NewDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsZoneName")

if err := client.CreateOrUpdateThenPoll(ctx, id, dnssecconfigs.DefaultCreateOrUpdateOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DnssecConfigsClient.Delete`

```go
ctx := context.TODO()
id := dnssecconfigs.NewDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsZoneName")

if err := client.DeleteThenPoll(ctx, id, dnssecconfigs.DefaultDeleteOperationOptions()); err != nil {
	// handle the error
}
```


### Example Usage: `DnssecConfigsClient.Get`

```go
ctx := context.TODO()
id := dnssecconfigs.NewDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsZoneName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `DnssecConfigsClient.ListByDnsZone`

```go
ctx := context.TODO()
id := dnssecconfigs.NewDnsZoneID("12345678-1234-9876-4563-123456789012", "example-resource-group", "dnsZoneName")

// alternatively `client.ListByDnsZone(ctx, id)` can be used to do batched pagination
items, err := client.ListByDnsZoneComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```
