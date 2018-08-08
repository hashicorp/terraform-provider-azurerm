# Basic VM Scale Set (Linux)

This example provisions a VM Scale Set behind a Load Balancer.

-> **NOTE:** The numbers prefixing the files in this directory are for readability purposes and have no impact on the ordering.

## Variables

* `prefix` - (Required) The Prefix used for all resources in this example.
* `location` - (Required) The Azure Region in which the resources in this example should exist.
* `admin_username` - (Required) The Admin Username used for all VM's in this Scale Set.
* `admin_password` - (Required) The Admin Password used for all VM's in this Scale Set.
* `instance_count` - (Optional) Number of VM instances (100 or less). Defaults to 3.
