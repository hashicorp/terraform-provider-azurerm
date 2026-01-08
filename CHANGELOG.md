## 4.57.0 (December 18, 2025)

**NOTE:** This release removes the Mobile Network (`azurerm_mobile_network*`) resources and data sources due to Azure having retired the service

FEATURES:

* **New Resource:** `azurerm_automation_runtime_environment` ([#30991](https://github.com/hashicorp/terraform-provider-azurerm/issues/30991))

ENHANCEMENTS:

* `azurerm_data_protection_backup_vault_customer_managed_key` - the `key_vault_key_id` property now supports keys from a Managed HSM vault ([#31365](https://github.com/hashicorp/terraform-provider-azurerm/issues/31365))
* `azurerm_kubernetes_cluster` - support for the `node_provisioning_profile` block ([#30517](https://github.com/hashicorp/terraform-provider-azurerm/issues/30517))
* `azurerm_log_analytics_cluster_customer_managed_key` - the `key_vault_key_id` property now supports keys from a Managed HSM vault ([#31375](https://github.com/hashicorp/terraform-provider-azurerm/issues/31375))
* `azurerm_mssql_database` - the `transparent_data_encryption_key_vault_key_id` property now supports keys from a Managed HSM vault ([#31373](https://github.com/hashicorp/terraform-provider-azurerm/issues/31373))

BUG FIXES:

* `azurerm_data_factory` - fix ID parsing errors when `customer_managed_key_identity_id` is an empty string ([#28621](https://github.com/hashicorp/terraform-provider-azurerm/issues/28621))
* `azurerm_eventhub` - `partition_count` can now be updated for dedicated clusters ([#30993](https://github.com/hashicorp/terraform-provider-azurerm/issues/30993))
* `azurerm_linux_function_app` - fix panic when deployed without all required permissions ([#31344](https://github.com/hashicorp/terraform-provider-azurerm/issues/31344))

## 4.56.0 (December 11, 2025)

ENHANCEMENTS:

* dependencies: `healthbot` - update to API version `2025-05-25` ([#31328](https://github.com/hashicorp/terraform-provider-azurerm/issues/31328))
* dependencies: `terraform-plugin-testing` - update to `v1.14.0`  ([#31334](https://github.com/hashicorp/terraform-provider-azurerm/issues/31334))
* Data Source: `azurerm_cognitive_account` - add support for new attributes ([#30778](https://github.com/hashicorp/terraform-provider-azurerm/issues/30778))
* `azurerm_cognitive_account` - add support for the `kind` property to rollback or upgrade from `OpenAI` to `AIServices` ([#31063](https://github.com/hashicorp/terraform-provider-azurerm/issues/31063))
* `azurerm_databricks_workspace_root_dbfs_customer_managed_key` - the `key_vault_key_id` property now supports keys from Managed HSM Vaults ([#31336](https://github.com/hashicorp/terraform-provider-azurerm/issues/31336))
* `azurerm_databricks_workspace_root_dbfs_customer_managed_key` - the `key_vault_key_id` property now supports versionless keys ([#31336](https://github.com/hashicorp/terraform-provider-azurerm/issues/31336))
* `azurerm_healthbot` - add support for the `C1` and `PES` SKUs ([#31328](https://github.com/hashicorp/terraform-provider-azurerm/issues/31328))
* `azurerm_lb` fix `ignore_changes` behaviour in updatable properties ([#31318](https://github.com/hashicorp/terraform-provider-azurerm/issues/31318))
* `azurerm_network_manager_network_group` - add support for the `member_type` property [GH-30672
* `azurerm_network_manager_static_member` - add support for using a subnet as the target resource ([#30672](https://github.com/hashicorp/terraform-provider-azurerm/issues/30672))
* `azurerm_virtual_network_gateway` - add support for the `ErGwScale` SKU ([#31082](https://github.com/hashicorp/terraform-provider-azurerm/issues/31082))

BUG FIXES:

* `azurerm_container_app_environment_certificate` - fix an issue that prevented creating the resource with an empty value for `certificate_password` ([#31335](https://github.com/hashicorp/terraform-provider-azurerm/issues/31335))
* `azurerm_databricks_workspace_root_dbfs_customer_managed_key` - fix a panic that occurred when the customer managed key was removed from the workspace outside of Terraform ([#31336](https://github.com/hashicorp/terraform-provider-azurerm/issues/31336))
* `azurerm_databricks_workspace_root_dbfs_customer_managed_key` - fix the timeout for the delete operation ([#31336](https://github.com/hashicorp/terraform-provider-azurerm/issues/31336))
* `azurerm_storage_blob_inventory_policy` - fix setting Resource Identity data ([#31313](https://github.com/hashicorp/terraform-provider-azurerm/issues/31313))

## 4.55.0 (December 04, 2025)

FEATURES:

* **New Data Source**: `azurerm_api_management_workspace` ([#30241](https://github.com/hashicorp/terraform-provider-azurerm/issues/30241))
* **New Resource**: `azurerm_cognitive_account_project` ([#30916](https://github.com/hashicorp/terraform-provider-azurerm/issues/30916))
* **New Resource**: `azurerm_log_analytics_workspace_table_custom_log` ([#30800](https://github.com/hashicorp/terraform-provider-azurerm/issues/30800))
* **New Resource**: `azurerm_mongo_cluster_user` ([#31205](https://github.com/hashicorp/terraform-provider-azurerm/issues/31205))
* **New Resource**: `azurerm_palo_alto_next_generation_firewall_virtual_hub_strata_cloud_manager` ([#30613](https://github.com/hashicorp/terraform-provider-azurerm/issues/30613))
* **New Resource**: `azurerm_palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager` ([#30613](https://github.com/hashicorp/terraform-provider-azurerm/issues/30613))
* **New List Resource**: `azurerm_private_dns_zone` ([#31157](https://github.com/hashicorp/terraform-provider-azurerm/issues/31157))

ENHANCEMENTS:

* dependencies: `containerregistry` - update to API version `2025-04-01` ([#30205](https://github.com/hashicorp/terraform-provider-azurerm/issues/30205))
* dependencies: `go-azure-helpers` - update to `v0.75.1` ([#31148](https://github.com/hashicorp/terraform-provider-azurerm/issues/31148))
* dependencies: `go-azure-sdk` - update to `v0.20251202.1181053` ([#31253](https://github.com/hashicorp/terraform-provider-azurerm/issues/31253))
* dependencies: `managedidentity` - upgrade API version to `2024-11-30` ([#30535](https://github.com/hashicorp/terraform-provider-azurerm/issues/30535))
* dependencies: `postgres` - update to API version `2025-08-01` ([#31162](https://github.com/hashicorp/terraform-provider-azurerm/issues/31162))
* `azurerm_cognitive_account` - update validation for `customer_managed_key.key_vault_key_id` to allow managed HSM keys as input ([#31147](https://github.com/hashicorp/terraform-provider-azurerm/issues/31147))
* `azurerm_container_app_environment` - extend validation for `workload_profile_type` for additional supported SKUs ([#30738](https://github.com/hashicorp/terraform-provider-azurerm/issues/30738))
* `azurerm_container_app_environment_certificate` - add support for the `certificate_key_vault` block ([#30510](https://github.com/hashicorp/terraform-provider-azurerm/issues/30510))
* `azurerm_data_factory` - update validation for `customer_managed_key_id` to allow managed HSM keys as input ([#31146](https://github.com/hashicorp/terraform-provider-azurerm/issues/31146))
* `azurerm_mongo_cluster` - support for new properties `customer_managed_key`, `data_api_mode_enabled`, `identity`, `restore`, `authentication_methods` and `storage_type` ([#31100](https://github.com/hashicorp/terraform-provider-azurerm/issues/31100))
* `azurerm_mysql_flexible_server` - add support for MySQL version `8.4` ([#31099](https://github.com/hashicorp/terraform-provider-azurerm/issues/31099))
* `azurerm_oracle_autonomous_database` - the `admin_password` property is no longer `ForceNew` ([#30966](https://github.com/hashicorp/terraform-provider-azurerm/issues/30966))
* `azurerm_postgresql_flexible_server` - update validation for `customer_managed_key.key_vault_key_id` and `customer_managed_key.geo_backup_key_vault_key_id` to allow managed HSM keys as input ([#31148](https://github.com/hashicorp/terraform-provider-azurerm/issues/31148))
* `azurerm_postgresql_flexible_server` - add support for PostgreSQL version `18` ([#31162](https://github.com/hashicorp/terraform-provider-azurerm/issues/31162))
* `azurerm_storage_encryption_scope` - update validation for `key_vault_key_id` to allow managed HSM keys as input ([#31145](https://github.com/hashicorp/terraform-provider-azurerm/issues/31145))

BUG FIXES:

* Data Source: `azurerm_ssh_public_key` - fix normalisation for `public_key` to avoid removing a literal `EOT` from the base64 encoded content ([#31249](https://github.com/hashicorp/terraform-provider-azurerm/issues/31249))
* `azurerm_data_protection_backup_vault` - poll delete request for completion ([#31202](https://github.com/hashicorp/terraform-provider-azurerm/issues/31202))
* `azurerm_function_app_hybrid_connection` - remove validation preventing resource import when using an elastic service plan SKU ([#31134](https://github.com/hashicorp/terraform-provider-azurerm/issues/31134))
* `azurerm_key_vault_key` - `not_before_date` and `expiration_date` are now set into state when empty, fixing an issue where drift was not detected ([#31192](https://github.com/hashicorp/terraform-provider-azurerm/issues/31192))
* `azurerm_key_vault_secret` - `not_before_date` and `expiration_date` are now set into state when empty, fixing an issue where drift was not detected ([#31192](https://github.com/hashicorp/terraform-provider-azurerm/issues/31192))
* `azurerm_kubernetes_cluster` - fix drift on `azure_policy_enabled` when updating cluster ([#30917](https://github.com/hashicorp/terraform-provider-azurerm/issues/30917))
* `azurerm_kubernetes_fleet_update_run` - fix a nil pointer dereference to prevent panics ([#31213](https://github.com/hashicorp/terraform-provider-azurerm/issues/31213))
* `azurerm_lb_nat_rule` - fix an issue that prevented changing `floating_ip_enabled` and `tcp_reset_enabled` from `true` to `false` ([#31244](https://github.com/hashicorp/terraform-provider-azurerm/issues/31244))
* `azurerm_lb_outbound_rule` - fix an issue that prevented changing `tcp_reset_enabled` from `true` to `false` ([#31244](https://github.com/hashicorp/terraform-provider-azurerm/issues/31244))
* `azurerm_lb_rule` - fix an issue that prevented changing `floating_ip_enabled` and `tcp_reset_enabled` from `true` to `false` ([#31244](https://github.com/hashicorp/terraform-provider-azurerm/issues/31244))
* `azurerm_private_endpoint` - ensure Resource Identity data is set on create to avoid `Missing Resource Identity After Create` errors ([#31246](https://github.com/hashicorp/terraform-provider-azurerm/issues/31246))
* `azurerm_resource_group` - fix poller for the `prevent_deletion_if_contains_resources` feature, resolving an Azure eventual consistency issue ([#31253](https://github.com/hashicorp/terraform-provider-azurerm/issues/31253))
* `azurerm_storage_account` - ensure Resource Identity data is set on create to avoid `Missing Resource Identity After Create` errors ([#31246](https://github.com/hashicorp/terraform-provider-azurerm/issues/31246))
* `azurerm_traffic_manager_profile` - fix an issue that prevented changing `traffic_view_enabled` from `true` to `false` ([#31066](https://github.com/hashicorp/terraform-provider-azurerm/issues/31066))

## 4.54.0 (November 19, 2025)

FEATURES:

* **New Action**: `azurerm_cdn_front_door_cache_purge`  ([#30765](https://github.com/hashicorp/terraform-provider-azurerm/issues/30765))
* **New Action**: `azurerm_data_protection_backup_instance_protect` ([#31085](https://github.com/hashicorp/terraform-provider-azurerm/issues/31085))
* **New Action**: `azurerm_managed_redis_databases_flush` ([#31132](https://github.com/hashicorp/terraform-provider-azurerm/issues/31132))
* **New Action**: `azurerm_mssql_execute_job` ([#31095](https://github.com/hashicorp/terraform-provider-azurerm/issues/31095))
* **New List Resource**: `azurerm_network_interface` ([#31012](https://github.com/hashicorp/terraform-provider-azurerm/issues/31012))
* **New List Resource**: `azurerm_network_profile` ([#31127](https://github.com/hashicorp/terraform-provider-azurerm/issues/31127))
* **New List Resource**: `azurerm_network_security_group` ([#31014](https://github.com/hashicorp/terraform-provider-azurerm/issues/31014))
* **New List Resource**: `azurerm_route_table` ([#31015](https://github.com/hashicorp/terraform-provider-azurerm/issues/31015))

ENHANCEMENTS:

* dependencies: `go-azure-sdk` - update to `v0.20251107.1191907` ([#31095](https://github.com/hashicorp/terraform-provider-azurerm/issues/31095))
* Data Source: `azurerm_container_app` - add support for the `template.cooldown_period_in_seconds` and `template.polling_interval_in_seconds` properties ([#29426](https://github.com/hashicorp/terraform-provider-azurerm/issues/29426))
* `azurerm_container_app` - add support for the `template.cooldown_period_in_seconds` and `template.polling_interval_in_seconds` properties ([#29426](https://github.com/hashicorp/terraform-provider-azurerm/issues/29426))
* `azurerm_linux_function_app` - add support for `dotnet_version` `10.0` ([#31007](https://github.com/hashicorp/terraform-provider-azurerm/issues/31007))
* `azurerm_linux_function_app_slot` - add support for `dotnet_version` `10.0` ([#31007](https://github.com/hashicorp/terraform-provider-azurerm/issues/31007))
* `azurerm_linux_web_app` - add support for `dotnet_version` `10.0` ([#31007](https://github.com/hashicorp/terraform-provider-azurerm/issues/31007))
* `azurerm_linux_web_app_slot` - add support for `dotnet_version` `10.0` ([#31007](https://github.com/hashicorp/terraform-provider-azurerm/issues/31007))
* `azurerm_managed_redis` - add support for `persistence_append_only_file_backup_frequency` and `persistence_redis_database_backup_frequency` properties  ([#30964](https://github.com/hashicorp/terraform-provider-azurerm/issues/30964))
* `azurerm_resource_group` - refactored from legacy SDK to use `go-azure-sdk` ([#30616](https://github.com/hashicorp/terraform-provider-azurerm/issues/30616))
* `azurerm_service_plan` - suppress casing difference on `sku_name` ([#30907](https://github.com/hashicorp/terraform-provider-azurerm/issues/30907))
* `azurerm_storage_share_directory` - Deprecate `storage_share_id` in favour of `storage_share_url` ([#28457](https://github.com/hashicorp/terraform-provider-azurerm/issues/28457))
* `azurerm_storage_share_file` - Deprecate `storage_share_id` in favour of `storage_share_url` ([#28457](https://github.com/hashicorp/terraform-provider-azurerm/issues/28457))
* `azurerm_windows_function_app` - add support for `dotnet_version` `v10.0` ([#31007](https://github.com/hashicorp/terraform-provider-azurerm/issues/31007))
* `azurerm_windows_function_app_slot` - add support for `dotnet_version` `v10.0` ([#31007](https://github.com/hashicorp/terraform-provider-azurerm/issues/31007))
* `azurerm_windows_web_app` - add support for `dotnet_version` `v10.0` ([#31007](https://github.com/hashicorp/terraform-provider-azurerm/issues/31007))
* `azurerm_windows_web_app_slot` - add support for `dotnet_version` `v10.0` ([#31007](https://github.com/hashicorp/terraform-provider-azurerm/issues/31007))

BUG FIXES:

* `azurerm_orchestrated_virtual_machine_scale_set` - Fix issue when using a specialized image ([#30889](https://github.com/hashicorp/terraform-provider-azurerm/issues/30889))
* `azurerm_virtual_network` - remove RO values from update to avoid issues with API payload size limitation ([#30945](https://github.com/hashicorp/terraform-provider-azurerm/issues/30945))

## 4.53.0 (November 14, 2025)

FEATURES:

* **New Resource**: `azurerm_api_management_workspace_certificate` ([#30628](https://github.com/hashicorp/terraform-provider-azurerm/issues/30628))
* **New Resource**: `azurerm_mongo_cluster_firewall_rule` ([#31062](https://github.com/hashicorp/terraform-provider-azurerm/issues/31062))

ENHANCEMENTS:

* dependencies: `automation` - update to API version `2024-10-23` ([#30890](https://github.com/hashicorp/terraform-provider-azurerm/issues/30890))
* dependencies: `go-azure-sdk` - update to `v0.20251029.1173336` ([#31051](https://github.com/hashicorp/terraform-provider-azurerm/issues/31051))
* dependencies: `managedredis` - update to API Version `2025-07-01` ([#31051](https://github.com/hashicorp/terraform-provider-azurerm/issues/31051))
* dependencies: `mongocluster` - update to API version `2025-09-01` ([#30982](https://github.com/hashicorp/terraform-provider-azurerm/issues/30982))
* `azurerm_api_management_backend` - add support for the `circuit_breaker_rule` block  ([#30471](https://github.com/hashicorp/terraform-provider-azurerm/issues/30471))
* `azurerm_dynatrace_monitor` - support for the `YEARLY` value in the `billing_cycle` property ([#31078](https://github.com/hashicorp/terraform-provider-azurerm/issues/31078))
* `azurerm_kubernetes_cluster_node_pool` - support for the `undrainable_node_behavior` and `max_unavailable` properties ([#30563](https://github.com/hashicorp/terraform-provider-azurerm/issues/30563))
* `azurerm_managed_disk` - support expanding Ultra Disks and Premium SSD v2 disk without downtime ([#30593](https://github.com/hashicorp/terraform-provider-azurerm/issues/30593))
* `azurerm_managed_redis` - add support for `public_network_access` ([#31051](https://github.com/hashicorp/terraform-provider-azurerm/issues/31051))
* `azurerm_storage_table_entity` - resource is now removed from state if it no longer exists in Azure ([#31064](https://github.com/hashicorp/terraform-provider-azurerm/issues/31064))
* `azurerm_synapse_spark_pool` - add support for `spark_version` `3.5` ([#30900](https://github.com/hashicorp/terraform-provider-azurerm/issues/30900))
* `data.azurerm_postgresql_flexible_server` - add support for `zone` and `high_availability` ([#31034](https://github.com/hashicorp/terraform-provider-azurerm/issues/31034))

BUG FIXES:

* `azurerm_dynatrace_monitor` -  the `phone_number` and `country` properties are no longer Required ([#31077](https://github.com/hashicorp/terraform-provider-azurerm/issues/31077))
* `azurerm_dynatrace_tag_rules` - the `log_rule.filtering_tag` property is no longer required ([#31065](https://github.com/hashicorp/terraform-provider-azurerm/issues/31065))
* `azurerm_dynatrace_tag_rules` - the `metric_rule.filtering_tag` property is no longer required ([#31065](https://github.com/hashicorp/terraform-provider-azurerm/issues/31065))
* `azurerm_kubernetes_cluster` - fix crash in use of `azure_active_directory_role_based_access_control` ([#31101](https://github.com/hashicorp/terraform-provider-azurerm/issues/31101))
* `azurerm_logic_app_workflow` - fix inaccurate error messages ([#30963](https://github.com/hashicorp/terraform-provider-azurerm/issues/30963))
* `azurerm_virtual_network_gateway` - fix validation for `policy_group.name` and `vpn_client_configuration.virtual_network_gateway_client_connection.policy_group_names` ([#30454](https://github.com/hashicorp/terraform-provider-azurerm/issues/30454))

## 4.52.0 (November 06, 2025)

**NOTE:** This release removes the `azurerm_spatial_anchors_account` resource and data source due to Azure having retired the service

FEATURES:

* **New Resource**: `azurerm_api_management_workspace_api_version_set` ([#30498](https://github.com/hashicorp/terraform-provider-azurerm/issues/30498))

ENHANCEMENTS:

* dependencies: `Go` updated to `v1.25.3` ([#31020](https://github.com/hashicorp/terraform-provider-azurerm/issues/31020))
* Data Source: `azurerm_application_gateway` - add support for the `backend_http_settings.dedicated_backend_connection_enabled` property ([#31033](https://github.com/hashicorp/terraform-provider-azurerm/issues/31033))
* `azurerm_application_gateway` - add support for the `backend_http_settings.dedicated_backend_connection_enabled` property ([#31033](https://github.com/hashicorp/terraform-provider-azurerm/issues/31033))
* `azurerm_machine_learning_datastore_blobstorage` - improve validation for `storage_container_id` ([#31002](https://github.com/hashicorp/terraform-provider-azurerm/issues/31002))
* `azurerm_machine_learning_datastore_datalake_gen2` - improve validation for `storage_container_id` ([#31002](https://github.com/hashicorp/terraform-provider-azurerm/issues/31002))
* `azurerm_windows_web_app` - add support for the `virtual_network_image_pull_enabled` property ([#30920](https://github.com/hashicorp/terraform-provider-azurerm/issues/30920))
* `azurerm_windows_web_app_slot` - add support for the `virtual_network_image_pull_enabled` property ([#30920](https://github.com/hashicorp/terraform-provider-azurerm/issues/30920))

BUG FIXES:

* `azurerm_container_registry_task` - prevent a panic by adding a nil check ([#31043](https://github.com/hashicorp/terraform-provider-azurerm/issues/31043))

## 4.51.0 (October 30, 2025)

FEATURES:

* **New Data Source**: `azurerm_oracle_resource_anchor` ([#30823](https://github.com/hashicorp/terraform-provider-azurerm/issues/30823))
* **New Resource**: `azurerm_network_manager_routing_rule` ([#30439](https://github.com/hashicorp/terraform-provider-azurerm/issues/30439))
* **New Resource**: `azurerm_oracle_resource_anchor` ([#30823](https://github.com/hashicorp/terraform-provider-azurerm/issues/30823))

ENHANCEMENTS:

* dependencies: `dashboard` - update to API version `2025-08-01` ([#30972](https://github.com/hashicorp/terraform-provider-azurerm/issues/30972))
* dependencies: `go-azure-sdk` - update to `v0.20251024.1223440` ([#30952](https://github.com/hashicorp/terraform-provider-azurerm/issues/30952))
* dependencies: `network` - update to API version `2025-01-01` ([#30904](https://github.com/hashicorp/terraform-provider-azurerm/issues/30904))
* `azurerm_cognitive_account` - add `TextAnalytics` to allowed `kind` validation for `network_acls.bypass` ([#30887](https://github.com/hashicorp/terraform-provider-azurerm/issues/30887))
* `azurerm_subnet_service_endpoint_storage_policy` - add support for the `/services/Azure/Databricks` value in the `definition.service_resources` property ([#30762](https://github.com/hashicorp/terraform-provider-azurerm/issues/30762))

BUG FIXES:

* Data Source: `azurerm_managed_redis` - fix a panic caused by a nested field access on a pointer without nil checking ([#30978](https://github.com/hashicorp/terraform-provider-azurerm/issues/30978))

## 4.50.0 (October 23, 2025)

FEATURES:

* **New Data Source**: `azurerm_managed_redis` ([#30060](https://github.com/hashicorp/terraform-provider-azurerm/issues/30060))
* **New Resource**: `azurerm_managed_redis` ([#30060](https://github.com/hashicorp/terraform-provider-azurerm/issues/30060))
* **New Resource**: `azurerm_managed_redis_geo_replication` ([#30060](https://github.com/hashicorp/terraform-provider-azurerm/issues/30060))

ENHANCEMENTS:

* dependencies: `go-azure-sdk` update to `v0.20251016.1163854` ([#30883](https://github.com/hashicorp/terraform-provider-azurerm/issues/30883))
* dependencies: `oracle` - update to API version `2025-09-01` ([#30796](https://github.com/hashicorp/terraform-provider-azurerm/issues/30796))
* Data Source: `azurerm_container_app_environment` - add support for the `public_network_access` property ([#30817](https://github.com/hashicorp/terraform-provider-azurerm/issues/30817))
* `azurerm_container_app_environment` - add support for the `public_network_access` property ([#30817](https://github.com/hashicorp/terraform-provider-azurerm/issues/30817))
* `azurerm_mssql_job_target_group` - the `job_target.job_credential_id` property is no longer required when `database_name` is not set to allow for authentication using a managed identity ([#30898](https://github.com/hashicorp/terraform-provider-azurerm/issues/30898))
* `azurerm_netapp_volume_resource` - support for Cross Zone Region replication through the `data_protection_replication` block ([#30872](https://github.com/hashicorp/terraform-provider-azurerm/issues/30872))
* `azurerm_search_service` - implement plan time error when `local_authentication_enabled = false` and `authentication_failure_mode` is set ([#30882](https://github.com/hashicorp/terraform-provider-azurerm/issues/30882))

BUG FIXES:

* `azurerm_mssql_database` - allow existing zero or null value for `auto_pause_delay_in_minutes` and `min_capacity` of non-serverless database ([#30924](https://github.com/hashicorp/terraform-provider-azurerm/issues/30924))

## 4.49.0 (October 16, 2025)

FEATURES:

* **New Data Source**: `azurerm_graph_services_account` ([#30697](https://github.com/hashicorp/terraform-provider-azurerm/issues/30697))
* **New Data Source**: `azurerm_oracle_exascale_database_storage_vault` ([#30043](https://github.com/hashicorp/terraform-provider-azurerm/issues/30043))
* **New Resource**: `azurerm_api_management_workspace_policy_fragment` ([#30678](https://github.com/hashicorp/terraform-provider-azurerm/issues/30678))
* **New Resource**: `azurerm_oracle_exascale_database_storage_vault` ([#30043](https://github.com/hashicorp/terraform-provider-azurerm/issues/30043))

ENHANCEMENTS:

* Data Source: `azurerm_data_protection_backup_vault` - add support for the `identity.identity_ids` property ([#29061](https://github.com/hashicorp/terraform-provider-azurerm/issues/29061))
* `azurerm_consumption_budget_management_group` - remove the maximum count validation for the `notification` block ([#29200](https://github.com/hashicorp/terraform-provider-azurerm/issues/29200))
* `azurerm_consumption_budget_resource_group` - remove the maximum count validation for the `notification` block ([#29200](https://github.com/hashicorp/terraform-provider-azurerm/issues/29200))
* `azurerm_consumption_budget_subscription` - remove the maximum count validation for the `notification` block ([#29200](https://github.com/hashicorp/terraform-provider-azurerm/issues/29200))
* `azurerm_data_protection_backup_vault` - add support for the `identity.identity_ids` property ([#29061](https://github.com/hashicorp/terraform-provider-azurerm/issues/29061))
* `azurerm_data_protection_backup_vault` - add support for `UserAssigned` and `SystemAssigned, UserAssigned` values to the `identity.type` property ([#29061](https://github.com/hashicorp/terraform-provider-azurerm/issues/29061))
* `azurerm_monitor_data_collection_rule` - improve validation for `data_sources.*.name` ([#30851](https://github.com/hashicorp/terraform-provider-azurerm/issues/30851))
* `azurerm_search_service` - support upgrading the `sku` based on tier  ([#30842](https://github.com/hashicorp/terraform-provider-azurerm/issues/30842))
* `azurerm_storage_queue` - support migrating from `storage_account_name` to `storage_account_id`  ([#30836](https://github.com/hashicorp/terraform-provider-azurerm/issues/30836))

BUG FIXES:

* `azurerm_application_insights` - fix an issue that caused `tags` to be removed when other properties were updated ([#30758](https://github.com/hashicorp/terraform-provider-azurerm/issues/30758))
* `azurerm_container_registry` - fix the `name` length validation to allow 50 rather than 49 ([#30858](https://github.com/hashicorp/terraform-provider-azurerm/issues/30858))
* `azurerm_function_app_flex_consumption` - the `auth_settings` block contents are now set into state when `auth_settings.enabled` is set to `false` ([#30781](https://github.com/hashicorp/terraform-provider-azurerm/issues/30781))
* `azurerm_linux_function_app` - the `auth_settings` block contents are now set into state when `auth_settings.enabled` is set to `false` ([#30781](https://github.com/hashicorp/terraform-provider-azurerm/issues/30781))
* `azurerm_linux_function_app_slot` - the `auth_settings` block contents are now set into state when `auth_settings.enabled` is set to `false` ([#30781](https://github.com/hashicorp/terraform-provider-azurerm/issues/30781))
* `azurerm_linux_web_app` - the `auth_settings` block contents are now set into state when `auth_settings.enabled` is set to `false` ([#30781](https://github.com/hashicorp/terraform-provider-azurerm/issues/30781))
* `azurerm_linux_web_app_slot` - the `auth_settings` block contents are now set into state when `auth_settings.enabled` is set to `false` ([#30781](https://github.com/hashicorp/terraform-provider-azurerm/issues/30781))
* `azurerm_mssql_database` - fix validation for `min_capacity` and `auto_pause_delay_in_minutes` being set on non-serverless SKUs ([#30856](https://github.com/hashicorp/terraform-provider-azurerm/issues/30856))
* `azurerm_signalr_service_custom_certificate` - remove unnecessary API requests and checks that could lead to a panic ([#30412](https://github.com/hashicorp/terraform-provider-azurerm/issues/30412))
* `azurerm_windows_function_app` - the `auth_settings` block contents are now set into state when `auth_settings.enabled` is set to `false` ([#30781](https://github.com/hashicorp/terraform-provider-azurerm/issues/30781))
* `azurerm_windows_function_app_slot` - the `auth_settings` block contents are now set into state when `auth_settings.enabled` is set to `false` ([#30781](https://github.com/hashicorp/terraform-provider-azurerm/issues/30781))
* `azurerm_windows_web_app` - the `auth_settings` block contents are now set into state when `auth_settings.enabled` is set to `false` ([#30781](https://github.com/hashicorp/terraform-provider-azurerm/issues/30781))
* `azurerm_windows_web_app_slot` - the `auth_settings` block contents are now set into state when `auth_settings.enabled` is set to `false` ([#30781](https://github.com/hashicorp/terraform-provider-azurerm/issues/30781))

## 4.48.0 (October 13, 2025)

FEATURES:

* **New Data Source**: `azurerm_oracle_autonomous_database_clone_from_backup` ([#29633](https://github.com/hashicorp/terraform-provider-azurerm/issues/29633))
* **New Data Source**: `azurerm_oracle_autonomous_database_clone_from_database` ([#29633](https://github.com/hashicorp/terraform-provider-azurerm/issues/29633))
* **New Resource**: `azurerm_oracle_autonomous_database_clone_from_backup` ([#29633](https://github.com/hashicorp/terraform-provider-azurerm/issues/29633))
* **New Resource**: `azurerm_oracle_autonomous_database_clone_from_database` ([#29633](https://github.com/hashicorp/terraform-provider-azurerm/issues/29633))

ENHANCEMENTS:

* dependencies: `containerapps` - update to API version `2025-07-01` ([#30801](https://github.com/hashicorp/terraform-provider-azurerm/issues/30801))
* dependencies: `containerservice` - update to API version `2025-07-01` ([#30719](https://github.com/hashicorp/terraform-provider-azurerm/issues/30719))
* dependencies: `go-azure-sdk` - update to `v0.20251007.1195632` ([#30799](https://github.com/hashicorp/terraform-provider-azurerm/issues/30799))
* dependencies: `guestconfiguration` - update to API version `2024-04-05` ([#30642](https://github.com/hashicorp/terraform-provider-azurerm/issues/30642))
* dependencies: `search` - update to API version `2025-05-01` ([#30314](https://github.com/hashicorp/terraform-provider-azurerm/issues/30314))
* `azurerm_kubernetes_cluster` - add support for `AzureLinux3` and `Ubuntu2204` to the `default_node_pool.os_sku` property ([#30719](https://github.com/hashicorp/terraform-provider-azurerm/issues/30719))
* `azurerm_kubernetes_cluster` - add support for the `ai_toolchain_operator_enabled` property ([#30713](https://github.com/hashicorp/terraform-provider-azurerm/issues/30713))
* `azurerm_kubernetes_cluster_node_pool` - add support for `AzureLinux3` and `Ubuntu2204` to the `os_sku` property ([#30719](https://github.com/hashicorp/terraform-provider-azurerm/issues/30719))
* `azurerm_linux_virtual_machine_scale_set` - add support for the `resilient_vm_creation_enabled` and `resilient_vm_deletion_enabled` properties ([#30204](https://github.com/hashicorp/terraform-provider-azurerm/issues/30204))
* `azurerm_network_watcher_flow_log` - changing the `target_resource_id` property no longer forces the resource to be replaced ([#30776](https://github.com/hashicorp/terraform-provider-azurerm/issues/30776))
* `azurerm_notification_hub_namespace` - add support for `replication_region` and `zone_redundancy_enabled` ([#30531](https://github.com/hashicorp/terraform-provider-azurerm/issues/30531))
* `azurerm_windows_virtual_machine_scale_set` - add support for the `resilient_vm_creation_enabled` and `resilient_vm_deletion_enabled` properties ([#30204](https://github.com/hashicorp/terraform-provider-azurerm/issues/30204))

BUG FIXES:

* `azurerm_eventhub_namespace` - `maximum_throughput_units` can be set to `0` when `auto_inflate_enabled` is disabled ([#30777](https://github.com/hashicorp/terraform-provider-azurerm/issues/30777))
* `azurerm_log_analytics_workspace` - fix the default value for `local_authentication_enabled` ([#30759](https://github.com/hashicorp/terraform-provider-azurerm/issues/30759))
* `azurerm_mssql_database` - add validation to ensure that `min_capacity` and `auto_pause_delay_in_minutes` can only be set on serverless dbs ([#30790](https://github.com/hashicorp/terraform-provider-azurerm/issues/30790))
* `azurerm_mssql_server` - the `azuread_administrator` block now updates in place rather than being deleted/recreated ([#30742](https://github.com/hashicorp/terraform-provider-azurerm/issues/30742))
* `azurerm_network_watcher_flow_log` - the `target_resource_id` property is now included in the update request payload resolving an issue where changing it failed to recreate or update the resource ([#30776](https://github.com/hashicorp/terraform-provider-azurerm/issues/30776))
* `azurerm_pim_eligible_role_assignment` - improve filter used during List requests to prevent timeouts ([#30705](https://github.com/hashicorp/terraform-provider-azurerm/issues/30705))
* `azurerm_postgresql_flexible_server_virtual_endpoint` - fix read error when in replica set in failover state ([#30789](https://github.com/hashicorp/terraform-provider-azurerm/issues/30789))

## 4.47.0 (October 02, 2025)

FEATURES:

* **New Resource**: `azurerm_api_management_workspace_policy` ([#30547](https://github.com/hashicorp/terraform-provider-azurerm/issues/30547))

ENHANCEMENTS:

* dependencies: `go-azure-sdk` - update to `v0.20250924.1155608` ([#30693](https://github.com/hashicorp/terraform-provider-azurerm/issues/30693))
* `azurerm_cognitive_account` - add support for value `AIServices` to `kind` property ([#30423](https://github.com/hashicorp/terraform-provider-azurerm/issues/30423))
* `azurerm_cognitive_account` - add the `project_management_enabled` property ([#30423](https://github.com/hashicorp/terraform-provider-azurerm/issues/30423))
* `azurerm_cognitive_account` - add the `network_injection` property ([#30423](https://github.com/hashicorp/terraform-provider-azurerm/issues/30423))
* `azurerm_palo_alto_local_rulestack_rule` - increase limit for `priority` to `1000000` ([#30712](https://github.com/hashicorp/terraform-provider-azurerm/issues/30712))
* `azurerm_stream_analytics_job` - add support for the `Msi` value in the `job_storage_account.authentication_mode` property ([#30728](https://github.com/hashicorp/terraform-provider-azurerm/issues/30728))

BUG FIXES:

* `azurerm_management_group_policy_remediation` - suppress casing difference on `policy_definition_reference_id` to avoid a perpetual diff as the API doesn't honour casing ([#30736](https://github.com/hashicorp/terraform-provider-azurerm/issues/30736))
* `azurerm_resource_group_policy_remediation` - suppress casing difference on `policy_definition_reference_id` to avoid a perpetual diff as the API doesn't honour casing ([#30736](https://github.com/hashicorp/terraform-provider-azurerm/issues/30736))
* `azurerm_resource_policy_remediation` - suppress casing difference on `policy_definition_reference_id` to avoid a perpetual diff as the API doesn't honour casing ([#30736](https://github.com/hashicorp/terraform-provider-azurerm/issues/30736))
* `azurerm_storage_account` - fix error that occurs around `queue_properties` when not specified ([#30746](https://github.com/hashicorp/terraform-provider-azurerm/issues/30746))
* `azurerm_subscription_policy_remediation` - suppress casing difference on `policy_definition_reference_id` to avoid a perpetual diff as the API doesn't honour casing ([#30736](https://github.com/hashicorp/terraform-provider-azurerm/issues/30736))

## 4.46.0 (September 25, 2025)

ENHANCEMENTS:

* dependencies: `frontdoor/webapplicationfirewallpolicies` - update to API version `2025-03-01` ([#29742](https://github.com/hashicorp/terraform-provider-azurerm/issues/29742))
* `azurerm_cdn_frontdoor_firewall_policy` - support for the `captcha_cookie_expiration_in_minutes` property and  the `CAPTCHA` value in the `custom_rule.action` property ([#29742](https://github.com/hashicorp/terraform-provider-azurerm/issues/29742))
* `azurerm_cdn_frontdoor_security_policy` - add update ability ([#30299](https://github.com/hashicorp/terraform-provider-azurerm/issues/30299))
* `azurerm_cognitive_account` - add support for `C2`, `C3`, `C4`, `D3`, and `S1`  to `sku_name` ([#30655](https://github.com/hashicorp/terraform-provider-azurerm/issues/30655))
* `azurerm_flex_function_app` - add support for the `http_concurrency` property ([#29678](https://github.com/hashicorp/terraform-provider-azurerm/issues/29678))
* `azurerm_kubernetes_cluster` - add support for the `api_server_access_profile.virtual_network_integration_enabled` and `api_server_access_profile.subnet_id` properties ([#30559](https://github.com/hashicorp/terraform-provider-azurerm/issues/30559))
* `azurerm_machine_learning_workspace` - add support for the `service_side_encryption_enabled` property ([#30478](https://github.com/hashicorp/terraform-provider-azurerm/issues/30478))
* `azurerm_mysql_flexible_server` - add support for the `managed_hsm_key_id` property ([#30502](https://github.com/hashicorp/terraform-provider-azurerm/issues/30502))
* `azurerm_netapp_volume` - add support for updating `protocols` ([#30643](https://github.com/hashicorp/terraform-provider-azurerm/issues/30643))
* `azurerm_netapp_volume_group_oracle` - add support for updating `protocols` ([#30643](https://github.com/hashicorp/terraform-provider-azurerm/issues/30643))
* `azurerm_netapp_volume_group_sap_hana` - add support for updating `protocols` ([#30643](https://github.com/hashicorp/terraform-provider-azurerm/issues/30643))
* `azurerm_postgresql_flexible_server` - add support for the `17` value in the `version` property ([#30683](https://github.com/hashicorp/terraform-provider-azurerm/issues/30683))
* `azurerm_storage_queue` - add support for the `storage_account_id` property ([#28752](https://github.com/hashicorp/terraform-provider-azurerm/issues/28752))

BUG FIXES:

* `azurerm_cdn_frontdoor_firewall_policy` - fix the read function so it now correctly marks the resource as gone ([#30704](https://github.com/hashicorp/terraform-provider-azurerm/issues/30704))

## 4.45.1 (September 22, 2025)

NOTES:

This release contains a Terraform Plugin SDK v2 version bump that prevents identity change validation from raising an error when prior identity is empty (all attributes are null).

BUG FIXES:

* dependencies: `hashicorp/terraform-plugin-sdk/v2` - update to `v2.38.1` ([#30667](https://github.com/hashicorp/terraform-provider-azurerm/issues/30667))
* `azurerm_network_interface` - ensure identity is set during non-refresh apply operations ([#30667](https://github.com/hashicorp/terraform-provider-azurerm/issues/30667))

## 4.45.0 (September 18, 2025)

FEATURES:

* **New Action**: `azurerm_virtual_machine_power` ([#30647](https://github.com/hashicorp/terraform-provider-azurerm/issues/30647))
* **New List Resource**: `azurerm_storage_account` ([#30614](https://github.com/hashicorp/terraform-provider-azurerm/issues/30614))
* **New List Resource**: `azurerm_virtual_network` ([#30614](https://github.com/hashicorp/terraform-provider-azurerm/issues/30614))

ENHANCEMENTS:

* dependencies: `go-azure-sdk` - update to `v0.20250908.1192604` ([#30644](https://github.com/hashicorp/terraform-provider-azurerm/issues/30644))
* `azurerm_kubernetes_cluster` - add support for the `network_profile.advanced_networking` block ([#30434](https://github.com/hashicorp/terraform-provider-azurerm/issues/30434))
* `azurerm_storage_account` - `expiration_action` supports `Block` ([#30599](https://github.com/hashicorp/terraform-provider-azurerm/issues/30599)) ([#30599](https://github.com/hashicorp/terraform-provider-azurerm/issues/30599))
* `azurerm_subnet` - add support for `sharing_scope` ([#30600](https://github.com/hashicorp/terraform-provider-azurerm/issues/30600))

## 4.44.0 (September 11, 2025)

ENHANCEMENTS:

* dependencies: `go-azure-sdk` - update to `v0.20250903.1204452` ([#30557](https://github.com/hashicorp/terraform-provider-azurerm/issues/30557))
* dependencies: `paloaltonetworks/firewalls` - update to API version `2025-05-23` ([#30587](https://github.com/hashicorp/terraform-provider-azurerm/issues/30587))
* `azurerm_kubernetes_cluster` - add support for `bootstrap_profile` ([#30532](https://github.com/hashicorp/terraform-provider-azurerm/issues/30532))
* `azurerm_kubernetes_cluster` - add support for the `none` value for `network_profile.outbound_type` ([#30532](https://github.com/hashicorp/terraform-provider-azurerm/issues/30532))
* `azurerm_netapp_volume` - add support for `accept_grow_capacity_pool_for_short_term_clone_split` ([#30494](https://github.com/hashicorp/terraform-provider-azurerm/issues/30494))
* `azurerm_service_plan` - add support for Premium V4 SKUs ([#30163](https://github.com/hashicorp/terraform-provider-azurerm/issues/30163))

BUG FIXES:

* `azurerm_application_insights_standard_web_test` - prevent Resource ID parsing errors when parsing `hidden-link` tags ([#28034](https://github.com/hashicorp/terraform-provider-azurerm/issues/28034))
* `azurerm_custom_ip_prefix` - fix an incorrect type-assertion that caused an error when attempting to read the API response ([#30537](https://github.com/hashicorp/terraform-provider-azurerm/issues/30537))
* `azurerm_monitor_activity_log_alert` - fix `name` validation ([#30590](https://github.com/hashicorp/terraform-provider-azurerm/issues/30590))
* `azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack` - fix import by parsing `network_virtual_appliance_id` insensitively ([#30597](https://github.com/hashicorp/terraform-provider-azurerm/issues/30597))
* `azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama` - fix import by parsing `network_virtual_appliance_id` insensitively ([#30597](https://github.com/hashicorp/terraform-provider-azurerm/issues/30597))
* `azurerm_postgresql_flexible_server` - fix a bug when setting the `source_server_id` property ([#30497](https://github.com/hashicorp/terraform-provider-azurerm/issues/30497))

## 4.43.0 (September 04, 2025)

FEATURES:

* **New Data Source:** `azurerm_oracle_autonomous_database_backup` ([#30201](https://github.com/hashicorp/terraform-provider-azurerm/issues/30201))
* **New Data Source:** `azurerm_oracle_autonomous_database_backups` ([#30201](https://github.com/hashicorp/terraform-provider-azurerm/issues/30201))
* **New Resource:** `azurerm_oracle_autonomous_database_backup` ([#30201](https://github.com/hashicorp/terraform-provider-azurerm/issues/30201))

ENHANCEMENTS:

* dependencies: `azurerm_api_management_backend` - upgrade API version to `2024-05-01` ([#30500](https://github.com/hashicorp/terraform-provider-azurerm/issues/30500))
* dependencies: `eventgrid` - upgrade to API version `2025-02-15` ([#30481](https://github.com/hashicorp/terraform-provider-azurerm/issues/30481))
* Data Source: `azurerm_dev_center_project_pool` - add support for the `single_sign_on_enabled` property ([#30440](https://github.com/hashicorp/terraform-provider-azurerm/issues/30440))
* `azurerm_dev_center_project_pool` - add support for the `single_sign_on_enabled` property ([#30440](https://github.com/hashicorp/terraform-provider-azurerm/issues/30440))
* `azurerm_management_group_policy_assignment` - `override.kind` can now be configured ([#30524](https://github.com/hashicorp/terraform-provider-azurerm/issues/30524))
* `azurerm_monitor_activity_log_alert` - add support for the `Security` value in the `recommendation_category` property ([#30192](https://github.com/hashicorp/terraform-provider-azurerm/issues/30192))
* `azurerm_postgresql_flexible_server_firewall_rule` - improve validation for the `start_ip_address` and `end_ip_address` properties to ensure the values are valid IPv4 addresses ([#30514](https://github.com/hashicorp/terraform-provider-azurerm/issues/30514))
* `azurerm_resource_group_policy_assignment` - `override.kind` can now be configured ([#30524](https://github.com/hashicorp/terraform-provider-azurerm/issues/30524))
* `azurerm_resource_policy_assignment` - `override.kind` can now be configured ([#30524](https://github.com/hashicorp/terraform-provider-azurerm/issues/30524))
* `azurerm_sentinel_automation_rule` - add support for the `action_incident_task` block ([#29295](https://github.com/hashicorp/terraform-provider-azurerm/issues/29295))
* `azurerm_subscription_policy_assignment` - `override.kind` can now be configured ([#30524](https://github.com/hashicorp/terraform-provider-azurerm/issues/30524))

BUG FIXES:

* `azurerm_flex_function_app` - fix `instance_memory_in_mb` update issue (GH-30489)
* `azurerm_kubernetes_cluster` - remove read-only field `NodeImageVersion` when cycling node pool ([#30416](https://github.com/hashicorp/terraform-provider-azurerm/issues/30416))
* `azurerm_kubernetes_cluster_node_pool` - remove read-only field `NodeImageVersion` when cycling node pool ([#30416](https://github.com/hashicorp/terraform-provider-azurerm/issues/30416))
* `azurerm_management_group_policy_set_definition` - fix an issue that caused API errors when `policy_definition_reference` blocks were added or removed ([#30493](https://github.com/hashicorp/terraform-provider-azurerm/issues/30493))
* `azurerm_policy_set_definition` - fix an issue that caused API errors when `policy_definition_reference` blocks were added or removed ([#30493](https://github.com/hashicorp/terraform-provider-azurerm/issues/30493))
* `azurerm_virtual_machine` - fix potential panic caused by the hash function for the `os_profile_linux_config` block ([#30456](https://github.com/hashicorp/terraform-provider-azurerm/issues/30456))

## 4.42.0 (August 28, 2025)

NOTES:

* This release contains a state migration that fixes a resource state parsing error for `azurerm_kusto_cluster` when the `language_extensions` property is defined. Users upgrading from a version older than `4.0.0` should upgrade directly to this release.

FEATURES:

* **New Data Source:** `azurerm_managed_disks` ([#30394](https://github.com/hashicorp/terraform-provider-azurerm/issues/30394))

ENHANCEMENTS:

* dependencies: `containerservice` - update api version to `2025-05-01` ([#30401](https://github.com/hashicorp/terraform-provider-azurerm/issues/30401))
* dependencies: `go-azure-sdk/resourcemanager` update to `v0.20250814.1105543` ([#30401](https://github.com/hashicorp/terraform-provider-azurerm/issues/30401))
* dependencies: `go-azure-sdk/sdk` update to `v0.20250814.1105543` ([#30401](https://github.com/hashicorp/terraform-provider-azurerm/issues/30401))
* `azurerm_iothub` - add support for `endpoint.subscription_id` property ([#27524](https://github.com/hashicorp/terraform-provider-azurerm/issues/27524))
* `azurerm_iothub_endpoint_cosmosdb_account` - add support for `endpoint.subscription_id` property ([#27524](https://github.com/hashicorp/terraform-provider-azurerm/issues/27524))
* `azurerm_iothub_endpoint_eventhub` - add support for `endpoint.subscription_id` property ([#27524](https://github.com/hashicorp/terraform-provider-azurerm/issues/27524))
* `azurerm_iothub_endpoint_servicebus_queue` - add support for `endpoint.subscription_id` property ([#27524](https://github.com/hashicorp/terraform-provider-azurerm/issues/27524))
* `azurerm_iothub_endpoint_servicebus_topic` - add support for `endpoint.subscription_id` property ([#27524](https://github.com/hashicorp/terraform-provider-azurerm/issues/27524))
* `azurerm_linux_virtual_machine` - add support for `os_managed_disk_id` property ([#30394](https://github.com/hashicorp/terraform-provider-azurerm/issues/30394))
* `azurerm_windows_virtual_machine` - add support for `os_managed_disk_id` property ([#30394](https://github.com/hashicorp/terraform-provider-azurerm/issues/30394))

BUG FIXES:

* `azurerm_kusto_cluster` - add a state migration for `language_extensions` to migrate from a list of strings to a list of objects (block) ([#30438](https://github.com/hashicorp/terraform-provider-azurerm/issues/30438))
* `azurerm_kusto_cluster` - fix an issue where removal of the `language_extensions` property was not applied to the API request ([#30449](https://github.com/hashicorp/terraform-provider-azurerm/issues/30449))
* `azurerm_linux_web_app` - normalize docker url ([#30368](https://github.com/hashicorp/terraform-provider-azurerm/issues/30368))

## 4.41.0 (August 21, 2025)

FEATURES:

* **New Resource**: `azurerm_network_manager_ipam_pool_static_cidr` ([#29501](https://github.com/hashicorp/terraform-provider-azurerm/issues/29501))
* **New Resource**: `azurerm_network_manager_routing_rule_collection` ([#29783](https://github.com/hashicorp/terraform-provider-azurerm/issues/29783))

ENHANCEMENTS:

* `azurerm_cdn_frontdoor_profile` - add support for the `log_scrubbing_rule` block ([#30115](https://github.com/hashicorp/terraform-provider-azurerm/issues/30115))
* `azurerm_monitor_diagnostic_setting` - update validation for `target_resource_id` to allow management group IDs as input ([#30447](https://github.com/hashicorp/terraform-provider-azurerm/issues/30447))
* `azurerm_netapp_account_encryption` - add support for `federated_client_id` and `cross_tenant_key_vault_resource_id`  ([#30373](https://github.com/hashicorp/terraform-provider-azurerm/issues/30373))
* `azurerm_netapp_pool` - add support for `custom_throughput_mibps`  ([#30404](https://github.com/hashicorp/terraform-provider-azurerm/issues/30404))

BUG FIXES:

* `azurerm_app_service_environment_v3` - fix drift on the `allow_new_private_endpoint_connections` property ([#30391](https://github.com/hashicorp/terraform-provider-azurerm/issues/30391))
* `azurerm_private_endpoint` - retry on `RetryableError` and `StorageAccountOperationInProgress` errors during LRO ([#28112](https://github.com/hashicorp/terraform-provider-azurerm/issues/28112))

## 4.40.0 (August 14, 2025)

FEATURES:

* **New Resource**: `azurerm_data_factory_customer_managed_key` ([#30341](https://github.com/hashicorp/terraform-provider-azurerm/issues/30341))

ENHANCEMENTS:

* `azurerm_eventgrid_system_topic` - suppress case difference on `source_resource_id` ([#30379](https://github.com/hashicorp/terraform-provider-azurerm/issues/30379))
* `azurerm_kubernetes_cluster` - add support for `gpu_profile` property ([#29954](https://github.com/hashicorp/terraform-provider-azurerm/issues/29954))
* `azurerm_load_test` - improved validation for the `encryption.identity.identity_id` property ([#30323](https://github.com/hashicorp/terraform-provider-azurerm/issues/30323))
* `azurerm_logic_app_standard` - refactored to leverage shared code with other `appservice` apps ([#30272](https://github.com/hashicorp/terraform-provider-azurerm/issues/30272))
* `azurerm_machine_learning_workspace` - support `provision_on_creation_enabled` property ([#30312](https://github.com/hashicorp/terraform-provider-azurerm/issues/30312))

## 4.39.0 (August 08, 2025)

FEATURES:

* **New Resource**: `azurerm_api_management_standalone_gateway` ([#30226](https://github.com/hashicorp/terraform-provider-azurerm/issues/30226))
* **New Resource**: `azurerm_eventgrid_partner_namespace` ([#30266](https://github.com/hashicorp/terraform-provider-azurerm/issues/30266))
* **New Resource**: `azurerm_postgresql_flexible_server_backup` ([#29201](https://github.com/hashicorp/terraform-provider-azurerm/issues/29201))

ENHANCEMENTS:

* dependencies: `cognitive` - update API version to `2025-06-01` ([#30302](https://github.com/hashicorp/terraform-provider-azurerm/issues/30302))
* dependencies: `machinelearning` - update API version to `2025-06-01` ([#30268](https://github.com/hashicorp/terraform-provider-azurerm/issues/30268))
* Data Source: `azurerm_oracle_db_system_shapes` - add support for the `zone` property ([#30071](https://github.com/hashicorp/terraform-provider-azurerm/issues/30071))
* Data Source: `azurerm_oracle_gi_versions` - add support for the `shape` and `zone` properties ([#30071](https://github.com/hashicorp/terraform-provider-azurerm/issues/30071))
* `azurerm_cognitive_deployment` - remove `model.format` validation ([#30276](https://github.com/hashicorp/terraform-provider-azurerm/issues/30276))
* `azurerm_eventgrid_system_topic` - add support for the `SystemAssigned, UserAssigned` value for the `identity.type` property ([#30339](https://github.com/hashicorp/terraform-provider-azurerm/issues/30339))
* `azurerm_linux_web_app` - add support for the `8.4` value in the `php_version` property ([#30281](https://github.com/hashicorp/terraform-provider-azurerm/issues/30281))
* `azurerm_linux_web_app_slot` - add support for the `8.4` value in the `php_version` property ([#30281](https://github.com/hashicorp/terraform-provider-azurerm/issues/30281))
* `azurerm_postgresql_flexible_server` - the `customer_managed_key.geo_backup_key_vault_key_id` now supports versionless IDs ([#30305](https://github.com/hashicorp/terraform-provider-azurerm/issues/30305))
* `azurerm_site_recovery_replicated_vm` - the `target_disk_type` property now supports the `StandardSSD_ZRS`, `Premium_ZRS` and `PremiumV2_LRS` values and the `target_replica_disk_type` now supports the `StandardSSD_ZRS` and `Premium_ZRS` properties ([#30291](https://github.com/hashicorp/terraform-provider-azurerm/issues/30291))

BUG FIXES:

* `azurerm_container_app_environment` - fix an issue where `identity` was not set to the update request payload ([#30311](https://github.com/hashicorp/terraform-provider-azurerm/issues/30311))
* `azurerm_function_app_flex_consumption` - the `maximum_instance_count` property now updates as expected ([#30342](https://github.com/hashicorp/terraform-provider-azurerm/issues/30342))
* `azurerm_kubernetes_cluster_node_pool` - add locks on `vnet_subnet_id` and `pod_subnet_id` to prevent conflicts while updating multiple node pools in parallel ([#29537](https://github.com/hashicorp/terraform-provider-azurerm/issues/29537))
* `azurerm_postgresql_flexible_server` - fix an issue where `administrator_password_wo` was not set to the update request payload ([#29475](https://github.com/hashicorp/terraform-provider-azurerm/issues/29475))

## 4.38.1 (July 31, 2025)

**NOTE:** This patch release addresses a critical problem in App Service and Logic Apps resources preventing all Long Running Operations from completing successfully. 

BUG FIXES:

* dependencies: `go-azure-sdk/sdk` update to `v0.20250731.1142049` ([#30282](https://github.com/hashicorp/terraform-provider-azurerm/issues/30282))

## 4.38.0 (July 30, 2025)

FEATURES:

* **New Data Source** : `azurerm_api_connection` ([#30178](https://github.com/hashicorp/terraform-provider-azurerm/issues/30178))
* **New Data Source**: `azurerm_log_analytics_workspace_table` ([#30261](https://github.com/hashicorp/terraform-provider-azurerm/issues/30261))
* **New Data Source**: `azurerm_mssql_failover_group` ([#29428](https://github.com/hashicorp/terraform-provider-azurerm/issues/29428))
* **New Data Source**: `azurerm_trusted_signing_account` ([#29293](https://github.com/hashicorp/terraform-provider-azurerm/issues/29293))
* **New Resource**: `azurerm_application_load_balancer_security_policy` ([#30128](https://github.com/hashicorp/terraform-provider-azurerm/issues/30128))
* **New Resource** : `azurerm_eventgrid_partner_registration` ([#29736](https://github.com/hashicorp/terraform-provider-azurerm/issues/29736))
* **New Resource**: `azurerm_mssql_managed_instance_start_stop_schedule` ([#26702](https://github.com/hashicorp/terraform-provider-azurerm/issues/26702))

ENHANCEMENTS:

* dependencies: `go-azure-sdk` update to `v0.20250728.1144148` ([#30254](https://github.com/hashicorp/terraform-provider-azurerm/issues/30254))
* dependencies: `go-azure-sdk` update to v0.20250716.1144812 ([#30171](https://github.com/hashicorp/terraform-provider-azurerm/issues/30171))
* dependencies: `golang.org/x/crypto` update to `v0.40.0` ([#30171](https://github.com/hashicorp/terraform-provider-azurerm/issues/30171))
* dependencies: `golang.org/x/mod` update to `v0.26.0` ([#30171](https://github.com/hashicorp/terraform-provider-azurerm/issues/30171))
* dependencies: `golang.org/x/net` update to `v0.42.0` ([#30171](https://github.com/hashicorp/terraform-provider-azurerm/issues/30171))
* dependencies: `golang.org/x/sync` update to `v0.16.0` ([#30171](https://github.com/hashicorp/terraform-provider-azurerm/issues/30171))
* dependencies: `golang.org/x/sys` update to `v0.34.0` ([#30171](https://github.com/hashicorp/terraform-provider-azurerm/issues/30171))
* dependencies: `golang.org/x/text` update to `v0.27.0` ([#30171](https://github.com/hashicorp/terraform-provider-azurerm/issues/30171))
* dependencies: `golang.org/x/tools` update to `v0.35.0` ([#30171](https://github.com/hashicorp/terraform-provider-azurerm/issues/30171))
* dependencies: `servicebus` - update to API version `2024-01-01` ([#30231](https://github.com/hashicorp/terraform-provider-azurerm/issues/30231))
* Data Source: `azurerm_databricks_workspace` - add support for the `custom_parameters` property ([#30214](https://github.com/hashicorp/terraform-provider-azurerm/issues/30214))
* Data Source: `azurerm_oracle_cloud_vm_cluster` - add support for the `file_system_configuration` block ([#30092](https://github.com/hashicorp/terraform-provider-azurerm/issues/30092))
* Data Source: `azurerm_oracle_exadata_infrastructure` - add support for the `defined_file_system_configuration` block ([#30092](https://github.com/hashicorp/terraform-provider-azurerm/issues/30092))
* `azurerm_batch_pool` - fix `start_task.0.task_retry_maximum` validation ([#30182](https://github.com/hashicorp/terraform-provider-azurerm/issues/30182))
* `azurerm_dev_center` - add support for the `project_catalog_item_sync_enabled` property ([#29274](https://github.com/hashicorp/terraform-provider-azurerm/issues/29274))
* `azurerm_dev_center_project_pool` - add support for the `managed_virtual_network_regions` property ([#30061](https://github.com/hashicorp/terraform-provider-azurerm/issues/30061))
* `azurerm_dynatrace_monitor` -  add support for  the `environment_properties` block ([#29251](https://github.com/hashicorp/terraform-provider-azurerm/issues/29251))
* `azurerm_image` - improve validation for `os_disk`, `data_disk` and `zone_resilient` ([#30222](https://github.com/hashicorp/terraform-provider-azurerm/issues/30222))
* `azurerm_managed_lustre_file_system` - add support for the `root_squash` block ([#29876](https://github.com/hashicorp/terraform-provider-azurerm/issues/29876))
* `azurerm_management_group_policy_assignment` - improve validation for the `name` property ([#30179](https://github.com/hashicorp/terraform-provider-azurerm/issues/30179))
* `azurerm_management_group_policy_set_definition` - now forces a new resource to be created when the number of `parameters` is decreased ([#29866](https://github.com/hashicorp/terraform-provider-azurerm/issues/29866))
* `azurerm_mongo_cluster` - add support for `version` 8.0 ([#29823](https://github.com/hashicorp/terraform-provider-azurerm/issues/29823))
* `azurerm_network_security_rule` - improve validation for source and destination properties  ([#29675](https://github.com/hashicorp/terraform-provider-azurerm/issues/29675))
* `azurerm_oracle_cloud_vm_cluster` - add support for the `file_system_configuration` block ([#30092](https://github.com/hashicorp/terraform-provider-azurerm/issues/30092))
* `azurerm_policy_set_definition` - now forces a new resource to be created when the number of `parameters` is decreased ([#29866](https://github.com/hashicorp/terraform-provider-azurerm/issues/29866))
* `azurerm_resource_group_policy_assignment` - improve validation for the `name` property ([#30179](https://github.com/hashicorp/terraform-provider-azurerm/issues/30179))
* `azurerm_resource_policy_assignment` - improve validation for the `name` property ([#30179](https://github.com/hashicorp/terraform-provider-azurerm/issues/30179))
* `azurerm_subnet` - add support for the `ip_address_pool` block ([#29840](https://github.com/hashicorp/terraform-provider-azurerm/issues/29840))
* `azurerm_subscription_policy_assignment` - improve validation for the `name` property ([#30179](https://github.com/hashicorp/terraform-provider-azurerm/issues/30179))
* `azurerm_video_indexer_account` - add support for the `public_network_access` property ([#29725](https://github.com/hashicorp/terraform-provider-azurerm/issues/29725))

BUG FIXES:

* Data Source: `azurerm_kusto_cluster` - fix returned error if cluster was not found ([#30232](https://github.com/hashicorp/terraform-provider-azurerm/issues/30232))
* `appservice` - now checks for deployment service availability before zip deployment  ([#30066](https://github.com/hashicorp/terraform-provider-azurerm/issues/30066))
* `azurerm_ai_foundry` - no longer crashes when the `key_vault_id` property is nil ([#30252](https://github.com/hashicorp/terraform-provider-azurerm/issues/30252))
* `azurerm_container_app_environment` - no longer panics when `log_analytics_workspace_id` is from another subscription ([#29829](https://github.com/hashicorp/terraform-provider-azurerm/issues/29829))
* `azurerm_eventhub` - fix perpetual diff with `message_retention` ([#30169](https://github.com/hashicorp/terraform-provider-azurerm/issues/30169))
* `azurerm_kusto_attached_database_configuration` - resource is now removed from state if it no longer exists ([#30232](https://github.com/hashicorp/terraform-provider-azurerm/issues/30232))
* `azurerm_kusto_cluster` - resource is now removed from state if it no longer exists ([#30232](https://github.com/hashicorp/terraform-provider-azurerm/issues/30232))
* `azurerm_kusto_cluster_customer_managed_key` - resource is now removed from state if it no longer exists ([#30232](https://github.com/hashicorp/terraform-provider-azurerm/issues/30232))
* `azurerm_kusto_cluster_principal_assignment` - resource is now removed from state if it no longer exists ([#30232](https://github.com/hashicorp/terraform-provider-azurerm/issues/30232))
* `azurerm_log_analytics_workspace_table` - the `retention_in_days` property can now be reset ([#29182](https://github.com/hashicorp/terraform-provider-azurerm/issues/29182))
* `azurerm_monitor_alert_prometheus_rule_group` - prevent an error caused by the request containing an empty string for the `rule.for` property when not set ([#30180](https://github.com/hashicorp/terraform-provider-azurerm/issues/30180))
* `azurerm_mssql_database` - the `max_size_gb` can now support `0.1` and `0.5` ([#28334](https://github.com/hashicorp/terraform-provider-azurerm/issues/28334))
* `azurerm_mssql_managed_instance` - `administrator_login` is now Computed, preventing resource recreation when `azure_active_directory_administrator.azuread_authentication_only_enabled` is `true` ([#30263](https://github.com/hashicorp/terraform-provider-azurerm/issues/30263))
* `azurerm_postgresql_flexible_server_virtual_endpoint` - no longer causes an error when `replica_server_id` is from another subscription ([#29270](https://github.com/hashicorp/terraform-provider-azurerm/issues/29270))
* `azurerm_role_management_policy` - fix perpetual diff on `activation_rules.approval_stage` ([#29084](https://github.com/hashicorp/terraform-provider-azurerm/issues/29084))
* `azurerm_service_plan` - fix an issue that prevented supported SKUs from specifying `zone_balancing_enabled` as `true` ([#30165](https://github.com/hashicorp/terraform-provider-azurerm/issues/30165))
* `azurerm_web_application_firewall_policy` - `js_challenge_cookie_expiration_in_minutes` is now set to default value if not returned from API ([#30245](https://github.com/hashicorp/terraform-provider-azurerm/issues/30245))

## 4.37.0 (July 17, 2025)

FEATURES:

* **New Data Source**: `azurerm_network_manager_ipam_pool` ([#30145](https://github.com/hashicorp/terraform-provider-azurerm/issues/30145))

ENHANCEMENTS:

* Data Source: `azurerm_virtual_machine_scale_set` - add support for the `auxiliary_mode` and `auxiliary_sku` properties ([#30159](https://github.com/hashicorp/terraform-provider-azurerm/issues/30159))
* `azurerm_container_app_environment` - add support for the `identity` block ([#29409](https://github.com/hashicorp/terraform-provider-azurerm/issues/29409))
* `azurerm_eventhub` - add  support for the `retention_description` block ([#29427](https://github.com/hashicorp/terraform-provider-azurerm/issues/29427))
* `azurerm_kubernetes_cluster` - add support for the `Daily` value in the `maintenance_window_auto_upgrade.frequency` property ([#30133](https://github.com/hashicorp/terraform-provider-azurerm/issues/30133))
* `azurerm_kubernetes_flux_configuration` - add support for the `git_repository.provider` property ([#30082](https://github.com/hashicorp/terraform-provider-azurerm/issues/30082))
* `azurerm_mssql_job_step` - the `job_credential_id` and `output_target.job_credential_id` properties are now optional ([#30031](https://github.com/hashicorp/terraform-provider-azurerm/issues/30031))
* `azurerm_orchestrated_virtual_machine_scale_set` - add support for `auxiliary_mode` and `auxiliary_sku`  ([#30102](https://github.com/hashicorp/terraform-provider-azurerm/issues/30102))
* `azurerm_storage_account` - add support for the `provisioned_billing_model_version` property ([#29043](https://github.com/hashicorp/terraform-provider-azurerm/issues/29043))
* `azurerm_vpn_gateway_connection` - add support for the `dpd_timeout_seconds` property ([#29434](https://github.com/hashicorp/terraform-provider-azurerm/issues/29434))

BUG FIXES:

* Data Source: `azurerm_virtual_machine_scale_set` - fix a panic caused by missing properties ([#30159](https://github.com/hashicorp/terraform-provider-azurerm/issues/30159))
* `azurerm_container_app_environment` - fix import for `workload_profile` ([#30139](https://github.com/hashicorp/terraform-provider-azurerm/issues/30139))
* `azurerm_mongo_cluster` - the `create_mode` property no longer causes ForceNews on import ([#29375](https://github.com/hashicorp/terraform-provider-azurerm/issues/29375))
* `azurerm_virtual_network` - suppress a perpetual diff on `address_space` when using `ip_address_pool` ([#30073](https://github.com/hashicorp/terraform-provider-azurerm/issues/30073))
* `azurerm_vpn_gateway_connection` - the `shared_key` is now Optional + Computed ([#30152](https://github.com/hashicorp/terraform-provider-azurerm/issues/30152))

## 4.36.0 (July 10, 2025)

FEATURES:

* **New Resource**: `azurerm_api_management_workspace` ([#30033](https://github.com/hashicorp/terraform-provider-azurerm/issues/30033))
* **New Resource**: `azurerm_network_manager_verifier_workspace_reachability_analysis_intent` ([#28956](https://github.com/hashicorp/terraform-provider-azurerm/issues/28956))

ENHANCEMENTS:

* dependencies: `kubernetesconfiguration` - update to API version `2024-11-01` ([#29896](https://github.com/hashicorp/terraform-provider-azurerm/issues/29896))
* dependencies: `oracle` - update to API version `2025-03-01` ([#29721](https://github.com/hashicorp/terraform-provider-azurerm/issues/29721))
* dependencies: `servicenetworking` - update to API version `2025-01-01` ([#30103](https://github.com/hashicorp/terraform-provider-azurerm/issues/30103))
* Data Source: `azurerm_container_registry` - add support for the `data_endpoint_host_names` property ([#30086](https://github.com/hashicorp/terraform-provider-azurerm/issues/30086))
* Data Source: `azurerm_dev_center_dev_box_definition` - add support for the `hibernate_support_enabled` property ([#29995](https://github.com/hashicorp/terraform-provider-azurerm/issues/29995))
* Data Source: `azurerm_marketplace_agreement` - add support for the `accepted` property ([#30118](https://github.com/hashicorp/terraform-provider-azurerm/issues/30118))
* Data Source: `azurerm_oracle_autonomous_database` - add support for `compute_model` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* Data Source: `azurerm_oracle_cloud_vm_cluster` - add support for `compute_model` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* Data Source: `azurerm_oracle_db_servers` - add support for the `compute_model` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* Data Source: `azurerm_oracle_db_system_shapes` - add support for the `display_name` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* Data Source: `azurerm_oracle_db_system_shapes` - add support for the `are_server_types_supported` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* Data Source: `azurerm_oracle_db_system_shapes` - add support for the `compute_model` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* Data Source: `azurerm_oracle_exadata_infrastructure` - add support for the `compute_model` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* Data Source: `azurerm_oracle_exadata_infrastructure` - add support for the `database_server_type` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* Data Source: `azurerm_oracle_exadata_infrastructure` - add support for the `storage_server_type` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* Data Source: `azurerm_private_dns_zone_virtual_network_link` - add support for the `resolution_policy` property ([#29861](https://github.com/hashicorp/terraform-provider-azurerm/issues/29861))
* `azurerm_api_management` - `sku_name` now supports V2 Tiers  ([#29657](https://github.com/hashicorp/terraform-provider-azurerm/issues/29657))
* `azurerm_container_registry` - add support for the `data_endpoint_host_names` property ([#30086](https://github.com/hashicorp/terraform-provider-azurerm/issues/30086))
* `azurerm_data_protection_backup_instance_disk` - support cross subscription snapshot resource group ([#30087](https://github.com/hashicorp/terraform-provider-azurerm/issues/30087))
* `azurerm_dev_center_dev_box_definition` - add support for the `hibernate_support_enabled` property ([#29995](https://github.com/hashicorp/terraform-provider-azurerm/issues/29995))
* `azurerm_kubernetes_cluster` - add support for the `custom_ca_trust_certificates_base64` property ([#29894](https://github.com/hashicorp/terraform-provider-azurerm/issues/29894))
* `azurerm_kubernetes_cluster` - support for the `web_app_routing.default_nginx_controller` property ([#29879](https://github.com/hashicorp/terraform-provider-azurerm/issues/29879))
* `azurerm_linux_virtual_machine_scale_set` - add support for the `network_interface.auxiliary_mode` and `network_interface.auxiliary_sku` properties ([#29724](https://github.com/hashicorp/terraform-provider-azurerm/issues/29724))
* `azurerm_linux_web_app` - support for the `vnet_image_pull_enabled` property ([#29452](https://github.com/hashicorp/terraform-provider-azurerm/issues/29452))
* `azurerm_linux_web_app_slot` - support for the `vnet_image_pull_enabled` property ([#29452](https://github.com/hashicorp/terraform-provider-azurerm/issues/29452))
* `azurerm_log_analytics_workspace` - now returns an error during planning when creating with/updating to a `Standard` or `Premium` SKU as this is no longer supported by Azure ([#30101](https://github.com/hashicorp/terraform-provider-azurerm/issues/30101))
* `azurerm_logic_app_workflow` - The `access_control.trigger.allowed_caller_ip_address_range` property is now optional ([#30041](https://github.com/hashicorp/terraform-provider-azurerm/issues/30041))
* `azurerm_machine_learning_datastore_blobstorage` - the `shared_access_signature` and `account_key` properties are now optional ([#30079](https://github.com/hashicorp/terraform-provider-azurerm/issues/30079))
* `azurerm_netapp_volume` - add support for the `cool_access` block ([#29915](https://github.com/hashicorp/terraform-provider-azurerm/issues/29915))
* `azurerm_oracle_autonomous_database` - Add support for `allowed_ips` ([#29412](https://github.com/hashicorp/terraform-provider-azurerm/issues/29412))
* `azurerm_oracle_exadata_infrastructure` - add support for the `database_server_type` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* `azurerm_oracle_exadata_infrastructure` - add support for the `storage_server_type` property ([#29801](https://github.com/hashicorp/terraform-provider-azurerm/issues/29801))
* `azurerm_private_dns_zone_virtual_network_link` - add support for the `resolution_policy` property ([#29861](https://github.com/hashicorp/terraform-provider-azurerm/issues/29861))
* `azurerm_public_ip_prefix` - add support for the `custom_ip_prefix_id` property ([#29851](https://github.com/hashicorp/terraform-provider-azurerm/issues/29851))
* `azurerm_service_plan` - allow updating `zone_balancing_enabled` without recreating the resource in supported configurations ([#29810](https://github.com/hashicorp/terraform-provider-azurerm/issues/29810))
* `azurerm_virtual_hub` - add support for the `branch_to_branch_traffic_enabled` property ([#29453](https://github.com/hashicorp/terraform-provider-azurerm/issues/29453))
* `azurerm_windows_virtual_machine_scale_set` - add support for the `network_interface.auxiliary_mode` and `network_interface.auxiliary_sku` properties ([#29724](https://github.com/hashicorp/terraform-provider-azurerm/issues/29724))

BUG FIXES:

* `azurerm_mobile_network_packet_core_control_plane` - the `site_ids` property is now marked as `ForceNew` ([#30056](https://github.com/hashicorp/terraform-provider-azurerm/issues/30056))
* `azurerm_mobile_network_slice` - the `single_network_slice_selection_assistance_information` property is now updated correctly ([#30057](https://github.com/hashicorp/terraform-provider-azurerm/issues/30057))
* `azurerm_private_dns_resolver_dns_forwarding_ruleset` - fix an issue where `private_dns_resolver_outbound_endpoint_ids` failed to update ([#30046](https://github.com/hashicorp/terraform-provider-azurerm/issues/30046))

## 4.35.0 (July 01, 2025)

FEATURES:
* **New Resource**: `azurerm_email_communication_service_domain_sennder_username` ([#29340](https://github.com/hashicorp/terraform-provider-azurerm/issues/29340))
* **New Resource**: `azurerm_management_group_policy_set_definition` ([#29863](https://github.com/hashicorp/terraform-provider-azurerm/issues/29863))

ENHANCEMENTS:
* **Data Source**: `azurerm_communication_service` - add support for the `immutable_resource_id` property ([#29912](https://github.com/hashicorp/terraform-provider-azurerm/issues/29912))
* `azurerm_cdn_endpoint` - block creation of all Azure CDN(classic) resources while allowing existing resources to be updated ([#29299](https://github.com/hashicorp/terraform-provider-azurerm/issues/29299))
* `azurerm_cdn_endpoint_custom_domain` - block creation of all Azure CDN(classic) resources while allowing existing resources to be updated ([#29299](https://github.com/hashicorp/terraform-provider-azurerm/issues/29299))
* `azurerm_cdn_profile` - block creation of all Azure CDN(classic) resources while allowing existing resources to be updated ([#29299](https://github.com/hashicorp/terraform-provider-azurerm/issues/29299))
* `azurerm_container_app_job` - add support for the `volume_mounts.sub_path` property ([#29883](https://github.com/hashicorp/terraform-provider-azurerm/issues/29883))
* `azurerm_container_app` - add support for the `cors` property ([#29785](https://github.com/hashicorp/terraform-provider-azurerm/issues/29785))
* `azurerm_data_protection_backup_policy_disk` - the `absolute_criteria` property now supports the `AllBackup`, `FirstOfMonth` and `FirstOfYear` values ([#29917](https://github.com/hashicorp/terraform-provider-azurerm/issues/29917))
* `azurerm_frontdoor` - block new resource creation while allowing existing resources to be updated ([#29257](https://github.com/hashicorp/terraform-provider-azurerm/issues/29257))
* `azurerm_frontdoor_custom_https_configuration` - block new resource creation while allowing existing resources to be updated ([#29257](https://github.com/hashicorp/terraform-provider-azurerm/issues/29257))
* `azurerm_frontdoor_firewall_policy` - block new resource creation while allowing existing resources to be updated ([#29257](https://github.com/hashicorp/terraform-provider-azurerm/issues/29257))
* `azurerm_frontdoor_rules_engine` - block new resource creation while allowing existing resources to be updated ([#29257](https://github.com/hashicorp/terraform-provider-azurerm/issues/29257))
* `azurerm_function_app_flex_consumption` - add support for the `vnet_route_all_enabled` property ([#29839](https://github.com/hashicorp/terraform-provider-azurerm/issues/29839))
* `azurerm_machine_learning_compute_cluster` - the `scale_settings` block and its sub-properties are no longer `ForceNew` ([#29878](https://github.com/hashicorp/terraform-provider-azurerm/issues/29878))
* `azurerm_machine_learning_compute_cluster` - the `tags` property is no longer `ForceNew` ([#29878](https://github.com/hashicorp/terraform-provider-azurerm/issues/29878))
* `azurerm_oracle_autonomous_database ` - add support for `long_term_backup_schedule` ([#29207](https://github.com/hashicorp/terraform-provider-azurerm/issues/29207)) 
* `azurerm_policy_set_definition` - add support for the `policy_definition_reference.version` property ([#29924](https://github.com/hashicorp/terraform-provider-azurerm/issues/29924))
* `azurerm_policy_set_definition` - migrate to use `go-azure-sdk` ([#29863](https://github.com/hashicorp/terraform-provider-azurerm/issues/29863)) 
* `azurerm_private_link_service` - add support for the `destination_ip_address` property ([#29395](https://github.com/hashicorp/terraform-provider-azurerm/issues/29395))
* `azurerm_purview_account` - add support for the `managed_event_hub_enabled` and `aws_external_id` properties ([#29732](https://github.com/hashicorp/terraform-provider-azurerm/issues/29732))
* `azurerm_virtual_network_gateway` - the `ip_configuration.public_ip_address_id` property is now optional ([#30038](https://github.com/hashicorp/terraform-provider-azurerm/issues/30038))
* `azurerm_windows_virtual_machine`: `os_disk.0.diff_disk_settings.0.placement` now supports `NvmeDisk` ([#29922](https://github.com/hashicorp/terraform-provider-azurerm/issues/29922))

BUG FIXES:
* `provider` - allow missing `subscription_id` when `use_cli` is `true` ([#29985](https://github.com/hashicorp/terraform-provider-azurerm/issues/29985))
* `azurerm_netapp_backup_policy` - the `weekly_backups_to_keep` and `monthly_backups_to_keep` properties can now be set to `0` ([#29920](https://github.com/hashicorp/terraform-provider-azurerm/issues/29920))

## 4.34.0 (June 20, 2025)

ENHANCEMENTS:
* dependencies: `containerservice` - update API version to `2025-02-01` ([#29761](https://github.com/hashicorp/terraform-provider-azurerm/issues/29761))
* `azurerm_network_manager_ipam_pool` - `display_name` is now optional ([#29842](https://github.com/hashicorp/terraform-provider-azurerm/issues/29842))

* `dependencies`: `go-azure-sdk` - update to `v0.20250613.1153526` ([#29871](https://github.com/hashicorp/terraform-provider-azurerm/issues/29871))
* `provider`: add support for `msi_api_version` property and `ARM_MSI_API_VERSION` env var. ([#29871](https://github.com/hashicorp/terraform-provider-azurerm/issues/29871))
* `azurerm_kusto_cluster_customer_managed_key` - add support for `managed_hsm_key_id` ([#29416](https://github.com/hashicorp/terraform-provider-azurerm/issues/29416))

FEATURES:
* **New Data Source**: `azurerm_dev_center_environment_type` ([#29782](https://github.com/hashicorp/terraform-provider-azurerm/issues/29782))
* **New Data Source**: `azurerm_dev_center_project_pool` ([#29778](https://github.com/hashicorp/terraform-provider-azurerm/issues/29778))

BUG FIXES:
* `azurerm_eventgrid_namespace` - validations for `maximum_session_expiry_in_hours` and `maximum_client_sessions_per_authentication_name` are now correct ([#29919](https://github.com/hashicorp/terraform-provider-azurerm/issues/29919))
* `azurerm_api_management_api_operation` - fix validation for the `url_template` property to allow parameters prefixed with `*` ([#29895](https://github.com/hashicorp/terraform-provider-azurerm/issues/29895))
* `azurerm_mysql_flexible_server` - reverted a change made to the validation of the `sku_name` property that caused errors for existing resources ([#29909](https://github.com/hashicorp/terraform-provider-azurerm/issues/29909))
* `azurerm_orchestrated_virtual_machine_scale_set` - prevent a panic when an empty `os_profile` block is present in the configuration  ([#29809](https://github.com/hashicorp/terraform-provider-azurerm/issues/29809))

## 4.33.0 (June 12, 2025)

FEATURES:

* **New Data Source**: `azurerm_dev_center_attached_network` ([#29793](https://github.com/hashicorp/terraform-provider-azurerm/issues/29793))
* **New Data Source**: `azurerm_dev_center_dev_box_definition` ([#29790](https://github.com/hashicorp/terraform-provider-azurerm/issues/29790))
* **New Data Source**: `azurerm_dev_center_catalog` ([#29794](https://github.com/hashicorp/terraform-provider-azurerm/issues/29794))
* **New Data Source**: `azurerm_dev_center_gallery` ([#29795](https://github.com/hashicorp/terraform-provider-azurerm/issues/29795))
* **New Data Source**: `azurerm_dev_center_network_connection` ([#29792](https://github.com/hashicorp/terraform-provider-azurerm/issues/29792))

ENHANCEMENTS:

* `azurem_netapp_volume_group_oracle_resource` - add support for `data_protection_replication ` including Cross-Region Replication (CRR) and Cross-Zone Replication (CZR) ([#29771](https://github.com/hashicorp/terraform-provider-azurerm/issues/29771))
* `azurerm_postgresql_flexible_server` - the `create_mode` property now supports the `ReviveDropped` value  ([#29814](https://github.com/hashicorp/terraform-provider-azurerm/issues/29814))
* `azurerm_postgresql_flexible_server` - add support for `SystemAssigned, UserAssigned` to the `identity.type` property ([#29320](https://github.com/hashicorp/terraform-provider-azurerm/issues/29320))

BUG FIXES:
* `azurerm_windows_function_app` - the `app_settings` property is no longer marked as sensitive ([#29834](https://github.com/hashicorp/terraform-provider-azurerm/issues/29834))
* `azurerm_mssql_server_vulnerability_assessment` - `storage_account_access_key` and `storage_container_sas_key` are no longer required to be set ([#29789](https://github.com/hashicorp/terraform-provider-azurerm/issues/29789))

## 4.32.0 (June 05, 2025)

FEATURES:

* **New Data Source**: `azurerm_dev_center_project` ([#29747](https://github.com/hashicorp/terraform-provider-azurerm/issues/29747))
* **New Data Source**: `azurerm_dev_center_project_environment_type` ([#29762](https://github.com/hashicorp/terraform-provider-azurerm/issues/29762))
* **New Resource**: `azurerm_qumulo_file_system` ([#28704](https://github.com/hashicorp/terraform-provider-azurerm/issues/28704))

ENHANCEMENTS:

* dependencies: `go-azure-sdk` - update to `v0.20250526.1224007` ([#29745](https://github.com/hashicorp/terraform-provider-azurerm/issues/29745))
* Data Source: `azurerm_netapp_volume` - export the `large_volume_enabled` property ([#29712](https://github.com/hashicorp/terraform-provider-azurerm/issues/29712))
* Data Source: `azurerm_vpn_gateway` - export the `ip_configuration` block ([#29186](https://github.com/hashicorp/terraform-provider-azurerm/issues/29186))
* `azurerm_kubernetes_cluster` - the `vm_size` property is now optional ([#29612](https://github.com/hashicorp/terraform-provider-azurerm/issues/29612))
* `azurerm_kubernetes_cluster_node_pool` - the `vm_size` property is now optional ([#29612](https://github.com/hashicorp/terraform-provider-azurerm/issues/29612))
* `azurerm_netapp_volume` - allow volumes made from snapshots to have a different pool than the original volume ([#29425](https://github.com/hashicorp/terraform-provider-azurerm/issues/29425))
* `azurerm_netapp_volume` - add support for the `large_volume_enabled` property ([#29712](https://github.com/hashicorp/terraform-provider-azurerm/issues/29712))
* `azurerm_postgresql_flexible_server` - add support for versionless key vault key IDs to the `customer_managed_key.key_vault_key_id` property ([#29741](https://github.com/hashicorp/terraform-provider-azurerm/issues/29741))
* `azurerm_virtual_network` - add support for the `ip_address_pool` block ([#29021](https://github.com/hashicorp/terraform-provider-azurerm/issues/29021))
* `azurerm_vpn_gateway` - export the `ip_configuration` block ([#29186](https://github.com/hashicorp/terraform-provider-azurerm/issues/29186))

BUG FIXES:

* Data Source: `azurerm_lb_backend_address_pool` - the `inbound_nat_rule_port_mapping.frontend_port` and `inbound_nat_rule_port_mapping.backend_port` are now set correctly ([#29791](https://github.com/hashicorp/terraform-provider-azurerm/issues/29791))
* `keyvault` - fix locking around the keyvault cache ([#28330](https://github.com/hashicorp/terraform-provider-azurerm/issues/28330))

## 4.31.0 (May 29, 2025)

FEATURES:

* **New Data Source**: `azurerm_dev_center` ([#29716](https://github.com/hashicorp/terraform-provider-azurerm/issues/29716))
* **New Resource**: `azurerm_network_manager_routing_configuration` ([#29310](https://github.com/hashicorp/terraform-provider-azurerm/issues/29310))

ENHANCEMENTS:

* dependencies: `azurerm_managed_lustre_file_system` - update to API version `2024-07-01` ([#29433](https://github.com/hashicorp/terraform-provider-azurerm/issues/29433))
* dependencies: `azurerm_mssql_server_vulnerability_assessment` - update to API version `2023-08-01-preview` ([#29373](https://github.com/hashicorp/terraform-provider-azurerm/issues/29373))
* dependencies: `azurerm_virtual_machine_scale_set_standby_pool` - update to API version `2025-03-01` ([#29649](https://github.com/hashicorp/terraform-provider-azurerm/issues/29649))
* dependencies: `compute` - partial update to API version `2024-11-01` ([#29666](https://github.com/hashicorp/terraform-provider-azurerm/issues/29666))
* dependencies: `videoindexer` - update to API version `2025-04-01` ([#29715](https://github.com/hashicorp/terraform-provider-azurerm/issues/29715))
* `azurerm_backup_protected_vm` - add support for the `BackupsSuspended` value to the `protection_state` property ([#29710](https://github.com/hashicorp/terraform-provider-azurerm/issues/29710))
* `azurerm_dashboard_grafana_managed_private_endpoint` - add support for the `privatelink_service_url` property ([#29466](https://github.com/hashicorp/terraform-provider-azurerm/issues/29466))
* `azurerm_dynatrace_tag_rules` - add support for the `sending_metrics_enabled` property ([#29499](https://github.com/hashicorp/terraform-provider-azurerm/issues/29499))
* `azurerm_function_app_flex_consumption` - add support for the `https_only` property ([#29024](https://github.com/hashicorp/terraform-provider-azurerm/issues/29024))
* `azurerm_mysql_flexible_server` - add support for the `MO_Standard_E96ads_v5` value to the `sku_name` property ([#29709](https://github.com/hashicorp/terraform-provider-azurerm/issues/29709))
* `azurerm_postgresql_flexible_server` - lock the source server when creating a replica server ([#29337](https://github.com/hashicorp/terraform-provider-azurerm/issues/29337))

BUG FIXES:

* `azurerm_api_management_product` - allow setting the `subscriptions_limit` property to `0` ([#28133](https://github.com/hashicorp/terraform-provider-azurerm/issues/28133))
* `azurerm_api_management_api` - add additional validation to catch when `api_type` is `websocket` but `service_url` is left empty ([#29624](https://github.com/hashicorp/terraform-provider-azurerm/issues/29624))
* `azurerm_batch_pool` - the `data_disks` property will now be correctly updated ([#29377](https://github.com/hashicorp/terraform-provider-azurerm/issues/29377))
* `azurerm_data_factory_dataset_binary` - fix incorrect casing of the `compression.type` property when sent to the API which caused compression to not be set ([#29273](https://github.com/hashicorp/terraform-provider-azurerm/issues/29273))
* `azurerm_cdn_frontdoor_rule` - fix shared schema validation of the `operator` property and use the correct package for validations ([#29482](https://github.com/hashicorp/terraform-provider-azurerm/issues/29482))
* `azurerm_hdinsight_hadoop_cluster` - changing the `script_action` property now forces a new resource to be created instead of silenty failing to update ([#28262](https://github.com/hashicorp/terraform-provider-azurerm/issues/28262))
* `azurerm_hbase_hadoop_cluster` - changing the `script_action` property now forces a new resource to be created instead of silenty failing to update ([#28262](https://github.com/hashicorp/terraform-provider-azurerm/issues/28262))
* `azurerm_interactive_query_hadoop_cluster` - changing the `script_action` property now forces a new resource to be created instead of silenty failing to update ([#28262](https://github.com/hashicorp/terraform-provider-azurerm/issues/28262))
* `azurerm_kafka_hadoop_cluster` - changing the `script_action` property now forces a new resource to be created instead of silenty failing to update ([#28262](https://github.com/hashicorp/terraform-provider-azurerm/issues/28262))
* `azurerm_linux_virtual_machine` - fix update for `identity` when VM has VMExtensions configured ([#29717](https://github.com/hashicorp/terraform-provider-azurerm/issues/29717))
* `azurerm_mongo_cluster` - connection strings conaining a `$` now get exported correctly ([#29669](https://github.com/hashicorp/terraform-provider-azurerm/issues/29669))
* `azurerm_mssql_virtual_machine` - `auto_patching` is now disabled when the block is not specified ([#29723](https://github.com/hashicorp/terraform-provider-azurerm/issues/29723))
* `azurerm_mssql_server_vulnerability_assessment` - `storage_account_access_key` or `storage_container_sas_key` property is now a `required` field ([#29373](https://github.com/hashicorp/terraform-provider-azurerm/issues/29373))
* `azurerm_network_interface` - `tags` can now be updated when NIC is attached to a private endpoint ([#29319](https://github.com/hashicorp/terraform-provider-azurerm/issues/29319))
* `azurerm_postgresql_flexible_server_configuration` - now checks the server state before restarting it ([#29221](https://github.com/hashicorp/terraform-provider-azurerm/issues/29221))
* `azurerm_search_service` - prevent a bug that cleared the `network_rule_bypass_option` property when only updating the `allowed_ips` property ([#29246](https://github.com/hashicorp/terraform-provider-azurerm/issues/29246))
* `azurerm_service_fabric_managed_cluster` - support for the `subnet_id` property ([#29216](https://github.com/hashicorp/terraform-provider-azurerm/issues/29216))
* `azurerm_spark_hadoop_cluster` - changing the `script_action` property now forces a new resource to be created instead of silenty failing to update ([#28262](https://github.com/hashicorp/terraform-provider-azurerm/issues/28262))

## 4.30.0 (May 22, 2025)

ENHANCEMENTS:

* dependencies: `go-azure-sdk` - update to `v0.20250520.1180806` ([#29665](https://github.com/hashicorp/terraform-provider-azurerm/issues/29665))
* Data Source: `azurerm_managed_disk` - add support for `location` ([#29513](https://github.com/hashicorp/terraform-provider-azurerm/issues/29513))
* `azurerm_dns_caa_record` - add support for the `contactemail` value in the `tag` property ([#29664](https://github.com/hashicorp/terraform-provider-azurerm/issues/29664))
* `azurerm_eventhub_namespace_schema_group` - add support for the `Json` value in the `schema_type` property ([#29641](https://github.com/hashicorp/terraform-provider-azurerm/issues/29641))
* `azurerm_function_app_flex_consumption` - add support for the `always_ready` block ([#29023](https://github.com/hashicorp/terraform-provider-azurerm/issues/29023))
* `azurerm_security_center_subscription_pricing` - add support for the `AI` value for the `resource_type` property ([#29631](https://github.com/hashicorp/terraform-provider-azurerm/issues/29631))

## 4.29.0 (May 16, 2025)

FEATURES: 

* **New Resource**: `azurerm_system_center_virtual_machine_manager_virtual_machine_instance_guest_agent` ([#28953](https://github.com/hashicorp/terraform-provider-azurerm/issues/28953))

ENHANCEMENTS:

* `azurerm_api_management_api` - fix `import` of resources ([#28193](https://github.com/hashicorp/terraform-provider-azurerm/issues/28193))
* `azurerm_app_configuration` - add support for `developer` tier to the `sku` property ([#29492](https://github.com/hashicorp/terraform-provider-azurerm/issues/29492))
* `azurerm_app_configuration` - the `sku` property can now be downgraded from `premium` to `standard` without recreating the resource ([#29492](https://github.com/hashicorp/terraform-provider-azurerm/issues/29492))
* `azurerm_key_vault_managed_hardware_security_module_key` - add support for the `import` value in the `key_opts` property ([#29524](https://github.com/hashicorp/terraform-provider-azurerm/issues/29524))
* `azurerm_netapp_pool` - add support for `cool_access_enabled` ([#29468](https://github.com/hashicorp/terraform-provider-azurerm/issues/29468))
* `azurerm_network_manager_deployment` - add support for the `Routing` value in the `scope_access` property ([#29536](https://github.com/hashicorp/terraform-provider-azurerm/issues/29536))
* `azurerm_private_endpoint_application_security_group_association` - resource is now removed from state if it no longer exist ([#29601](https://github.com/hashicorp/terraform-provider-azurerm/issues/29601))
* `azurerm_virtual_machine_implicit_data_disk_from_source` - the `disk_size_gb` property can now be increased without recreating the resource ([#29239](https://github.com/hashicorp/terraform-provider-azurerm/issues/29239))
* `azurerm_web_application_firewall_policy` - add support for the `JSChallenge` in the `action` property ([#29614](https://github.com/hashicorp/terraform-provider-azurerm/issues/29614))

BUG FIXES:

* `azurerm_api_management_api` - no longer returns an error on the  `oauth2_authorization` and `openid_authentication` properties when updating ([#29042](https://github.com/hashicorp/terraform-provider-azurerm/issues/29042))
* `azurerm_route_map` - the validation for the `name` now allows numbers ([#29519](https://github.com/hashicorp/terraform-provider-azurerm/issues/29519))

## 4.28.0 (May 09, 2025)

FEATURES:

* **New Resource**: `azurerm_nginx_api_key` ([#28919](https://github.com/hashicorp/terraform-provider-azurerm/issues/28919))
* **New Data Source**: `azurerm_nginx_api_key` ([#28919](https://github.com/hashicorp/terraform-provider-azurerm/issues/28919))

ENHANCEMENTS:

* dependencies: `azurerm_mssql_database` - Update to API version `2023-08-01-preview/replicationlinks` ([#28705](https://github.com/hashicorp/terraform-provider-azurerm/issues/28705))
* dependencies: `azurerm_mssql_server_security_alert_policy` - update to API version `2023-08-01-preview/serversecurityalertpolicies` ([#29363](https://github.com/hashicorp/terraform-provider-azurerm/issues/29363))
* dependencies: `eventhub` - update to API version `2024-01-01` ([#29397](https://github.com/hashicorp/terraform-provider-azurerm/issues/29397))
* dependencies: `azurerm_shared_image_version` - update to API version `2024-03-01` ([#28954](https://github.com/hashicorp/terraform-provider-azurerm/issues/28954))
* `azurerm_ai_foundry_project` - add support for the `primary_user_assigned_identity` property ([#29197](https://github.com/hashicorp/terraform-provider-azurerm/issues/29197))
* `azurerm_storage_account_static_website` - the `index_document` property now has validation for length and excluding slashes ([#29431](https://github.com/hashicorp/terraform-provider-azurerm/issues/29431))

BUG FIXES:

* `azurerm_application_insights` - the `workspace_id` is now `Computed` ([#29396](https://github.com/hashicorp/terraform-provider-azurerm/issues/29396))
* `azurerm_batch_pool` - prevent error when `certificate` is not used ([#29443](https://github.com/hashicorp/terraform-provider-azurerm/issues/29443))
* `azurerm_nginx_deployment` - add support for the `web_application_firewall` property ([#27454](https://github.com/hashicorp/terraform-provider-azurerm/issues/27454))
* `azurerm_postgresql_flexible_server_virtual_endpoint` - is no longer removed from state when a fail-over occurs ([#29424](https://github.com/hashicorp/terraform-provider-azurerm/issues/29424))
* `azurerm_servicebus_queue` - no longer waits on resource creation ([#29435](https://github.com/hashicorp/terraform-provider-azurerm/issues/29435))
* `azurerm_virtual_network_gateway` - prevent a panic when `vpn_client_configuration` is removed from from the configuration ([#29456](https://github.com/hashicorp/terraform-provider-azurerm/issues/29456))
* `azurerm_web_pubsub_custom_certificate` - no longer crashes when `custom_certificate_id` is in a different subscription ([#29410](https://github.com/hashicorp/terraform-provider-azurerm/issues/29410))
* `azurerm_windows_web_app` - fix perpetual diff around incorrect default for `always_on` and ignore default values for `logs.0.application_logs` ([#29150](https://github.com/hashicorp/terraform-provider-azurerm/issues/29150))
* `azurerm_windows_web_app_slot` - fix perpetual diff around incorrect default for `always_on` and ignore default values for `logs.0.application_logs` ([#29150](https://github.com/hashicorp/terraform-provider-azurerm/issues/29150))

## 4.27.0 (April 25, 2025)

FEATURES:

* **New Resource**: `azurerm_eventgrid_partner_configuration` ([#28676](https://github.com/hashicorp/terraform-provider-azurerm/issues/28676))

ENHANCEMENTS:

* dependencies: update `go-azure-sdk` to `v0.20250409.1192141` ([#29307](https://github.com/hashicorp/terraform-provider-azurerm/issues/29307))
* dependencies: `containerapps` - update to API version  `2025-01-01` ([#29296](https://github.com/hashicorp/terraform-provider-azurerm/issues/29296))
* dependencies: `netapp` - update to API version `2025-01-01` ([#29382](https://github.com/hashicorp/terraform-provider-azurerm/issues/29382))
* dependencies: `operationalinsights` - partial update to API version `2023-09-01` ([#29283](https://github.com/hashicorp/terraform-provider-azurerm/issues/29283))
* `azurerm_cdn_frontdoor_origin` - support `managedEnvironments` value for `private_link.target_type` ([#28239](https://github.com/hashicorp/terraform-provider-azurerm/issues/28239))
* `azurerm_cdn_frontdoor_origin` - add support for the `web_secondary` `Gateway` values in the `private_link.target_type` property ([#29380](https://github.com/hashicorp/terraform-provider-azurerm/issues/29380))
* `azurerm_cognitive_deployment` - add support for the `Cohere` value in the `model.format` property ([#29143](https://github.com/hashicorp/terraform-provider-azurerm/issues/29143))
* `azurerm_container_app_environment`: add support for cross subscription `log_analytics_workspace_id` ([#28740](https://github.com/hashicorp/terraform-provider-azurerm/issues/28740))
* `azurerm_dev_center_project` - add support for the `identity` property ([#29278](https://github.com/hashicorp/terraform-provider-azurerm/issues/29278))
* `azurerm_dynatrace_tag_rules` - the `log_rule` and `metric_rule` blocks and their properties are no longer `ForceNew` ([#29298](https://github.com/hashicorp/terraform-provider-azurerm/issues/29298))
* `azurerm_monitor_data_collection_endpoint` - add support for the `metrics_ingestion_endpoint` attribute ([#29292](https://github.com/hashicorp/terraform-provider-azurerm/issues/29292))
* `azurerm_mysql_flexible_server` - support for the `log_on_disk_enabled` property ([#28929](https://github.com/hashicorp/terraform-provider-azurerm/issues/28929))
* `azurerm_subnet` - add support for the `Microsoft.PowerAutomate/hostedRpa` value in the `delegation.service_delegation.name` property ([#29271](https://github.com/hashicorp/terraform-provider-azurerm/issues/29271))
* `azurerm_subnet` - add support for the `Microsoft.Network/applicationGateways` value in the `delegation.service_delegation.name` property ([#29361](https://github.com/hashicorp/terraform-provider-azurerm/issues/29361))
* `azurerm_virtual_network` - add support for the `Microsoft.PowerAutomate/hostedRpa` value in the `subnet.delegation.service_delegation.name` property ([#29271](https://github.com/hashicorp/terraform-provider-azurerm/issues/29271))
* `azurerm_virtual_network` - add support for the `Microsoft.Network/applicationGateways` value in the `subnet.delegation.service_delegation.name` property ([#29361](https://github.com/hashicorp/terraform-provider-azurerm/issues/29361))

BUG FIXES:

* provider: ensure `x-ms-correlation-request-id` header is only set once during list operations ([#28974](https://github.com/hashicorp/terraform-provider-azurerm/issues/28974))
* `azurerm_app_configuration_feature` - suppress casing differences for `configuration_store_id` to prevent resource recreation ([#29285](https://github.com/hashicorp/terraform-provider-azurerm/issues/29285))
* `azurerm_app_configuration_key` - suppress casing differences for `configuration_store_id` to prevent resource recreation ([#29285](https://github.com/hashicorp/terraform-provider-azurerm/issues/29285))
* `azurerm_container_app_environment` - updates are now made using the `PATCH` method, preventing errors due to missing properties in the request ([#29317](https://github.com/hashicorp/terraform-provider-azurerm/issues/29317))
* `azurerm_eventhub_namespace` - remove max items from network/ip rules as they can be increased above upon request ([#29333](https://github.com/hashicorp/terraform-provider-azurerm/issues/29333))
* `azurerm_kusto_iothub_data_connection` - update `event_system_properties` validation and documentation to be more flexible ([#29314](https://github.com/hashicorp/terraform-provider-azurerm/issues/29314))
* `azurerm_linux_web_app` - correctly read `backup.schedule.start_time` into state ([#29254](https://github.com/hashicorp/terraform-provider-azurerm/issues/29254))
* `azurerm_netapp_volume` - update validation for `storage_quota_in_gb` to allow values from `50` to `102400` ([#29341](https://github.com/hashicorp/terraform-provider-azurerm/issues/29341))
* `azurerm_postgresql_flexible_server` - downgrading `version` forces a new resource to be created ([#28559](https://github.com/hashicorp/terraform-provider-azurerm/issues/28559))
* `azurerm_postgresql_flexible_server` - downgrading `storage_mb` forces a new resource to be created ([#29309](https://github.com/hashicorp/terraform-provider-azurerm/issues/29309))
* `azurerm_private_endpoint` - `private_dns_zone_group.private_dns_zone_ids` can now be updated correctly ([#29329](https://github.com/hashicorp/terraform-provider-azurerm/issues/29329))
* `azurerm_search_shared_private_link_service` - add locks to prevent conflicts when creating multiple instances ([#29294](https://github.com/hashicorp/terraform-provider-azurerm/issues/29294))

## 4.26.0 (April 04, 2025)

BREAKING CHANGES:

* feature: The Provider `feature` configuration item `virtual_machines.graceful_shutdown` is now not used due to a breaking change in the `compute` API. This feature block setting is now deprecated and ignored if set and will be removed in v5.0 of the provider. ([#29185](https://github.com/hashicorp/terraform-provider-azurerm/issues/29185))
* `azurerm_linux_virtual_machine` - the `vm_agent_platform_updates_enabled` property is now read-only due to a recent API breaking change ([#29211](https://github.com/hashicorp/terraform-provider-azurerm/issues/29211))
* `azurerm_windows_virtual_machine` - the `vm_agent_platform_updates_enabled` property is now read-only due to a recent API breaking change ([#29211](https://github.com/hashicorp/terraform-provider-azurerm/issues/29211))

FEATURES:

* **New Data Source**: `azurerm_role_assignments` ([#29214](https://github.com/hashicorp/terraform-provider-azurerm/issues/29214))

ENHANCEMENTS:

* dependencies: `azurerm_sentinel_automation_rule` - update to API version `2024-09-01` ([#29240](https://github.com/hashicorp/terraform-provider-azurerm/issues/29240))
* dependencies: `devcenter` - update to API version `2025-02-01` ([#29240](https://github.com/hashicorp/terraform-provider-azurerm/issues/29240))
* dependencies: `recoveryservices` - partial update to API version `2024-10-01` ([#29240](https://github.com/hashicorp/terraform-provider-azurerm/issues/29240))
* Data Source: `azurerm_mssql_server` - export the `express_vulnerability_assessment_enabled` property ([#29168](https://github.com/hashicorp/terraform-provider-azurerm/issues/29168))
* `azurerm_dashboard_grafana` - `grafana_major_version` is no longer ForceNew ([#29212](https://github.com/hashicorp/terraform-provider-azurerm/issues/29212))
* `azurerm_data_factory_linked_service_sftp` - add support for SSH authentication and Key Vault secret references ([#28690](https://github.com/hashicorp/terraform-provider-azurerm/issues/28690))
* `azurerm_databricks_workspace` - resources using managed resource groups that contain UC can now be deleted with the `force_delete` Provider Feature flag ([#29095](https://github.com/hashicorp/terraform-provider-azurerm/issues/29095))
* `azurerm_mssql_server` - add support for the `express_vulnerability_assessment_enabled` property ([#29168](https://github.com/hashicorp/terraform-provider-azurerm/issues/29168))
* `azurerm_mysql_flexible_server` - deprecate `public_network_access_enabled` in favor of `public_network_access` ([#28890](https://github.com/hashicorp/terraform-provider-azurerm/issues/28890))
* `azurerm_netapp_volume` - `service_level` can now be updated ([#29209](https://github.com/hashicorp/terraform-provider-azurerm/issues/29209))
* `azurerm_nginx_deployment` - `frontend_public`, `frontend_private`, and `network_interface` are no longer `ForceNew` ([#28577](https://github.com/hashicorp/terraform-provider-azurerm/issues/28577))
* `azurerm_orchestrated_virtual_machine_scale_set` - add support for the `upgrade_mode` and `rolling_upgrade_policy` properties ([#28354](https://github.com/hashicorp/terraform-provider-azurerm/issues/28354))
* `azurerm_static_webapp` - mark `app_settings` sensitive in schema ([#28689](https://github.com/hashicorp/terraform-provider-azurerm/issues/28689))

BUG FIXES:

* `azurerm_linux_virtual_machine` - `license_type` can now be updated to None ([#28786](https://github.com/hashicorp/terraform-provider-azurerm/issues/28786))
* `azurerm_mysql_flexible_server` - prevent a panic when `customer_managed_key` is nil ([#29225](https://github.com/hashicorp/terraform-provider-azurerm/issues/29225))
* `azurerm_traffic_manager_nested_endpoint` - remove `Computed` from `priority` property as these are assigned dynamically by the API ([#29217](https://github.com/hashicorp/terraform-provider-azurerm/issues/29217))

## 4.25.0 (March 28, 2025)

ENHANCEMENTS:
 
* dependencies: `go-azure-helpers` - update  to `0.72.0` ([#29206](https://github.com/hashicorp/terraform-provider-azurerm/issues/29206))
* dependencies: `redisenterprise` - update to API version `2024-10-01` ([#29073](https://github.com/hashicorp/terraform-provider-azurerm/issues/29073))
* dependencies: `servicefabricmanaged` - update to API version `2024-04-01` ([#29199](https://github.com/hashicorp/terraform-provider-azurerm/issues/29199))
* Data Source: `azurerm_virtual_hub_connection` - add support for the `static_vnet_propagate_static_routes` property ([#28560](https://github.com/hashicorp/terraform-provider-azurerm/issues/28560))
* `azurerm_cosmosdb_account` - add support for the `DeleteAllItemsByPartitionKey` value in the `capabilities` property ([#29126](https://github.com/hashicorp/terraform-provider-azurerm/issues/29126))
* `azurerm_hdinsight_spark_cluster_resource` - add support for the `zones` property ([#28149](https://github.com/hashicorp/terraform-provider-azurerm/issues/28149))
* `azurerm_linux_function_app` - add support for Python version `3.13` ([#29131](https://github.com/hashicorp/terraform-provider-azurerm/issues/29131))
* `azurerm_linux_function_app_slot` - add support for Python version `3.13` ([#29131](https://github.com/hashicorp/terraform-provider-azurerm/issues/29131))
* `azurerm_linux_web_app` - add support for Python version `3.13` ([#29131](https://github.com/hashicorp/terraform-provider-azurerm/issues/29131))
* `azurerm_linux_web_app_slot` - add support for Python version `3.13` ([#29131](https://github.com/hashicorp/terraform-provider-azurerm/issues/29131))
* `azurerm_log_analytics_workspace` - add support for the `LACluster` SKU ([#29137](https://github.com/hashicorp/terraform-provider-azurerm/issues/29137))
* `azurerm_managed_disk` - allow disk expansion without downtime for all `storage_account_type` ([#28730](https://github.com/hashicorp/terraform-provider-azurerm/issues/28730))
* `azurerm_mssql_job_agent` - add support for the `identity` and `sku` properties ([#29090](https://github.com/hashicorp/terraform-provider-azurerm/issues/29090))
* `azurerm_network_manager` - `scope_accesses` is now optional ([#28781](https://github.com/hashicorp/terraform-provider-azurerm/issues/28781))
* `azurerm_oracle_cloud_vm_cluster` - add support for the `system_version` property ([#29093](https://github.com/hashicorp/terraform-provider-azurerm/issues/29093))
* `azurerm_powerbi_embedded `- add support for `A7` and `A8` values for `sku_name` ([#29153](https://github.com/hashicorp/terraform-provider-azurerm/issues/29153))
* `azurerm_virtual_hub_connection` - add support for the `static_vnet_propagate_static_routes` property ([#28560](https://github.com/hashicorp/terraform-provider-azurerm/issues/28560))

BUG FIXES

* Data source: `azurerm_container_app_environment` - prevent an error when the log analytics workspace is in a different subscription ([#28647](https://github.com/hashicorp/terraform-provider-azurerm/issues/28647))
* `azurerm_kubernetes_cluster_node_pool` - fix issue where `kubelet_disk_type` couldn't be updated, updating this will now rotate the node pool ([#29135](https://github.com/hashicorp/terraform-provider-azurerm/issues/29135))
* `azurerm_linux_virtual_machine` - fix issue where a user assigned identity couldn't be removed from the resource ([#29157](https://github.com/hashicorp/terraform-provider-azurerm/issues/29157))
* `azurerm_linux_virtual_machine_scale_set` - fix issue where a user assigned identity couldn't be removed from the resource ([#29157](https://github.com/hashicorp/terraform-provider-azurerm/issues/29157))
* `azurerm_log_analytics_workspace` - prevent an error when the workspace is in a soft-deleted state and linked to a log analytics cluster ([#29137](https://github.com/hashicorp/terraform-provider-azurerm/issues/29137))
* `azurerm_postgresql_flexible_server_virtual_endpoint` - add a lock on the replica server to prevent a race condition ([#29071](https://github.com/hashicorp/terraform-provider-azurerm/issues/29071))
* `azurerm_signalr_service` - set `location` in payload when updating to prevent an API error ([#29184](https://github.com/hashicorp/terraform-provider-azurerm/issues/29184))
* `azurerm_storage_account_queue_properties` - prevent a panic when the storage account is removed out of band ([#28371](https://github.com/hashicorp/terraform-provider-azurerm/issues/28371))
* `azurerm_storage_account_static_website` - prevent a panic when the storage account is removed out of band ([#28371](https://github.com/hashicorp/terraform-provider-azurerm/issues/28371))
* `azurerm_stream_analytics_job` - update validation to notify users if `content_storage_policy` hasn't been correctly set to setup `job_storage_account` ([#29158](https://github.com/hashicorp/terraform-provider-azurerm/issues/29158))

## 4.24.0 (March 21, 2025)

FEATURES:

* **New Resource**: `azurerm_servicebus_namespace_customer_managed_key` ([#28888](https://github.com/hashicorp/terraform-provider-azurerm/issues/28888))
* **New Resource**: `azurerm_stream_analytics_job_storage_account` ([#29113](https://github.com/hashicorp/terraform-provider-azurerm/issues/29113))
* **New Resource**: `azurerm_web_pubsub_socketio` ([#28992](https://github.com/hashicorp/terraform-provider-azurerm/issues/28992))

ENHANCEMENTS:

* dependencies: `hashicorp/go-azure-sdk` - update to `v0.20250314.1213156` ([#29081](https://github.com/hashicorp/terraform-provider-azurerm/issues/29081))
* dependencies: `loganalytics` - partial update to API version `2023-03-01` ([#28977](https://github.com/hashicorp/terraform-provider-azurerm/issues/28977))
* dependencies: `monitor` - partial update to API version `2023-03-01` ([#28977](https://github.com/hashicorp/terraform-provider-azurerm/issues/28977))
* dependencies: `postgresql` - partial update to API version `2024-08-01` ([#28964](https://github.com/hashicorp/terraform-provider-azurerm/issues/28964))
* Data Source: `azurerm_linux_function_app` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))
* Data Source: `azurerm_linux_web_app` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))
* Data Source: `azurerm_windows_function_app` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))
* `azurerm_ai_services` - add support for the `network_acls.bypass` property ([#28569](https://github.com/hashicorp/terraform-provider-azurerm/issues/28569))
* `azurerm_dashboard_grafana` - add support for `grafana_major_version` `11` ([#28884](https://github.com/hashicorp/terraform-provider-azurerm/issues/28884))
* `azurerm_kubernetes_cluster_node_pool` - remove call to retrieve the parent cluster in the read ([#29088](https://github.com/hashicorp/terraform-provider-azurerm/issues/29088))
* `azurerm_linux_function_app` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))
* `azurerm_linux_function_app` - set `pre_warmed_instance_count` on create ([#28739](https://github.com/hashicorp/terraform-provider-azurerm/issues/28739))
* `azurerm_linux_function_app_slot` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))
* `azurerm_linux_web_app` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))
* `azurerm_linux_web_app_slot` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))
* `azurerm_redis_cache` - tighten validation for `sku_name`, `family`, `capacity` ([#29079](https://github.com/hashicorp/terraform-provider-azurerm/issues/29079))
* `azurerm_windows_function_app` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))
* `azurerm_windows_function_app` - set `pre_warmed_instance_count` on create ([#28739](https://github.com/hashicorp/terraform-provider-azurerm/issues/28739))
* `azurerm_windows_function_app_slot` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))
* `azurerm_windows_web_app` - add support for node version `~22` ([#29082](https://github.com/hashicorp/terraform-provider-azurerm/issues/29082))
* `azurerm_windows_web_app` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))
* `azurerm_windows_web_app_slot` - add support for node version `~22` ([#29082](https://github.com/hashicorp/terraform-provider-azurerm/issues/29082))
* `azurerm_windows_web_app_slot` - add support for the `virtual_network_backup_restore_enabled` property ([#29012](https://github.com/hashicorp/terraform-provider-azurerm/issues/29012))

BUG FIXES:

* `azurerm_app_configuration` - the `encryption` block can now be removed ([#28173](https://github.com/hashicorp/terraform-provider-azurerm/issues/28173))
* `azurerm_cdn_frontdoor_origin_group` - `health_probe` no longer resets during update unless specified ([#29094](https://github.com/hashicorp/terraform-provider-azurerm/issues/29094))
* `azurerm_cognitive_account` - `customer_managed_key` can now be removed ([#28368](https://github.com/hashicorp/terraform-provider-azurerm/issues/28368))
* `azurerm_container_group` - `dns_name_label_reuse_policy` is now marked as ForceNew ([#29040](https://github.com/hashicorp/terraform-provider-azurerm/issues/29040))
* `azurerm_disk_encryption_set` - prevent crash when retrieving Key Vault details when updating ([#29018](https://github.com/hashicorp/terraform-provider-azurerm/issues/29018))
* `azurerm_express_route_circuit` - fix issue where `bandwidth_in_mbps` isn't updated correctly ([#28822](https://github.com/hashicorp/terraform-provider-azurerm/issues/28822))
* `azurerm_key_vault_secret` - revert CustomizeDiff logic to recreate the resource when `expiration_date` is removed ([#28920](https://github.com/hashicorp/terraform-provider-azurerm/issues/28920))
* `azurerm_kubernetes_cluster` - `fips_enabled` can be updated by cycling the default node pool ([#29096](https://github.com/hashicorp/terraform-provider-azurerm/issues/29096))
* `azurerm_monitor_diagnostic_setting` - the `enabled_log` block can now be removed ([#28485](https://github.com/hashicorp/terraform-provider-azurerm/issues/28485))
* `azurerm_mssql_database` - fix validation for `auto_pause_delay_in_minutes` ([#28670](https://github.com/hashicorp/terraform-provider-azurerm/issues/28670))
* `azurerm_mssql_server` - fix an issue where the provider would incorrectly error during plan operations if `administrator_login` or `administrator_login_password` were added to `lifecycle.ignore_changes` ([#29107](https://github.com/hashicorp/terraform-provider-azurerm/issues/29107))


## 4.23.0 (March 13, 2025)

NOTES:

* `azurerm_key_vault_secret` - resource now supports the `value_wo`[write-only argument](https://developer.hashicorp.com/terraform/language/v1.11.x/resources/ephemeral#write-only-arguments) ([#28947](https://github.com/hashicorp/terraform-provider-azurerm/issues/28947))

FEATURES:

* **New Resource**: `azurerm_network_manager_ipam_pool` ([#28695](https://github.com/hashicorp/terraform-provider-azurerm/issues/28695))

ENHANCEMENTS:

* dependencies: update `Go` version to `1.24.1` ([#28999](https://github.com/hashicorp/terraform-provider-azurerm/issues/28999))
* dependencies: `hashicorp/go-azure-sdk` - update to `v0.20250310.1130319` ([#29009](https://github.com/hashicorp/terraform-provider-azurerm/issues/29009))
* `azurerm_cognitive_deployment` - add support for `DataZoneBatch` in the `sku.name` property ([#28973](https://github.com/hashicorp/terraform-provider-azurerm/issues/28973))
* `azurerm_mongo_cluster` - add support for `M10`, `M20`, and `M200` compute tiers ([#29026](https://github.com/hashicorp/terraform-provider-azurerm/issues/29026))


BUG FIXES:

* `azurerm_linux_function_app` - fix validation for `site_config.application_stack.node_version` to allow `22` ([#28988](https://github.com/hashicorp/terraform-provider-azurerm/issues/28988))
* `azurerm_postgresql_flexible_server` - fix validation for `customer_managed_key.key_vault_key_id` and `customer_managed_key.geo_backup_key_id` to disallow versionless keys preventing unclear error messages ([#28981](https://github.com/hashicorp/terraform-provider-azurerm/issues/28981))
* `azurerm_web_pubsub_hub` - validation for the `auth.managed_identity_id` now supports token audience as a valid input ([#28495](https://github.com/hashicorp/terraform-provider-azurerm/issues/28495))

## 4.22.0 (March 07, 2025)

FEATURES:

* **New Data Source**: `azurerm_extended_location_custom_location` ([#28066](https://github.com/hashicorp/terraform-provider-azurerm/issues/28066))
* **New Resource**: `azurerm_system_center_virtual_machine_manager_virtual_machine_instance` ([#27622](https://github.com/hashicorp/terraform-provider-azurerm/issues/27622))

ENHANCEMENTS: 

* dependencies: `containers` - update API version to  `2024-09-01` ([#28598](https://github.com/hashicorp/terraform-provider-azurerm/issues/28598))
* dependencies: `hashicorp/go-azure-sdk` - update to `v0.20250227.1125644` ([#28902](https://github.com/hashicorp/terraform-provider-azurerm/issues/28902))
* dependencies: `signalr` - update API version to `2024-03-01` ([#28940](https://github.com/hashicorp/terraform-provider-azurerm/issues/28940))
* Data Source: `azurerm_container_app` - add support for the `template.volume.mount_options` property ([#28619](https://github.com/hashicorp/terraform-provider-azurerm/issues/28619))
* Data Source: `azurerm_storage_account_queue_properties` - now gets the parent account directly rather than searching the list of all accounts when the Resource Manager ID is available ([#28617](https://github.com/hashicorp/terraform-provider-azurerm/issues/28617))
* Data Source: `azurerm_storage_account_static_website` - now gets the parent account directly rather than searching the list of all accounts when the Resource Manager ID is available ([#28617](https://github.com/hashicorp/terraform-provider-azurerm/issues/28617))
* Data Source: `azurerm_storage_containers` - now gets the parent account directly rather than searching the list of all accounts when the Resource Manager ID is available ([#28617](https://github.com/hashicorp/terraform-provider-azurerm/issues/28617))
* `azurerm_api_connection` - `display_name` and `parameter_values` are no longer `ForceNew` ([#28721](https://github.com/hashicorp/terraform-provider-azurerm/issues/28721))
* `azurerm_cdn_frontdoor_firewall_policy` - add support for the `log_scrubbing` properties ([#28834](https://github.com/hashicorp/terraform-provider-azurerm/issues/28834))
* `azurerm_container_app` - add support for the `template.volume.mount_options` property ([#28619](https://github.com/hashicorp/terraform-provider-azurerm/issues/28619))
* `azurerm_container_app_job` - add support for the `template.volume.mount_options` property ([#28619](https://github.com/hashicorp/terraform-provider-azurerm/issues/28619))
* `azurerm_extended_custom_location` - deprecated in favour of `azurerm_extended_location_custom_location` ([#28066](https://github.com/hashicorp/terraform-provider-azurerm/issues/28066))
* `azurerm_mongo_cluster` - add support for the `connection_strings` attribute ([#28880](https://github.com/hashicorp/terraform-provider-azurerm/issues/28880))
* `azurerm_storage_account` - now gets the parent account directly rather than searching the list of all accounts when the Resource Manager ID is available ([#28617](https://github.com/hashicorp/terraform-provider-azurerm/issues/28617))
* `azurerm_storage_account_queue_properties` - now gets the parent account directly rather than searching the list of all accounts when the Resource Manager ID is available ([#28617](https://github.com/hashicorp/terraform-provider-azurerm/issues/28617))
* `azurerm_storage_account_static_website` - now gets the parent account directly rather than searching the list of all accounts when the Resource Manager ID is available ([#28617](https://github.com/hashicorp/terraform-provider-azurerm/issues/28617))
* `azurerm_workloads_sap_discovery_virtual_instance` - add support for the `managed_resources_network_access_type` property ([#28881](https://github.com/hashicorp/terraform-provider-azurerm/issues/28881))
* `azurerm_workloads_sap_single_node_virtual_instance` - add support for the `managed_resources_network_access_type` property ([#28881](https://github.com/hashicorp/terraform-provider-azurerm/issues/28881))
* `azurerm_workloads_sap_three_tier_virtual_instance` - add support for the `managed_resources_network_access_type` property ([#28881](https://github.com/hashicorp/terraform-provider-azurerm/issues/28881))

BUG FIXES:

* `azurerm_api_management_api` - split create/update methods ([#28271](https://github.com/hashicorp/terraform-provider-azurerm/issues/28271))
* `azurerm_express_route_circuit` - `allow_classic_operations` is now set when resource is created ([#28748](https://github.com/hashicorp/terraform-provider-azurerm/issues/28748))
* `azurerm_key_vault_certificate` - set partial when updating key vault certificate ([#28848](https://github.com/hashicorp/terraform-provider-azurerm/issues/28848))
* `azurerm_managed_disk` - always set `network_access_policy` into state to allow Terraform to detect drift ([#28934](https://github.com/hashicorp/terraform-provider-azurerm/issues/28934))
* `azurerm_mssql_managed_instance` - fix an issue that prevented using values only known during apply for `administrator_login_password` ([#28843](https://github.com/hashicorp/terraform-provider-azurerm/issues/28843))
* `azurerm_mssql_server` - prevent panic by removing function call on a value that may be unknown ([#28949](https://github.com/hashicorp/terraform-provider-azurerm/issues/28949))


## 4.21.1 (February 28, 2025)

BUG FIXES:

* `azurerm_mssql_server` - prevent panic by checking if `administrator_login` exists in the raw config map ([#28909](https://github.com/hashicorp/terraform-provider-azurerm/issues/28909))


## 4.21.0 (February 27, 2025)

NOTES:

* The `azurerm_mssql_job_credential` resource now supports the `password_wo` [write-only argument](https://developer.hashicorp.com/terraform/language/v1.11.x/resources/ephemeral#write-only-arguments)
* The `azurerm_mssql_server` resource now supports the `administrator_login_password_wo` [write-only argument](https://developer.hashicorp.com/terraform/language/v1.11.x/resources/ephemeral#write-only-arguments)
* The `azurerm_mysql_flexible_server` resource now supports the `administrator_password_wo` [write-only argument](https://developer.hashicorp.com/terraform/language/v1.11.x/resources/ephemeral#write-only-arguments)
* The `azurerm_postgresql_flexible_server` resource now supports the `administrator_password_wo` [write-only argument](https://developer.hashicorp.com/terraform/language/v1.11.x/resources/ephemeral#write-only-arguments)
* The `azurerm_postgresql_server` resource now supports the `administrator_login_password_wo` [write-only argument](https://developer.hashicorp.com/terraform/language/v1.11.x/resources/ephemeral#write-only-arguments)

FEATURES:

* **New Resource**: `azurerm_linux_function_app_flex_consumption` ([#28199](https://github.com/hashicorp/terraform-provider-azurerm/issues/28199))
* **New Resource**: `azurerm_network_manager_verifier_workspace` ([#28754](https://github.com/hashicorp/terraform-provider-azurerm/issues/28754))

ENHANCEMENTS:

* dependencies: `azurerm_kubernetes_cluster_trusted_access_role_binding` - update API version to `2024-05-01` ([#28853](https://github.com/hashicorp/terraform-provider-azurerm/issues/28853))
* dependencies: `desktopvirtualization` - update API version to `2024-04-03` ([#28771](https://github.com/hashicorp/terraform-provider-azurerm/issues/28771))
* dependencies: `kusto` - update API version to `2024-04-13` ([#28685](https://github.com/hashicorp/terraform-provider-azurerm/issues/28685))
* dependencies: `redis` - update API version to `2024-11-01` ([#28696](https://github.com/hashicorp/terraform-provider-azurerm/issues/28696))
* dependencies: `workloads` - update API version to `2024-09-01` ([#28825](https://github.com/hashicorp/terraform-provider-azurerm/issues/28825))
* `azurerm_fluid_relay_server` - fix `versionless_id` support for `key_vault_key_id` ([#28864](https://github.com/hashicorp/terraform-provider-azurerm/issues/28864))
* `azurerm_kubernetes_cluster` - add support for the `upgrade_override_setting` property ([#27962](https://github.com/hashicorp/terraform-provider-azurerm/issues/27962))
* `azurerm_kusto_cluster_principal_assignment` - add support for `AllDatabaseMonitor` role type ([#28685](https://github.com/hashicorp/terraform-provider-azurerm/issues/28685))
* `azurerm_linux_function_app` - correctly update `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` when changed in `app_settings` ([#28859](https://github.com/hashicorp/terraform-provider-azurerm/issues/28859))
* `azurerm_linux_function_app_slot` - correctly update `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` when changed in `app_settings` ([#28859](https://github.com/hashicorp/terraform-provider-azurerm/issues/28859))
* `azurerm_linux_web_app` - add support for Node Version `22` ([#28840](https://github.com/hashicorp/terraform-provider-azurerm/issues/28840))
* `azurerm_linux_web_app_slot` - add support for Node Version `22` ([#28840](https://github.com/hashicorp/terraform-provider-azurerm/issues/28840))
* `azurerm_logic_app_standard` - add support for the `vnet_content_share_enabled` property ([#28879](https://github.com/hashicorp/terraform-provider-azurerm/issues/28879))
* `azurerm_mssql_job_credential` - add support for the `password_wo` and `password_wo_version` properties ([#28808](https://github.com/hashicorp/terraform-provider-azurerm/issues/28808))
* `azurerm_mssql_managed_instance` - add support for the `database_format` and `hybrid_secondary_usage` properties ([#28248](https://github.com/hashicorp/terraform-provider-azurerm/issues/28248))
* `azurerm_mssql_server` - add support for the `administrator_login_password_wo` and `administrator_login_password_wo_version` properties ([#28818](https://github.com/hashicorp/terraform-provider-azurerm/issues/28818))
* `azurerm_mysql_flexible_server` - add support for the `administrator_password_wo` and `administrator_password_wo_version` properties ([#28799](https://github.com/hashicorp/terraform-provider-azurerm/issues/28799))
* `azurerm_postgresql_flexible_server` - add support for the `administrator_password_wo` and `administrator_password_wo_version` properties ([#28857](https://github.com/hashicorp/terraform-provider-azurerm/issues/28857))
* `azurerm_postgresql_server` - add support for the `administrator_login_password_wo` and `administrator_login_password_wo_version` properties ([#28856](https://github.com/hashicorp/terraform-provider-azurerm/issues/28856))
* `azurerm_service_plan` - add support for the `I1mv2`, `I2mv2`, `I3mv2`, `I4mv2`, `I5mv2` skus ([#28316](https://github.com/hashicorp/terraform-provider-azurerm/issues/28316))
* `azurerm_servicebus_namespace` - split create/update functions ([#28539](https://github.com/hashicorp/terraform-provider-azurerm/issues/28539))
* `azurerm_storage_account` - nested attributes in `immutability_policy` can now be updated ([#28122](https://github.com/hashicorp/terraform-provider-azurerm/issues/28122))
* `azurerm_windows_function_app` - correctly update `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` when changed in `app_settings` ([#28859](https://github.com/hashicorp/terraform-provider-azurerm/issues/28859))
* `azurerm_windows_function_app_slot` - correctly update `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` when changed in `app_settings` ([#28859](https://github.com/hashicorp/terraform-provider-azurerm/issues/28859))


BUG FIXES:

* `azurerm_key_vault_secret` - recreate the resource if `expiration_date` is removed after having been set ([#28494](https://github.com/hashicorp/terraform-provider-azurerm/issues/28494))
* `azurerm_log_analytics_cluster_customer_managed_key` - fix error due to read-only property included in request payload during create/update/delete operations ([#28862](https://github.com/hashicorp/terraform-provider-azurerm/issues/28862))
* `azurerm_log_analytics_cluster_customer_managed_key` - remove resource from state when deleted outside of Terraform ([#28862](https://github.com/hashicorp/terraform-provider-azurerm/issues/28862))
* `azurerm_log_analytics_cluster_customer_managed_key` - fix resource delete function ([#28862](https://github.com/hashicorp/terraform-provider-azurerm/issues/28862))
* `azurerm_security_center_pricing` - updating `subplan` now recreates the resource to work around API behaviour that enables certain settings on updated ([#27805](https://github.com/hashicorp/terraform-provider-azurerm/issues/27805))
* `azurerm_windows_web_app` - fix change detection for `tomcat_version` ([#28842](https://github.com/hashicorp/terraform-provider-azurerm/issues/28842))

## 4.20.0 (February 20, 2025)

FEATURES:

* **New Data Source**: `azurerm_dynatrace_monitor` ([#28381](https://github.com/hashicorp/terraform-provider-azurerm/issues/28381))
* **New Resource**: `azurerm_data_protection_backup_vault_customer_managed_key` ([#28679](https://github.com/hashicorp/terraform-provider-azurerm/issues/28679))

ENHANCEMENTS:

* dependencies: `hashicorp/terraform-plugin-sdk/v2` - update to `v2.36.0` ([#28788](https://github.com/hashicorp/terraform-provider-azurerm/issues/28788))
* dependencies: `azurerm_data_factory_pipeline` - update to use `hashicorp/go-azure-sdk` ([#28768](https://github.com/hashicorp/terraform-provider-azurerm/issues/28768))
* Data Source: `azurerm_logic_app_standard` - add support for the `ftp_publish_basic_authentication_enabled` and `scm_publish_basic_authentication_enabled` properties ([#28763](https://github.com/hashicorp/terraform-provider-azurerm/issues/28763))
* `azurerm_logic_app_standard` - add support for the `ftp_publish_basic_authentication_enabled` and `scm_publish_basic_authentication_enabled` properties ([#28763](https://github.com/hashicorp/terraform-provider-azurerm/issues/28763))
* `azurerm_pim_active_role_assignment` - add support for Azure RBAC conditions ([#27947](https://github.com/hashicorp/terraform-provider-azurerm/issues/27947))
* `azurerm_storage_container` - add support for migrating from deprecated `storage_account_name` to  `storage_account_id` ([#28784](https://github.com/hashicorp/terraform-provider-azurerm/issues/28784))
* `azurerm_storage_share` - add support for migrating from deprecated `storage_account_name` to  `storage_account_id` ([#28784](https://github.com/hashicorp/terraform-provider-azurerm/issues/28784))
* `azurerm_storage_table` - add attribute `resource_manager_id` ([#28809](https://github.com/hashicorp/terraform-provider-azurerm/issues/28809))
* `azurerm_windows_function_app` - add support for node `~22` ([#28815](https://github.com/hashicorp/terraform-provider-azurerm/issues/28815))
* `azurerm_windows_function_app_slot` - add support for node `~22` ([#28815](https://github.com/hashicorp/terraform-provider-azurerm/issues/28815))

BUG FIXES:

* Data Source: `azurerm_container_app` - add missing `ingress.client_certificate_mode` property that caused an error when retrieving data ([#28793](https://github.com/hashicorp/terraform-provider-azurerm/issues/28793))
* `azurerm_data_factory_pipeline` - fix error when unmarshaling the headers for a web activity ([#28768](https://github.com/hashicorp/terraform-provider-azurerm/issues/28768))
* `azurerm_mssql_virtual_machine` - fix an issue that prevented users from using values only known during apply as the value for `auto_backup.encryption_password` ([#28223](https://github.com/hashicorp/terraform-provider-azurerm/issues/28223))

## 4.19.0 (February 14, 2025)

FEATURES:

* **New Data Source**: `azurerm_stack_hci_storage_path` ([#28602](https://github.com/hashicorp/terraform-provider-azurerm/issues/28602))
* **New Resource**: `azurerm_ai_foundry` ([#27424](https://github.com/hashicorp/terraform-provider-azurerm/issues/27424))
* **New Resource**: `azurerm_ai_foundry_project` ([#27424](https://github.com/hashicorp/terraform-provider-azurerm/issues/27424))
* **New Resource**: `azurerm_mssql_job_step` ([#28691](https://github.com/hashicorp/terraform-provider-azurerm/issues/28691))
* **New Resource**: `azurerm_netapp_volume_group_oracle` ([#28391](https://github.com/hashicorp/terraform-provider-azurerm/issues/28391))
* **New Resource**: `azurerm_virtual_machine_scale_set_standby_pool` ([#28441](https://github.com/hashicorp/terraform-provider-azurerm/issues/28441))

ENHANCEMENTS:

* dependencies: `hashicorp/go-azure-sdk` update to `v0.20250213.1092825` ([#28767](https://github.com/hashicorp/terraform-provider-azurerm/issues/28767))
* dependencies: `sentinel` partial update to `2023-12-01-preview` ([#28195](https://github.com/hashicorp/terraform-provider-azurerm/issues/28195))
* Data Source: `azurerm_app_configuration` - add support for the `data_plane_proxy_authentication_mode` and `data_plane_proxy_private_link_delegation_enabled` properties ([#28712](https://github.com/hashicorp/terraform-provider-azurerm/issues/28712))
* `azurerm_app_configuration` - add support for the `data_plane_proxy_authentication_mode` and `data_plane_proxy_private_link_delegation_enabled` properties ([#28712](https://github.com/hashicorp/terraform-provider-azurerm/issues/28712))
* `azurerm_container_app` - add support for the `client_certificate_mode` property ([#28523](https://github.com/hashicorp/terraform-provider-azurerm/issues/28523))
* `azurerm_cdn_frontdoor_firewall_policy` - add support for `JSChallenge` for `custom` rules ([#28717](https://github.com/hashicorp/terraform-provider-azurerm/issues/28717))
* `azurerm_express_route_circuit` - add support for the `rate_limiting_enabled` property ([#28659](https://github.com/hashicorp/terraform-provider-azurerm/issues/28659))
* `azurerm_mssql_managed_instance_failover_group` - add support for `secondary_type` ([#28633](https://github.com/hashicorp/terraform-provider-azurerm/issues/28633))
* `azurerm_sentinal_alert_rule_scheduled` - increase combined limit of `entity_mapping` and `sentinal_entity_mapping` to 10 ([#28195](https://github.com/hashicorp/terraform-provider-azurerm/issues/28195))
* `azurerm_service_plan` - support for `premium_plan_auto_scale_enabled` ([#28524](https://github.com/hashicorp/terraform-provider-azurerm/issues/28524))

BUG FIXES:

* `azurerm_cdn_frontdoor_firewall_policy` - fixed issue where the `js_challenge_cookie_expiration_in_minutes` policies `default` value caused `Standard_AzureFrontDoor` skus to receive a `BadRequest` error ([#28726](https://github.com/hashicorp/terraform-provider-azurerm/issues/28726))
* `azurerm_servicebus_topic` - prevent perma diff when provisioning a partitioned topic within a non-partitioned namespace ([#26680](https://github.com/hashicorp/terraform-provider-azurerm/issues/26680))
* `azurerm_linux_function_app` - will no longer plan when `site_config.0.cors` is the default value ([#28703](https://github.com/hashicorp/terraform-provider-azurerm/issues/28703))
* `azurerm_linux_function_app_slot` - fix issue where `site_config.0.elastic_instance_minimum` was not being set ([#28725](https://github.com/hashicorp/terraform-provider-azurerm/issues/28725))
* `azurerm_linux_web_app` - will no longer plan when `site_config.0.cors` is the default value ([#28703](https://github.com/hashicorp/terraform-provider-azurerm/issues/28703))
* `azurerm_postgresql_flexible_server_virtual_endpoint` - allow `source_server_id` and `replica_server_id` to reference the same server ([#28733](https://github.com/hashicorp/terraform-provider-azurerm/issues/28733))
* `azurerm_windows_function_app` - will no longer plan when `site_config.0.cors` is the default value ([#28703](https://github.com/hashicorp/terraform-provider-azurerm/issues/28703))
* `azurerm_windows_function_app_slot` - fix issue where `site_config.0.elastic_instance_minimum` was not being set ([#28725](https://github.com/hashicorp/terraform-provider-azurerm/issues/28725))
* `azurerm_windows_web_app` - will no longer plan when `site_config.0.cors` is the default value ([#28703](https://github.com/hashicorp/terraform-provider-azurerm/issues/28703))

## 4.18.0 (February 07, 2025)

ENHANCEMENTS:

* dependencies: `appconfiguration` - update to API version `2024-05-01` ([#28700](https://github.com/hashicorp/terraform-provider-azurerm/issues/28700))
* dependencies: update `azurerm_cdn_frontdoor_rule` to API version `2024-02-01` ([#28308](https://github.com/hashicorp/terraform-provider-azurerm/issues/28308))
* dependencies: update `azurerm_cdn_frontdoor_ruleset` to API version `2024-02-01` ([#28308](https://github.com/hashicorp/terraform-provider-azurerm/issues/28308))
* dependencies: update `go-azure-sdk` to `v0.20250131.1134653` ([#28674](https://github.com/hashicorp/terraform-provider-azurerm/issues/28674))
* Data Source: `azurerm_cdn_frontdoor_firewall_policy` - add support for `js_challenge_cookie_expiration_in_minutes` policy ([#28284](https://github.com/hashicorp/terraform-provider-azurerm/issues/28284))
* Data Source: `azurerm_nginx_configuration` - add support for the `protected_file.content_hash` property ([#28532](https://github.com/hashicorp/terraform-provider-azurerm/issues/28532))
* `azurerm_cdn_frontdoor_firewall_policy` - add support for `js_challenge_cookie_expiration_in_minutes` policy ([#28284](https://github.com/hashicorp/terraform-provider-azurerm/issues/28284))
* `azurerm_cdn_frontdoor_firewall_policy` - add support for `JSChallenge` `action` type in the `managed_rule` `override` block ([#28308](https://github.com/hashicorp/terraform-provider-azurerm/issues/28308))
* `azurerm_container_app` - add support for the `volume_mounts.sub_path` property ([#27533](https://github.com/hashicorp/terraform-provider-azurerm/issues/27533))
* `azurerm_nginx_configuration` - add support for the `protected_file.content_hash` property ([#28532](https://github.com/hashicorp/terraform-provider-azurerm/issues/28532))
* `azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack` - add support for the `marketplace_offer_id` and `plan_id` properties ([#28537](https://github.com/hashicorp/terraform-provider-azurerm/issues/28537))
* `azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama` - add support for the `marketplace_offer_id` and `plan_id` properties ([#28537](https://github.com/hashicorp/terraform-provider-azurerm/issues/28537))
* `azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack` - add support for the `marketplace_offer_id` and `plan_id` properties ([#28537](https://github.com/hashicorp/terraform-provider-azurerm/issues/28537))
* `azurerm_palo_alto_next_generation_firewall_virtual_network_panorama` - add support for the `marketplace_offer_id` and `plan_id` properties ([#28537](https://github.com/hashicorp/terraform-provider-azurerm/issues/28537))
* `azurerm_route_server` - add support for the `hub_routing_preference` property ([#28363](https://github.com/hashicorp/terraform-provider-azurerm/issues/28363))

BUG FIXES:

* `azurerm_logic_app_action_http` - fix issue where `queries` would be set to an empty map instead of null when omitted from the configuration ([#28447](https://github.com/hashicorp/terraform-provider-azurerm/issues/28447))
* `azurerm_machine_learning_compute_cluster` - allow resource creation when `node_public_ip_enabled` is `false` and `subnet_resource_id` has not been specified ([#28673](https://github.com/hashicorp/terraform-provider-azurerm/issues/28673))
* `azurerm_network_watcher_flow_log` - prevent panic when removing the `traffic_analytics` block ([#28416](https://github.com/hashicorp/terraform-provider-azurerm/issues/28416))
* `azurerm_oracle_autonomous_database` - fix incorrect type for the `supported_regions_to_clone_to` property ([#28536](https://github.com/hashicorp/terraform-provider-azurerm/issues/28536))

## 4.17.0 (January 31, 2025)

FEATURES:

* **New Data Source**: `azurerm_api_management_subscription` ([#27824](https://github.com/hashicorp/terraform-provider-azurerm/issues/27824))
* **New Resource**: `azurerm_cognitive_account_rai_policy` ([#28013](https://github.com/hashicorp/terraform-provider-azurerm/issues/28013))
* **New Resource**: `azurerm_mssql_job_target_group` ([#28492](https://github.com/hashicorp/terraform-provider-azurerm/issues/28492))

ENHANCEMENTS:

* dependencies: `network` - update to use `2024-05-01` ([#28146](https://github.com/hashicorp/terraform-provider-azurerm/issues/28146))
* dependencies: `privatedns` - update to use `2024-06-01` ([#28599](https://github.com/hashicorp/terraform-provider-azurerm/issues/28599))
* dependencies: `storage` - update to use `2023-05-01` ([#27760](https://github.com/hashicorp/terraform-provider-azurerm/issues/27760))
* Data Source: `azure_communication_service` - add support for the `hostname` property ([#28620](https://github.com/hashicorp/terraform-provider-azurerm/issues/28620))
* `azurerm_api_management` - `capacity` now has a max limit of 50 ([#28648](https://github.com/hashicorp/terraform-provider-azurerm/issues/28648))
* `azurerm_backup_protected_vm` - add support for feature `vm_backup_suspend_protection_and_retain_data_on_destroy` ([#27950](https://github.com/hashicorp/terraform-provider-azurerm/issues/27950))
* `azurerm_cognitive_account` - support for the `bypass` property ([#28221](https://github.com/hashicorp/terraform-provider-azurerm/issues/28221))
* `azure_communication_service` - add support for the `hostname` property ([#28620](https://github.com/hashicorp/terraform-provider-azurerm/issues/28620))
* `azurerm_container_app_environment` - add support for Azure Monitor as a log destination ([#26047](https://github.com/hashicorp/terraform-provider-azurerm/issues/26047))
* `azurerm_mssql_elasticpool`- add support for `MOPRMS` pool type and update validation for `PRMS` and `Gen5` pool types ([#28453](https://github.com/hashicorp/terraform-provider-azurerm/issues/28453))
* `azurerm_mssql_managed_instance_transparent_data_encryption` - support for the `managed_hsm_key_id` property ([#28480](https://github.com/hashicorp/terraform-provider-azurerm/issues/28480))
* `azurerm_stream_analytics_output_cosmosdb` - support for the `authentication_mode` property ([#28372](https://github.com/hashicorp/terraform-provider-azurerm/issues/28372))
* `azurerm_stream_analytics_stream_input_blob` - add support for `authentication_mode` ([#27853](https://github.com/hashicorp/terraform-provider-azurerm/issues/27853))

BUG FIXES:

* `azurerm_container_app` - update the validation regex for the resource's name ([#28528](https://github.com/hashicorp/terraform-provider-azurerm/issues/28528))
* `azurerm_kubernetes_cluster` - parse `oms_agent.log_analytics_workspace_id` insensitively to handle inconsistent casing ([#28575](https://github.com/hashicorp/terraform-provider-azurerm/issues/28575))
* `azurerm_kubernetes_flux_configuration` - fix issue where removing `post_build` from a `kustomization` resulted in an error from the API ([#28590](https://github.com/hashicorp/terraform-provider-azurerm/issues/28590))
* `azurerm_linux_virtual_machine_scale_set` - prevent crash caused by ommited `extensions_to_provision_after_vm_creation` block ([#28549](https://github.com/hashicorp/terraform-provider-azurerm/issues/28549))
* `azurerm_log_analytics_storage_insights` - use subscription from workspace ID when building the resource ID ([#28469](https://github.com/hashicorp/terraform-provider-azurerm/issues/28469))
* `azurerm_orchestrated_virtual_machine_scale_set` - prevent crash caused by ommited `extensions_to_provision_after_vm_creation` block ([#28549](https://github.com/hashicorp/terraform-provider-azurerm/issues/28549))
* `azurerm_virtual_machine` - parse `os_disk` insensitively to handle inconsistent casing ([#28592](https://github.com/hashicorp/terraform-provider-azurerm/issues/28592))
* `azurerm_windows_virtual_machine_scale_set` - Prevent crash caused by ommited `extensions_to_provision_after_vm_creation` block ([#28549](https://github.com/hashicorp/terraform-provider-azurerm/issues/28549))

## 4.16.0 (January 16, 2025)

**NOTE:** This release contains a breaking change reverting `redisenterprise` API version from `2024-10-01` to `2024-06-01-preview` as not all regions are currently supported in the `2024-10-01` version 

BREAKING CHANGES:

* dependencies - `redisenterprise` API version reverted from `2024-10-01` to `2024-06-01-preview` ([#28516](https://github.com/hashicorp/terraform-provider-azurerm/issues/28516))

FEATURES:

* **New Resource**: `azurerm_container_registry_credential_set` ([#27528](https://github.com/hashicorp/terraform-provider-azurerm/issues/27528))
* **New Resource**: `azurerm_mssql_job` ([#28456](https://github.com/hashicorp/terraform-provider-azurerm/issues/28456))
* **New Resource**: `azurerm_mssql_job_schedule` ([#28456](https://github.com/hashicorp/terraform-provider-azurerm/issues/28456))

ENHANCEMENTS:

* dependencies - update `hashicorp/go-azure-sdk` to `v0.20250115.1141151` ([#28519](https://github.com/hashicorp/terraform-provider-azurerm/issues/28519))
* dependencies - `costmanagement` update to use `2023-08-01` ([#27680](https://github.com/hashicorp/terraform-provider-azurerm/issues/27680))
* dependencies - `postgresql` - partial upgrade to API version `2024-08-01` ([#28474](https://github.com/hashicorp/terraform-provider-azurerm/issues/28474))
* `azurerm_container_app` – support for the `termination_grace_period_seconds` property ([#28307](https://github.com/hashicorp/terraform-provider-azurerm/issues/28307))
* `azurerm_cost_anomaly_alert` - add support for the `notification_email` property ([#27680](https://github.com/hashicorp/terraform-provider-azurerm/issues/27680))
* `azurerm_data_protection_backup_vault` - support for `immutability` property ([#27859](https://github.com/hashicorp/terraform-provider-azurerm/issues/27859))
* `azurerm_databricks_workspace` - fix `ignore_changes` support ([#28527](https://github.com/hashicorp/terraform-provider-azurerm/issues/28527))
* `azurerm_kubernetes_cluster_node_pool` - add support for the `temporary_name_for_rotation` property to allow node pool rotation ([#27791](https://github.com/hashicorp/terraform-provider-azurerm/issues/27791))
* `azurerm_linux_function_app` - add  support for node `22` and java `17` support for `JBOSSEAP` ([#28472](https://github.com/hashicorp/terraform-provider-azurerm/issues/28472))
* `azurerm_linux_web_app` - add  support for node `22` and java `17` support for `JBOSSEAP` ([#28472](https://github.com/hashicorp/terraform-provider-azurerm/issues/28472))
* `azurerm_windows_function_app` - add  support for node `22` and java `17` support for `JBOSSEAP` ([#28472](https://github.com/hashicorp/terraform-provider-azurerm/issues/28472))


BUG FIXES:

* `azurerm_logic_app_standard` - fix setting `public_network_access` for conflicting API properties ([#28465](https://github.com/hashicorp/terraform-provider-azurerm/issues/28465))
* `azurerm_redis_cache` - `data_persistence_authentication_method` can now be unset ([#27932](https://github.com/hashicorp/terraform-provider-azurerm/issues/27932))
* `azurerm_mssql_database` - fix bug where verifying TDE might fail to return an error on failure ([#28505](https://github.com/hashicorp/terraform-provider-azurerm/issues/28505))
* `azurerm_mssql_database` - fix several potential bugs where retry functions could return false negatives for actual errors ([#28505](https://github.com/hashicorp/terraform-provider-azurerm/issues/28505))
* `azurerm_private_endpoint` - fix a bug where reading Private DNS could error and exit the Read of the resource early without raising an error ([#28505](https://github.com/hashicorp/terraform-provider-azurerm/issues/28505))


## 4.15.0 (January 10, 2025)

FEATURES:

* **New Data Source**: `azurerm_kubernetes_fleet_manager` ([#28278](https://github.com/hashicorp/terraform-provider-azurerm/issues/28278))
* **New Resource**: `azurerm_arc_kubernetes_provisioned_cluster` ([#28216](https://github.com/hashicorp/terraform-provider-azurerm/issues/28216))
* **New Resource**: `azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint` ([#27874](https://github.com/hashicorp/terraform-provider-azurerm/issues/27874))
* **New Resource** `azurerm_machine_learning_workspace_network_outbound_rule_service_tag` ([#27931](https://github.com/hashicorp/terraform-provider-azurerm/issues/27931))
* **New Resource** `azurerm_dynatrace_tag_rules` ([#27985](https://github.com/hashicorp/terraform-provider-azurerm/issues/27985))

ENHANCEMENTS:

* dependencies - update tool Go version and bump `go-git` version to `5.13.0` ([#28425](https://github.com/hashicorp/terraform-provider-azurerm/issues/28425))
* dependencies - update `hashicorp/go-azure-sdk` to `v0.20241212.1154051` ([#28360](https://github.com/hashicorp/terraform-provider-azurerm/issues/28360))
* dependencies - `frontdoor` - partial update to use `2024-02-01` API ([#28233](https://github.com/hashicorp/terraform-provider-azurerm/issues/28233))
* dependencies - `postgresql` - update to `2024-08-01` ([#28380](https://github.com/hashicorp/terraform-provider-azurerm/issues/28380))
* dependencies - `redisenterprise` - update to `2024-10-01` and support for new skus ([#28280](https://github.com/hashicorp/terraform-provider-azurerm/issues/28280))
* Data Source: `azurerm_healthcare_dicom_service` - add support for the `data_partitions_enabled`, `cors`, `encryption_key_url` and `storage` properties ([#27375](https://github.com/hashicorp/terraform-provider-azurerm/issues/27375))
* Data Source: `azurerm_nginx_deployment` - add support for the `dataplane_api_endpoint` property ([#28379](https://github.com/hashicorp/terraform-provider-azurerm/issues/28379)) 
* Data Source: `azurerm_static_web_app` - add  support for the `repository_url` and `repository_branch` properties ([#27401](https://github.com/hashicorp/terraform-provider-azurerm/issues/27401))
* `azurerm_billing_account_cost_management_export` - add support for the `file_format` property ([#27122](https://github.com/hashicorp/terraform-provider-azurerm/issues/27122))
* `azurerm_cdn_frontdoor_profile` - add support for the `identity` property ([#28281](https://github.com/hashicorp/terraform-provider-azurerm/issues/28281))
* `azurerm_cognitive_deployment` - `DataZoneProvisionedManaged` and `GlobalProvisionedManaged` skus are now supported ([#28404](https://github.com/hashicorp/terraform-provider-azurerm/issues/28404))
* `azurerm_databricks_access_connector` - `SystemAssigned,UserAssigned` identity is now supported ([#28442](https://github.com/hashicorp/terraform-provider-azurerm/issues/28442))
* `azurerm_healthcare_dicom_service` - add support for the `data_partitions_enabled`, `cors`, `encryption_key_url` and `storage` properties ([#27375](https://github.com/hashicorp/terraform-provider-azurerm/issues/27375))
* `azurerm_kubernetes_flux_configuration` - add support for the `post_build` and `wait` properties ([#25695](https://github.com/hashicorp/terraform-provider-azurerm/issues/25695))
* `azurerm_linux_virtual_machine` - export the `os_disk.0.id` attribute ([#28352](https://github.com/hashicorp/terraform-provider-azurerm/issues/28352))
* `azurerm_netapp_volume` - make the `network_features` property Optional/Computed ([#28390](https://github.com/hashicorp/terraform-provider-azurerm/issues/28390))
* `azurerm_nginx_deployment` - add support for the `dataplane_api_endpoint` property ([#28379](https://github.com/hashicorp/terraform-provider-azurerm/issues/28379)) 
* `azurerm_resource_group_cost_management_export` - add support for the `file_format` property ([#27122](https://github.com/hashicorp/terraform-provider-azurerm/issues/27122))
* `azurerm_site_recovery_replicated_vm` - support for the `network_interface.recovery_load_balancer_backend_address_pool_ids` property ([#28398](https://github.com/hashicorp/terraform-provider-azurerm/issues/28398))
* `azurerm_static_web_app` - add  support for the `repository_url`, `repository_branch` and `repository_token` properties ([#27401](https://github.com/hashicorp/terraform-provider-azurerm/issues/27401))
* `azurerm_subscription_cost_management_export` - add support for the `file_format` property ([#27122](https://github.com/hashicorp/terraform-provider-azurerm/issues/27122))
* `azurerm_virtual_network` - support for the `private_endpoint_vnet_policies` property ([#27830](https://github.com/hashicorp/terraform-provider-azurerm/issues/27830))
* `azurerm_windows_virtual_machine` - export the `os_disk.0.id` attribute ([#28352](https://github.com/hashicorp/terraform-provider-azurerm/issues/28352))
* `azurerm_mssql_managed_instance` - support for new property `azure_active_directory_administrator` ([#24801](https://github.com/hashicorp/terraform-provider-azurerm/issues/24801))

BUG FIXES:

* `azurerm_api_management` - update the `capacity` property to allow increasing the apim scalability to `31` ([#28427](https://github.com/hashicorp/terraform-provider-azurerm/issues/28427))
* `azurerm_automation_software_update_configuration` remove deprecated misspelled attribute `error_meesage` ([#28312](https://github.com/hashicorp/terraform-provider-azurerm/issues/28312))
* `azurerm_batch_pool` - support for new block `security_profile` ([#28069](https://github.com/hashicorp/terraform-provider-azurerm/issues/28069))
* `azurerm_log_analytics_data_export_rule` - now creates successfully without returning `404` ([#27876](https://github.com/hashicorp/terraform-provider-azurerm/issues/27876))
* `azurerm_mongo_cluster` - remove CustomizeDiff logic for `administrator_password` to allow the input to be generated by the `random_password` resource ([#28215](https://github.com/hashicorp/terraform-provider-azurerm/issues/28215))
* `azurerm_mongo_cluster` - valdation updated so the resource now creates successfully when using `create_mode` `GeoReplica` ([#28269](https://github.com/hashicorp/terraform-provider-azurerm/issues/28269))
* `azurerm_mssql_managed_instance` - allow system and user assigned identities, fix update failure ([#28319](https://github.com/hashicorp/terraform-provider-azurerm/issues/28319))
* `azurerm_storage_account` - fix error handling for `static_website` and `queue_properties` availability checks ([#28279](https://github.com/hashicorp/terraform-provider-azurerm/issues/28279))



## 4.14.0 (December 12, 2024)

BREAKING CHANGES:

* `nginx` - update api version to `2024-09-01-preview`, this API no longer supports certain properties which have had to be deprecated in the provider for the upgrade ([#27776](https://github.com/hashicorp/terraform-provider-azurerm/issues/27776))
* Data Source: `azurerm_nginx_configuration` - the `protected_file.content` property will not be populated and has been deprecated ([#27776](https://github.com/hashicorp/terraform-provider-azurerm/issues/27776))
* Data Source: `azurerm_nginx_deployment` - the `managed_resource_group` property will not be populated and has been deprecated ([#27776](https://github.com/hashicorp/terraform-provider-azurerm/issues/27776))
* `azurerm_network_function_collector_policy` - the API doesn't preserve the ordering of the `ipfx_ingestion.source_resource_ids` property causing non-empty plans after apply, this property's type has been changed from a list to a set to prevent Terraform from continually trying to recreate this resource. If this property is being referenced anywhere you will need to update your config to convert it to a list before referencing it ([#27915](https://github.com/hashicorp/terraform-provider-azurerm/issues/27915))
* `azurerm_nginx_deployment` - the `managed_resource_group` property is no longer supported and has been deprecated ([#27776](https://github.com/hashicorp/terraform-provider-azurerm/issues/27776))

FEATURES:

* **New Resource**: `azurerm_cognitive_account_rai_blocklist` ([#28043](https://github.com/hashicorp/terraform-provider-azurerm/issues/28043))
* **New Resource**: `azurerm_fabric_capacity` ([#28080](https://github.com/hashicorp/terraform-provider-azurerm/issues/28080))

ENHANCEMENTS:

* dependencies - update `go-azure-sdk` to `v0.20241206.1180327` ([#28211](https://github.com/hashicorp/terraform-provider-azurerm/issues/28211))
* `nginx` - update api version to `2024-11-01-preview` ([#28227](https://github.com/hashicorp/terraform-provider-azurerm/issues/28227))
* `azurerm_linux_function_app` - add support for  preview  value `21` for `java_version` ([#26304](https://github.com/hashicorp/terraform-provider-azurerm/issues/26304))
* `azurerm_linux_function_app_slot` - support `1.3` for `site_config.minimum_tls_version` and `site_config.scm_minimum_tls_version` ([#28016](https://github.com/hashicorp/terraform-provider-azurerm/issues/28016))
* `azurerm_linux_web_app` - add support for  preview  value `21` for `java_version` ([#26304](https://github.com/hashicorp/terraform-provider-azurerm/issues/26304))
* `azurerm_orchestrated_virtual_machine_scale_set` - support hot patching for `2025-datacenter-azure-edition-core-smalldisk` ([#28160](https://github.com/hashicorp/terraform-provider-azurerm/issues/28160))
* `azurerm_search_service` - add support for the `network_rule_bypass_option` property ([#28139](https://github.com/hashicorp/terraform-provider-azurerm/issues/28139))
* `azurerm_windows_function_app` - add support for  preview  value `21` for `java_version` ([#26304](https://github.com/hashicorp/terraform-provider-azurerm/issues/26304))
* `azurerm_windows_function_app_slot` - support `1.3` for `site_config.minimum_tls_version` and `site_config.scm_minimum_tls_version` ([#28016](https://github.com/hashicorp/terraform-provider-azurerm/issues/28016))
* `azurerm_windows_virtual_machine` - support hot patching for `2025-datacenter-azure-edition-core-smalldisk` ([#28160](https://github.com/hashicorp/terraform-provider-azurerm/issues/28160))
* `azurerm_windows_web_app` - add support for  preview  value `21` for `java_version` ([#26304](https://github.com/hashicorp/terraform-provider-azurerm/issues/26304))

BUG FIXES:

* `azurerm_management_group` - fix regression where subscription ID can't be parsed correctly anymore ([#28228](https://github.com/hashicorp/terraform-provider-azurerm/issues/28228))

## 4.13.0 (December 05, 2024)

ENHANCEMENTS:

* `azurerm_cognitive_deployment` - support for the `dynamic_throttling_enabled` property ([#28100](https://github.com/hashicorp/terraform-provider-azurerm/issues/28100))
* `azurerm_key_vault_managed_hardware_security_module_key` - the `key_type` property now supports `oct-HSM` ([#28171](https://github.com/hashicorp/terraform-provider-azurerm/issues/28171))
* `azurerm_machine_learning_datastore_datalake_gen2` - can now be used with storage account in a different subscription ([#28123](https://github.com/hashicorp/terraform-provider-azurerm/issues/28123))
* `azurerm_network_watcher_flow_log` - `target_resource_id` supports subnets and network interfaces ([#28177](https://github.com/hashicorp/terraform-provider-azurerm/issues/28177))

BUG:

* Data Source: `azurerm_logic_app_standard` - update the `identity` property to support User Assigned Identities ([#28158](https://github.com/hashicorp/terraform-provider-azurerm/issues/28158))
* `azurerm_cdn_frontdoor_origin_group` - update validation of the `interval_in_seconds` property to match API behaviour ([#28143](https://github.com/hashicorp/terraform-provider-azurerm/issues/28143))
* `azurerm_container_group` - retrieve log analytics workspace key from config when updating resource ([#28025](https://github.com/hashicorp/terraform-provider-azurerm/issues/28025))
* `azurerm_mssql_elasticpool` - fix sku tier and family validation that prevented the creation of Hyperscale PRMS pools ([#28178](https://github.com/hashicorp/terraform-provider-azurerm/issues/28178))
* `azurerm_search_service` -  the `partition_count` property can now be up to `3` when using basic sku ([#28105](https://github.com/hashicorp/terraform-provider-azurerm/issues/28105))

## 4.12.0 (November 28, 2024)

FEATURES:

* **New Data Source**: `azurerm_mssql_managed_database` ([#27026](https://github.com/hashicorp/terraform-provider-azurerm/issues/27026))

BUG FIXES:

* `azurerm_application_insights_api_key` - fix condition that nil checks the list of available API keys to prevent an indefinate loop when keys created outside of Terraform are present ([#28037](https://github.com/hashicorp/terraform-provider-azurerm/issues/28037))
* `azurerm_data_factory_linked_service_azure_sql_database` - send `tenant_id` only if it has been specified ([#28120](https://github.com/hashicorp/terraform-provider-azurerm/issues/28120))
* `azurerm_eventgrid_event_subscription` - fix crash when flattening `advanced_filter` ([#28110](https://github.com/hashicorp/terraform-provider-azurerm/issues/28110))
* `azurerm_virtual_network_gateway` - fix crash issue when specifying `root_certificate ` or `revoked_certificate` ([#28099](https://github.com/hashicorp/terraform-provider-azurerm/issues/28099))

ENHANCEMENTS:

* dependencies - update `go-azure-sdk` to `v0.20241128.1112539` ([#28137](https://github.com/hashicorp/terraform-provider-azurerm/issues/28137))
* `containerapps` - update api version to `2024-03-01` ([#28074](https://github.com/hashicorp/terraform-provider-azurerm/issues/28074))
* `Search` - update api version to `2024-06-01-preview` ([#27803](https://github.com/hashicorp/terraform-provider-azurerm/issues/27803))
* Data Source: `azurerm_logic_app_standard` - add support for the `public_network_access` property ([#27913](https://github.com/hashicorp/terraform-provider-azurerm/issues/27913))
* Data Source: `azurerm_search_service` - add support for the `customer_managed_key_encryption_compliance_status` property ([#27478](https://github.com/hashicorp/terraform-provider-azurerm/issues/27478))
* `azurerm_container_registry_task` - add validation on `cpu` as well as on `agent_pool_name`and `agent_setting` ([#28098](https://github.com/hashicorp/terraform-provider-azurerm/issues/28098))
* `azurerm_databricks_workspace` - add support for the `enhanced_security_compliance` block ([#26606](https://github.com/hashicorp/terraform-provider-azurerm/issues/26606))
* `azurerm_eventhub` - deprecate `namespace_name` and `resource_group_name` in favour of `namespace_id` ([#28055](https://github.com/hashicorp/terraform-provider-azurerm/issues/28055))
* `azurerm_logic_app_standard` - add support for the `public_network_access` property ([#27913](https://github.com/hashicorp/terraform-provider-azurerm/issues/27913))
* `azurerm_search_service` - add support for the `customer_managed_key_encryption_compliance_status` property ([#27478](https://github.com/hashicorp/terraform-provider-azurerm/issues/27478))
* `azurerm_cosmosdb_account` - add support for value `EnableNoSQLFullTextSearch` in the  `capabilities.name` property ([#28114](https://github.com/hashicorp/terraform-provider-azurerm/issues/28114))

## 4.11.0 (November 22, 2024)

NOTES:
* New [ephemeral resources](https://developer.hashicorp.com/terraform/language/v1.10.x/resources/ephemeral) `azurerm_key_vault_certificate` and `azurerm_key_vault_secret` now support [ephemeral values](https://developer.hashicorp.com/terraform/language/v1.10.x/values/variables#exclude-values-from-state)

FEATURES:

* **New Ephemeral Resource**: `azurerm_key_vault_certificate` ([#28083](https://github.com/hashicorp/terraform-provider-azurerm/issues/28083))
* **New Ephemeral Resource**: `azurerm_key_vault_secret` ([#28083](https://github.com/hashicorp/terraform-provider-azurerm/issues/28083))
* **New Resource**: `azurerm_eventgrid_namespace` ([#27682](https://github.com/hashicorp/terraform-provider-azurerm/issues/27682))

ENHANCEMENTS:

* dependencies: update `hashicorp/go-azure-sdk` to `v0.20241118.1115603` ([#28075](https://github.com/hashicorp/terraform-provider-azurerm/issues/28075))
* `batch` - upgrade api version to `2024-07-01` ([#27982](https://github.com/hashicorp/terraform-provider-azurerm/issues/27982))
* `containerregistry` - upgrade api version to `2023-11-01-preview` ([#27983](https://github.com/hashicorp/terraform-provider-azurerm/issues/27983))
* `azurerm_application_gateway` - `1.1` is now accepted as a valid `rule_set_version` in the `waf_configuration` block ([#28039](https://github.com/hashicorp/terraform-provider-azurerm/issues/28039))
* `azurerm_arc_machine` - add support for the `identity` and `tags` properties ([#27987](https://github.com/hashicorp/terraform-provider-azurerm/issues/27987))
* `azurerm_container_app` - `secret.name` now accepts up to 253 characters and `.` ([#27935](https://github.com/hashicorp/terraform-provider-azurerm/issues/27935))
* `azurerm_network_manager` - `scope_accesses` now accepts `Routing` ([#28033](https://github.com/hashicorp/terraform-provider-azurerm/issues/28033))
* `azurerm_network_watcher_flow_log` - add support for the `target_resource_id` property ([#26015](https://github.com/hashicorp/terraform-provider-azurerm/issues/26015))
* `azurerm_role_assignment` - `condition_version` will be defaulted to `2.0` when `condition` has been set ([#27189](https://github.com/hashicorp/terraform-provider-azurerm/issues/27189))
* `azurerm_subnet` - `Informatica.DataManagement/organizations` is a valid `service_delegation` ([#27993](https://github.com/hashicorp/terraform-provider-azurerm/issues/27993))
* `azurerm_virtual_network` - `Informatica.DataManagement/organizations` is a valid `service_delegation` ([#27993](https://github.com/hashicorp/terraform-provider-azurerm/issues/27993))
* `azurerm_web_application_firewall_policy` - `1.1` is now accepted as a valid `version` for `Microsoft_BotManagerRuleSet` rule types ([#28039](https://github.com/hashicorp/terraform-provider-azurerm/issues/28039))

BUG FIXES:

* `azurerm_api_management` - `public_ip_address_id` is no longer required when `zone` has been set ([#27976](https://github.com/hashicorp/terraform-provider-azurerm/issues/27976))
* `azurerm_api_management_diagnostic` - raise and error when `operation_name_format` is used with and `identity` that is not `applicationinsights` ([#27630](https://github.com/hashicorp/terraform-provider-azurerm/issues/27630))
* `azurerm_api_management_api_diagnostic` - raise and error when `operation_name_format` is used with and `identity` that is not `applicationinsights` ([#27630](https://github.com/hashicorp/terraform-provider-azurerm/issues/27630))
* `azurerm_application_gateway` - `rewrite_rule_set` can be supplied when using `Basic` sku ([#28011](https://github.com/hashicorp/terraform-provider-azurerm/issues/28011))
* `azurerm_container_registry_token_password` - correctly mark as gone if container registry token doesn't exist ([#27232](https://github.com/hashicorp/terraform-provider-azurerm/issues/27232))
* `azurerm_kusto_cluster` - `allowed_fqdn` and `allowed_ip_ranges` can now be set to empty lists ([#27529](https://github.com/hashicorp/terraform-provider-azurerm/issues/27529))
* `azurerm_linux_function_app_slot` - create content settings when using a consumpton plan ([#25412](https://github.com/hashicorp/terraform-provider-azurerm/issues/25412))
* `azurerm_virtual_network_gatway` - updating `ip_configuration` now recreates the resource ([#27828](https://github.com/hashicorp/terraform-provider-azurerm/issues/27828))


## 4.10.0 (November 14, 2024)

BREAKING CHANGES:

* dependencies - update `cognitive` to `2024-10-01`, due to a behavioural change in this version of the API, the `primary_access_key` and `secondary_access_key` can not be retrieved if `local_authentication_enabled` has been set to `false`. These properties that may have had values previously will now be empty. This has affected the `azurerm_ai_services` and `azurerm_cognitive_account` resources as well as the `azurerm_cognitive_account` data source ([#27851](https://github.com/hashicorp/terraform-provider-azurerm/issues/27851))

FEATURES:

* **New Data Source**: `azurerm_key_vault_managed_hardware_security_module_key` ([#27827](https://github.com/hashicorp/terraform-provider-azurerm/issues/27827))
* **New Resource**: `azurerm_netapp_backup_vault` ([#27188](https://github.com/hashicorp/terraform-provider-azurerm/issues/27188))
* **New Resource**: `azurerm_netapp_backup_policy` ([#27188](https://github.com/hashicorp/terraform-provider-azurerm/issues/27188))

ENHANCEMENTS:

* dependencies: update `terraform-plugin-framework` to version `v1.13.0` ([#27936](https://github.com/hashicorp/terraform-provider-azurerm/issues/27936))
* dependencies: update `terraform-plugin-framework-validators` to version `v0.14.0` ([#27936](https://github.com/hashicorp/terraform-provider-azurerm/issues/27936))
* dependencies: update `terraform-plugin-go` to version `v0.25.0` ([#27936](https://github.com/hashicorp/terraform-provider-azurerm/issues/27936))
* dependencies: update `terraform-plugin-mux` to version `v0.17.0` ([#27936](https://github.com/hashicorp/terraform-provider-azurerm/issues/27936))
* dependencies: update `terraform-plugin-sdk/v2` to version `v2.35.0` ([#27936](https://github.com/hashicorp/terraform-provider-azurerm/issues/27936))
* Data Source: `azurerm_bastion_host` - add support for the `zones` property ([#27909](https://github.com/hashicorp/terraform-provider-azurerm/issues/27909))
* `azurerm_application_gateway` - support more values for the `status_code` property ([#27535](https://github.com/hashicorp/terraform-provider-azurerm/issues/27535))
* `azurerm_bastion_host` - support for the `zones` property ([#27909](https://github.com/hashicorp/terraform-provider-azurerm/issues/27909))
* `azurerm_communication_service` - support for `usgov` region ([#27919](https://github.com/hashicorp/terraform-provider-azurerm/issues/27919))
* `azurerm_email_communication_service` - support for `usgov` region added ([#27919](https://github.com/hashicorp/terraform-provider-azurerm/issues/27919))
* `azurerm_linux_function_app` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_linux_function_app_slot` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_linux_web_app` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_linux_web_app_slot` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_windows_web_app` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_windows_web_app_slot` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_windows_function_app` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_windows_function_app_slot` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))

BUG FIXES:

* `azurerm_log_analytics_workspace_table` - use the subscription from workspace ID ([#27590](https://github.com/hashicorp/terraform-provider-azurerm/issues/27590))
* `azurerm_traffic_manager_external_endpoint` - the value for `priority` will be dynamically assigned by the API ([#27966](https://github.com/hashicorp/terraform-provider-azurerm/issues/27966))
* `azurerm_traffic_manager_azure_endpoint` - the value for `priority` will be dynamically assigned by the API ([#27966](https://github.com/hashicorp/terraform-provider-azurerm/issues/27966))


## 4.9.0 (November 08, 2024)

FEATURES:

* **New Resource**: `azurerm_dynatrace_monitor` ([#27432](https://github.com/hashicorp/terraform-provider-azurerm/issues/27432))
* **New Resource**: `azurerm_dashboard_grafana_managed_private_endpoint` ([#27781](https://github.com/hashicorp/terraform-provider-azurerm/issues/27781))
* **New Resource**: `azurerm_data_protection_backup_instance_mysql_flexible_server` ([#27464](https://github.com/hashicorp/terraform-provider-azurerm/issues/27464))
* **New Resource**: `azurerm_mongo_cluster` ([#27636](https://github.com/hashicorp/terraform-provider-azurerm/issues/27636))
* **New Resource**: `azurerm_stack_hci_network_interface` ([#26888](https://github.com/hashicorp/terraform-provider-azurerm/issues/26888))

ENHANCEMENTS:

* dependencies - update `go-azure-sdk` to `v0.20241104.1140654` ([#27896](https://github.com/hashicorp/terraform-provider-azurerm/issues/27896))
* dependencies - update `go-azure-helpers` to `v0.71.0` ([#27897](https://github.com/hashicorp/terraform-provider-azurerm/issues/27897))
* dependencies - update `golang-jwt` to `v4.5.1` ([#27938](https://github.com/hashicorp/terraform-provider-azurerm/issues/27938))
* `storage` - allow `azurerm_storage_account` to be used in Data Plane restrictive environments ([#27818](https://github.com/hashicorp/terraform-provider-azurerm/issues/27818))
* `azurerm_cognitive_deployment` - `sku.0.name` now supports `DataZoneStandard` ([#27926](https://github.com/hashicorp/terraform-provider-azurerm/issues/27926))
* `azurerm_mssql_managed_database` - support for the `tags` property ([#27857](https://github.com/hashicorp/terraform-provider-azurerm/issues/27857))
* `azurerm_oracle_cloud_vm_cluster` - support for the `domain`, `scan_listener_port_tcp`, `scan_listener_port_tcp_ssl` and `zone_id` properties ([#27808](https://github.com/hashicorp/terraform-provider-azurerm/issues/27808))
* `azurerm_public_ip_prefix` - support for the `sku_tier` property ([#27882](https://github.com/hashicorp/terraform-provider-azurerm/issues/27882))
* `azurerm_public_ip` - support for the `domain_name_label_scope` property ([#27748](https://github.com/hashicorp/terraform-provider-azurerm/issues/27748))
* `azurerm_subnet` - `default_outbound_access_enabled` can now be updated ([#27858](https://github.com/hashicorp/terraform-provider-azurerm/issues/27858))
* `azurerm_storage_container` - support for the `storage_account_id` property ([#27733](https://github.com/hashicorp/terraform-provider-azurerm/issues/27733))
* `azurerm_storage_share` - support for the `storage_account_id` property ([#27733](https://github.com/hashicorp/terraform-provider-azurerm/issues/27733))

## 4.8.0 (October 31, 2024)

FEATURES:

* **New Data Source**: `azurerm_virtual_network_peering` ([#27530](https://github.com/hashicorp/terraform-provider-azurerm/issues/27530))
* **New Resource**: `azurerm_machine_learning_workspace_network_outbound_rule_fqdn` ([#27384](https://github.com/hashicorp/terraform-provider-azurerm/issues/27384))
* **New Resource**: `azurerm_stack_hci_extension` ([#26929](https://github.com/hashicorp/terraform-provider-azurerm/issues/26929))
* **New Resource**: `azurerm_stack_hci_marketplace_gallery_image` ([#27532](https://github.com/hashicorp/terraform-provider-azurerm/issues/27532))
* **New Resource**: `azurerm_trusted_signing_account` ([#27720](https://github.com/hashicorp/terraform-provider-azurerm/issues/27720))

ENHANCEMENTS:

* `mysql` - upgrade api version to `2023-12-30` ([#27767](https://github.com/hashicorp/terraform-provider-azurerm/issues/27767))
* `network` - upgrade api version to `2024-03-01 ` ([#27746](https://github.com/hashicorp/terraform-provider-azurerm/issues/27746))
* `azurerm_cosmosdb_account`: support for CMK through `managed_hsm_key_id` property ([#26521](https://github.com/hashicorp/terraform-provider-azurerm/issues/26521))
* `azurerm_cosmosdb_account` - support further versions for `mongo_server_version` ([#27763](https://github.com/hashicorp/terraform-provider-azurerm/issues/27763))
* `azurerm_container_app_environment` - changing the `log_analytics_workspace_id` property no longer creates a new resource ([#27794](https://github.com/hashicorp/terraform-provider-azurerm/issues/27794))
* `azurerm_data_factory_linked_service_azure_sql_database` - add support for the `credential_name` property ([#27629](https://github.com/hashicorp/terraform-provider-azurerm/issues/27629))
* `azurerm_key_vault_key` - `expiration_date` only recreates the resource when it is removed from the config file ([#27813](https://github.com/hashicorp/terraform-provider-azurerm/issues/27813))
* `azurerm_kubernetes_cluster` - fix issue where`maintenance_window_auto_upgrade`/`maintenance_window_auto_upgrade`/`maintenance_window_node_os ` might not be read into state ([#26915](https://github.com/hashicorp/terraform-provider-azurerm/issues/26915))
* `azurerm_kubernetes_cluster` - support for the `backend_pool_type` property ([#27596](https://github.com/hashicorp/terraform-provider-azurerm/issues/27596))
* `azurerm_kubernetes_cluster` - support for the `daemonset_eviction_for_empty_nodes_enabled`, `daemonset_eviction_for_occupied_nodes_enabled`, and `ignore_daemonsets_utilization_enabled` properties ([#27588](https://github.com/hashicorp/terraform-provider-azurerm/issues/27588))
* `azurerm_load_test` - `description` can now be updated ([#27800](https://github.com/hashicorp/terraform-provider-azurerm/issues/27800))
* `azurerm_oracle_cloud_vm_cluster` - export the `ocid` property ([#27785](https://github.com/hashicorp/terraform-provider-azurerm/issues/27785))
* `azurerm_orchestrated_virtual_machine_scale_set` - add support for `sku_profile` block ([#27599](https://github.com/hashicorp/terraform-provider-azurerm/issues/27599))
* `azurerm_web_application_firewall_policy` - add support for `policy_settings.0.file_upload_enforcement` ([#27774](https://github.com/hashicorp/terraform-provider-azurerm/issues/27774))

BUG FIXES:

* `azurerm_automation_hybrid_runbook_worker_group` - correctly mark resource as gone if it's absent when reading it ([#27797](https://github.com/hashicorp/terraform-provider-azurerm/issues/27797))
* `azurerm_automation_hybrid_runbook_worker` - correctly mark resource as gone if it's absent when reading it ([#27797](https://github.com/hashicorp/terraform-provider-azurerm/issues/27797))
* `azurerm_automation_python3_package` - correctly mark resource as gone if it's absent when reading it ([#27797](https://github.com/hashicorp/terraform-provider-azurerm/issues/27797))
* `azurerm_data_protection_backup_vault` - prevent panic when checking value of `cross_region_restore_enabled` ([#27762](https://github.com/hashicorp/terraform-provider-azurerm/issues/27762))
* `azurerm_role_management_policy` - fix panic when unmarshalling the policy into a specific type ([#27731](https://github.com/hashicorp/terraform-provider-azurerm/issues/27731))
* `azurerm_security_center_subscription_pricing` - correctly type assert the `additional_extension_properties` property when building the payload ([#27721](https://github.com/hashicorp/terraform-provider-azurerm/issues/27721))
* `azurerm_synapse_workspace_aad_admin` - will no correctly delete when using `azurerm_synapse_workspace_aad_admin` with `azurerm_synapse_workspace` ([#27606](https://github.com/hashicorp/terraform-provider-azurerm/issues/27606))
* `azurerm_windows_function_app_slot` - fixed panic in state migration ([#27700](https://github.com/hashicorp/terraform-provider-azurerm/issues/27700))

## 4.7.0 (October 24, 2024)

FEATURES:

* **New Data Source**: `azurerm_oracle_adbs_character_sets` ([#27698](https://github.com/hashicorp/terraform-provider-azurerm/issues/27698))
* **New Data Source**: `azurerm_oracle_adbs_national_character_sets` ([#27698](https://github.com/hashicorp/terraform-provider-azurerm/issues/27698))
* **New Data Source**: `azurerm_oracle_autonomous_database` ([#27696](https://github.com/hashicorp/terraform-provider-azurerm/issues/27696))
* **New Data Source**: `azurerm_oracle_db_nodes` ([#27698](https://github.com/hashicorp/terraform-provider-azurerm/issues/27698))
* **New Data Source**: `azurerm_oracle_db_system_shapes` ([#27698](https://github.com/hashicorp/terraform-provider-azurerm/issues/27698))
* **New Data Source**: `azurerm_oracle_gi_versions` ([#27698](https://github.com/hashicorp/terraform-provider-azurerm/issues/27698))
* **New Resource**: `azurerm_dev_center_project_pool` ([#27706](https://github.com/hashicorp/terraform-provider-azurerm/issues/27706))
* **New Resource**: `azurerm_oracle_autonomous_database` ([#27696](https://github.com/hashicorp/terraform-provider-azurerm/issues/27696))
* **New Resource**: `azurerm_video_indexer_account` ([#27632](https://github.com/hashicorp/terraform-provider-azurerm/issues/27632))

ENHANCEMENTS:

* dependencies - update `go-azure-sdk` to `v0.20241021.1074254` ([#27713](https://github.com/hashicorp/terraform-provider-azurerm/issues/27713))
* `newrelic` - upgrade api version to `2024-03-01`  ([#27135](https://github.com/hashicorp/terraform-provider-azurerm/issues/27135))
* `cosmosdb` - upgrade api version to `2024-08-15` ([#27659](https://github.com/hashicorp/terraform-provider-azurerm/issues/27659))
* `azurerm_application_gateway` - support for the new `Basic` SKU value ([#27440](https://github.com/hashicorp/terraform-provider-azurerm/issues/27440))
* `azurerm_consumption_budget_management_group` - the property `notification.threshold_type` can now be updated ([#27511](https://github.com/hashicorp/terraform-provider-azurerm/issues/27511))
* `azurerm_consumption_budget_resource_group` - the property `notification.threshold_type` can now be updated ([#27511](https://github.com/hashicorp/terraform-provider-azurerm/issues/27511))
* `azurerm_container_app` - add support for the `template.container.readiness_probe.initial_delay` and `template.container.startup_probe.initial_delay` properties ([#27551](https://github.com/hashicorp/terraform-provider-azurerm/issues/27551))
* `azurerm_mssql_managed_instance` - the `storage_account_type` property can now be updated ([#27737](https://github.com/hashicorp/terraform-provider-azurerm/issues/27737))

BUG FIXES:

* `azurerm_automation_software_update_configuration` - correct validation to not allow `5` and allow `-1` ([#25574](https://github.com/hashicorp/terraform-provider-azurerm/issues/25574))
* `azurerm_cosmosdb_sql_container` - fix recreation logic for `partition_key_version` ([#27692](https://github.com/hashicorp/terraform-provider-azurerm/issues/27692))
* `azurerm_mssql_database` - updating short term retention policy now works as expected ([#27714](https://github.com/hashicorp/terraform-provider-azurerm/issues/27714))
* `azurerm_network_watcher_flow_log` - fix issue where `tags` were not being updated ([#27389](https://github.com/hashicorp/terraform-provider-azurerm/issues/27389))
* `azurerm_postgresql_flexible_server_virtual_endpoint` - retrieve and parse `replica_server_id` for cross-region scenarios as well as remove custom poller for the delete operation ([#27509](https://github.com/hashicorp/terraform-provider-azurerm/issues/27509))

## 4.6.0 (October 18, 2024)

FEATURES:

* **New Resource**: `azurerm_dev_center_attached_network` ([#27638](https://github.com/hashicorp/terraform-provider-azurerm/issues/27638))
* **New Resource**: `azurerm_oracle_cloud_vm_cluster` ([#27678](https://github.com/hashicorp/terraform-provider-azurerm/issues/27678))
* **New Resource**: `azurerm_oracle_exadata_infrastructure` ([#27678](https://github.com/hashicorp/terraform-provider-azurerm/issues/27678))
* **New Data Source**: `azurerm_oracle_cloud_vm_cluster` ([#27678](https://github.com/hashicorp/terraform-provider-azurerm/issues/27678))
* **New Data Source**: `azurerm_oracle_db_servers` ([#27678](https://github.com/hashicorp/terraform-provider-azurerm/issues/27678))
* **New Data Source**: `azurerm_oracle_exadata_infrastructure` ([#27678](https://github.com/hashicorp/terraform-provider-azurerm/issues/27678))

ENHANCEMENTS:

* `redisenterprise` - upgrade api version to `2024-06-01-preview`  ([#27597](https://github.com/hashicorp/terraform-provider-azurerm/issues/27597))
* `azurerm_app_configuration` - support for premium sku ([#27674](https://github.com/hashicorp/terraform-provider-azurerm/issues/27674))
* `azurerm_container_app` - support for the `max_inactive_revisions` property ([#27598](https://github.com/hashicorp/terraform-provider-azurerm/issues/27598))
* `azurerm_kubernetes_cluster` - remove lock on subnets ([#27583](https://github.com/hashicorp/terraform-provider-azurerm/issues/27583))
* `azurerm_nginx_deployment` - allow updates for `sku` ([#27604](https://github.com/hashicorp/terraform-provider-azurerm/issues/27604))
* `azurerm_fluid_relay_server` - support for the `customer_managed_key` property ([#27581](https://github.com/hashicorp/terraform-provider-azurerm/issues/27581))
* `azurerm_linux_virtual_machine` - support the `UBUNTU_PRO` value for the `license_type` property ([#27534](https://github.com/hashicorp/terraform-provider-azurerm/issues/27534))


BUGS:

* `azurerm_api_management_api_diagnostic` - do not set `OperationNameFormat` when the `identifier` property is `azuremonitor` ([#27456](https://github.com/hashicorp/terraform-provider-azurerm/issues/27456))
* `azurerm_api_management` - prevent a panic ([#27649](https://github.com/hashicorp/terraform-provider-azurerm/issues/27649))
* `azurerm_mssql_database` - make `short_term_retention_policy.backup_interval_in_hours` computed ([#27656](https://github.com/hashicorp/terraform-provider-azurerm/issues/27656))

## 4.5.0 (October 10, 2024)

FEATURES:

* **New Resource**: `azurerm_stack_hci_virtual_hard_disk` ([#27474](https://github.com/hashicorp/terraform-provider-azurerm/issues/27474))

ENHANCEMENTS:

* `azurerm_bastion_host` - support for the `Premium` SKU and `session_recording_enabled` property ([#27278](https://github.com/hashicorp/terraform-provider-azurerm/issues/27278))
* `azurerm_log_analytics_cluster` - the `size_gb` property now supports all of 100, 200, 300, 400, 500, 1000, 2000, 5000, 10000, 25000, and 50000 ([#27616](https://github.com/hashicorp/terraform-provider-azurerm/issues/27616))
* `azurerm_mssql_elasticpool` - allow `PRMS` for the `family` property ([#27615](https://github.com/hashicorp/terraform-provider-azurerm/issues/27615))


BUG FIXES:

* `azurerm_mssql_database` - now creates successfully when elastic pool is hyperscale ([#27505](https://github.com/hashicorp/terraform-provider-azurerm/issues/27505))
* `azurerm_postgresql_flexible_server_configuration` - now locks to prevent conflicts when deploying multiple ([#27355](https://github.com/hashicorp/terraform-provider-azurerm/issues/27355))


## 4.4.0 (October 04, 2024)

ENHANCEMENTS: 

* dependencies - update `github.com/hashicorp/go-azure-sdk` to `v0.20240923.1151247` ([#27491](https://github.com/hashicorp/terraform-provider-azurerm/issues/27491))
* `azurerm_site_recovery_replicated_vm` - support for the `target_virtual_machine_size` property ([#27480](https://github.com/hashicorp/terraform-provider-azurerm/issues/27480))

BUG FIXES:

* `azurerm_app_service_certificate` - `key_vault_secret_id` can now be versionless ([#27537](https://github.com/hashicorp/terraform-provider-azurerm/issues/27537))
* `azurerm_linux_virtual_machine_scale_set` - prevent crash when `auto_upgrade_minor_version_enabled` is nil ([#27353](https://github.com/hashicorp/terraform-provider-azurerm/issues/27353))
* `azurerm_role_assignment` - correctly parse ID when it's a root or provider scope ([#27237](https://github.com/hashicorp/terraform-provider-azurerm/issues/27237))
* `azurerm_storage_blob` - `source_content` is now ForceNew ([#27508](https://github.com/hashicorp/terraform-provider-azurerm/issues/27508))
* `azurerm_virtual_network_gateway_connection` - revert `shared_key` to Optional and Computed ([#27560](https://github.com/hashicorp/terraform-provider-azurerm/issues/27560))

## 4.3.0 (September 19, 2024)

FEATURES:

* **New Resource**: `azurerm_advisor_suppression` ([#26177](https://github.com/hashicorp/terraform-provider-azurerm/issues/26177))
* **New Resource**: `azurerm_data_protection_backup_policy_mysql_flexible_server` ([#26955](https://github.com/hashicorp/terraform-provider-azurerm/issues/26955))
* **New Resource**: `azurerm_key_vault_managed_hardware_security_module_key_rotation_policy` ([#27306](https://github.com/hashicorp/terraform-provider-azurerm/issues/27306))
* **New Resource**: `azurerm_stack_hci_deployment_setting` ([#25646](https://github.com/hashicorp/terraform-provider-azurerm/issues/25646))
* **New Resource**: `azurerm_stack_hci_storage_path` ([#26509](https://github.com/hashicorp/terraform-provider-azurerm/issues/26509))
* **New Data Source**: `azurerm_vpn_server_configuration` ([#27054](https://github.com/hashicorp/terraform-provider-azurerm/issues/27054))

ENHANCEMENTS: 

* `managementgroups` - migrate to `hashicorp/go-azure-sdk` ([#26430](https://github.com/hashicorp/terraform-provider-azurerm/issues/26430))
* `nginx` - upgrade api version to `2024-06-01-preview`  ([#27345](https://github.com/hashicorp/terraform-provider-azurerm/issues/27345))
* `azurerm_linux[windows]_web[function]_app[app_slot]` - upgrade api version from `2023-01-01` to `2023-12-01` ([#27196](https://github.com/hashicorp/terraform-provider-azurerm/issues/27196))
* `azurerm_cosmosdb_account` - support for the capability `EnableNoSQLVectorSearch` ([#27357](https://github.com/hashicorp/terraform-provider-azurerm/issues/27357))azurerm_container_app_custom_domain - fix parsing the certificate ID error #25972
* `azurerm_container_app_custom_domain` - support other certificate types ([#25972](https://github.com/hashicorp/terraform-provider-azurerm/issues/25972))
* `azurerm_linux_virtual_machine_scale_set` - the `zones` property can now be updated without creating a new resource ([#27288](https://github.com/hashicorp/terraform-provider-azurerm/issues/27288))
* `azurerm_orchestrated_virtual_machine_scale_set` - the `zones` property can now be updated without creating a new resource ([#27288](https://github.com/hashicorp/terraform-provider-azurerm/issues/27288))
* `azurerm_role_management_policy` - support for resource scope ([#27205](https://github.com/hashicorp/terraform-provider-azurerm/issues/27205))
* `azurerm_spring_cloud_gateway` - changing the `environment_variables` and `sensitive_environment_variables` properties no longer creates a new resource ([#27404](https://github.com/hashicorp/terraform-provider-azurerm/issues/27404))
* `azurerm_static_web_app` - support for the `public_network_access_enabled` property ([#26345](https://github.com/hashicorp/terraform-provider-azurerm/issues/26345))
* `azurerm_shared_image` - support for the `disk_controller_type_nvme_enabled` property ([#26370](https://github.com/hashicorp/terraform-provider-azurerm/issues/26370))
* `azurerm_storage_blob` - changing the `source` property no longer creates a new resource ([#27394](https://github.com/hashicorp/terraform-provider-azurerm/issues/27394))
* `azurerm_storage_object_replication` - changing the `rules.x. source_container_name` and `rules.x. destination_container_name` properties no longer creates a new resource ([#27394](https://github.com/hashicorp/terraform-provider-azurerm/issues/27394))
* `azurerm_windows_virtual_machine_scale_set` - the `zones` property can now be updated without creating a new resource ([#27288](https://github.com/hashicorp/terraform-provider-azurerm/issues/27288)) 

BUG FIXES:

* `azurerm_application_insights` - fix crash when read for `DataVolumeCap` is `nil` ([#27352](https://github.com/hashicorp/terraform-provider-azurerm/issues/27352))
* `azurerm_container_app` - relax validation on the ingress traffic property ([#27396](https://github.com/hashicorp/terraform-provider-azurerm/issues/27396))
* `azurerm_log_analytics_workspace_table` - will now correctly set `total_retention_in_days` when `sku` is `Basic` ([#27420](https://github.com/hashicorp/terraform-provider-azurerm/issues/27420))

## 4.2.0 (September 12, 2024)

FEATURES:

* **New Resource**: `azurerm_arc_machine` ([#26647](https://github.com/hashicorp/terraform-provider-azurerm/issues/26647))
* **New Resource**: `azurerm_arc_machine_automanage_configuration_assignment` ([#26657](https://github.com/hashicorp/terraform-provider-azurerm/issues/26657)) 

ENHANCEMENTS:

* `network/bastionhosts` - upgrade api version from `2023-11-01` to `2024-01-01` ([#27277](https://github.com/hashicorp/terraform-provider-azurerm/issues/27277))
* `recoveryservices` - upgrade `recoveryservicessiterecovery` from `2022-10-0`1 to `2024-04-01` ([#27281](https://github.com/hashicorp/terraform-provider-azurerm/issues/27281))
* `azurerm_data_protection_backup_vault` - support for the `property cross_region_restore_enabled` property ([#27197](https://github.com/hashicorp/terraform-provider-azurerm/issues/27197))
* `azurem_mssql_managed_instance` - support for the `service_principal_type` property ([#27240](https://github.com/hashicorp/terraform-provider-azurerm/issues/27240))

BUG FIXES:

* `azurerm_cosmosdb_account` - fix crash during state migration ([#27302](https://github.com/hashicorp/terraform-provider-azurerm/issues/27302))
* `azurerm_servicebus_queue` - fix defaults of the `default_message_ttl` and `auto_delete_on_idle` properties ([#27305](https://github.com/hashicorp/terraform-provider-azurerm/issues/27305))

## 4.1.0 (September 05, 2024)

ENHANCEMENTS:

* dependencies - bump `hashicorp/go-azure-sdk` to `v0.20240903.1111904` ([#27268](https://github.com/hashicorp/terraform-provider-azurerm/issues/27268))
* Virtual Machine Scale Sets - upgrade api version from `2024-03-01` to `2024-07-01` ([#27230](https://github.com/hashicorp/terraform-provider-azurerm/issues/27230))
* `hdinsights` - update the HDInsights Node definition validation of VM sizes to include new V5 types ([#27270](https://github.com/hashicorp/terraform-provider-azurerm/issues/27270))
* `azurerm_api_management_logger` - support for the `application_insights.connection_string` property ([#27137](https://github.com/hashicorp/terraform-provider-azurerm/issues/27137))
* `azurerm_bot_service_azure_bot` - will now send the value for the `developer_app_insights_api_key` property ([#27280](https://github.com/hashicorp/terraform-provider-azurerm/issues/27280))
* `azurerm_netapp_volume` - support for the `smb3_protocol_encryption_enabled` property ([#27228](https://github.com/hashicorp/terraform-provider-azurerm/issues/27228))
* `azurerm_subnet` - support `Microsoft.DevOpsInfrastructure` as delegation service ([#27259](https://github.com/hashicorp/terraform-provider-azurerm/issues/27259))

BUG FIXES:

* `azurerm_mysql_flexible_server` - correctly set `source_server_id` in the state file ([#27295](https://github.com/hashicorp/terraform-provider-azurerm/issues/27295))
* `azurerm_cosmosdb_account` - the `ip_range_filter` property now supports IPV4 addresses ([#27208](https://github.com/hashicorp/terraform-provider-azurerm/issues/27208))
* `azurerm_cosmosdb_account` - added state migration for `ip_range_filter` underlying type change from `string` to `set` ([#27276](https://github.com/hashicorp/terraform-provider-azurerm/issues/27276))
* `azurerm_linux_virtual_machine` - the `admin_ssh_key.public_key` property now supports ed25519 ssh keys ([#27202](https://github.com/hashicorp/terraform-provider-azurerm/issues/27202))
* `azurerm_sentinel_automation_rule` - no longer panics when using `condition_json`  ([#27269](https://github.com/hashicorp/terraform-provider-azurerm/issues/27269))
* `azurerm_kubernetes_cluster` -  the `host_encryption_enabled` and `node_public_ip_enabled` properties are now set correctly ([#27218](https://github.com/hashicorp/terraform-provider-azurerm/issues/27218))

## 4.0.1 (August 23, 2024)

BUG FIXES:

* provider: fix a validation bug that prevents `terraform validate` from working when `subscription_id` is not specified ([#27178](https://github.com/hashicorp/terraform-provider-azurerm/issues/27178))
* `azurerm_cognitive_deployment` - fixed replacement of `scale` block with `sku` ([#27173](https://github.com/hashicorp/terraform-provider-azurerm/issues/27173))
* `azurerm_kubernetes_cluster` - prevent a panic ([#27183](https://github.com/hashicorp/terraform-provider-azurerm/issues/27183))
* `azurerm_kubernetes_cluster_node_pool` - prevent a panic caused by renamed `enable_*` properties ([#27164](https://github.com/hashicorp/terraform-provider-azurerm/issues/27164))
* `azurerm_sentinel_data_connector_microsoft_threat_intelligence` - prevent error by removing deprecated property `bing_safety_phishing_url_lookback_date` ([#27171](https://github.com/hashicorp/terraform-provider-azurerm/issues/27171))

## 4.0.0 (August 22, 2024)

NOTES:

* **Major Version**: Version 4.0 of the Azure Provider is a major version - some behaviours have changed and some deprecated fields/resources have been removed - please refer to [the 4.0 upgrade guide for more information](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/4.0-upgrade-guide).
* When upgrading to v4.0 of the AzureRM Provider, we recommend upgrading to the latest version of Terraform Core ([which can be found here](https://www.terraform.io/downloads)).

ENHANCEMENTS:

* Data Source: `azurerm_shared_image` - add support for the `trusted_launch_supported`, `trusted_launch_enabled`, `confidential_vm_supported`, `confidential_vm_enabled`, `accelerated_network_support_enabled` and `hibernation_enabled` properties ([#26975](https://github.com/hashicorp/terraform-provider-azurerm/issues/26975))
* dependencies: updating `hashicorp/go-azure-sdk` to `v0.20240819.1075239` ([#27107](https://github.com/hashicorp/terraform-provider-azurerm/issues/27107))
* `applicationgateways` - updating to use `2023-11-01` ([#26776](https://github.com/hashicorp/terraform-provider-azurerm/issues/26776))
* `containerregistry` - updating to use `2023-06-01-preview` ([#23393](https://github.com/hashicorp/terraform-provider-azurerm/issues/23393))
* `containerservice` - updating to `2024-05-01` ([#27105](https://github.com/hashicorp/terraform-provider-azurerm/issues/27105))
* `mssql` - updating to use `hashicorp/go-azure-sdk` and `023-08-01-preview` ([#27073](https://github.com/hashicorp/terraform-provider-azurerm/issues/27073))
* `mssqlmanagedinstance` - updating to use `hashicorp/go-azure-sdk` and `2023-08-01-preview` ([#26872](https://github.com/hashicorp/terraform-provider-azurerm/issues/26872))
* `azurerm_image` - add support for the `disk_encryption_set_id` property to the `data_disk` block ([#27015](https://github.com/hashicorp/terraform-provider-azurerm/issues/27015))
* `azurerm_log_analytics_workspace_table` - add support for more `total_retention_in_days` and `retention_in_days` values ([#27053](https://github.com/hashicorp/terraform-provider-azurerm/issues/27053))
* `azurerm_mssql_elasticpool` - add support for the `HS_MOPRMS` and `MOPRMS` skus ([#27085](https://github.com/hashicorp/terraform-provider-azurerm/issues/27085))
* `azurerm_netapp_pool` - allow `1` as a valid value for `size_in_tb` ([#27095](https://github.com/hashicorp/terraform-provider-azurerm/issues/27095))
* `azurerm_notification_hub` - add support for the `browser_credential` property ([#27058](https://github.com/hashicorp/terraform-provider-azurerm/issues/27058))
* `azurerm_redis_cache` - add support for the `access_keys_authentication_enabled` property ([#27039](https://github.com/hashicorp/terraform-provider-azurerm/issues/27039))
* `azurerm_role_assignment` - add support for the `/`, `/providers/Microsoft.Capacity` and `/providers/Microsoft.BillingBenefits` scopes ([#26663](https://github.com/hashicorp/terraform-provider-azurerm/issues/26663))
* `azurerm_shared_image` - add support for the `hibernation_enabled` property ([#26975](https://github.com/hashicorp/terraform-provider-azurerm/issues/26975))
* `azurerm_storage_account` - support `queue_encryption_key_type` and `table_encryption_key_type` for more storage account kinds ([#27112](https://github.com/hashicorp/terraform-provider-azurerm/issues/27112))
* `azurerm_web_application_firewall_policy` - add support for the `request_body_enforcement` property ([#27094](https://github.com/hashicorp/terraform-provider-azurerm/issues/27094))

BUG FIXES:

* `azurerm_ip_group_cidr` - fixed the position of the CIDR check to correctly refresh the resource when it's no longer present ([#27103](https://github.com/hashicorp/terraform-provider-azurerm/issues/27103))
* `azurerm_monitor_diagnostic_setting` - add further polling to work around an eventual consistency issue when creating the resource ([#27088](https://github.com/hashicorp/terraform-provider-azurerm/issues/27088))
* `azurerm_storage_account` - prevent API error by populating `infrastructure_encryption_enabled` when updating `customer_managed_key` ([#26971](https://github.com/hashicorp/terraform-provider-azurerm/issues/26971))
* `azurerm_storage_blob_inventory_policy` - the `filter` property can now be set when `scope` is `container` ([#27113](https://github.com/hashicorp/terraform-provider-azurerm/issues/27113))
* `azurerm_virtual_network_dns_servers` - moved locks to prevent the creation of subnets with stale data ([#27036](https://github.com/hashicorp/terraform-provider-azurerm/issues/27036))
* `azurerm_virtual_network_gateway_connection` - allow `0` as a valid value for `ipsec_policy.sa_datasize` ([#27056](https://github.com/hashicorp/terraform-provider-azurerm/issues/27056))

---

For information on changes between the v3.116.0 and v3.0.0 releases, please see [the previous v3.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v3.md).

For information on changes between the v2.99.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v1.44.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
