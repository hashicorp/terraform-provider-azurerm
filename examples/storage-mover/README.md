# Storage Mover examples

This example creates:

- `azurerm_storage_mover`
- `azurerm_storage_mover_target_endpoint` (Azure Blob container)
- `azurerm_storage_mover_source_endpoint` (NFS)

## Using the locally built provider

1. Build the provider: `make build` (from repo root)
2. Configure Terraform to use the local binary in `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "hashicorp/azurerm" = "/path/to/go/bin"   # or $GOPATH/bin
  }
  direct {}
}
```

3. Set Azure credentials (e.g. `ARM_CLIENT_ID`, `ARM_CLIENT_SECRET`, `ARM_TENANT_ID`, `ARM_SUBSCRIPTION_ID`)
4. Run: `terraform init && terraform plan` (or `terraform apply`)

## Multi-cloud connector endpoint

`azurerm_storage_mover_multi_cloud_connector_endpoint` requires:

- `multi_cloud_connector_id`: a `Microsoft.HybridConnectivity/publicCloudConnectors` resource ID
- `aws_s3_bucket_id`: a `Microsoft.AwsConnector/s3Buckets` resource ID

Those resources are created via Azure Arc Multicloud Connector and cannot be created by Terraform here. Uncomment the resource and variables in `main.tf` / `variables.tf` and set the IDs to test.

## Validation

To test that invalid inputs are rejected, copy the validation test file and run plan:

```bash
cp validation-test-invalid.tf.example validation-test-invalid.tf
terraform plan -input=false
```

Expected: plan fails with validation errors. Remove `validation-test-invalid.tf` to use the example normally.

The provider validates:

- **Storage Mover**: name non-empty
- **Target endpoint**: name 1–64 chars (alphanumeric, `_`, `-`); storage account ID; container name 3–63 chars, lowercase
- **Source endpoint**: name 1–64 chars; host non-empty; `nfs_version` one of `NFSauto`, `NFSv4`, `NFSv3`
- **NFS file share endpoint**: name 1–64 chars; storage mover ID; storage account ID; file share name 3–63 chars, lowercase alphanumeric and hyphens
- **Project**: name 1–64 chars; storage mover ID (Microsoft.StorageMover)
- **Job definition**: name 1–64 chars; project ID (Microsoft.StorageMover/.../projects/...); `copy_mode` one of `Additive`, `Mirror`
- **Multi-cloud connector endpoint** (on hold): name 1–64 chars; multi_cloud_connector_id and aws_s3_bucket_id must be the correct Azure resource ID types

Invalid IDs or formats fail at plan time.
