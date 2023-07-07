
## `github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/cluster` Documentation

The `cluster` SDK allows for interaction with the Azure Resource Manager Service `azurestackhci` (API Version `2023-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/cluster"
```


### Client Initialization

```go
client := cluster.NewClusterClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `ClusterClient.CreateIdentity`

```go
ctx := context.TODO()
id := cluster.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

if err := client.CreateIdentityThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `ClusterClient.ExtendSoftwareAssuranceBenefit`

```go
ctx := context.TODO()
id := cluster.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := cluster.SoftwareAssuranceChangeRequest{
	// ...
}


if err := client.ExtendSoftwareAssuranceBenefitThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `ClusterClient.UploadCertificate`

```go
ctx := context.TODO()
id := cluster.NewClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterValue")

payload := cluster.UploadCertificateRequest{
	// ...
}


if err := client.UploadCertificateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
