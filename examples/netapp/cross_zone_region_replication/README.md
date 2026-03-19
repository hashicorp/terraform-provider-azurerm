## Example: NetApp Volume Cross-Zone-Region Replication

This example demonstrates how to configure cross-zone-region replication for Azure NetApp Files. The configuration creates a primary volume with an availability zone and sets up two destination volumes:
1. A cross-zone replication destination in a different availability zone within the same region
2. A cross-region replication destination in a different region

### Variables

* `prefix` - (Required) The prefix used for all resources in this example.
* `location` - (Required) The Azure Region for the primary volume and cross-zone destination (e.g., `East US`).
* `alt_location` - (Required) The Azure Region for the cross-region destination volume (e.g., `West US`).

### Notes

* The primary volume must have an availability zone assigned if CZR is being used.
* Cross-zone-region replication supports up to two destination replications per source volume.
* When deploying multiple destination volumes, ensure they are created sequentially using `depends_on` to avoid conflicts when authorizing replication on the source volume. The `example_cross_region` volume depends on both `example_primary` and `example_cross_zone` to ensure proper ordering.
* For more information, see [Manage cross-zone-region replication for Azure NetApp Files](https://learn.microsoft.com/azure/azure-netapp-files/cross-zone-region-replication-configure).
