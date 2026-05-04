## Example: NetApp Files Volume Bucket (Object REST API) - Inline credentials

This example provisions a NetApp Files volume with an S3-compatible bucket and generates an inline credential pair: the `access_key` / `secret_key` are returned directly by the Azure NetApp Files API and persisted in Terraform state as sensitive attributes on `azurerm_netapp_volume_bucket_credentials`.

~> **Recommended:** For any non-throwaway workload use the [Azure Key Vault-backed example](../volume_bucket_akv) instead. The Key Vault variant avoids storing the bucket access / secret keys in Terraform state and matches the Azure NetApp Files security guidance for the Object REST API.

The Object REST API feature is in preview and must be registered on the subscription before buckets can be created:

```bash
az feature register --namespace Microsoft.NetApp --name ANFEnableObjectRESTAPI
```

### Consuming the credentials

The generated keys are exposed as sensitive resource attributes and re-exported via [`outputs.tf`](outputs.tf). After `terraform apply` you can read them with:

```bash
terraform output -raw bucket_access_key
terraform output -raw bucket_secret_key
```

You can also reference them directly from another resource or provider, e.g. to configure an S3-compatible client:

```hcl
provider "aws" {
  alias      = "anf_bucket"
  access_key = azurerm_netapp_volume_bucket_credentials.example.access_key
  secret_key = azurerm_netapp_volume_bucket_credentials.example.secret_key

  endpoints {
    s3 = "https://${azurerm_netapp_volume_bucket.example.server_ip_address}"
  }
}
```

~> **Note:** Because the keys live in Terraform state, anyone with access to the state file can read them. Use a state backend with encryption and restricted access, and prefer the Key Vault example for production.

### Rotation

`key_pair_expiry_days` is `ForceNew`. To rotate the credentials, taint the resource (or change `key_pair_expiry_days`) and re-apply - this generates a new key pair and immediately invalidates the previous one.

### Variables

* `prefix` - (Required) The prefix used for all resources in this example.

* `location` - (Required) The Azure Region in which the resources in this example should be created.
