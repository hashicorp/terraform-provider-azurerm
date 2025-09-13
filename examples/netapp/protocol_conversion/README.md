## Example: NetApp Volume Protocol Conversion between NFSv3 and NFSv4.1 (Provider v4)

This example demonstrates how to convert a NetApp volume between NFSv3 and NFSv4.1 protocols using Azure NetApp Files with **Terraform Provider AzureRM v4.x**.

> **Note**: This example uses provider v4 field naming conventions (`protocols_enabled`). For provider v5, use the `protocol` field instead.

### Prerequisites

Before running this example, you need to register the protocol conversion feature:

```bash
az feature register --namespace Microsoft.NetApp --name ANFProtocolTypeNFSConversion
az feature show --namespace Microsoft.NetApp --name ANFProtocolTypeNFSConversion
```

Wait until the registration state is "Registered" before proceeding.

### Variables

* `prefix` - (Optional) The prefix used for all resources in this example. Defaults to "example".

* `location` - (Optional) The Azure Region in which the resources in this example should be created. Defaults to "East US".

* `resource_group_name` - (Optional) The name of the resource group. Defaults to "rg-netapp-protocol-conversion-1".

* `protocol_type` - (Optional) The NFS protocol type (NFSv3 or NFSv4.1). Defaults to "NFSv3".

### Usage

1. **Initial setup** - Create a volume with NFSv3:
   ```bash
   terraform apply -var="protocol_type=NFSv3"
   ```

2. **Convert to NFSv4.1**:
   ```bash
   terraform apply -var="protocol_type=NFSv4.1"
   ```

3. **Convert back to NFSv3**:
   ```bash
   terraform apply -var="protocol_type=NFSv3"
   ```

### Important Considerations

- You cannot convert an NFSv4.1 volume with Kerberos enabled to NFSv3
- You cannot change the NFS version of a dual-protocol volume
- Converting from NFSv4.1 to NFSv3 will cause all advanced NFSv4.1 features such as ACLs and file locking to become unavailable
- The operation requires unmounting the volume from all clients before conversion
- The validation ensures that export policy protocols match the volume protocol during conversion

### Provider Version Compatibility

This example is designed for **Terraform Provider AzureRM v4.x**:

- **Export Policy Field**: Uses `protocols_enabled` (v4 convention with deprecation warning)
- **Provider Version**: `~>4.0`
- **Migration Path**: When upgrading to provider v5.x, change `protocols_enabled` to `protocol` in export policy rules