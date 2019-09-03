package azurerm

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/common"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	// NOTE: as part of migrating Data Sources/Resources into Packages - we should be able to use
	// the Service Registration interface to gradually migrate Data Sources/Resources over to the
	// new pattern.
	// However this requires that the following be done first:
	//  1. (DONE) Migrating the top level functions into the internal package
	//	2. (DONE) Finish migrating the SDK Clients into Packages
	//	3. Switch the remaining resources over to the new Storage SDK
	//		(so we can remove `getBlobStorageClientForStorageAccount` from `config.go`)
	//	4. Making the SDK Clients public in the ArmClient prior to moving
	//  5. Introducing a parent struct which becomes a nested field in `config.go`
	//  	for those properties, to ease migration (probably internal/common/clients.go)
	//	6. Migrating references from the `ArmClient` to the new parent client
	//
	// For the moment/until that's done, we'll have to continue defining these inline
	supportedServices := []common.ServiceRegistration{}

	dataSources := map[string]*schema.Resource{
		"azurerm_api_management":                         dataSourceApiManagementService(),
		"azurerm_api_management_api":                     dataSourceApiManagementApi(),
		"azurerm_api_management_group":                   dataSourceApiManagementGroup(),
		"azurerm_api_management_product":                 dataSourceApiManagementProduct(),
		"azurerm_api_management_user":                    dataSourceArmApiManagementUser(),
		"azurerm_app_service_plan":                       dataSourceAppServicePlan(),
		"azurerm_app_service":                            dataSourceArmAppService(),
		"azurerm_application_insights":                   dataSourceArmApplicationInsights(),
		"azurerm_application_security_group":             dataSourceArmApplicationSecurityGroup(),
		"azurerm_automation_variable_bool":               dataSourceArmAutomationVariableBool(),
		"azurerm_automation_variable_datetime":           dataSourceArmAutomationVariableDateTime(),
		"azurerm_automation_variable_int":                dataSourceArmAutomationVariableInt(),
		"azurerm_automation_variable_string":             dataSourceArmAutomationVariableString(),
		"azurerm_availability_set":                       dataSourceArmAvailabilitySet(),
		"azurerm_azuread_application":                    dataSourceArmAzureADApplication(),
		"azurerm_azuread_service_principal":              dataSourceArmActiveDirectoryServicePrincipal(),
		"azurerm_batch_account":                          dataSourceArmBatchAccount(),
		"azurerm_batch_certificate":                      dataSourceArmBatchCertificate(),
		"azurerm_batch_pool":                             dataSourceArmBatchPool(),
		"azurerm_builtin_role_definition":                dataSourceArmBuiltInRoleDefinition(),
		"azurerm_cdn_profile":                            dataSourceArmCdnProfile(),
		"azurerm_client_config":                          dataSourceArmClientConfig(),
		"azurerm_kubernetes_service_versions":            dataSourceArmKubernetesServiceVersions(),
		"azurerm_container_registry":                     dataSourceArmContainerRegistry(),
		"azurerm_cosmosdb_account":                       dataSourceArmCosmosDbAccount(),
		"azurerm_data_lake_store":                        dataSourceArmDataLakeStoreAccount(),
		"azurerm_dev_test_lab":                           dataSourceArmDevTestLab(),
		"azurerm_dev_test_virtual_network":               dataSourceArmDevTestVirtualNetwork(),
		"azurerm_dns_zone":                               dataSourceArmDnsZone(),
		"azurerm_eventhub_namespace":                     dataSourceEventHubNamespace(),
		"azurerm_express_route_circuit":                  dataSourceArmExpressRouteCircuit(),
		"azurerm_firewall":                               dataSourceArmFirewall(),
		"azurerm_image":                                  dataSourceArmImage(),
		"azurerm_hdinsight_cluster":                      dataSourceArmHDInsightSparkCluster(),
		"azurerm_maps_account":                           dataSourceArmMapsAccount(),
		"azurerm_key_vault_access_policy":                dataSourceArmKeyVaultAccessPolicy(),
		"azurerm_key_vault_key":                          dataSourceArmKeyVaultKey(),
		"azurerm_key_vault_secret":                       dataSourceArmKeyVaultSecret(),
		"azurerm_key_vault":                              dataSourceArmKeyVault(),
		"azurerm_kubernetes_cluster":                     dataSourceArmKubernetesCluster(),
		"azurerm_lb":                                     dataSourceArmLoadBalancer(),
		"azurerm_lb_backend_address_pool":                dataSourceArmLoadBalancerBackendAddressPool(),
		"azurerm_log_analytics_workspace":                dataSourceLogAnalyticsWorkspace(),
		"azurerm_logic_app_workflow":                     dataSourceArmLogicAppWorkflow(),
		"azurerm_managed_disk":                           dataSourceArmManagedDisk(),
		"azurerm_management_group":                       dataSourceArmManagementGroup(),
		"azurerm_monitor_action_group":                   dataSourceArmMonitorActionGroup(),
		"azurerm_monitor_diagnostic_categories":          dataSourceArmMonitorDiagnosticCategories(),
		"azurerm_monitor_log_profile":                    dataSourceArmMonitorLogProfile(),
		"azurerm_mssql_elasticpool":                      dataSourceArmMsSqlElasticpool(),
		"azurerm_network_interface":                      dataSourceArmNetworkInterface(),
		"azurerm_network_security_group":                 dataSourceArmNetworkSecurityGroup(),
		"azurerm_network_watcher":                        dataSourceArmNetworkWatcher(),
		"azurerm_notification_hub_namespace":             dataSourceNotificationHubNamespace(),
		"azurerm_notification_hub":                       dataSourceNotificationHub(),
		"azurerm_platform_image":                         dataSourceArmPlatformImage(),
		"azurerm_policy_definition":                      dataSourceArmPolicyDefinition(),
		"azurerm_public_ip":                              dataSourceArmPublicIP(),
		"azurerm_public_ips":                             dataSourceArmPublicIPs(),
		"azurerm_recovery_services_vault":                dataSourceArmRecoveryServicesVault(),
		"azurerm_recovery_services_protection_policy_vm": dataSourceArmRecoveryServicesProtectionPolicyVm(),
		"azurerm_redis_cache":                            dataSourceArmRedisCache(),
		"azurerm_resource_group":                         dataSourceArmResourceGroup(),
		"azurerm_role_definition":                        dataSourceArmRoleDefinition(),
		"azurerm_route_table":                            dataSourceArmRouteTable(),
		"azurerm_scheduler_job_collection":               dataSourceArmSchedulerJobCollection(),
		"azurerm_servicebus_namespace":                   dataSourceArmServiceBusNamespace(),
		"azurerm_shared_image_gallery":                   dataSourceArmSharedImageGallery(),
		"azurerm_shared_image_version":                   dataSourceArmSharedImageVersion(),
		"azurerm_shared_image":                           dataSourceArmSharedImage(),
		"azurerm_snapshot":                               dataSourceArmSnapshot(),
		"azurerm_sql_server":                             dataSourceSqlServer(),
		"azurerm_sql_database":                           dataSourceSqlDatabase(),
		"azurerm_stream_analytics_job":                   dataSourceArmStreamAnalyticsJob(),
		"azurerm_storage_account_sas":                    dataSourceArmStorageAccountSharedAccessSignature(),
		"azurerm_storage_account":                        dataSourceArmStorageAccount(),
		"azurerm_subnet":                                 dataSourceArmSubnet(),
		"azurerm_subscription":                           dataSourceArmSubscription(),
		"azurerm_subscriptions":                          dataSourceArmSubscriptions(),
		"azurerm_traffic_manager_geographical_location":  dataSourceArmTrafficManagerGeographicalLocation(),
		"azurerm_user_assigned_identity":                 dataSourceArmUserAssignedIdentity(),
		"azurerm_virtual_machine":                        dataSourceArmVirtualMachine(),
		"azurerm_virtual_network_gateway":                dataSourceArmVirtualNetworkGateway(),
		"azurerm_virtual_network_gateway_connection":     dataSourceArmVirtualNetworkGatewayConnection(),
		"azurerm_virtual_network":                        dataSourceArmVirtualNetwork(),
	}

	resources := map[string]*schema.Resource{
		"azurerm_analysis_services_server":                           resourceArmAnalysisServicesServer(),
		"azurerm_api_management":                                     resourceArmApiManagementService(),
		"azurerm_api_management_api":                                 resourceArmApiManagementApi(),
		"azurerm_api_management_api_operation":                       resourceArmApiManagementApiOperation(),
		"azurerm_api_management_api_operation_policy":                resourceArmApiManagementApiOperationPolicy(),
		"azurerm_api_management_api_policy":                          resourceArmApiManagementApiPolicy(),
		"azurerm_api_management_api_schema":                          resourceArmApiManagementApiSchema(),
		"azurerm_api_management_api_version_set":                     resourceArmApiManagementApiVersionSet(),
		"azurerm_api_management_authorization_server":                resourceArmApiManagementAuthorizationServer(),
		"azurerm_api_management_backend":                             resourceArmApiManagementBackend(),
		"azurerm_api_management_certificate":                         resourceArmApiManagementCertificate(),
		"azurerm_api_management_group":                               resourceArmApiManagementGroup(),
		"azurerm_api_management_group_user":                          resourceArmApiManagementGroupUser(),
		"azurerm_api_management_logger":                              resourceArmApiManagementLogger(),
		"azurerm_api_management_openid_connect_provider":             resourceArmApiManagementOpenIDConnectProvider(),
		"azurerm_api_management_product":                             resourceArmApiManagementProduct(),
		"azurerm_api_management_product_api":                         resourceArmApiManagementProductApi(),
		"azurerm_api_management_product_group":                       resourceArmApiManagementProductGroup(),
		"azurerm_api_management_product_policy":                      resourceArmApiManagementProductPolicy(),
		"azurerm_api_management_property":                            resourceArmApiManagementProperty(),
		"azurerm_api_management_subscription":                        resourceArmApiManagementSubscription(),
		"azurerm_api_management_user":                                resourceArmApiManagementUser(),
		"azurerm_app_service_active_slot":                            resourceArmAppServiceActiveSlot(),
		"azurerm_app_service_certificate":                            resourceArmAppServiceCertificate(),
		"azurerm_app_service_custom_hostname_binding":                resourceArmAppServiceCustomHostnameBinding(),
		"azurerm_app_service_plan":                                   resourceArmAppServicePlan(),
		"azurerm_app_service_slot":                                   resourceArmAppServiceSlot(),
		"azurerm_app_service":                                        resourceArmAppService(),
		"azurerm_application_gateway":                                resourceArmApplicationGateway(),
		"azurerm_application_insights_api_key":                       resourceArmApplicationInsightsAPIKey(),
		"azurerm_application_insights":                               resourceArmApplicationInsights(),
		"azurerm_application_insights_web_test":                      resourceArmApplicationInsightsWebTests(),
		"azurerm_application_security_group":                         resourceArmApplicationSecurityGroup(),
		"azurerm_automation_account":                                 resourceArmAutomationAccount(),
		"azurerm_automation_credential":                              resourceArmAutomationCredential(),
		"azurerm_automation_dsc_configuration":                       resourceArmAutomationDscConfiguration(),
		"azurerm_automation_dsc_nodeconfiguration":                   resourceArmAutomationDscNodeConfiguration(),
		"azurerm_automation_module":                                  resourceArmAutomationModule(),
		"azurerm_automation_runbook":                                 resourceArmAutomationRunbook(),
		"azurerm_automation_schedule":                                resourceArmAutomationSchedule(),
		"azurerm_automation_variable_bool":                           resourceArmAutomationVariableBool(),
		"azurerm_automation_variable_datetime":                       resourceArmAutomationVariableDateTime(),
		"azurerm_automation_variable_int":                            resourceArmAutomationVariableInt(),
		"azurerm_automation_variable_string":                         resourceArmAutomationVariableString(),
		"azurerm_autoscale_setting":                                  resourceArmAutoScaleSetting(),
		"azurerm_availability_set":                                   resourceArmAvailabilitySet(),
		"azurerm_azuread_application":                                resourceArmActiveDirectoryApplication(),
		"azurerm_azuread_service_principal_password":                 resourceArmActiveDirectoryServicePrincipalPassword(),
		"azurerm_azuread_service_principal":                          resourceArmActiveDirectoryServicePrincipal(),
		"azurerm_batch_account":                                      resourceArmBatchAccount(),
		"azurerm_batch_application":                                  resourceArmBatchApplication(),
		"azurerm_batch_certificate":                                  resourceArmBatchCertificate(),
		"azurerm_batch_pool":                                         resourceArmBatchPool(),
		"azurerm_cdn_endpoint":                                       resourceArmCdnEndpoint(),
		"azurerm_cdn_profile":                                        resourceArmCdnProfile(),
		"azurerm_cognitive_account":                                  resourceArmCognitiveAccount(),
		"azurerm_connection_monitor":                                 resourceArmConnectionMonitor(),
		"azurerm_container_group":                                    resourceArmContainerGroup(),
		"azurerm_container_registry_webhook":                         resourceArmContainerRegistryWebhook(),
		"azurerm_container_registry":                                 resourceArmContainerRegistry(),
		"azurerm_container_service":                                  resourceArmContainerService(),
		"azurerm_cosmosdb_account":                                   resourceArmCosmosDbAccount(),
		"azurerm_cosmosdb_cassandra_keyspace":                        resourceArmCosmosDbCassandraKeyspace(),
		"azurerm_cosmosdb_mongo_collection":                          resourceArmCosmosDbMongoCollection(),
		"azurerm_cosmosdb_mongo_database":                            resourceArmCosmosDbMongoDatabase(),
		"azurerm_cosmosdb_sql_container":                             resourceArmCosmosDbSQLContainer(),
		"azurerm_cosmosdb_sql_database":                              resourceArmCosmosDbSQLDatabase(),
		"azurerm_cosmosdb_table":                                     resourceArmCosmosDbTable(),
		"azurerm_data_factory":                                       resourceArmDataFactory(),
		"azurerm_data_factory_dataset_mysql":                         resourceArmDataFactoryDatasetMySQL(),
		"azurerm_data_factory_dataset_postgresql":                    resourceArmDataFactoryDatasetPostgreSQL(),
		"azurerm_data_factory_dataset_sql_server_table":              resourceArmDataFactoryDatasetSQLServerTable(),
		"azurerm_data_factory_linked_service_data_lake_storage_gen2": resourceArmDataFactoryLinkedServiceDataLakeStorageGen2(),
		"azurerm_data_factory_linked_service_mysql":                  resourceArmDataFactoryLinkedServiceMySQL(),
		"azurerm_data_factory_linked_service_postgresql":             resourceArmDataFactoryLinkedServicePostgreSQL(),
		"azurerm_data_factory_linked_service_sql_server":             resourceArmDataFactoryLinkedServiceSQLServer(),
		"azurerm_data_factory_pipeline":                              resourceArmDataFactoryPipeline(),
		"azurerm_data_lake_analytics_account":                        resourceArmDataLakeAnalyticsAccount(),
		"azurerm_data_lake_analytics_firewall_rule":                  resourceArmDataLakeAnalyticsFirewallRule(),
		"azurerm_data_lake_store_file":                               resourceArmDataLakeStoreFile(),
		"azurerm_data_lake_store_firewall_rule":                      resourceArmDataLakeStoreFirewallRule(),
		"azurerm_data_lake_store":                                    resourceArmDataLakeStore(),
		"azurerm_databricks_workspace":                               resourceArmDatabricksWorkspace(),
		"azurerm_ddos_protection_plan":                               resourceArmDDoSProtectionPlan(),
		"azurerm_dev_test_lab":                                       resourceArmDevTestLab(),
		"azurerm_dev_test_schedule":                                  resourceArmDevTestLabSchedules(),
		"azurerm_dev_test_linux_virtual_machine":                     resourceArmDevTestLinuxVirtualMachine(),
		"azurerm_dev_test_policy":                                    resourceArmDevTestPolicy(),
		"azurerm_dev_test_virtual_network":                           resourceArmDevTestVirtualNetwork(),
		"azurerm_dev_test_windows_virtual_machine":                   resourceArmDevTestWindowsVirtualMachine(),
		"azurerm_devspace_controller":                                resourceArmDevSpaceController(),
		"azurerm_dns_a_record":                                       resourceArmDnsARecord(),
		"azurerm_dns_aaaa_record":                                    resourceArmDnsAAAARecord(),
		"azurerm_dns_caa_record":                                     resourceArmDnsCaaRecord(),
		"azurerm_dns_cname_record":                                   resourceArmDnsCNameRecord(),
		"azurerm_dns_mx_record":                                      resourceArmDnsMxRecord(),
		"azurerm_dns_ns_record":                                      resourceArmDnsNsRecord(),
		"azurerm_dns_ptr_record":                                     resourceArmDnsPtrRecord(),
		"azurerm_dns_srv_record":                                     resourceArmDnsSrvRecord(),
		"azurerm_dns_txt_record":                                     resourceArmDnsTxtRecord(),
		"azurerm_dns_zone":                                           resourceArmDnsZone(),
		"azurerm_eventgrid_domain":                                   resourceArmEventGridDomain(),
		"azurerm_eventgrid_event_subscription":                       resourceArmEventGridEventSubscription(),
		"azurerm_eventgrid_topic":                                    resourceArmEventGridTopic(),
		"azurerm_eventhub_authorization_rule":                        resourceArmEventHubAuthorizationRule(),
		"azurerm_eventhub_consumer_group":                            resourceArmEventHubConsumerGroup(),
		"azurerm_eventhub_namespace_authorization_rule":              resourceArmEventHubNamespaceAuthorizationRule(),
		"azurerm_eventhub_namespace":                                 resourceArmEventHubNamespace(),
		"azurerm_eventhub":                                           resourceArmEventHub(),
		"azurerm_express_route_circuit_authorization":                resourceArmExpressRouteCircuitAuthorization(),
		"azurerm_express_route_circuit_peering":                      resourceArmExpressRouteCircuitPeering(),
		"azurerm_express_route_circuit":                              resourceArmExpressRouteCircuit(),
		"azurerm_firewall_application_rule_collection":               resourceArmFirewallApplicationRuleCollection(),
		"azurerm_firewall_nat_rule_collection":                       resourceArmFirewallNatRuleCollection(),
		"azurerm_firewall_network_rule_collection":                   resourceArmFirewallNetworkRuleCollection(),
		"azurerm_firewall":                                           resourceArmFirewall(),
		"azurerm_function_app":                                       resourceArmFunctionApp(),
		"azurerm_hdinsight_hadoop_cluster":                           resourceArmHDInsightHadoopCluster(),
		"azurerm_hdinsight_hbase_cluster":                            resourceArmHDInsightHBaseCluster(),
		"azurerm_hdinsight_interactive_query_cluster":                resourceArmHDInsightInteractiveQueryCluster(),
		"azurerm_hdinsight_kafka_cluster":                            resourceArmHDInsightKafkaCluster(),
		"azurerm_hdinsight_ml_services_cluster":                      resourceArmHDInsightMLServicesCluster(),
		"azurerm_hdinsight_rserver_cluster":                          resourceArmHDInsightRServerCluster(),
		"azurerm_hdinsight_spark_cluster":                            resourceArmHDInsightSparkCluster(),
		"azurerm_hdinsight_storm_cluster":                            resourceArmHDInsightStormCluster(),
		"azurerm_image":                                              resourceArmImage(),
		"azurerm_iot_dps":                                            resourceArmIotDPS(),
		"azurerm_iot_dps_certificate":                                resourceArmIotDPSCertificate(),
		"azurerm_iothub_consumer_group":                              resourceArmIotHubConsumerGroup(),
		"azurerm_iothub":                                             resourceArmIotHub(),
		"azurerm_iothub_shared_access_policy":                        resourceArmIotHubSharedAccessPolicy(),
		"azurerm_key_vault_access_policy":                            resourceArmKeyVaultAccessPolicy(),
		"azurerm_key_vault_certificate":                              resourceArmKeyVaultCertificate(),
		"azurerm_key_vault_key":                                      resourceArmKeyVaultKey(),
		"azurerm_key_vault_secret":                                   resourceArmKeyVaultSecret(),
		"azurerm_key_vault":                                          resourceArmKeyVault(),
		"azurerm_kubernetes_cluster":                                 resourceArmKubernetesCluster(),
		"azurerm_kusto_cluster":                                      resourceArmKustoCluster(),
		"azurerm_lb_backend_address_pool":                            resourceArmLoadBalancerBackendAddressPool(),
		"azurerm_lb_nat_pool":                                        resourceArmLoadBalancerNatPool(),
		"azurerm_lb_nat_rule":                                        resourceArmLoadBalancerNatRule(),
		"azurerm_lb_probe":                                           resourceArmLoadBalancerProbe(),
		"azurerm_lb_outbound_rule":                                   resourceArmLoadBalancerOutboundRule(),
		"azurerm_lb_rule":                                            resourceArmLoadBalancerRule(),
		"azurerm_lb":                                                 resourceArmLoadBalancer(),
		"azurerm_local_network_gateway":                              resourceArmLocalNetworkGateway(),
		"azurerm_log_analytics_solution":                             resourceArmLogAnalyticsSolution(),
		"azurerm_log_analytics_linked_service":                       resourceArmLogAnalyticsLinkedService(),
		"azurerm_log_analytics_workspace_linked_service":             resourceArmLogAnalyticsWorkspaceLinkedService(),
		"azurerm_log_analytics_workspace":                            resourceArmLogAnalyticsWorkspace(),
		"azurerm_logic_app_action_custom":                            resourceArmLogicAppActionCustom(),
		"azurerm_logic_app_action_http":                              resourceArmLogicAppActionHTTP(),
		"azurerm_logic_app_trigger_custom":                           resourceArmLogicAppTriggerCustom(),
		"azurerm_logic_app_trigger_http_request":                     resourceArmLogicAppTriggerHttpRequest(),
		"azurerm_logic_app_trigger_recurrence":                       resourceArmLogicAppTriggerRecurrence(),
		"azurerm_logic_app_workflow":                                 resourceArmLogicAppWorkflow(),
		"azurerm_managed_disk":                                       resourceArmManagedDisk(),
		"azurerm_management_group":                                   resourceArmManagementGroup(),
		"azurerm_management_lock":                                    resourceArmManagementLock(),
		"azurerm_maps_account":                                       resourceArmMapsAccount(),
		"azurerm_mariadb_configuration":                              resourceArmMariaDbConfiguration(),
		"azurerm_mariadb_database":                                   resourceArmMariaDbDatabase(),
		"azurerm_mariadb_firewall_rule":                              resourceArmMariaDBFirewallRule(),
		"azurerm_mariadb_server":                                     resourceArmMariaDbServer(),
		"azurerm_mariadb_virtual_network_rule":                       resourceArmMariaDbVirtualNetworkRule(),
		"azurerm_media_services_account":                             resourceArmMediaServicesAccount(),
		"azurerm_metric_alertrule":                                   resourceArmMetricAlertRule(),
		"azurerm_monitor_autoscale_setting":                          resourceArmMonitorAutoScaleSetting(),
		"azurerm_monitor_action_group":                               resourceArmMonitorActionGroup(),
		"azurerm_monitor_activity_log_alert":                         resourceArmMonitorActivityLogAlert(),
		"azurerm_monitor_diagnostic_setting":                         resourceArmMonitorDiagnosticSetting(),
		"azurerm_monitor_log_profile":                                resourceArmMonitorLogProfile(),
		"azurerm_monitor_metric_alert":                               resourceArmMonitorMetricAlert(),
		"azurerm_monitor_metric_alertrule":                           resourceArmMonitorMetricAlertRule(),
		"azurerm_mssql_elasticpool":                                  resourceArmMsSqlElasticPool(),
		"azurerm_mysql_configuration":                                resourceArmMySQLConfiguration(),
		"azurerm_mysql_database":                                     resourceArmMySqlDatabase(),
		"azurerm_mysql_firewall_rule":                                resourceArmMySqlFirewallRule(),
		"azurerm_mysql_server":                                       resourceArmMySqlServer(),
		"azurerm_mysql_virtual_network_rule":                         resourceArmMySqlVirtualNetworkRule(),
		"azurerm_network_connection_monitor":                         resourceArmNetworkConnectionMonitor(),
		"azurerm_network_ddos_protection_plan":                       resourceArmNetworkDDoSProtectionPlan(),
		"azurerm_network_interface":                                  resourceArmNetworkInterface(),
		"azurerm_network_interface_application_gateway_backend_address_pool_association": resourceArmNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation(),
		"azurerm_network_interface_application_security_group_association":               resourceArmNetworkInterfaceApplicationSecurityGroupAssociation(),
		"azurerm_network_interface_backend_address_pool_association":                     resourceArmNetworkInterfaceBackendAddressPoolAssociation(),
		"azurerm_network_interface_nat_rule_association":                                 resourceArmNetworkInterfaceNatRuleAssociation(),
		"azurerm_network_packet_capture":                                                 resourceArmNetworkPacketCapture(),
		"azurerm_network_profile":                                                        resourceArmNetworkProfile(),
		"azurerm_network_security_group":                                                 resourceArmNetworkSecurityGroup(),
		"azurerm_network_security_rule":                                                  resourceArmNetworkSecurityRule(),
		"azurerm_network_watcher":                                                        resourceArmNetworkWatcher(),
		"azurerm_notification_hub_authorization_rule":                                    resourceArmNotificationHubAuthorizationRule(),
		"azurerm_notification_hub_namespace":                                             resourceArmNotificationHubNamespace(),
		"azurerm_notification_hub":                                                       resourceArmNotificationHub(),
		"azurerm_packet_capture":                                                         resourceArmPacketCapture(),
		"azurerm_policy_assignment":                                                      resourceArmPolicyAssignment(),
		"azurerm_policy_definition":                                                      resourceArmPolicyDefinition(),
		"azurerm_policy_set_definition":                                                  resourceArmPolicySetDefinition(),
		"azurerm_postgresql_configuration":                                               resourceArmPostgreSQLConfiguration(),
		"azurerm_postgresql_database":                                                    resourceArmPostgreSQLDatabase(),
		"azurerm_postgresql_firewall_rule":                                               resourceArmPostgreSQLFirewallRule(),
		"azurerm_postgresql_server":                                                      resourceArmPostgreSQLServer(),
		"azurerm_postgresql_virtual_network_rule":                                        resourceArmPostgreSQLVirtualNetworkRule(),
		"azurerm_private_dns_zone":                                                       resourceArmPrivateDnsZone(),
		"azurerm_private_dns_a_record":                                                   resourceArmPrivateDnsARecord(),
		"azurerm_private_dns_cname_record":                                               resourceArmPrivateDnsCNameRecord(),
		"azurerm_public_ip":                                                              resourceArmPublicIp(),
		"azurerm_public_ip_prefix":                                                       resourceArmPublicIpPrefix(),
		"azurerm_recovery_network_mapping":                                               resourceArmRecoveryServicesNetworkMapping(),
		"azurerm_recovery_replicated_vm":                                                 resourceArmRecoveryServicesReplicatedVm(),
		"azurerm_recovery_services_fabric":                                               resourceArmRecoveryServicesFabric(),
		"azurerm_recovery_services_protected_vm":                                         resourceArmRecoveryServicesProtectedVm(),
		"azurerm_recovery_services_protection_container":                                 resourceArmRecoveryServicesProtectionContainer(),
		"azurerm_recovery_services_protection_container_mapping":                         resourceArmRecoveryServicesProtectionContainerMapping(),
		"azurerm_recovery_services_protection_policy_vm":                                 resourceArmRecoveryServicesProtectionPolicyVm(),
		"azurerm_recovery_services_replication_policy":                                   resourceArmRecoveryServicesReplicationPolicy(),
		"azurerm_recovery_services_vault":                                                resourceArmRecoveryServicesVault(),
		"azurerm_redis_cache":                                                            resourceArmRedisCache(),
		"azurerm_redis_firewall_rule":                                                    resourceArmRedisFirewallRule(),
		"azurerm_relay_namespace":                                                        resourceArmRelayNamespace(),
		"azurerm_resource_group":                                                         resourceArmResourceGroup(),
		"azurerm_role_assignment":                                                        resourceArmRoleAssignment(),
		"azurerm_role_definition":                                                        resourceArmRoleDefinition(),
		"azurerm_route_table":                                                            resourceArmRouteTable(),
		"azurerm_route":                                                                  resourceArmRoute(),
		"azurerm_scheduler_job_collection":                                               resourceArmSchedulerJobCollection(),
		"azurerm_scheduler_job":                                                          resourceArmSchedulerJob(),
		"azurerm_search_service":                                                         resourceArmSearchService(),
		"azurerm_security_center_contact":                                                resourceArmSecurityCenterContact(),
		"azurerm_security_center_subscription_pricing":                                   resourceArmSecurityCenterSubscriptionPricing(),
		"azurerm_security_center_workspace":                                              resourceArmSecurityCenterWorkspace(),
		"azurerm_service_fabric_cluster":                                                 resourceArmServiceFabricCluster(),
		"azurerm_servicebus_namespace_authorization_rule":                                resourceArmServiceBusNamespaceAuthorizationRule(),
		"azurerm_servicebus_namespace":                                                   resourceArmServiceBusNamespace(),
		"azurerm_servicebus_queue_authorization_rule":                                    resourceArmServiceBusQueueAuthorizationRule(),
		"azurerm_servicebus_queue":                                                       resourceArmServiceBusQueue(),
		"azurerm_servicebus_subscription_rule":                                           resourceArmServiceBusSubscriptionRule(),
		"azurerm_servicebus_subscription":                                                resourceArmServiceBusSubscription(),
		"azurerm_servicebus_topic_authorization_rule":                                    resourceArmServiceBusTopicAuthorizationRule(),
		"azurerm_servicebus_topic":                                                       resourceArmServiceBusTopic(),
		"azurerm_shared_image_gallery":                                                   resourceArmSharedImageGallery(),
		"azurerm_shared_image_version":                                                   resourceArmSharedImageVersion(),
		"azurerm_shared_image":                                                           resourceArmSharedImage(),
		"azurerm_signalr_service":                                                        resourceArmSignalRService(),
		"azurerm_snapshot":                                                               resourceArmSnapshot(),
		"azurerm_sql_active_directory_administrator":                                     resourceArmSqlAdministrator(),
		"azurerm_sql_database":                                                           resourceArmSqlDatabase(),
		"azurerm_sql_elasticpool":                                                        resourceArmSqlElasticPool(),
		"azurerm_sql_failover_group":                                                     resourceArmSqlFailoverGroup(),
		"azurerm_sql_firewall_rule":                                                      resourceArmSqlFirewallRule(),
		"azurerm_sql_server":                                                             resourceArmSqlServer(),
		"azurerm_sql_virtual_network_rule":                                               resourceArmSqlVirtualNetworkRule(),
		"azurerm_storage_account":                                                        resourceArmStorageAccount(),
		"azurerm_storage_blob":                                                           resourceArmStorageBlob(),
		"azurerm_storage_container":                                                      resourceArmStorageContainer(),
		"azurerm_storage_queue":                                                          resourceArmStorageQueue(),
		"azurerm_storage_share":                                                          resourceArmStorageShare(),
		"azurerm_storage_share_directory":                                                resourceArmStorageShareDirectory(),
		"azurerm_storage_table":                                                          resourceArmStorageTable(),
		"azurerm_storage_table_entity":                                                   resourceArmStorageTableEntity(),
		"azurerm_stream_analytics_job":                                                   resourceArmStreamAnalyticsJob(),
		"azurerm_stream_analytics_function_javascript_udf":                               resourceArmStreamAnalyticsFunctionUDF(),
		"azurerm_stream_analytics_output_blob":                                           resourceArmStreamAnalyticsOutputBlob(),
		"azurerm_stream_analytics_output_mssql":                                          resourceArmStreamAnalyticsOutputSql(),
		"azurerm_stream_analytics_output_eventhub":                                       resourceArmStreamAnalyticsOutputEventHub(),
		"azurerm_stream_analytics_output_servicebus_queue":                               resourceArmStreamAnalyticsOutputServiceBusQueue(),
		"azurerm_stream_analytics_stream_input_blob":                                     resourceArmStreamAnalyticsStreamInputBlob(),
		"azurerm_stream_analytics_stream_input_eventhub":                                 resourceArmStreamAnalyticsStreamInputEventHub(),
		"azurerm_stream_analytics_stream_input_iothub":                                   resourceArmStreamAnalyticsStreamInputIoTHub(),
		"azurerm_subnet_network_security_group_association":                              resourceArmSubnetNetworkSecurityGroupAssociation(),
		"azurerm_subnet_route_table_association":                                         resourceArmSubnetRouteTableAssociation(),
		"azurerm_subnet":                                                                 resourceArmSubnet(),
		"azurerm_template_deployment":                                                    resourceArmTemplateDeployment(),
		"azurerm_traffic_manager_endpoint":                                               resourceArmTrafficManagerEndpoint(),
		"azurerm_traffic_manager_profile":                                                resourceArmTrafficManagerProfile(),
		"azurerm_user_assigned_identity":                                                 resourceArmUserAssignedIdentity(),
		"azurerm_virtual_machine_data_disk_attachment":                                   resourceArmVirtualMachineDataDiskAttachment(),
		"azurerm_virtual_machine_extension":                                              resourceArmVirtualMachineExtensions(),
		"azurerm_virtual_machine_scale_set":                                              resourceArmVirtualMachineScaleSet(),
		"azurerm_virtual_machine":                                                        resourceArmVirtualMachine(),
		"azurerm_virtual_network_gateway_connection":                                     resourceArmVirtualNetworkGatewayConnection(),
		"azurerm_virtual_network_gateway":                                                resourceArmVirtualNetworkGateway(),
		"azurerm_virtual_network_peering":                                                resourceArmVirtualNetworkPeering(),
		"azurerm_virtual_network":                                                        resourceArmVirtualNetwork(),
		"azurerm_virtual_wan":                                                            resourceArmVirtualWan(),
		"azurerm_web_application_firewall_policy":                                        resourceArmWebApplicationFirewallPolicy(),
	}

	for _, service := range supportedServices {
		log.Printf("[DEBUG] Registering Data Sources for %q..", service.Name())
		for k, v := range service.SupportedDataSources() {
			if existing := dataSources[k]; existing != nil {
				panic(fmt.Sprintf("An existing Data Source exists for %q", k))
			}

			dataSources[k] = v
		}

		log.Printf("[DEBUG] Registering Resources for %q..", service.Name())
		for k, v := range service.SupportedResources() {
			if existing := resources[k]; existing != nil {
				panic(fmt.Sprintf("An existing Resource exists for %q", k))
			}

			resources[k] = v
		}
	}

	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", ""),
				Description: "The Subscription ID which should be used.",
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
				Description: "The Client ID which should be used.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
				Description: "The Tenant ID which should be used.",
			},

			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_ENVIRONMENT", "public"),
				Description: "The Cloud Environment which should be used. Possible values are public, usgovernment, german, and china. Defaults to public.",
			},

			// Client Certificate specific fields
			"client_certificate_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", ""),
				Description: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate.",
			},

			"client_certificate_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", ""),
				Description: "The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate",
			},

			// Client Secret specific fields
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
				Description: "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			// Managed Service Identity specific fields
			"use_msi": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_MSI", false),
				Description: "Allowed Managed Service Identity be used for Authentication.",
			},
			"msi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_MSI_ENDPOINT", ""),
				Description: "The path to a custom endpoint for Managed Service Identity - in most circumstances this should be detected automatically. ",
			},

			// Managed Tracking GUID for User-agent
			"partner_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.UUIDOrEmpty,
				DefaultFunc:  schema.EnvDefaultFunc("ARM_PARTNER_ID", ""),
				Description:  "A GUID/UUID that is registered with Microsoft to facilitate partner resource usage attribution.",
			},

			"disable_correlation_request_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DISABLE_CORRELATION_REQUEST_ID", false),
				Description: "This will disable the x-ms-correlation-request-id header.",
			},

			// Advanced feature flags
			"skip_credentials_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SKIP_CREDENTIALS_VALIDATION", false),
				Description: "This will cause the AzureRM Provider to skip verifying the credentials being used are valid.",
			},

			"skip_provider_registration": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SKIP_PROVIDER_REGISTRATION", false),
				Description: "Should the AzureRM Provider skip registering all of the Resource Providers that it supports, if they're not already registered?",
			},
		},

		DataSourcesMap: dataSources,
		ResourcesMap:   resources,
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		builder := &authentication.Builder{
			SubscriptionID:     d.Get("subscription_id").(string),
			ClientID:           d.Get("client_id").(string),
			ClientSecret:       d.Get("client_secret").(string),
			TenantID:           d.Get("tenant_id").(string),
			Environment:        d.Get("environment").(string),
			MsiEndpoint:        d.Get("msi_endpoint").(string),
			ClientCertPassword: d.Get("client_certificate_password").(string),
			ClientCertPath:     d.Get("client_certificate_path").(string),

			// Feature Toggles
			SupportsClientCertAuth:         true,
			SupportsClientSecretAuth:       true,
			SupportsManagedServiceIdentity: d.Get("use_msi").(bool),
			SupportsAzureCliToken:          true,

			// Doc Links
			ClientSecretDocsLink: "https://www.terraform.io/docs/providers/azurerm/auth/service_principal_client_secret.html",
		}

		config, err := builder.Build()
		if err != nil {
			return nil, fmt.Errorf("Error building AzureRM Client: %s", err)
		}

		partnerId := d.Get("partner_id").(string)
		skipProviderRegistration := d.Get("skip_provider_registration").(bool)
		disableCorrelationRequestID := d.Get("disable_correlation_request_id").(bool)

		client, err := getArmClient(config, skipProviderRegistration, partnerId, disableCorrelationRequestID)
		if err != nil {
			return nil, err
		}

		client.StopContext = p.StopContext()

		// replaces the context between tests
		p.MetaReset = func() error {
			client.StopContext = p.StopContext()
			return nil
		}

		skipCredentialsValidation := d.Get("skip_credentials_validation").(bool)
		if !skipCredentialsValidation {
			// List all the available providers and their registration state to avoid unnecessary
			// requests. This also lets us check if the provider credentials are correct.
			ctx := client.StopContext
			providerList, err := client.resource.ProvidersClient.List(ctx, nil, "")
			if err != nil {
				return nil, fmt.Errorf("Unable to list provider registration status, it is possible that this is due to invalid "+
					"credentials or the service principal does not have permission to use the Resource Manager API, Azure "+
					"error: %s", err)
			}

			if !skipProviderRegistration {
				availableResourceProviders := providerList.Values()
				requiredResourceProviders := requiredResourceProviders()

				err := ensureResourceProvidersAreRegistered(ctx, *client.resource.ProvidersClient, availableResourceProviders, requiredResourceProviders)
				if err != nil {
					return nil, fmt.Errorf("Error ensuring Resource Providers are registered: %s", err)
				}
			}
		}

		return client, nil
	}
}

// ignoreCaseStateFunc is a StateFunc from helper/schema that converts the
// supplied value to lower before saving to state for consistency.
func ignoreCaseStateFunc(val interface{}) string {
	return strings.ToLower(val.(string))
}

func userDataDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return userDataStateFunc(old) == new
}

func userDataStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		s = base64Encode(s)
		hash := sha1.Sum([]byte(s))
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}

func base64EncodedStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		return base64Encode(s)
	default:
		return ""
	}
}

// base64Encode encodes data if the input isn't already encoded using
// base64.StdEncoding.EncodeToString. If the input is already base64 encoded,
// return the original input unchanged.
func base64Encode(data string) string {
	// Check whether the data is already Base64 encoded; don't double-encode
	if isBase64Encoded(data) {
		return data
	}
	// data has not been encoded encode and return
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func isBase64Encoded(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}
