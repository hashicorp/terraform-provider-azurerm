
## `github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapdiskconfigurations` Documentation

The `sapdiskconfigurations` SDK allows for interaction with Azure Resource Manager `workloads` (API Version `2023-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapdiskconfigurations"
```


### Client Initialization

```go
client := sapdiskconfigurations.NewSAPDiskConfigurationsClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `SAPDiskConfigurationsClient.SAPDiskConfigurations`

```go
ctx := context.TODO()
id := sapdiskconfigurations.NewLocationID("12345678-1234-9876-4563-123456789012", "locationName")

payload := sapdiskconfigurations.SAPDiskConfigurationsRequest{
	// ...
}


read, err := client.SAPDiskConfigurations(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
