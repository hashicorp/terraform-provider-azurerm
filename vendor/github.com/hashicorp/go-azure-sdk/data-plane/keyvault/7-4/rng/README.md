
## `github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/rng` Documentation

The `rng` SDK allows for interaction with <unknown source data type> `keyvault` (API Version `7.4`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/data-plane/keyvault/7.4/rng"
```


### Client Initialization

```go
client := rng.NewRNGClientWithBaseURI("")
client.Client.Authorizer = authorizer
```


### Example Usage: `RNGClient.GetRandomBytes`

```go
ctx := context.TODO()

payload := rng.GetRandomBytesRequest{
	// ...
}


read, err := client.GetRandomBytes(ctx, payload)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```
