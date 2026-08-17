
## `github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets` Documentation

The `buckets` SDK allows for interaction with Azure Resource Manager `netapp` (API Version `2026-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
```


### Client Initialization

```go
client := buckets.NewBucketsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `BucketsClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := buckets.NewBucketID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "bucketName")

payload := buckets.Bucket{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BucketsClient.Delete`

```go
ctx := context.TODO()
id := buckets.NewBucketID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "bucketName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BucketsClient.GenerateAkvCredentials`

```go
ctx := context.TODO()
id := buckets.NewBucketID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "bucketName")

payload := buckets.BucketCredentialsExpiry{
	// ...
}


if err := client.GenerateAkvCredentialsThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `BucketsClient.GenerateCredentials`

```go
ctx := context.TODO()
id := buckets.NewBucketID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "bucketName")

payload := buckets.BucketCredentialsExpiry{
	// ...
}


read, err := client.GenerateCredentials(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BucketsClient.Get`

```go
ctx := context.TODO()
id := buckets.NewBucketID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "bucketName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `BucketsClient.List`

```go
ctx := context.TODO()
id := buckets.NewVolumeID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName")

// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `BucketsClient.RefreshCertificate`

```go
ctx := context.TODO()
id := buckets.NewBucketID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "bucketName")

if err := client.RefreshCertificateThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `BucketsClient.Update`

```go
ctx := context.TODO()
id := buckets.NewBucketID("12345678-1234-9876-4563-123456789012", "example-resource-group", "netAppAccountName", "capacityPoolName", "volumeName", "bucketName")

payload := buckets.BucketPatch{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
