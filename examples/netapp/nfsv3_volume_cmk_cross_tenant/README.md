## Example: NetApp Account Encryption with Cross-Tenant Customer-Managed Keys

This example demonstrates how to configure NetApp Account Encryption using a customer-managed key from a key vault that exists in a different Azure Active Directory tenant (cross-tenant scenario).

### Overview

In enterprise environments, it's common to have a centralized key management setup where encryption keys are stored in a key vault that belongs to a different Entra ID tenant than where the NetApp resources are deployed. The `federated_client_id` parameter enables this cross-tenant access by specifying the client ID of a multi-tenant Entra ID application that has been granted access to the remote key vault.

### Example Scenario

This example assumes:
- **Tenant A** (Current): Where you're deploying NetApp resources
- **Tenant B** (Remote): Where the key vault and encryption keys already exist
- The key vault in Tenant B has been configured to allow access from your multi-tenant application

### Prerequisites

#### 1. Cross-tenant Key Vault Setup (Done by Tenant B administrators)
The remote tenant administrators must:
- Create/configure the key vault with appropriate encryption keys

#### 2. Managed Identity, Multi-tenant Entra ID Application and permissions
Please see this document [Configure cross-tenant customer-managed keys for Azure NetApp Files volume encryption](https://learn.microsoft.com/en-us/azure/azure-netapp-files/customer-managed-keys-cross-tenant), but only the parts described below since some of the items described on it are already being deployed with this example configuration file:

- [Understand cross-tenant customer-managed keys](https://learn.microsoft.com/en-us/azure/azure-netapp-files/customer-managed-keys-cross-tenant) for overview.
- [Configure cross-tenant customer-managed keys for Azure NetApp Files](https://learn.microsoft.com/en-us/azure/azure-netapp-files/customer-managed-keys-cross-tenant#understand-cross-tenant-customer-managed-keys), up to step 3, step 4 is deployed through this Terraform config example.
- [Authorize access to the key vault](https://learn.microsoft.com/en-us/azure/azure-netapp-files/customer-managed-keys-cross-tenant#authorize-access-to-the-key-vault), to allow access to the remote tenant's Azure Key Vault. **Note** that step 3 of this section is needed only during the `terraform apply` phase when it start to wait for the user to authorize the private endpoint connection in the remote tenant Azure Key Vault.

### Required Variables

- `prefix`: The prefix used for all resources in this example
- `location`: The Azure region where all resources should be created
- `remote_subscription_id`: The subscription ID of the remote tenant where the key vault exists
- `cross_tenant_key_vault_name`: Name of the existing key vault in the remote tenant
- `cross_tenant_resource_group_name`: Resource group where the remote key vault exists
- `cross_tenant_key_name`: Name of the encryption key in the remote key vault
- `cross_tenant_key_vault_resource_id`: Full resource ID of the cross-tenant key vault (mandatory for proper validation)
- `federated_client_id`: Client ID of the multi-tenant Entra ID application
- `user_assigned_identity_id`: Resource ID of the pre-created user-assigned managed identity used by the NetApp account

### Optional Variables

These variables have default values but can be customized:

- `private_endpoint_manual_approval`: Whether the private endpoint connection requires manual approval (default: `true`)
- `private_endpoint_approval_wait_time`: Time to wait (in minutes) for private endpoint approval (default: `15`)

### Usage

1. **Coordinate with Remote Tenant**: Ensure the key vault, keys, and permissions are set up in the remote tenant
2. **Set Variables**: Configure all required variables with the remote tenant resource details
3. **Wait for Approval**: Terraform will always wait for a fixed time period (specified with `private_endpoint_approval_wait_time`) before proceeding, to allow for manual approval of the private endpoint connection.
4. **Deploy**: Run `terraform plan` and `terraform apply`
5. **Approve Connection**: Watch for private endpoint approval timing (see "Private Endpoint Approval" section below)

### Example terraform.tfvars

```hcl
prefix = "cmk-crosstenant"
location = "westeurope"

# Remote tenant details  
remote_subscription_id = "fb6dc404-2676-49d7-be5b-9509a7f6b09b"

# Remote tenant key vault details
cross_tenant_key_vault_name = "centralized-encryption-kv"
cross_tenant_resource_group_name = "security-rg"
cross_tenant_key_name = "netapp-encryption-key"

# Multi-tenant app client ID
federated_client_id = "12345678-1234-1234-1234-123456789012"

# Existing managed identity ID (in your tenant)
user_assigned_identity_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/anf-ct-cmk-cross-tenant-identity"

# Cross-tenant key vault resource ID (mandatory for proper validation)
cross_tenant_key_vault_resource_id = "/subscriptions/11111111-1111-1111-1111-111111111111/resourceGroups/security-rg/providers/Microsoft.KeyVault/vaults/centralized-encryption-kv"

# Require manual approval for cross-tenant private endpoint
private_endpoint_manual_approval = true

# Wait time configuration
private_endpoint_approval_wait_time = 15
```

### Private Endpoint Approval

Terraform will wait for the specified time period (using `private_endpoint_approval_wait_time`) before continuing. During this wait, you'll see messages like:

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

#### Wait Method

The configuration uses a time-based wait. Terraform will wait for the specified number of minutes for manual approval of the private endpoint connection. Skipping the wait is not supported, as it may cause the NetApp volume creation to fail if the private endpoint is not approved in time.
