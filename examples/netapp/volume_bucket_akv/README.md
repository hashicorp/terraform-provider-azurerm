## Example: NetApp Files Volume Bucket (Object REST API) - Key Vault backed (recommended)

This example provisions a NetApp Files volume with an S3-compatible bucket whose server certificate is sourced from Azure Key Vault and whose generated credentials are persisted to a second Azure Key Vault.

This is the **recommended** way to use the Azure NetApp Files Object REST API: the bucket access key / secret key never leave Azure Key Vault and are never written to Terraform state, which matches the security guidance in the public documentation.

The Object REST API feature is in preview and must be registered on the subscription before buckets can be created:

```bash
az feature register --namespace Microsoft.NetApp --name ANFEnableObjectRESTAPI
```

### What this example provisions

The example uses two separate Key Vaults (per the Azure NetApp Files recommendation): one read-mostly vault for the certificate and one write vault for the credentials. The Azure NetApp Files first-party service principal (`2b6fb936-77b9-4775-b03e-37edae8ab84b`) is granted the documented certificate / secret permissions on the corresponding vault via `azurerm_key_vault_access_policy`.

The bucket server certificate is created here as a self-signed certificate solely for example purposes - replace it with an imported CA-signed certificate for production. Its Subject Alternative Name must match the `server_fqdn` value passed to the bucket.

The companion `azurerm_netapp_volume_bucket_credentials` resource is configured with `store_in_key_vault = true`, which causes Azure NetApp Files to write the generated `access_key_id` / `secret_access_key` to the credentials vault as a JSON secret instead of returning them inline.

### Consuming the credentials

The bucket credentials are stored in the credentials Key Vault as a single secret with the name supplied via `key_vault.0.credentials_secret_name`. The secret value is JSON, e.g.:

```json
{
  "access_key_id": "...",
  "secret_access_key": "..."
}
```

Read it from Terraform with `azurerm_key_vault_secret`, or out-of-band with the Azure CLI:

```bash
az keyvault secret show \
  --vault-name <prefix>-creds-kv \
  --name <prefix>-bucket-creds \
  --query value -o tsv
```

### Rotation

`key_pair_expiry_days` and `store_in_key_vault` are `ForceNew`. To rotate the credentials, taint the resource (or change `key_pair_expiry_days`) and re-apply - this generates a new key pair, immediately invalidates the previous one, and overwrites the secret in the credentials Key Vault.

### Variables

* `prefix` - (Required) The prefix used for all resources in this example.

* `location` - (Required) The Azure Region in which the resources in this example should be created.

* `server_fqdn` - (Required) The DNS name that will resolve to the bucket endpoint IP address. Must match the Subject Alternative Name of the bucket server certificate.
