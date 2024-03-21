// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2023-02-01/workspaces"
	mlworkspace "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-10-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	resourcesParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDatabricksWorkspace() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceDatabricksWorkspaceCreateUpdate,
		Read:   resourceDatabricksWorkspaceRead,
		Update: resourceDatabricksWorkspaceCreateUpdate,
		Delete: resourceDatabricksWorkspaceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := workspaces.ParseWorkspaceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"standard",
					"premium",
					"trial",
				}, false),
			},

			"managed_resource_group_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"customer_managed_key_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"managed_disk_identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"tenant_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"infrastructure_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"network_security_group_rules_required": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(workspaces.RequiredNsgRulesAllRules),
					string(workspaces.RequiredNsgRulesNoAzureDatabricksRules),
					string(workspaces.RequiredNsgRulesNoAzureServiceRules),
				}, false),
			},

			"load_balancer_backend_address_pool_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
			},

			"custom_parameters": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"machine_learning_workspace_id": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: mlworkspace.ValidateWorkspaceID,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"nat_gateway_name": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"no_public_ip": {
							Type:         pluginsdk.TypeBool,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"public_ip_name": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"public_subnet_name": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"public_subnet_network_security_group_association_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"private_subnet_name": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"private_subnet_network_security_group_association_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"virtual_network_id": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: commonids.ValidateVirtualNetworkID,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"storage_account_name": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							Computed:     true,
							ValidateFunc: storageValidate.StorageAccountName,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						// Per Service Team: This field is actually changeable so the ForceNew is no longer required, however we agreed to not change the current behavior for consistency purposes
						"storage_account_sku_name": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"vnet_address_prefix": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: workspaceCustomParametersString(),
						},
					},
				},
			},

			"managed_resource_group_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"workspace_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"workspace_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"disk_encryption_set_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_account_identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"tenant_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			_, customerEncryptionEnabled := d.GetChange("customer_managed_key_enabled")
			_, infrastructureEncryptionEnabled := d.GetChange("infrastructure_encryption_enabled")
			_, publicNetworkAccess := d.GetChange("public_network_access_enabled")
			_, requireNsgRules := d.GetChange("network_security_group_rules_required")
			_, backendPool := d.GetChange("load_balancer_backend_address_pool_id")
			_, managedServicesCMK := d.GetChange("managed_services_cmk_key_vault_key_id")
			_, managedDiskCMK := d.GetChange("managed_disk_cmk_key_vault_key_id")

			oldSku, newSku := d.GetChange("sku")

			// Disabling Public Network Access means that this is a Private Endpoint Workspace
			// Having a Load Balancer Backend Address Pool means the this is a Secure Cluster Connectivity Workspace
			// You cannot have a Private Enpoint Workspace and a Secure Cluster Connectivity Workspace definitions in
			// the same workspace configuration...
			if !publicNetworkAccess.(bool) {
				if requireNsgRules.(string) == string(workspaces.RequiredNsgRulesAllRules) {
					return fmt.Errorf("having 'network_security_group_rules_required' set to %q and 'public_network_access_enabled' set to 'false' is an invalid configuration", string(workspaces.RequiredNsgRulesAllRules))
				}
				if backendPool.(string) != "" {
					return fmt.Errorf("having 'load_balancer_backend_address_pool_id' defined and having 'public_network_access_enabled' set to 'false' is an invalid configuration")
				}
			}

			if d.HasChange("sku") {
				if newSku == "trial" {
					log.Printf("[DEBUG] recreate databricks workspace, cannot be migrated to %s", newSku)
					d.ForceNew("sku")
				} else {
					log.Printf("[DEBUG] databricks workspace can be upgraded from %s to %s", oldSku, newSku)
				}
			}

			if (customerEncryptionEnabled.(bool) || infrastructureEncryptionEnabled.(bool) || managedServicesCMK.(string) != "" || managedDiskCMK.(string) != "") && !strings.EqualFold("premium", newSku.(string)) {
				return fmt.Errorf("'customer_managed_key_enabled', 'infrastructure_encryption_enabled', 'managed_disk_cmk_key_vault_key_id' and 'managed_services_cmk_key_vault_key_id' are only available with a 'premium' workspace 'sku', got %q", newSku)
			}

			return nil
		}),
	}

	if !features.FourPointOhBeta() {
		// NOTE: Added to support cross subscription cmk's in 3.x and to migrate the resource from
		// using the keys data plane URL as the fields value to using the keys resource id
		// instead in 4.0...
		resource.Schema["managed_services_cmk_key_vault_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			Deprecated:   "`managed_services_cmk_key_vault_id` will be removed in favour of the property `managed_services_cmk_key_vault_key_resource_id` in version 4.0 of the AzureRM Provider.",
		}

		resource.Schema["managed_services_cmk_key_vault_key_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyVaultValidate.KeyVaultChildID,
			Deprecated:   "`managed_services_cmk_key_vault_key_id` will be removed in favour of the property `managed_services_cmk_key_vault_key_resource_id` in version 4.0 of the AzureRM Provider.",
		}

		resource.Schema["managed_disk_cmk_key_vault_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			Deprecated:   "`managed_disk_cmk_key_vault_id` will be removed in favour of the property `managed_disk_cmk_key_vault_key_resource_id` in version 4.0 of the AzureRM Provider.",
		}

		resource.Schema["managed_disk_cmk_key_vault_key_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyVaultValidate.KeyVaultChildID,
			Deprecated:   "`managed_disk_cmk_key_vault_key_id` will be removed in favour of the property `managed_disk_cmk_key_vault_key_resource_id` in version 4.0 of the AzureRM Provider.",
		}

		// Old Reference...
		resource.Schema["managed_disk_cmk_rotation_to_latest_version_enabled"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			RequiredWith: []string{"managed_disk_cmk_key_vault_key_id"},
		}
	} else {
		// NOTE: These fields maybe versioned or versionless...
		resource.Schema["managed_services_cmk_key_vault_key_resource_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.Any(commonids.ValidateKeyVaultKeyID, commonids.ValidateKeyVaultKeyVersionID),
		}

		resource.Schema["managed_disk_cmk_key_vault_key_resource_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.Any(commonids.ValidateKeyVaultKeyID, commonids.ValidateKeyVaultKeyVersionID),
		}

		// TODO: Make sure I updated this reference in the code below, see // Old Reference above...
		resource.Schema["managed_disk_cmk_rotation_to_latest_version_enabled"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			RequiredWith: []string{"managed_disk_cmk_key_vault_key_resource_id"},
		}
	}

	return resource
}

func resourceDatabricksWorkspaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workspaces.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_databricks_workspace", id.ID())
		}
	}

	var backendPoolName, loadBalancerId string
	skuName := d.Get("sku").(string)
	managedResourceGroupName := d.Get("managed_resource_group_name").(string)
	location := location.Normalize(d.Get("location").(string))
	backendPool := d.Get("load_balancer_backend_address_pool_id").(string)
	expandedTags := tags.Expand(d.Get("tags").(map[string]interface{}))

	if backendPool != "" {
		backendPoolId, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(backendPool)
		if err != nil {
			return err
		}

		// Generate the load balancer ID from the Backend Address Pool Id...
		lbId := loadbalancers.NewLoadBalancerID(backendPoolId.SubscriptionId, backendPoolId.ResourceGroupName, backendPoolId.LoadBalancerName)

		backendPoolName = backendPoolId.BackendAddressPoolName
		loadBalancerId = lbId.ID()

		locks.ByID(backendPoolId.ID())
		defer locks.UnlockByID(backendPoolId.ID())

		locks.ByID(lbId.ID())
		defer locks.UnlockByID(lbId.ID())

		// check to make sure the load balancer exists as referred to by the Backend Address Pool...
		plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: backendPoolId.SubscriptionId, ResourceGroupName: backendPoolId.ResourceGroupName, LoadBalancerName: backendPoolId.LoadBalancerName}
		lb, err := lbClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
		if err != nil {
			if response.WasNotFound(lb.HttpResponse) {
				return fmt.Errorf("%s was not found", lbId)
			}
			return fmt.Errorf("retrieving %s: %+v", lbId, err)
		}
	}

	if managedResourceGroupName == "" {
		// no managed resource group name was provided, we use the default pattern
		log.Printf("[DEBUG][azurerm_databricks_workspace] no managed resource group id was provided, we use the default pattern.")
		managedResourceGroupName = fmt.Sprintf("databricks-rg-%s", id.ResourceGroupName)
	}

	managedResourceGroupID := resourcesParse.NewResourceGroupID(subscriptionId, managedResourceGroupName).ID()
	customerEncryptionEnabled := d.Get("customer_managed_key_enabled").(bool)
	infrastructureEncryptionEnabled := d.Get("infrastructure_encryption_enabled").(bool)
	publicNetowrkAccessRaw := d.Get("public_network_access_enabled").(bool)
	publicNetworkAccess := workspaces.PublicNetworkAccessDisabled
	if publicNetowrkAccessRaw {
		publicNetworkAccess = workspaces.PublicNetworkAccessEnabled
	}
	requireNsgRules := d.Get("network_security_group_rules_required").(string)
	customParamsRaw := d.Get("custom_parameters").([]interface{})
	customParams, pubSubAssoc, priSubAssoc := expandWorkspaceCustomParameters(customParamsRaw, customerEncryptionEnabled, infrastructureEncryptionEnabled, backendPoolName, loadBalancerId)

	if len(customParamsRaw) > 0 && customParamsRaw[0] != nil {
		config := customParamsRaw[0].(map[string]interface{})
		pubSub := config["public_subnet_name"].(string)
		priSub := config["private_subnet_name"].(string)

		if config["virtual_network_id"].(string) == "" && (pubSub != "" || priSub != "") {
			return fmt.Errorf("'public_subnet_name' and/or 'private_subnet_name' cannot be defined if 'virtual_network_id' is not set")
		}
		if config["virtual_network_id"].(string) != "" && (pubSub == "" || priSub == "") {
			return fmt.Errorf("'public_subnet_name' and 'private_subnet_name' must both have values if 'virtual_network_id' is set")
		}
		if pubSub != "" && pubSubAssoc == nil {
			return fmt.Errorf("you must define a value for 'public_subnet_network_security_group_association_id' if 'public_subnet_name' is set")
		}
		if priSub != "" && priSubAssoc == nil {
			return fmt.Errorf("you must define a value for 'private_subnet_network_security_group_association_id' if 'private_subnet_name' is set")
		}
	}

	// Set up customer-managed keys for managed services encryption (e.g. notebook)
	setEncrypt := false
	encrypt := &workspaces.WorkspacePropertiesEncryption{}
	encrypt.Entities = workspaces.EncryptionEntitiesDefinition{}

	// TODO: Remove in 4.0
	var managedServicesKeyIdRaw string
	var managedServicesKeyVaultId string
	var managedDiskKeyIdRaw string
	var managedDiskKeyVaultId string

	// NOTE: Keep in 4.0
	var serviceKeyIdRaw string
	var diskKeyIdRaw string
	var servicesKeyId string
	var diskKeyId string

	if !features.FourPointOhBeta() {
		// set default subscription as current subscription for key vault look-up...
		managedServicesResourceSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)
		managedDiskResourceSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

		if v, ok := d.GetOk("managed_services_cmk_key_vault_key_id"); ok {
			managedServicesKeyIdRaw = v.(string)
		}

		if v, ok := d.GetOk("managed_services_cmk_key_vault_id"); ok {
			managedServicesKeyVaultId = v.(string)
		}

		if v, ok := d.GetOk("managed_disk_cmk_key_vault_key_id"); ok {
			managedDiskKeyIdRaw = v.(string)
		}

		if v, ok := d.GetOk("managed_disk_cmk_key_vault_id"); ok {
			managedDiskKeyVaultId = v.(string)
		}

		if managedServicesKeyVaultId != "" {
			// If they passed the 'managed_cmk_key_vault_id' parse the Key Vault ID
			// to extract the correct key vault subscription for the exists call...
			v, err := commonids.ParseKeyVaultID(managedServicesKeyVaultId)
			if err != nil {
				return fmt.Errorf("parsing %q as a Key Vault ID: %+v", managedServicesKeyVaultId, err)
			}

			managedServicesResourceSubscriptionId = commonids.NewSubscriptionID(v.SubscriptionId)
		}

		if managedServicesKeyIdRaw != "" {
			setEncrypt = true
			key, err := keyVaultParse.ParseNestedItemID(managedServicesKeyIdRaw)
			if err != nil {
				return err
			}

			// make sure the key vault exists
			_, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, managedServicesResourceSubscriptionId, key.KeyVaultBaseUrl)
			if err != nil {
				return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed services Key Vault in subscription %q at URL %q: %+v", managedServicesResourceSubscriptionId, key.KeyVaultBaseUrl, err)
			}

			encrypt.Entities.ManagedServices = &workspaces.EncryptionV2{
				KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
				KeyVaultProperties: &workspaces.EncryptionV2KeyVaultProperties{
					KeyName:     key.Name,
					KeyVersion:  key.Version,
					KeyVaultUri: key.KeyVaultBaseUrl,
				},
			}
		}

		if managedDiskKeyVaultId != "" {
			// If they passed the 'managed_disk_cmk_key_vault_id' parse the Key Vault ID
			// to extract the correct key vault subscription for the exists call...
			v, err := commonids.ParseKeyVaultID(managedDiskKeyVaultId)
			if err != nil {
				return fmt.Errorf("parsing %q as a Key Vault ID: %+v", managedDiskKeyVaultId, err)
			}

			managedDiskResourceSubscriptionId = commonids.NewSubscriptionID(v.SubscriptionId)
		}

		if managedDiskKeyIdRaw != "" {
			setEncrypt = true
			key, err := keyVaultParse.ParseNestedItemID(managedDiskKeyIdRaw)
			if err != nil {
				return err
			}

			// make sure the key vault exists
			_, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, managedDiskResourceSubscriptionId, key.KeyVaultBaseUrl)
			if err != nil {
				return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed disk Key Vault in subscription %q at URL %q: %+v", managedDiskResourceSubscriptionId, key.KeyVaultBaseUrl, err)
			}

			encrypt.Entities.ManagedDisk = &workspaces.ManagedDiskEncryption{
				KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
				KeyVaultProperties: workspaces.ManagedDiskEncryptionKeyVaultProperties{
					KeyName:     key.Name,
					KeyVersion:  key.Version,
					KeyVaultUri: key.KeyVaultBaseUrl,
				},
			}
		}
	} else {
		// Migrate to new 4.0 resource ID fields...
		if v, ok := d.GetOk("managed_services_cmk_key_vault_key_resource_id"); ok {
			serviceKeyIdRaw = v.(string)
		}

		if v, ok := d.GetOk("managed_disk_cmk_key_vault_key_resource_id"); ok {
			diskKeyIdRaw = v.(string)
		}

		if serviceKeyIdRaw != "" {
			setEncrypt = true

			// NOTE: The key ID may or may not be versionless...
			key, err := parseWorkspaceManagedCmkKeyWithOptionalVersion(serviceKeyIdRaw)
			if err != nil {
				return err
			}

			// Make sure the key vault exists
			keyVaultId := commonids.NewKeyVaultID(key.SubscriptionId, key.ResourceGroupName, key.VaultName)
			resp, err := keyVaultsClient.VaultsClient.Get(ctx, keyVaultId)
			if err != nil {
				return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed services Key Vault %q in subscription %q: %+v", keyVaultId, key.SubscriptionId, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", keyVaultId)
			}

			encrypt.Entities.ManagedServices = &workspaces.EncryptionV2{
				KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
				KeyVaultProperties: &workspaces.EncryptionV2KeyVaultProperties{
					KeyName:     key.KeyName,
					KeyVaultUri: *model.Properties.VaultUri,
				},
			}

			if key.VersionName == "" {
				servicesKeyId = commonids.NewKeyVaultKeyID(key.SubscriptionId, key.ResourceGroupName, key.VaultName, key.KeyName).ID()
			} else {
				encrypt.Entities.ManagedServices.KeyVaultProperties.KeyVersion = key.VersionName
				servicesKeyId = key.ID()
			}
		}

		if diskKeyIdRaw != "" {
			setEncrypt = true

			// NOTE: The key ID may or may not be versionless...
			key, err := parseWorkspaceManagedCmkKeyWithOptionalVersion(diskKeyIdRaw)
			if err != nil {
				return err
			}

			// Make sure the key vault exists
			keyVaultId := commonids.NewKeyVaultID(key.SubscriptionId, key.ResourceGroupName, key.VaultName)
			resp, err := keyVaultsClient.VaultsClient.Get(ctx, keyVaultId)
			if err != nil {
				return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed disk Key Vault %q in subscription %q: %+v", keyVaultId, key.SubscriptionId, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", keyVaultId)
			}

			encrypt.Entities.ManagedDisk = &workspaces.ManagedDiskEncryption{
				KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
				KeyVaultProperties: workspaces.ManagedDiskEncryptionKeyVaultProperties{
					KeyName:     key.KeyName,
					KeyVaultUri: *model.Properties.VaultUri,
				},
			}

			if key.VersionName == "" {
				diskKeyId = commonids.NewKeyVaultKeyID(key.SubscriptionId, key.ResourceGroupName, key.VaultName, key.KeyName).ID()
			} else {
				encrypt.Entities.ManagedDisk.KeyVaultProperties.KeyVersion = key.VersionName
				diskKeyId = key.ID()
			}
		}
	}

	if rotationEnabled := d.Get("managed_disk_cmk_rotation_to_latest_version_enabled").(bool); rotationEnabled {
		encrypt.Entities.ManagedDisk.RotationToLatestKeyVersionEnabled = utils.Bool(rotationEnabled)
	}

	// Including the Tags in the workspace parameters will update the tags on
	// the workspace only
	workspace := workspaces.Workspace{
		Sku: &workspaces.Sku{
			Name: skuName,
		},
		Location: location,
		Properties: workspaces.WorkspaceProperties{
			PublicNetworkAccess:    &publicNetworkAccess,
			ManagedResourceGroupId: managedResourceGroupID,
			Parameters:             customParams,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if requireNsgRules != "" {
		requiredNsgRulesConst := workspaces.RequiredNsgRules(requireNsgRules)
		workspace.Properties.RequiredNsgRules = &requiredNsgRulesConst
	}

	if setEncrypt {
		workspace.Properties.Encryption = encrypt
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, workspace); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	// Only call Update(e.g. PATCH) if it is not a new resource and the Tags have changed
	// this will cause the updated tags to be propagated to all of the connected
	// workspace resources.
	// TODO: can be removed once https://github.com/Azure/azure-sdk-for-go/issues/14571 is fixed
	if !d.IsNewResource() && d.HasChange("tags") {
		workspaceUpdate := workspaces.WorkspaceUpdate{
			Tags: expandedTags,
		}

		err := client.UpdateThenPoll(ctx, id, workspaceUpdate)
		if err != nil {
			return fmt.Errorf("updating %s Tags: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	// I have to set the custom_parameters so I can pass the public and private
	// subnet NSG association along with the backend Pool Id since they are not
	// returned in the read from Azure...
	custom, backendPoolReadId := flattenWorkspaceCustomParameters(customParams, pubSubAssoc, priSubAssoc)
	d.Set("load_balancer_backend_address_pool_id", backendPoolReadId)

	if err := d.Set("custom_parameters", custom); err != nil {
		return fmt.Errorf("setting `custom_parameters`: %+v", err)
	}

	// Always set these even if they are empty to keep the state file
	// consistent with the configuration file...
	if !features.FourPointOhBeta() {
		d.Set("managed_services_cmk_key_vault_key_id", managedServicesKeyIdRaw)
		d.Set("managed_disk_cmk_key_vault_key_id", managedDiskKeyIdRaw)
		d.Set("managed_services_cmk_key_vault_id", managedServicesKeyVaultId)
		d.Set("managed_disk_cmk_key_vault_id", managedDiskKeyVaultId)
	} else {
		d.Set("managed_services_cmk_key_vault_key_resource_id", servicesKeyId)
		d.Set("managed_disk_cmk_key_vault_key_resource_id", diskKeyId)
	}

	return resourceDatabricksWorkspaceRead(d, meta)
}

func resourceDatabricksWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku", sku.Name)
		}

		managedResourceGroupID, err := resourcesParse.ResourceGroupIDInsensitively(model.Properties.ManagedResourceGroupId)
		if err != nil {
			return err
		}
		d.Set("managed_resource_group_id", model.Properties.ManagedResourceGroupId)
		d.Set("managed_resource_group_name", managedResourceGroupID.ResourceGroup)

		publicNetworkAccess := model.Properties.PublicNetworkAccess
		if publicNetworkAccess != nil {
			d.Set("public_network_access_enabled", *publicNetworkAccess != workspaces.PublicNetworkAccessDisabled)
			if *publicNetworkAccess == workspaces.PublicNetworkAccessDisabled {
				if model.Properties.RequiredNsgRules != nil {
					d.Set("network_security_group_rules_required", string(*model.Properties.RequiredNsgRules))
				}
			}
		}

		var cmkEnabled, infraEnabled bool
		if parameters := model.Properties.Parameters; parameters != nil {
			if parameters.PrepareEncryption != nil {
				cmkEnabled = parameters.PrepareEncryption.Value
			}
			d.Set("customer_managed_key_enabled", cmkEnabled)

			if parameters.RequireInfrastructureEncryption != nil {
				infraEnabled = parameters.RequireInfrastructureEncryption.Value
			}
			d.Set("infrastructure_encryption_enabled", infraEnabled)

			// The subnet associations only exist in the statefile, so we need to do a Get before we Set
			// with what has come back from the Azure response...
			customParamsRaw := d.Get("custom_parameters").([]interface{})
			_, pubSubAssoc, priSubAssoc := expandWorkspaceCustomParameters(customParamsRaw, cmkEnabled, infraEnabled, "", "")

			custom, backendPoolReadId := flattenWorkspaceCustomParameters(parameters, pubSubAssoc, priSubAssoc)
			if err := d.Set("custom_parameters", custom); err != nil {
				return fmt.Errorf("setting `custom_parameters`: %+v", err)
			}

			d.Set("load_balancer_backend_address_pool_id", backendPoolReadId)
		}

		if err := d.Set("storage_account_identity", flattenWorkspaceManagedIdentity(model.Properties.StorageAccountIdentity)); err != nil {
			return fmt.Errorf("setting `storage_account_identity`: %+v", err)
		}

		if err := d.Set("managed_disk_identity", flattenWorkspaceManagedIdentity(model.Properties.ManagedDiskIdentity)); err != nil {
			return fmt.Errorf("setting `managed_disk_identity`: %+v", err)
		}

		var workspaceUrl string
		if model.Properties.WorkspaceUrl != nil {
			workspaceUrl = *model.Properties.WorkspaceUrl
		}
		d.Set("workspace_url", workspaceUrl)

		var workspaceId string
		if model.Properties.WorkspaceId != nil {
			workspaceId = *model.Properties.WorkspaceId
		}
		d.Set("workspace_id", workspaceId)

		var encryptDiskRotationEnabled bool

		if !features.FourPointOhBeta() {
			var managedServicesKeyName string
			var managedServicesKeyVersion string
			var managedServicesKeyVaultURI string
			var managedServicesKeyId string
			var managedDiskKeyName string
			var managedDiskKeyVersion string
			var managedDiskKeyVaultURI string
			var managedDiskKeyId string

			// customer managed key for managed services
			if encryption := model.Properties.Encryption; encryption != nil {
				if encryptionProps := encryption.Entities.ManagedServices; encryptionProps != nil {
					managedServicesKeyName = encryptionProps.KeyVaultProperties.KeyName
					managedServicesKeyVersion = encryptionProps.KeyVaultProperties.KeyVersion
					managedServicesKeyVaultURI = encryptionProps.KeyVaultProperties.KeyVaultUri

					if managedServicesKeyVaultURI != "" {
						key, err := keyVaultParse.NewNestedItemID(managedServicesKeyVaultURI, keyVaultParse.NestedItemTypeKey, managedServicesKeyName, managedServicesKeyVersion)
						if err == nil {
							managedServicesKeyId = key.ID()
						}
					}
				}
			}
			d.Set("managed_services_cmk_key_vault_key_id", managedServicesKeyId)

			// customer managed key for managed disk
			if encryption := model.Properties.Encryption; encryption != nil {
				if encryptionProps := encryption.Entities.ManagedDisk; encryptionProps != nil {
					managedDiskKeyName = encryptionProps.KeyVaultProperties.KeyName
					managedDiskKeyVersion = encryptionProps.KeyVaultProperties.KeyVersion
					managedDiskKeyVaultURI = encryptionProps.KeyVaultProperties.KeyVaultUri
					encryptDiskRotationEnabled = *encryptionProps.RotationToLatestKeyVersionEnabled
				}

				if managedDiskKeyVaultURI != "" {
					key, err := keyVaultParse.NewNestedItemID(managedDiskKeyVaultURI, keyVaultParse.NestedItemTypeKey, managedDiskKeyName, managedDiskKeyVersion)
					if err == nil {
						managedDiskKeyId = key.ID()
					}
				}
			}
			d.Set("managed_disk_cmk_key_vault_key_id", managedDiskKeyId)

			var managedServicesKeyVaultId string
			if v, ok := d.GetOk("managed_services_cmk_key_vault_id"); ok {
				managedServicesKeyVaultId = v.(string)
			}
			d.Set("managed_services_cmk_key_vault_id", managedServicesKeyVaultId)

			var managedDiskKeyVaultId string
			if v, ok := d.GetOk("managed_disk_cmk_key_vault_id"); ok {
				managedDiskKeyVaultId = v.(string)
			}
			d.Set("managed_disk_cmk_key_vault_id", managedDiskKeyVaultId)
		} else {
			// Set new 4.0 values...
			var servicesKeyIdRaw string
			var diskKeyIdRaw string
			var servicesKeyId string
			var diskKeyId string

			// NOTE: I have to pull this from state else I won't know what subscription
			// to use in the call to the NewKeyVaultKeyID/NewKeyVaultKeyVersionID function...
			if v := d.Get("managed_services_cmk_key_vault_key_resource_id"); v != nil {
				servicesKeyIdRaw = v.(string)
			}

			if v := d.Get("managed_disk_cmk_key_vault_key_resource_id"); v != nil {
				diskKeyIdRaw = v.(string)
			}

			// customer managed key for managed services
			if servicesKeyIdRaw != "" {
				key, err := parseWorkspaceManagedCmkKeyWithOptionalVersion(servicesKeyIdRaw)
				if err == nil {
					if key.VersionName == "" {
						servicesKeyId = commonids.NewKeyVaultKeyID(key.SubscriptionId, key.ResourceGroupName, key.VaultName, key.KeyName).ID()
					} else {
						servicesKeyId = key.ID()
					}
				}
			}
			d.Set("managed_services_cmk_key_vault_key_resource_id", servicesKeyId)

			// customer managed key for managed disk
			if diskKeyIdRaw != "" {
				key, err := parseWorkspaceManagedCmkKeyWithOptionalVersion(diskKeyIdRaw)
				if err == nil {
					if key.VersionName == "" {
						diskKeyId = commonids.NewKeyVaultKeyID(key.SubscriptionId, key.ResourceGroupName, key.VaultName, key.KeyName).ID()
					} else {
						diskKeyId = key.ID()
					}
				}

				if encryption := model.Properties.Encryption; encryption != nil {
					if encryptionProps := encryption.Entities.ManagedDisk; encryptionProps != nil {
						encryptDiskRotationEnabled = *encryptionProps.RotationToLatestKeyVersionEnabled
					}
				}
			}
			d.Set("managed_disk_cmk_key_vault_key_resource_id", diskKeyId)
		}

		d.Set("managed_disk_cmk_rotation_to_latest_version_enabled", encryptDiskRotationEnabled)

		var encryptDiskEncryptionSetId string
		if model.Properties.DiskEncryptionSetId != nil {
			encryptDiskEncryptionSetId = *model.Properties.DiskEncryptionSetId
		}
		d.Set("disk_encryption_set_id", encryptDiskEncryptionSetId)

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func parseWorkspaceManagedCmkKeyWithOptionalVersion(key string) (*commonids.KeyVaultKeyVersionId, error) {
	var result commonids.KeyVaultKeyVersionId
	if key == "" {
		return nil, fmt.Errorf("parsing %q as a Key Vault Key Resource ID: string empty", key)
	}

	keyId, err := commonids.ParseKeyVaultKeyID(key)
	if err == nil {
		result.SubscriptionId = keyId.SubscriptionId
		result.ResourceGroupName = keyId.ResourceGroupName
		result.VaultName = keyId.VaultName
		result.KeyName = keyId.KeyName
	} else {
		// Try parsing resource ID as a versioned key vault key...
		keyVersionId, err := commonids.ParseKeyVaultKeyVersionID(key)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a Key Vault Key Resource ID: %+v", key, err)
		}

		result.SubscriptionId = keyVersionId.SubscriptionId
		result.ResourceGroupName = keyVersionId.ResourceGroupName
		result.VaultName = keyVersionId.VaultName
		result.KeyName = keyVersionId.KeyName
		result.VersionName = keyVersionId.VersionName
	}

	return &result, nil
}

func resourceDatabricksWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func flattenWorkspaceManagedIdentity(input *workspaces.ManagedIdentityConfiguration) []interface{} {
	if input == nil {
		return nil
	}

	e := make(map[string]interface{})

	if t := input.PrincipalId; t != nil {
		e["principal_id"] = *t
	}

	if t := input.TenantId; t != nil {
		e["tenant_id"] = *t
	}

	if t := input.Type; t != nil {
		e["type"] = *t
	}

	if len(e) != 0 {
		return []interface{}{e}
	}

	return []interface{}{e}
}

func flattenWorkspaceCustomParameters(input *workspaces.WorkspaceCustomParameters, publicSubnetAssociation, privateSubnetAssociation *string) ([]interface{}, string) {
	if input == nil {
		return nil, ""
	}

	var backendAddressPoolId, backendName, loadBalancerId string
	parameters := make(map[string]interface{})

	if publicSubnetAssociation != nil && *publicSubnetAssociation != "" {
		parameters["public_subnet_network_security_group_association_id"] = *publicSubnetAssociation
	}

	if privateSubnetAssociation != nil && *privateSubnetAssociation != "" {
		parameters["private_subnet_network_security_group_association_id"] = *privateSubnetAssociation
	}

	if v := input.LoadBalancerBackendPoolName; v != nil {
		backendName = v.Value
	}

	if v := input.LoadBalancerId; v != nil {
		loadBalancerId = v.Value
	}

	if v := input.NatGatewayName; v != nil {
		parameters["nat_gateway_name"] = v.Value
	}

	if v := input.PublicIPName; v != nil {
		parameters["public_ip_name"] = v.Value
	}

	if v := input.StorageAccountName; v != nil {
		parameters["storage_account_name"] = v.Value
	}

	if v := input.StorageAccountSkuName; v != nil {
		parameters["storage_account_sku_name"] = v.Value
	}

	if v := input.VnetAddressPrefix; v != nil {
		parameters["vnet_address_prefix"] = v.Value
	}

	if v := input.AmlWorkspaceId; v != nil {
		parameters["machine_learning_workspace_id"] = v.Value
	}

	if v := input.EnableNoPublicIP; v != nil {
		parameters["no_public_ip"] = v.Value
	}

	if v := input.CustomPrivateSubnetName; v != nil {
		parameters["private_subnet_name"] = v.Value
	}

	if v := input.CustomPublicSubnetName; v != nil {
		parameters["public_subnet_name"] = v.Value
	}

	if v := input.CustomVirtualNetworkId; v != nil {
		parameters["virtual_network_id"] = v.Value
	}

	lbId, err := loadbalancers.ParseLoadBalancerIDInsensitively(loadBalancerId)

	if err == nil {
		backendId := loadbalancers.NewLoadBalancerBackendAddressPoolID(lbId.SubscriptionId, lbId.ResourceGroupName, lbId.LoadBalancerName, backendName)
		backendAddressPoolId = backendId.ID()
	}

	return []interface{}{parameters}, backendAddressPoolId
}

func expandWorkspaceCustomParameters(input []interface{}, customerManagedKeyEnabled, infrastructureEncryptionEnabled bool, backendAddressPoolName, loadBalancerId string) (workspaceCustomParameters *workspaces.WorkspaceCustomParameters, publicSubnetAssociation, privateSubnetAssociation *string) {
	if len(input) == 0 || input[0] == nil {
		// This will be hit when there are no custom params set but we still
		// need to pass the customerManagedKeyEnabled and infrastructureEncryptionEnabled
		// flags anyway...
		parameters := workspaces.WorkspaceCustomParameters{}

		parameters.PrepareEncryption = &workspaces.WorkspaceCustomBooleanParameter{
			Value: customerManagedKeyEnabled,
		}

		parameters.RequireInfrastructureEncryption = &workspaces.WorkspaceCustomBooleanParameter{
			Value: infrastructureEncryptionEnabled,
		}

		return &parameters, nil, nil
	}

	config := input[0].(map[string]interface{})
	var pubSubnetAssoc, priSubnetAssoc *string
	parameters := workspaces.WorkspaceCustomParameters{}

	if v, ok := config["public_subnet_network_security_group_association_id"].(string); ok && v != "" {
		pubSubnetAssoc = &v
	}

	if v, ok := config["private_subnet_network_security_group_association_id"].(string); ok && v != "" {
		priSubnetAssoc = &v
	}

	if backendAddressPoolName != "" {
		parameters.LoadBalancerBackendPoolName = &workspaces.WorkspaceCustomStringParameter{
			Value: backendAddressPoolName,
		}
	}

	if loadBalancerId != "" {
		parameters.LoadBalancerId = &workspaces.WorkspaceCustomStringParameter{
			Value: loadBalancerId,
		}
	}

	if v, ok := config["nat_gateway_name"].(string); ok && v != "" {
		parameters.NatGatewayName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["public_ip_name"].(string); ok && v != "" {
		parameters.PublicIPName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["storage_account_name"].(string); ok && v != "" {
		parameters.StorageAccountName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["storage_account_sku_name"].(string); ok && v != "" {
		parameters.StorageAccountSkuName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["vnet_address_prefix"].(string); ok && v != "" {
		parameters.VnetAddressPrefix = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["machine_learning_workspace_id"].(string); ok && v != "" {
		parameters.AmlWorkspaceId = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["no_public_ip"].(bool); ok {
		parameters.EnableNoPublicIP = &workspaces.WorkspaceCustomBooleanParameter{
			Value: v,
		}
	}

	if v, ok := config["public_subnet_name"].(string); ok && v != "" {
		parameters.CustomPublicSubnetName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	parameters.PrepareEncryption = &workspaces.WorkspaceCustomBooleanParameter{
		Value: customerManagedKeyEnabled,
	}

	parameters.RequireInfrastructureEncryption = &workspaces.WorkspaceCustomBooleanParameter{
		Value: infrastructureEncryptionEnabled,
	}

	if v, ok := config["private_subnet_name"].(string); ok && v != "" {
		parameters.CustomPrivateSubnetName = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	if v, ok := config["virtual_network_id"].(string); ok && v != "" {
		parameters.CustomVirtualNetworkId = &workspaces.WorkspaceCustomStringParameter{
			Value: v,
		}
	}

	return &parameters, pubSubnetAssoc, priSubnetAssoc
}

func workspaceCustomParametersString() []string {
	return []string{
		"custom_parameters.0.machine_learning_workspace_id", "custom_parameters.0.no_public_ip",
		"custom_parameters.0.public_subnet_name", "custom_parameters.0.private_subnet_name", "custom_parameters.0.virtual_network_id",
		"custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id",
		"custom_parameters.0.nat_gateway_name", "custom_parameters.0.public_ip_name", "custom_parameters.0.storage_account_name", "custom_parameters.0.storage_account_sku_name",
		"custom_parameters.0.vnet_address_prefix",
	}
}
