
## `github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/getprivatednszonesuffix` Documentation

The `getprivatednszonesuffix` SDK allows for interaction with Azure Resource Manager `mysql` (API Version `2022-01-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2022-01-01/getprivatednszonesuffix"
```


### Client Initialization

```go
client := getprivatednszonesuffix.NewGetPrivateDnsZoneSuffixClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
```


### Example Usage: `GetPrivateDnsZoneSuffixClient.Execute`

```go
ctx := context.TODO()


read, err := client.Execute(ctx)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
