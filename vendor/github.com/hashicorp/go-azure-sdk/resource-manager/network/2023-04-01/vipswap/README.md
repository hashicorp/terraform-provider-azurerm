
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/vipswap` Documentation

The `vipswap` SDK allows for interaction with the Azure Resource Manager Service `network` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/vipswap"
```


### Client Initialization

```go
client := vipswap.NewVipSwapClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VipSwapClient.Create`

```go
ctx := context.TODO()
id := vipswap.NewCloudServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupValue", "cloudServiceValue")

payload := vipswap.SwapResource{
	// ...
}


if err := client.CreateThenPoll(ctx, id, payload); err != nil {
	// handle the error
}
```


### Example Usage: `VipSwapClient.Get`

```go
ctx := context.TODO()
id := vipswap.NewCloudServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupValue", "cloudServiceValue")

read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `VipSwapClient.List`

```go
ctx := context.TODO()
id := vipswap.NewCloudServiceID("12345678-1234-9876-4563-123456789012", "resourceGroupValue", "cloudServiceValue")

read, err := client.List(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
