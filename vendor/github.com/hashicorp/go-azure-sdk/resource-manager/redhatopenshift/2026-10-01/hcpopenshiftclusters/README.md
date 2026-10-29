
## `github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/hcpopenshiftclusters` Documentation

The `hcpopenshiftclusters` SDK allows for interaction with Azure Resource Manager `redhatopenshift` (API Version `2026-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
import "github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2026-10-01/hcpopenshiftclusters"
```


### Client Initialization

```go
client := hcpopenshiftclusters.NewHcpOpenShiftClustersClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `HcpOpenShiftClustersClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := hcpopenshiftclusters.NewHcpOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName")

payload := hcpopenshiftclusters.HcpOpenShiftCluster{
	// ...
}


if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HcpOpenShiftClustersClient.Delete`

```go
ctx := context.TODO()
id := hcpopenshiftclusters.NewHcpOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName")

if err := client.DeleteThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `HcpOpenShiftClustersClient.Get`

```go
ctx := context.TODO()
id := hcpopenshiftclusters.NewHcpOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `HcpOpenShiftClustersClient.ListByResourceGroup`

```go
ctx := context.TODO()
id := commonids.NewResourceGroupID("12345678-1234-9876-4563-123456789012", "example-resource-group")

// alternatively `client.ListByResourceGroup(ctx, id)` can be used to do batched pagination
items, err := client.ListByResourceGroupComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HcpOpenShiftClustersClient.ListBySubscription`

```go
ctx := context.TODO()
id := commonids.NewSubscriptionID("12345678-1234-9876-4563-123456789012")

// alternatively `client.ListBySubscription(ctx, id)` can be used to do batched pagination
items, err := client.ListBySubscriptionComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `HcpOpenShiftClustersClient.RequestAdminCredential`

```go
ctx := context.TODO()
id := hcpopenshiftclusters.NewHcpOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName")

payload := hcpopenshiftclusters.HcpOpenShiftClusterAdminCredentialRequest{
	// ...
}


if err := client.RequestAdminCredentialThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `HcpOpenShiftClustersClient.RevokeCredentials`

```go
ctx := context.TODO()
id := hcpopenshiftclusters.NewHcpOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName")

if err := client.RevokeCredentialsThenPoll(ctx, id); err != nil {
	// handle the error
}
```


### Example Usage: `HcpOpenShiftClustersClient.Update`

```go
ctx := context.TODO()
id := hcpopenshiftclusters.NewHcpOpenShiftClusterID("12345678-1234-9876-4563-123456789012", "example-resource-group", "hcpOpenShiftClusterName")

payload := hcpopenshiftclusters.HcpOpenShiftClusterUpdate{
	// ...
}


if err := client.UpdateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```
