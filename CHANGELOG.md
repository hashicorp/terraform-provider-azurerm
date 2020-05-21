## 2.11.0 (Unreleased)

DEPENDENCIES: 

* updating `github.com/Azure/azure-sdk-for-go` to `v42.1.0` [GH-6725]
* updating `network` to `2020-03-01` [GH-6727]

FEATURES:
* **Opt-In/Experimental Enhanced Validation for Locations:** This allows validating that the `location` field being specified is a valid Azure Region within the Azure Environment being used - which can be caught via `terraform plan` rather than `terraform apply`. This can be enabled by setting the Environment Variable `ARM_PROVIDER_ENHANCED_VALIDATION` to `true` and will be enabled by default in a future release of the AzureRM Provider [GH-6927]

IMPROVEMENTS:

* Data Source: `azurerm_platform_image` - support for `version` filter [GH-6948]
* `azurerm_api_management_api_version_set` - updating the validation for the `name` field [GH-6947]
* `azurerm_app_service` - the `ip_restriction` block now supports the `action` property [GH-6967]
* `azurerm_databricks_workspace` - exposing `workspace_id` and `workspace_url` [GH-6973]
* `azurerm_netapp_volume` - support the `mount_ip_addresses` property [GH-5526]
* `azurerm_redis_cache` - support new maxmemory policies `allkeys-lfu` & `volatile-lfu` [GH-7031]
* `azurerm_storage_account` - allowing the value `PATCH` for `allowed_methods` within the `cors_rule` block within the `blob_properties` block [GH-6964]

BUG FIXES:

* Data Source: `azurerm_api_management_group` - raising an error when the Group cannot be found [GH-7024]
* Data Source: `azurerm_image` - raising an error when the Image cannot be found [GH-7024]
* Data Source: `azurerm_data_lake_store` - raising an error when Data Lake Store cannot be found [GH-7024]
* Data Source: `azurerm_data_share_account` - raising an error when Data Share Account cannot be found [GH-7024]
* Data Source: `azurerm_hdinsight_cluster` - raising an error when the HDInsight Cluster cannot be found [GH-7024]
* Data Source: `azurerm_healthcare_service` - raising an error when the HealthCare Service cannot be found [GH-7024]
* Data Source: `azurerm_healthcare_service` - ensuring all blocks are set in the response [GH-7024]
* Data Source: `azurerm_firewall` - raising an error when the Firewall cannot be found [GH-7024]
* Data Source: `azurerm_maintenance_configuration` - raising an error when the Maintenance Configuration cannot be found [GH-7024]
* Data Source: `azurerm_private_endpoint_connection` - raising an error when the Private Endpoint Connection cannot be found [GH-7024]
* Data Source: `azurerm_shared_image_version` - raising an error when the Image Version cannot be found [GH-7024]
* Data Source: `azurerm_shared_image_versions` - raising an error when Image Versions cannot be found [GH-7024]
* Data Source: `azurerm_user_assigned_identity` - raising an error when the User Assigned Identity cannot be found [GH-7024]
* `azurerm_api_management_subscription` - fix the export of `primary_key` and `secondary_key` [GH-6938]
* `azurerm_eventgrid_event_subscription` - correctly parsing the ID [GH-6958]
* `azurerm_healthcare_service` - ensuring all blocks are set in the response [GH-7024]
* `azurerm_linux_virtual_machine` - allowing name to end with a capital letter [GH-7023]
* `azurerm_linux_virtual_machine_scale_set` - allowing name to end with a capital [GH-7023]
* `azurerm_management_group` - workaround for 403 bug in service response [GH-6668]
* `azurerm_windows_virtual_machine` - allowing name to end with a capital [GH-7023]
* `azurerm_windows_virtual_machine_scale_set` - allowing name to end with a capital [GH-7023]

---

For information on changes between the v2.10.0 and v2.0.0 releases, please see [the v2.10.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/v2.10.0/CHANGELOG.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.44.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/v1.44.0/CHANGELOG.md).
