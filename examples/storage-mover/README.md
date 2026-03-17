# Storage Mover examples (SMB mount endpoint)

Creates: `azurerm_storage_mover`, target endpoint, source endpoint, **SMB mount endpoint**, project, job definition.

## Usage

1. Build: `make build` (from repo root)
2. Use local provider via `dev_overrides` in `~/.terraformrc`
3. Set ARM_* env vars, then: `terraform init && terraform plan`

Set `smb_host` and `smb_share_name` (and optionally credentials) for apply.

## Validation

Copy `validation-test-invalid.tf.example` to `validation-test-invalid.tf` and run `terraform plan` to verify invalid inputs are rejected.
