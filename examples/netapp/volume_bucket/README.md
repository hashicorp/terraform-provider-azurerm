## Example: NetApp Files Volume Buckets (Object REST API) - inline certificate

This example provisions a NetApp Files volume with two S3-compatible buckets that expose the volume through the Azure NetApp Files Object REST API. The bucket server certificate is supplied inline through `server.certificate_pem`.

The first bucket on a volume establishes the bucket server (FQDN and certificate) and is created with the `azurerm_netapp_volume_bucket_with_server` resource. Every subsequent bucket reuses that server configuration and is created with the server-less `azurerm_netapp_volume_bucket` resource. Declaring a `server` block on more than one bucket would overwrite the shared server configuration.

The Object REST API feature is in preview and must be registered on the subscription before buckets can be created:

```bash
az feature register --namespace Microsoft.NetApp --name ANFEnableObjectRESTAPI
```

### What this example provisions

The bucket server certificate is created here as a self-signed certificate solely for example purposes - replace it with a CA-signed certificate for production. Its Subject Alternative Name must match the `server_fqdn` value passed to the first bucket.

The example does not generate bucket credentials. To mint access keys, configure a `key_vault` block on the bucket and use the `azurerm_netapp_volume_bucket_credentials` action (see the `volume_bucket_akv` example).

### Variables

* `prefix` - (Required) The prefix used for all resources in this example.

* `location` - (Required) The Azure Region in which the resources in this example should be created.

* `server_fqdn` - (Required) The DNS name that will resolve to the bucket endpoint IP address. Must match the Subject Alternative Name of the bucket server certificate.
