// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-10-01-preview/accessconnector"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2024-05-01/workspaces"
	mlworkspace "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	resourcesParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatabricksWorkspace() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceDatabricksWorkspaceCreate,
		Read:   resourceDatabricksWorkspaceRead,
		Update: resourceDatabricksWorkspaceUpdate,
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
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				// NOTE: O+C We set a value for this if omitted so this should remain Computed
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

			"default_storage_firewall_enabled": {
				Type:         pluginsdk.TypeBool,
				Optional:     true,
				RequiredWith: []string{"access_connector_id"},
			},

			"access_connector_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				RequiredWith: []string{"default_storage_firewall_enabled"},
			},

			"network_security_group_rules_required": {
				Type:     pluginsdk.TypeString,
				Optional: true,
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
				// NOTE: O+C The API populates these and since many are ForceNew there doesn't appear to be a need to remove this once set to use the defaults
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
							Default:      true,
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

						"storage_account_sku_name": {
							Type:         pluginsdk.TypeString,
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

			"managed_services_cmk_key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.KeyVaultChildID,
			},

			"managed_services_cmk_key_vault_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateKeyVaultID,
			},

			"managed_disk_cmk_key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.KeyVaultChildID,
			},
			"managed_disk_cmk_key_vault_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateKeyVaultID,
			},

			"managed_disk_cmk_rotation_to_latest_version_enabled": {
				Type:         pluginsdk.TypeBool,
				Optional:     true,
				RequiredWith: []string{"managed_disk_cmk_key_vault_key_id"},
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

			"enhanced_security_compliance": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"automatic_cluster_update_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"compliance_security_profile_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"compliance_security_profile_standards": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(workspaces.ComplianceStandardHIPAA),
									string(workspaces.ComplianceStandardPCIDSS),
								}, false),
							},
						},
						"enhanced_security_monitoring_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
				_, customerEncryptionEnabled := d.GetChange("customer_managed_key_enabled")
				_, defaultStorageFirewallEnabled := d.GetChange("default_storage_firewall_enabled")
				_, infrastructureEncryptionEnabled := d.GetChange("infrastructure_encryption_enabled")
				_, publicNetworkAccess := d.GetChange("public_network_access_enabled")
				_, requireNsgRules := d.GetChange("network_security_group_rules_required")
				_, backendPool := d.GetChange("load_balancer_backend_address_pool_id")
				_, managedServicesCMK := d.GetChange("managed_services_cmk_key_vault_key_id")
				_, managedDiskCMK := d.GetChange("managed_disk_cmk_key_vault_key_id")
				_, enhancedSecurityCompliance := d.GetChange("enhanced_security_compliance")

				oldSku, newSku := d.GetChange("sku")

				// Disabling Public Network Access means that this is a Private Endpoint Workspace
				// Having a Load Balancer Backend Address Pool means the this is a Secure Cluster Connectivity Workspace
				// You cannot have a Private Enpoint Workspace and a Secure Cluster Connectivity Workspace definitions in
				// the same workspace configuration...
				if !publicNetworkAccess.(bool) {
					if requireNsgRules.(string) == string(workspaces.RequiredNsgRulesAllRules) {
						return fmt.Errorf("having `network_security_group_rules_required` set to %q and `public_network_access_enabled` set to `false` is an invalid configuration", string(workspaces.RequiredNsgRulesAllRules))
					}
					if backendPool.(string) != "" {
						return fmt.Errorf("having `load_balancer_backend_address_pool_id` defined and having `public_network_access_enabled` set to `false` is an invalid configuration")
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

				if (customerEncryptionEnabled.(bool) || defaultStorageFirewallEnabled.(bool) || len(enhancedSecurityCompliance.([]interface{})) > 0 || infrastructureEncryptionEnabled.(bool) || managedServicesCMK.(string) != "" || managedDiskCMK.(string) != "") && !strings.EqualFold("premium", newSku.(string)) {
					return fmt.Errorf("`customer_managed_key_enabled`, `default_storage_firewall_enabled`, `enhanced_security_compliance`, `infrastructure_encryption_enabled`, `managed_disk_cmk_key_vault_key_id` and `managed_services_cmk_key_vault_key_id` are only available with a `premium` workspace `sku`, got %q", newSku)
				}

				return nil
			}),

			// Once compliance security profile has been enabled, disabling it will force a workspace replacement
			pluginsdk.ForceNewIfChange("enhanced_security_compliance.0.compliance_security_profile_enabled", func(ctx context.Context, old, new, meta interface{}) bool {
				return old.(bool) && !new.(bool)
			}),

			// Once a compliance standard is enabled, disabling it will force a workspace replacement
			pluginsdk.ForceNewIfChange("enhanced_security_compliance.0.compliance_security_profile_standards", func(ctx context.Context, old, new, meta interface{}) bool {
				removedStandards := old.(*pluginsdk.Set).Difference(new.(*pluginsdk.Set))
				return removedStandards.Len() > 0
			}),

			// Compliance security profile requires automatic cluster update and enhanced security monitoring to be enabled
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
				_, complianceSecurityProfileEnabled := d.GetChange("enhanced_security_compliance.0.compliance_security_profile_enabled")
				_, automaticClusterUpdateEnabled := d.GetChange("enhanced_security_compliance.0.automatic_cluster_update_enabled")
				_, enhancedSecurityMonitoringEnabled := d.GetChange("enhanced_security_compliance.0.enhanced_security_monitoring_enabled")

				if complianceSecurityProfileEnabled.(bool) && (!automaticClusterUpdateEnabled.(bool) || !enhancedSecurityMonitoringEnabled.(bool)) {
					return fmt.Errorf("`automatic_cluster_update_enabled` and `enhanced_security_monitoring_enabled` must be set to true when `compliance_security_profile_enabled` is set to true")
				}

				return nil
			}),

			// compliance standards cannot be specified without enabling compliance profile
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
				_, complianceSecurityProfileEnabled := d.GetChange("enhanced_security_compliance.0.compliance_security_profile_enabled")
				_, complianceStandards := d.GetChange("enhanced_security_compliance.0.compliance_security_profile_standards")

				if !complianceSecurityProfileEnabled.(bool) && complianceStandards.(*pluginsdk.Set).Len() > 0 {
					return fmt.Errorf("`compliance_security_profile_standards` cannot be set when `compliance_security_profile_enabled` is false")
				}

				return nil
			}),
		),
	}

	return resource
}

func resourceDatabricksWorkspaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	acClient := meta.(*clients.Client).DataBricks.AccessConnectorClient
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subnetsClient := meta.(*clients.Client).Network.Subnets
	keyVaultsClient := meta.(*clients.Client).KeyVault
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workspaces.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_databricks_workspace", id.ID())
	}

	var backendPoolName, loadBalancerId string
	skuName := d.Get("sku").(string)
	managedResourceGroupName := d.Get("managed_resource_group_name").(string)
	location := location.Normalize(d.Get("location").(string))
	backendPool := d.Get("load_balancer_backend_address_pool_id").(string)

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
	defaultStorageFirewallEnabledRaw := d.Get("default_storage_firewall_enabled").(bool)
	defaultStorageFirewallEnabled := workspaces.DefaultStorageFirewallDisabled
	if defaultStorageFirewallEnabledRaw {
		defaultStorageFirewallEnabled = workspaces.DefaultStorageFirewallEnabled
	}
	publicNetworkAccessRaw := d.Get("public_network_access_enabled").(bool)
	publicNetworkAccess := workspaces.PublicNetworkAccessDisabled
	if publicNetworkAccessRaw {
		publicNetworkAccess = workspaces.PublicNetworkAccessEnabled
	}
	requireNsgRules := d.Get("network_security_group_rules_required").(string)
	customParamsRaw := d.Get("custom_parameters").([]interface{})
	customParams, pubSubAssoc, priSubAssoc := expandWorkspaceCustomParameters(customParamsRaw, customerEncryptionEnabled, infrastructureEncryptionEnabled, backendPoolName, loadBalancerId)

	if len(customParamsRaw) > 0 && customParamsRaw[0] != nil {
		config := customParamsRaw[0].(map[string]interface{})
		pubSub := config["public_subnet_name"].(string)
		priSub := config["private_subnet_name"].(string)
		vnetID := config["virtual_network_id"].(string)

		if config["virtual_network_id"].(string) == "" && (pubSub != "" || priSub != "") {
			return fmt.Errorf("`public_subnet_name` and/or `private_subnet_name` cannot be defined if `virtual_network_id` is not set")
		}
		if config["virtual_network_id"].(string) != "" && (pubSub == "" || priSub == "") {
			return fmt.Errorf("`public_subnet_name` and `private_subnet_name` must both have values if `virtual_network_id` is set")
		}
		if pubSub != "" && pubSubAssoc == nil {
			return fmt.Errorf("you must define a value for `public_subnet_network_security_group_association_id` if `public_subnet_name` is set")
		}
		if priSub != "" && priSubAssoc == nil {
			return fmt.Errorf("you must define a value for `private_subnet_network_security_group_association_id` if `private_subnet_name` is set")
		}

		if subnetDelegationErr := checkSubnetDelegations(ctx, subnetsClient, vnetID, pubSub, priSub); subnetDelegationErr != nil {
			return subnetDelegationErr
		}
	}

	// Set up customer-managed keys for managed services encryption (e.g. notebook)
	setEncrypt := false
	encrypt := &workspaces.WorkspacePropertiesEncryption{}
	encrypt.Entities = workspaces.EncryptionEntitiesDefinition{}

	var servicesKeyId string
	var servicesKeyVaultId string
	var diskKeyId string
	var diskKeyVaultId string

	if v, ok := d.GetOk("managed_services_cmk_key_vault_key_id"); ok {
		servicesKeyId = v.(string)
	}

	if v, ok := d.GetOk("managed_services_cmk_key_vault_id"); ok {
		servicesKeyVaultId = v.(string)
	}

	if v, ok := d.GetOk("managed_disk_cmk_key_vault_key_id"); ok {
		diskKeyId = v.(string)
	}

	if v, ok := d.GetOk("managed_disk_cmk_key_vault_id"); ok {
		diskKeyVaultId = v.(string)
	}

	// set default subscription as current subscription for key vault look-up...
	servicesResourceSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)
	diskResourceSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	if servicesKeyVaultId != "" {
		// If they passed the 'managed_cmk_key_vault_id' parse the Key Vault ID
		// to extract the correct key vault subscription for the exists call...
		v, err := commonids.ParseKeyVaultID(servicesKeyVaultId)
		if err != nil {
			return fmt.Errorf("parsing %q as a Key Vault ID: %+v", servicesKeyVaultId, err)
		}

		servicesResourceSubscriptionId = commonids.NewSubscriptionID(v.SubscriptionId)
	}

	if servicesKeyId != "" {
		setEncrypt = true
		key, err := keyVaultParse.ParseNestedItemID(servicesKeyId)
		if err != nil {
			return err
		}

		// make sure the key vault exists
		_, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, servicesResourceSubscriptionId, key.KeyVaultBaseUrl)
		if err != nil {
			return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed services Key Vault in subscription %q at URL %q: %+v", servicesResourceSubscriptionId, key.KeyVaultBaseUrl, err)
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

	if diskKeyVaultId != "" {
		// If they passed the 'managed_disk_cmk_key_vault_id' parse the Key Vault ID
		// to extract the correct key vault subscription for the exists call...
		v, err := commonids.ParseKeyVaultID(diskKeyVaultId)
		if err != nil {
			return fmt.Errorf("parsing %q as a Key Vault ID: %+v", diskKeyVaultId, err)
		}

		diskResourceSubscriptionId = commonids.NewSubscriptionID(v.SubscriptionId)
	}

	if diskKeyId != "" {
		setEncrypt = true
		key, err := keyVaultParse.ParseNestedItemID(diskKeyId)
		if err != nil {
			return err
		}

		// make sure the key vault exists
		_, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, diskResourceSubscriptionId, key.KeyVaultBaseUrl)
		if err != nil {
			return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed disk Key Vault in subscription %q at URL %q: %+v", diskResourceSubscriptionId, key.KeyVaultBaseUrl, err)
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

	if rotationEnabled := d.Get("managed_disk_cmk_rotation_to_latest_version_enabled").(bool); rotationEnabled {
		encrypt.Entities.ManagedDisk.RotationToLatestKeyVersionEnabled = pointer.To(rotationEnabled)
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

	if defaultStorageFirewallEnabledRaw {
		accessConnectorProperties := workspaces.WorkspacePropertiesAccessConnector{}
		accessConnectorIdRaw := d.Get("access_connector_id").(string)
		accessConnectorId, err := accessconnector.ParseAccessConnectorID(accessConnectorIdRaw)
		if err != nil {
			return fmt.Errorf("parsing Access Connector ID %s: %+v", accessConnectorIdRaw, err)
		}

		accessConnector, err := acClient.Get(ctx, *accessConnectorId)
		if err != nil {
			return fmt.Errorf("retrieving Access Connector %s: %+v", accessConnectorId.AccessConnectorName, err)
		}

		if accessConnector.Model.Identity != nil {
			accIdentityId := ""
			for raw := range accessConnector.Model.Identity.IdentityIds {
				id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
				if err != nil {
					return fmt.Errorf("parsing %q as a User Assigned Identity ID: %+v", raw, err)
				}
				accIdentityId = id.ID()
				break
			}

			accessConnectorProperties.Id = *accessConnector.Model.Id
			accessConnectorProperties.IdentityType = workspaces.IdentityType(accessConnector.Model.Identity.Type)
			accessConnectorProperties.UserAssignedIdentityId = &accIdentityId
		}

		workspace.Properties.AccessConnector = &accessConnectorProperties
		workspace.Properties.DefaultStorageFirewall = &defaultStorageFirewallEnabled
	}

	if !d.IsNewResource() && d.HasChange("default_storage_firewall_enabled") {
		workspace.Properties.DefaultStorageFirewall = &defaultStorageFirewallEnabled
	}

	if requireNsgRules != "" {
		requiredNsgRulesConst := workspaces.RequiredNsgRules(requireNsgRules)
		workspace.Properties.RequiredNsgRules = &requiredNsgRulesConst
	}

	if setEncrypt {
		workspace.Properties.Encryption = encrypt
	}

	enhancedSecurityCompliance := d.Get("enhanced_security_compliance")
	workspace.Properties.EnhancedSecurityCompliance = expandWorkspaceEnhancedSecurity(enhancedSecurityCompliance.([]interface{}))

	if err := client.CreateOrUpdateThenPoll(ctx, id, workspace); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
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
	d.Set("managed_services_cmk_key_vault_key_id", servicesKeyId)
	d.Set("managed_disk_cmk_key_vault_key_id", diskKeyId)
	d.Set("managed_services_cmk_key_vault_id", servicesKeyVaultId)
	d.Set("managed_disk_cmk_key_vault_id", diskKeyVaultId)

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

	var encryptDiskRotationEnabled bool
	var servicesKeyVaultId string
	if v, ok := d.GetOk("managed_services_cmk_key_vault_id"); ok {
		servicesKeyVaultId = v.(string)
	}

	var diskKeyVaultId string
	if v, ok := d.GetOk("managed_disk_cmk_key_vault_id"); ok {
		diskKeyVaultId = v.(string)
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

		if defaultStorageFirewall := model.Properties.DefaultStorageFirewall; defaultStorageFirewall != nil {
			d.Set("default_storage_firewall_enabled", *defaultStorageFirewall != workspaces.DefaultStorageFirewallDisabled)
			if model.Properties.AccessConnector != nil {
				d.Set("access_connector_id", model.Properties.AccessConnector.Id)
			}
		}

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
		if model.Properties.WorkspaceURL != nil {
			workspaceUrl = *model.Properties.WorkspaceURL
		}
		d.Set("workspace_url", workspaceUrl)

		var workspaceId string
		if model.Properties.WorkspaceId != nil {
			workspaceId = *model.Properties.WorkspaceId
		}
		d.Set("workspace_id", workspaceId)

		// customer managed key for managed services
		var servicesKeyId string
		if encryption := model.Properties.Encryption; encryption != nil {
			if encryptionProps := encryption.Entities.ManagedServices; encryptionProps != nil {
				if encryptionProps.KeyVaultProperties.KeyVaultUri != "" {
					key, err := keyVaultParse.NewNestedItemID(encryptionProps.KeyVaultProperties.KeyVaultUri, keyVaultParse.NestedItemTypeKey, encryptionProps.KeyVaultProperties.KeyName, encryptionProps.KeyVaultProperties.KeyVersion)
					if err == nil {
						servicesKeyId = key.ID()
					}
				}
			}
		}

		// customer managed key for managed disk
		var diskKeyId string
		if encryption := model.Properties.Encryption; encryption != nil {
			if encryptionProps := encryption.Entities.ManagedDisk; encryptionProps != nil {
				if encryptionProps.KeyVaultProperties.KeyVaultUri != "" {
					key, err := keyVaultParse.NewNestedItemID(encryptionProps.KeyVaultProperties.KeyVaultUri, keyVaultParse.NestedItemTypeKey, encryptionProps.KeyVaultProperties.KeyName, encryptionProps.KeyVaultProperties.KeyVersion)
					if err == nil {
						diskKeyId = key.ID()
					}
				}

				encryptDiskRotationEnabled = *encryptionProps.RotationToLatestKeyVersionEnabled
			}
		}

		d.Set("enhanced_security_compliance", flattenWorkspaceEnhancedSecurity(model.Properties.EnhancedSecurityCompliance))

		var encryptDiskEncryptionSetId string
		if model.Properties.DiskEncryptionSetId != nil {
			encryptDiskEncryptionSetId = *model.Properties.DiskEncryptionSetId
		}
		d.Set("disk_encryption_set_id", encryptDiskEncryptionSetId)

		// Always set these even if they are empty to keep the state file
		// consistent with the configuration file...
		d.Set("managed_services_cmk_key_vault_key_id", servicesKeyId)
		d.Set("managed_services_cmk_key_vault_id", servicesKeyVaultId)
		d.Set("managed_disk_cmk_key_vault_key_id", diskKeyId)
		d.Set("managed_disk_cmk_key_vault_id", diskKeyVaultId)
		d.Set("managed_disk_cmk_rotation_to_latest_version_enabled", encryptDiskRotationEnabled)

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceDatabricksWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	deleteOperationOptions := workspaces.DefaultDeleteOperationOptions()
	if meta.(*clients.Client).Features.DatabricksWorkspace.ForceDelete {
		deleteOperationOptions.ForceDeletion = pointer.To(true)
	}

	if err = client.DeleteThenPoll(ctx, *id, deleteOperationOptions); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func resourceDatabricksWorkspaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	acClient := meta.(*clients.Client).DataBricks.AccessConnectorClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model is nil", id)
	}

	model := *existing.Model

	props := model.Properties

	if d.HasChange("sku") {
		if model.Sku == nil {
			model.Sku = &workspaces.Sku{}
		}
		model.Sku.Name = d.Get("sku").(string)
	}

	if d.HasChange("tags") {
		model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("customer_managed_key_enabled") {
		if props.Parameters == nil {
			props.Parameters = &workspaces.WorkspaceCustomParameters{}
		}
		props.Parameters.PrepareEncryption = &workspaces.WorkspaceCustomBooleanParameter{
			Value: d.Get("customer_managed_key_enabled").(bool),
		}
	}

	if d.HasChange("infrastructure_encryption_enabled") {
		if props.Parameters == nil {
			props.Parameters = &workspaces.WorkspaceCustomParameters{}
		}
		props.Parameters.RequireInfrastructureEncryption = &workspaces.WorkspaceCustomBooleanParameter{
			Value: d.Get("infrastructure_encryption_enabled").(bool),
		}
	}

	if d.HasChange("default_storage_firewall_enabled") {
		defaultStorageFirewallEnabled := workspaces.DefaultStorageFirewallDisabled
		defaultStorageFirewallEnabledRaw := d.Get("default_storage_firewall_enabled").(bool)

		if defaultStorageFirewallEnabledRaw {
			defaultStorageFirewallEnabled = workspaces.DefaultStorageFirewallEnabled

			accessConnectorProperties := workspaces.WorkspacePropertiesAccessConnector{}
			accessConnectorIdRaw := d.Get("access_connector_id").(string)
			accessConnectorId, err := accessconnector.ParseAccessConnectorID(accessConnectorIdRaw)
			if err != nil {
				return err
			}

			accessConnector, err := acClient.Get(ctx, *accessConnectorId)
			if err != nil {
				return fmt.Errorf("retrieving Access Connector %s: %+v", accessConnectorId.AccessConnectorName, err)
			}

			if accessConnector.Model.Identity != nil {
				accIdentityId := ""
				for raw := range accessConnector.Model.Identity.IdentityIds {
					identityId, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
					if err != nil {
						return err
					}
					accIdentityId = identityId.ID()
					break
				}

				accessConnectorProperties.Id = *accessConnector.Model.Id
				accessConnectorProperties.IdentityType = workspaces.IdentityType(accessConnector.Model.Identity.Type)
				accessConnectorProperties.UserAssignedIdentityId = &accIdentityId
			}

			props.AccessConnector = &accessConnectorProperties
		}

		props.DefaultStorageFirewall = &defaultStorageFirewallEnabled
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccessRaw := d.Get("public_network_access_enabled").(bool)
		publicNetworkAccess := workspaces.PublicNetworkAccessDisabled
		if publicNetworkAccessRaw {
			publicNetworkAccess = workspaces.PublicNetworkAccessEnabled
		}
		props.PublicNetworkAccess = &publicNetworkAccess
	}

	if d.HasChange("network_security_group_rules_required") {
		props.RequiredNsgRules = pointer.To(workspaces.RequiredNsgRules(d.Get("network_security_group_rules_required").(string)))
	}

	if d.HasChange("custom_parameters") {
		if props.Parameters == nil {
			props.Parameters = &workspaces.WorkspaceCustomParameters{}
		}

		if customParams := d.Get("custom_parameters").([]interface{}); len(customParams) > 0 && customParams[0] != nil {
			config := customParams[0].(map[string]interface{})
			var pubSubnetAssoc, priSubnetAssoc *string

			pubSub := config["public_subnet_name"].(string)
			priSub := config["private_subnet_name"].(string)

			if v, ok := config["public_subnet_network_security_group_association_id"].(string); ok {
				pubSubnetAssoc = &v
			}

			if v, ok := config["private_subnet_network_security_group_association_id"].(string); ok {
				priSubnetAssoc = &v
			}

			if config["virtual_network_id"].(string) == "" && (pubSub != "" || priSub != "") {
				return fmt.Errorf("`public_subnet_name` and/or `private_subnet_name` cannot be defined if `virtual_network_id` is not set")
			}
			if config["virtual_network_id"].(string) != "" && (pubSub == "" || priSub == "") {
				return fmt.Errorf("`public_subnet_name` and `private_subnet_name` must both have values if `virtual_network_id` is set")
			}
			if pubSub != "" && pubSubnetAssoc == nil {
				return fmt.Errorf("you must define a value for `public_subnet_network_security_group_association_id` if `public_subnet_name` is set")
			}
			if priSub != "" && priSubnetAssoc == nil {
				return fmt.Errorf("you must define a value for `private_subnet_network_security_group_association_id` if `private_subnet_name` is set")
			}

			if v, ok := config["nat_gateway_name"].(string); ok && v != "" {
				props.Parameters.NatGatewayName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["public_ip_name"].(string); ok && v != "" {
				props.Parameters.PublicIPName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["storage_account_name"].(string); ok && v != "" {
				props.Parameters.StorageAccountName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["storage_account_sku_name"].(string); ok && v != "" {
				props.Parameters.StorageAccountSkuName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["vnet_address_prefix"].(string); ok && v != "" {
				props.Parameters.VnetAddressPrefix = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["machine_learning_workspace_id"].(string); ok && v != "" {
				props.Parameters.AmlWorkspaceId = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["no_public_ip"].(bool); ok {
				props.Parameters.EnableNoPublicIP = &workspaces.WorkspaceNoPublicIPBooleanParameter{
					Value: v,
				}
			}

			if v, ok := config["public_subnet_name"].(string); ok && v != "" {
				props.Parameters.CustomPublicSubnetName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["private_subnet_name"].(string); ok && v != "" {
				props.Parameters.CustomPrivateSubnetName = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}

			if v, ok := config["virtual_network_id"].(string); ok && v != "" {
				props.Parameters.CustomVirtualNetworkId = &workspaces.WorkspaceCustomStringParameter{
					Value: v,
				}
			}
		}
	}

	// Set up customer-managed keys for managed services encryption (e.g. notebook)
	setEncrypt := false
	encrypt := &workspaces.WorkspacePropertiesEncryption{}
	encrypt.Entities = workspaces.EncryptionEntitiesDefinition{}

	var servicesKeyId string
	var servicesKeyVaultId string
	var diskKeyId string
	var diskKeyVaultId string

	if v, ok := d.GetOk("managed_services_cmk_key_vault_key_id"); ok {
		servicesKeyId = v.(string)
	}

	if v, ok := d.GetOk("managed_services_cmk_key_vault_id"); ok {
		servicesKeyVaultId = v.(string)
	}

	if v, ok := d.GetOk("managed_disk_cmk_key_vault_key_id"); ok {
		diskKeyId = v.(string)
	}

	if v, ok := d.GetOk("managed_disk_cmk_key_vault_id"); ok {
		diskKeyVaultId = v.(string)
	}

	// set default subscription as current subscription for key vault look-up...
	servicesResourceSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)
	diskResourceSubscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	if servicesKeyVaultId != "" {
		// If they passed the 'managed_cmk_key_vault_id' parse the Key Vault ID
		// to extract the correct key vault subscription for the exists call...
		v, err := commonids.ParseKeyVaultID(servicesKeyVaultId)
		if err != nil {
			return err
		}

		servicesResourceSubscriptionId = commonids.NewSubscriptionID(v.SubscriptionId)
	}

	if servicesKeyId != "" {
		setEncrypt = true
		key, err := keyVaultParse.ParseNestedItemID(servicesKeyId)
		if err != nil {
			return err
		}

		// make sure the key vault exists
		_, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, servicesResourceSubscriptionId, key.KeyVaultBaseUrl)
		if err != nil {
			return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed services Key Vault in subscription %q at URL %q: %+v", servicesResourceSubscriptionId, key.KeyVaultBaseUrl, err)
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

	if diskKeyVaultId != "" {
		// If they passed the 'managed_disk_cmk_key_vault_id' parse the Key Vault ID
		// to extract the correct key vault subscription for the exists call...
		v, err := commonids.ParseKeyVaultID(diskKeyVaultId)
		if err != nil {
			return err
		}

		diskResourceSubscriptionId = commonids.NewSubscriptionID(v.SubscriptionId)
	}

	if diskKeyId != "" {
		setEncrypt = true
		key, err := keyVaultParse.ParseNestedItemID(diskKeyId)
		if err != nil {
			return err
		}

		// make sure the key vault exists
		_, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, diskResourceSubscriptionId, key.KeyVaultBaseUrl)
		if err != nil {
			return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed disk Key Vault in subscription %q at URL %q: %+v", diskResourceSubscriptionId, key.KeyVaultBaseUrl, err)
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

	if rotationEnabled := d.Get("managed_disk_cmk_rotation_to_latest_version_enabled").(bool); rotationEnabled {
		encrypt.Entities.ManagedDisk.RotationToLatestKeyVersionEnabled = pointer.To(rotationEnabled)
	}

	if setEncrypt {
		props.Encryption = encrypt
	}

	enhancedSecurityCompliance := d.Get("enhanced_security_compliance")
	props.EnhancedSecurityCompliance = expandWorkspaceEnhancedSecurity(enhancedSecurityCompliance.([]interface{}))

	model.Properties = props

	if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if d.HasChange("tags") {
		workspaceUpdate := workspaces.WorkspaceUpdate{
			Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		}

		err := client.UpdateThenPoll(ctx, *id, workspaceUpdate)
		if err != nil {
			return fmt.Errorf("updating %s Tags: %+v", id, err)
		}
	}

	return resourceDatabricksWorkspaceRead(d, meta)
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
		parameters.EnableNoPublicIP = &workspaces.WorkspaceNoPublicIPBooleanParameter{
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

func flattenWorkspaceEnhancedSecurity(input *workspaces.EnhancedSecurityComplianceDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	enhancedSecurityCompliance := make(map[string]interface{})

	if v := input.AutomaticClusterUpdate; v != nil {
		enhancedSecurityCompliance["automatic_cluster_update_enabled"] = pointer.From(v.Value) != workspaces.AutomaticClusterUpdateValueDisabled
	}

	if v := input.EnhancedSecurityMonitoring; v != nil {
		enhancedSecurityCompliance["enhanced_security_monitoring_enabled"] = pointer.From(v.Value) != workspaces.EnhancedSecurityMonitoringValueDisabled
	}

	if v := input.ComplianceSecurityProfile; v != nil {
		enhancedSecurityCompliance["compliance_security_profile_enabled"] = pointer.From(v.Value) != workspaces.ComplianceSecurityProfileValueDisabled

		standards := pluginsdk.NewSet(pluginsdk.HashString, nil)
		for _, s := range pointer.From(v.ComplianceStandards) {
			if s == workspaces.ComplianceStandardNONE {
				continue
			}
			standards.Add(string(s))
		}

		enhancedSecurityCompliance["compliance_security_profile_standards"] = standards
	}

	return []interface{}{enhancedSecurityCompliance}
}

func expandWorkspaceEnhancedSecurity(input []interface{}) *workspaces.EnhancedSecurityComplianceDefinition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})

	automaticClusterUpdateEnabled := workspaces.AutomaticClusterUpdateValueDisabled
	if enabled, ok := config["automatic_cluster_update_enabled"].(bool); ok && enabled {
		automaticClusterUpdateEnabled = workspaces.AutomaticClusterUpdateValueEnabled
	}

	enhancedSecurityMonitoringEnabled := workspaces.EnhancedSecurityMonitoringValueDisabled
	if enabled, ok := config["enhanced_security_monitoring_enabled"].(bool); ok && enabled {
		enhancedSecurityMonitoringEnabled = workspaces.EnhancedSecurityMonitoringValueEnabled
	}

	complianceSecurityProfileEnabled := workspaces.ComplianceSecurityProfileValueDisabled
	if enabled, ok := config["compliance_security_profile_enabled"].(bool); ok && enabled {
		complianceSecurityProfileEnabled = workspaces.ComplianceSecurityProfileValueEnabled
	}

	complianceStandards := []workspaces.ComplianceStandard{}
	if standardSet, ok := config["compliance_security_profile_standards"].(*pluginsdk.Set); ok {
		for _, s := range standardSet.List() {
			complianceStandards = append(complianceStandards, workspaces.ComplianceStandard(s.(string)))
		}
	}

	if complianceSecurityProfileEnabled == workspaces.ComplianceSecurityProfileValueEnabled && len(complianceStandards) == 0 {
		complianceStandards = append(complianceStandards, workspaces.ComplianceStandardNONE)
	}

	return &workspaces.EnhancedSecurityComplianceDefinition{
		AutomaticClusterUpdate: &workspaces.AutomaticClusterUpdateDefinition{
			Value: &automaticClusterUpdateEnabled,
		},
		EnhancedSecurityMonitoring: &workspaces.EnhancedSecurityMonitoringDefinition{
			Value: &enhancedSecurityMonitoringEnabled,
		},
		ComplianceSecurityProfile: &workspaces.ComplianceSecurityProfileDefinition{
			Value:               &complianceSecurityProfileEnabled,
			ComplianceStandards: &complianceStandards,
		},
	}
}

func checkSubnetDelegations(ctx context.Context, client *subnets.SubnetsClient, vnetID, publicSubnetName, privateSubnetName string) error {
	requiredDelegationService := "Microsoft.Databricks/workspaces"

	if vnetID == "" || (publicSubnetName == "" && privateSubnetName == "") {
		return nil
	}

	id, err := commonids.ParseVirtualNetworkID(vnetID)
	if err != nil {
		return err
	}

	if publicSubnetName != "" {
		subnetID := commonids.NewSubnetID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName, publicSubnetName)
		resp, err := client.Get(ctx, subnetID, subnets.DefaultGetOperationOptions())
		if err != nil || resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("failed to check public subnet delegation for %s: %s", publicSubnetName, err)
		}
		if resp.Model.Properties.Delegations == nil {
			return fmt.Errorf("required public subnet delegation to %s on %s not found", requiredDelegationService, publicSubnetName)
		}

		if delegations := resp.Model.Properties.Delegations; delegations == nil {
			return fmt.Errorf("required public subnet delegation to %s on %s not found", requiredDelegationService, publicSubnetName)
		} else {
			found := false
			for _, v := range *delegations {
				if v.Properties == nil {
					continue
				}
				if pointer.From(v.Properties.ServiceName) == requiredDelegationService {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("required public subnet delegation to %s on %s not found", requiredDelegationService, publicSubnetName)
			}
		}
	}

	if privateSubnetName != "" {
		subnetID := commonids.NewSubnetID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName, privateSubnetName)
		resp, err := client.Get(ctx, subnetID, subnets.DefaultGetOperationOptions())
		if err != nil || resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("failed to check private subnet delegation for %s: %s", privateSubnetName, err)
		}

		if delegations := resp.Model.Properties.Delegations; delegations == nil {
			return fmt.Errorf("required private subnet delegation to %s on %s not found", requiredDelegationService, privateSubnetName)
		} else {
			found := false
			for _, v := range *delegations {
				if v.Properties == nil {
					continue
				}
				if pointer.From(v.Properties.ServiceName) == requiredDelegationService {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("required private subnet delegation to %s on %s not found", requiredDelegationService, privateSubnetName)
			}
		}
	}

	return nil
}
