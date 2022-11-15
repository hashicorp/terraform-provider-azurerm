## Example: Private Endpoint

This example provisions a Private Endpoint which connects to an application gateway within Azure.

The `AllowApplicationGatewayPrivateLink` feature must be registered on the subscription:

```bash
az feature register --name AllowApplicationGatewayPrivateLink --namespace Microsoft.Network
```

### Variables

* `prefix` - (Required) The prefix used for all resources in this example.

* `location` - (Required) The Azure Region in which all resources in this example should be created.
