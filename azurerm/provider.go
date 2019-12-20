package azurerm

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/provider"
)

func Provider() terraform.ResourceProvider {
	dataSources := map[string]*schema.Resource{
		"azurerm_application_security_group":         dataSourceArmApplicationSecurityGroup(),
		"azurerm_builtin_role_definition":            dataSourceArmBuiltInRoleDefinition(),
		"azurerm_client_config":                      dataSourceArmClientConfig(),
		"azurerm_image":                              dataSourceArmImage(),
		"azurerm_key_vault_access_policy":            dataSourceArmKeyVaultAccessPolicy(),
		"azurerm_key_vault_key":                      dataSourceArmKeyVaultKey(),
		"azurerm_key_vault_secret":                   dataSourceArmKeyVaultSecret(),
		"azurerm_key_vault":                          dataSourceArmKeyVault(),
		"azurerm_platform_image":                     dataSourceArmPlatformImage(),
		"azurerm_proximity_placement_group":          dataSourceArmProximityPlacementGroup(),
		"azurerm_resources":                          dataSourceArmResources(),
		"azurerm_resource_group":                     dataSourceArmResourceGroup(),
		"azurerm_storage_account_blob_container_sas": dataSourceArmStorageAccountBlobContainerSharedAccessSignature(),
		"azurerm_storage_account_sas":                dataSourceArmStorageAccountSharedAccessSignature(),
		"azurerm_storage_account":                    dataSourceArmStorageAccount(),
		"azurerm_storage_management_policy":          dataSourceArmStorageManagementPolicy(),
	}

	resources := map[string]*schema.Resource{
		"azurerm_application_gateway":               resourceArmApplicationGateway(),
		"azurerm_application_security_group":        resourceArmApplicationSecurityGroup(),
		"azurerm_autoscale_setting":                 resourceArmAutoScaleSetting(),
		"azurerm_bastion_host":                      resourceArmBastionHost(),
		"azurerm_connection_monitor":                resourceArmConnectionMonitor(),
		"azurerm_dashboard":                         resourceArmDashboard(),
		"azurerm_image":                             resourceArmImage(),
		"azurerm_key_vault_access_policy":           resourceArmKeyVaultAccessPolicy(),
		"azurerm_key_vault_certificate":             resourceArmKeyVaultCertificate(),
		"azurerm_key_vault_key":                     resourceArmKeyVaultKey(),
		"azurerm_key_vault_secret":                  resourceArmKeyVaultSecret(),
		"azurerm_key_vault":                         resourceArmKeyVault(),
		"azurerm_management_lock":                   resourceArmManagementLock(),
		"azurerm_marketplace_agreement":             resourceArmMarketplaceAgreement(),
		"azurerm_resource_group":                    resourceArmResourceGroup(),
		"azurerm_shared_image_gallery":              resourceArmSharedImageGallery(),
		"azurerm_shared_image_version":              resourceArmSharedImageVersion(),
		"azurerm_shared_image":                      resourceArmSharedImage(),
		"azurerm_storage_account":                   resourceArmStorageAccount(),
		"azurerm_storage_account_network_rules":     resourceArmStorageAccountNetworkRules(),
		"azurerm_storage_blob":                      resourceArmStorageBlob(),
		"azurerm_storage_container":                 resourceArmStorageContainer(),
		"azurerm_storage_data_lake_gen2_filesystem": resourceArmStorageDataLakeGen2FileSystem(),
		"azurerm_storage_management_policy":         resourceArmStorageManagementPolicy(),
		"azurerm_storage_queue":                     resourceArmStorageQueue(),
		"azurerm_storage_share":                     resourceArmStorageShare(),
		"azurerm_storage_share_directory":           resourceArmStorageShareDirectory(),
		"azurerm_storage_table":                     resourceArmStorageTable(),
		"azurerm_storage_table_entity":              resourceArmStorageTableEntity(),
		"azurerm_template_deployment":               resourceArmTemplateDeployment(),
	}

	return provider.AzureProvider(dataSources, resources)
}
