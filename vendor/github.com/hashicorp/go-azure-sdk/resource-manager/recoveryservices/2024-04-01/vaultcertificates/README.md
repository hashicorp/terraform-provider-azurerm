
## `github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-04-01/vaultcertificates` Documentation

The `vaultcertificates` SDK allows for interaction with Azure Resource Manager `recoveryservices` (API Version `2024-04-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-04-01/vaultcertificates"
```


### Client Initialization

```go
client := vaultcertificates.NewVaultCertificatesClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `VaultCertificatesClient.Create`

```go
ctx := context.TODO()
id := vaultcertificates.NewCertificateID("12345678-1234-9876-4563-123456789012", "example-resource-group", "vaultName", "certificateName")

payload := vaultcertificates.CertificateRequest{
	// ...
}


read, err := client.Create(ctx, id, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
