
## `github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-10-01/vpnserverconfigurations` Documentation

The `vpnserverconfigurations` SDK allows for interaction with Azure Resource Manager `network` (API Version `2024-10-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-10-01/vpnserverconfigurations"
```


### Client Initialization

```go
client := vpnserverconfigurations.NewVpnServerConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VpnServerConfigurationsClient.ListRadiusSecrets`

```go
ctx := context.TODO()
id := vpnserverconfigurations.NewVpnServerConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationName")

// alternatively `client.ListRadiusSecrets(ctx, id)` can be used to do batched pagination
items, err := client.ListRadiusSecretsComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `VpnServerConfigurationsClient.UpdateTags`

```go
ctx := context.TODO()
id := vpnserverconfigurations.NewVpnServerConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vpnServerConfigurationName")

payload := vpnserverconfigurations.TagsObject{
	// ...
}


read, err := client.UpdateTags(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
