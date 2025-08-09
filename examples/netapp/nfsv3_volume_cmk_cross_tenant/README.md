# NetApp Account Encryption with Cross-Tenant Customer-Managed Keys

This example demonstrates how to configure NetApp Account Encryption using a customer-managed key from a key vault that exists in a different Azure Active Directory tenant (cross-tenant scenario).

## Overview

In enterprise environments, it's common to have a centralized key management setup where encryption keys are stored in a key vault that belongs to a different Azure AD tenant than where the NetApp resources are deployed. The `federated_client_id` parameter enables this cross-tenant access by specifying the client ID of a multi-tenant Azure AD application that has been granted access to the remote key vault.

## Real-World Scenario

This example assumes:
- **Tenant A** (Current): Where you're deploying NetApp resources
- **Tenant B** (Remote): Where the key vault and encryption keys already exist
- The key vault in Tenant B has been configured to allow access from your multi-tenant application

## Prerequisites

### 1. Multi-tenant Azure AD Application
You need to create a multi-tenant Azure AD application that can be used to access the key vault across tenants.

### 2. Cross-tenant Key Vault Setup (Done by Tenant B administrators)
The remote tenant administrators must:
- Create/configure the key vault with appropriate encryption keys
- Grant access permissions to your multi-tenant application
- Optionally configure private endpoint approval policies

### 3. Network Connectivity
- Private endpoint connection from your VNet to the remote key vault
- Appropriate DNS resolution for the remote key vault

## Required Variables

- `cross_tenant_key_vault_name`: Name of the existing key vault in the remote tenant
- `cross_tenant_resource_group_name`: Resource group where the remote key vault exists
- `cross_tenant_key_name`: Name of the encryption key in the remote key vault
- `federated_client_id`: Client ID of the multi-tenant Azure AD application
- `user_assigned_identity_id`: Resource ID of the pre-created user-assigned managed identity used by the NetApp account
- `cross_tenant_key_vault_resource_id`: Full resource ID of the cross-tenant key vault (recommended for proper validation)
- `private_endpoint_manual_approval`: Whether private endpoint requires manual approval (typically true for cross-tenant)

## Understanding `is_manual_connection`

The `is_manual_connection` parameter in the private endpoint configuration determines:

- **`false`**: Automatic approval - The connection is approved automatically (requires appropriate RBAC permissions)
- **`true`**: Manual approval - The connection requires manual approval by the target resource owner

In cross-tenant scenarios, `is_manual_connection = true` is often required because:
1. You typically don't have automatic approval rights in the remote tenant
2. Security policies in enterprise environments often require manual approval for cross-tenant connections
3. The remote tenant administrators need to review and approve the connection request

## Usage

