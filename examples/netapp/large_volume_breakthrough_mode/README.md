## Example: NetApp Files Large Volume with Breakthrough Mode

This example provisions a NetApp Files large volume running in Breakthrough Mode, which places the volume on dedicated capacity delivering higher throughput and allowing the volume to grow up to 2,400 TiB.

Breakthrough Mode requires the `ANFBreakthroughMode` and `ANFLargeVolumes` features to be registered on the subscription, can only be used together with `large_volume_enabled`, cannot be combined with `cool_access` at creation time, and cannot be changed after the volume is created.

### Variables

* `prefix` - (Required) The prefix used for all resources in this example.

* `location` - (Required) The Azure Region in which the resources in this example should be created.

* `pool_size_in_tb` - (Optional) The size of the capacity pool in TiB. Defaults to `4`.

* `storage_quota_in_gb` - (Optional) The size of the Breakthrough Mode volume in GiB, valid values are between `2400` and `2457600`. Defaults to `2400`.

* `throughput_in_mibps` - (Optional) The throughput assigned to the volume in MiB/s. Defaults to `38.4`.
