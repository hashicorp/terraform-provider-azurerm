## Example: NetApp Volume Protocol Conversion between NFSv3 and NFSv4.1

This example demonstrates how to convert a NetApp volume between NFSv3 and NFSv4.1 protocols and vice-versa.

### Variables

* `prefix` - (Optional) The prefix used for all resources in this example. Defaults to "example".

* `location` - (Optional) The Azure Region in which the resources in this example should be created. Defaults to "westus3".

* `resource_group_name` - (Optional) The name of the resource group. Defaults to "example-netapp-protocol-conversion-rg".

* `protocol_type` - (Optional) The NFS protocol type (NFSv3 or NFSv4.1). Defaults to "NFSv3".

### Usage example to demonstrate the conversion

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