1. **Coordinate with Remote Tenant**: Ensure the key vault, keys, and permissions are set up in the remote tenant
2. **Set Variables**: Configure all required variables with the remote tenant resource details
3. **Choose Wait Method**: Set `private_endpoint_approval_wait_method` to control how Terraform waits for approval:
   - `"time"`: Terraform will wait for a fixed time period (specify with `private_endpoint_approval_wait_time`)
   - `"none"`: No wait (you'll need to manage approval timing manually)
4. **Deploy**: Run `terraform plan` and `terraform apply`
5. **Approve Connection**: Watch for private endpoint approval timing (see "Private Endpoint Approval" section below)

## Example terraform.tfvars

```hcl
prefix = "mycompany-netapp"
location = "West Europe"

# Remote tenant key vault details
cross_tenant_key_vault_name = "centralized-encryption-kv"
cross_tenant_resource_group_name = "security-rg"
cross_tenant_key_name = "netapp-encryption-key"

# Multi-tenant app client ID
federated_client_id = "12345678-1234-1234-1234-123456789012"

# Existing managed identity ID (in your tenant)
user_assigned_identity_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/anf-ct-cmk-cross-tenant-identity"

# Cross-tenant key vault resource ID (recommended for proper validation)
cross_tenant_key_vault_resource_id = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/security-rg/providers/Microsoft.KeyVault/vaults/centralized-encryption-kv"

# Require manual approval for cross-tenant private endpoint
private_endpoint_manual_approval = true

# Wait method configuration
private_endpoint_approval_wait_method = "time"  # or "none"
private_endpoint_approval_wait_time   = 15      # minutes (only used with "time" method)
```

## Private Endpoint Approval

### When Using "time" Wait Method

When `private_endpoint_approval_wait_method = "time"` is configured, Terraform will wait for the specified time period before continuing. During this wait, you'll see messages like:

```
time_sleep.private_endpoint_approval_time_wait[0]: Still creating... [2m30s elapsed]
time_sleep.private_endpoint_approval_time_wait[0]: Still creating... [2m40s elapsed]
```

**Action Required:** When you see these "Still creating..." messages, coordinate with the remote tenant administrators to:

1. Go to the Azure Portal in the **remote tenant**
2. Navigate to the Key Vault: `<your-cross-tenant-key-vault-name>`
3. Go to **Networking** → **Private endpoint connections**
4. Find the pending connection request
5. **Approve** the private endpoint connection

The approval should happen as soon as possible after seeing the "Still creating..." messages. The wait time provides a buffer for this manual coordination process.

### When Using "none" Wait Method

With `private_endpoint_approval_wait_method = "none"`, Terraform will not wait and will immediately proceed. You must ensure the private endpoint is approved beforehand or the NetApp volume creation may fail.

### Identifying the Connection Request

In the remote tenant's Key Vault, look for a private endpoint connection request with:
- **Source subscription**: Your local subscription ID
- **Resource group**: Your local resource group name
- **Connection name**: Will contain your prefix (e.g., `anf-ct-cmk-pe-ct-akv`)

## Important Notes

- The key vault and keys must already exist in the remote tenant
- You need read permissions on the remote key vault to use data sources
- Private endpoint connections may require approval from remote tenant administrators
- Ensure proper DNS resolution for the remote key vault private endpoint
- The `federated_client_id` must be registered and approved in the remote tenant
- The `cross_tenant_key_vault_resource_id` parameter is recommended for cross-tenant scenarios to ensure proper validation by Azure APIs

## Security Considerations

- Use managed identities where possible
- Limit key vault permissions to minimum required (Get, Encrypt, Decrypt, WrapKey, UnwrapKey)
- Monitor cross-tenant access and audit logs
- Implement proper network security controls for private endpoints

## Troubleshooting

### Private Endpoint Creation Fails with "InternalServerError"

If you encounter an error like:
```
Error: creating Private Endpoint: polling after CreateOrUpdate: polling failed: the Azure API returned the following error:
Status: "InternalServerError"
Code: ""
Message: "An error occurred."
```

This typically indicates one of the following issues:

1. **Target Key Vault Doesn't Exist**: Verify that the key vault exists in the specified subscription and resource group in the remote tenant.

2. **Network Policies**: The target key vault may have network policies that prevent private endpoint creation from external tenants.

3. **Insufficient Permissions**: Your service principal may not have the required permissions to create a private endpoint to the target resource.

4. **Key Vault Configuration**: The target key vault may not be configured to accept private endpoints or may have specific firewall rules.

**Resolution Steps:**
1. Verify the key vault exists: Check with remote tenant administrators
2. Confirm key vault network settings allow private endpoints
3. Ensure the key vault firewall is configured to allow the connection
4. Try setting `is_manual_connection = true` if not already set
5. Consider creating the private endpoint manually first to test connectivity

### DNS Resolution Issues

If the private endpoint is created but DNS resolution fails:
1. Verify the private DNS zone is correctly configured
2. Check that the virtual network link is properly established
3. Ensure DNS settings point to the private IP address

### Key Access Issues

If NetApp account encryption fails to use the key:
1. Verify the `federated_client_id` has been granted access in the remote tenant
2. Check that the managed identity has the required permissions
3. Confirm the key vault access policies are correctly configured
4. Verify that the `cross_tenant_key_vault_resource_id` parameter is set correctly
