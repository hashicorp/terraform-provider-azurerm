// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/accessconnector"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	mlworkspace "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/subnets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -test-params "premium"

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

		Importer: pluginsdk.ImporterValidatingIdentity(&workspaces.WorkspaceId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&workspaces.WorkspaceId{}),
		},

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
				ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
			},

			"managed_disk_cmk_key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
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
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice(validate.PossibleValuesForComplianceStandard(), false),
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

			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
				// Neither of these arguments can be removed once set
				for _, k := range []string{"managed_disk_cmk_key_vault_key_id", "managed_services_cmk_key_vault_key_id"} {
					o, n := d.GetChange(k)

					if o.(string) != "" && n.(string) == "" {
						// Check RawConfig to prevent replacments on `(known after apply)` values
						rawConfig := d.GetRawConfig()
						if rawConfig.IsNull() || !rawConfig.IsKnown() {
							return nil
						}
						rawValues := rawConfig.AsValueMap()

						if !rawValues[k].IsNull() {
							return nil
						}

						if err := d.ForceNew(k); err != nil {
							return err
						}
					}
				}

				return nil
			}),
		),
	}

	if !features.SixPointOh() {
		resource.Schema["managed_services_cmk_key_vault_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			DiffSuppressFunc: func(_, o, n string, _ *pluginsdk.ResourceData) bool {
				// Suppress removal diff for 5.x since that does not require an update
				return o != "" && n == ""
			},
			Deprecated: "`managed_services_cmk_key_vault_id` has been deprecated and will be removed in v6.0 of the AzureRM provider. This property is no longer required for cross-subscription scenarios.",
		}

		resource.Schema["managed_disk_cmk_key_vault_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			DiffSuppressFunc: func(_, o, n string, _ *pluginsdk.ResourceData) bool {
				// Suppress removal diff for 5.x since that does not require an update
				return o != "" && n == ""
			},
			Deprecated: "`managed_disk_cmk_key_vault_id` has been deprecated and will be removed in v6.0 of the AzureRM provider. This property is no longer required for cross-subscription scenarios.",
		}
	}

	return resource
}

func resourceDatabricksWorkspaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	acClient := meta.(*clients.Client).DataBricks.AccessConnectorClient
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subnetsClient := meta.(*clients.Client).Network.Subnets

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workspaces.NewWorkspaceID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
	managedResourceGroupName := d.Get("managed_resource_group_name").(string)

	if backendPool := d.Get("load_balancer_backend_address_pool_id").(string); backendPool != "" {
		backendPoolId, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(backendPool)
		if err != nil {
			return err
		}

		lbId := loadbalancers.NewProviderLoadBalancerID(backendPoolId.SubscriptionId, backendPoolId.ResourceGroupName, backendPoolId.LoadBalancerName)
		backendPoolName = backendPoolId.BackendAddressPoolName
		loadBalancerId = lbId.ID()

		locks.ByID(backendPoolId.ID())
		defer locks.UnlockByID(backendPoolId.ID())

		locks.ByID(lbId.ID())
		defer locks.UnlockByID(lbId.ID())

		// check to make sure the load balancer exists as referred to by the Backend Address Pool...
		lb, err := lbClient.Get(ctx, lbId, loadbalancers.GetOperationOptions{})
		if err != nil {
			if response.WasNotFound(lb.HttpResponse) {
				return fmt.Errorf("%s was not found", lbId)
			}
			return fmt.Errorf("retrieving %s: %+v", lbId, err)
		}
	}

	if managedResourceGroupName == "" {
		log.Printf("[DEBUG] no managed resource group name was provided, using the default pattern")
		managedResourceGroupName = fmt.Sprintf("databricks-rg-%s", id.ResourceGroupName)
	}

	publicNetworkAccess := workspaces.PublicNetworkAccessDisabled
	if d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = workspaces.PublicNetworkAccessEnabled
	}

	customParamsRaw := d.Get("custom_parameters").([]interface{})
	customParams, pubSubAssoc, priSubAssoc := expandWorkspaceCustomParameters(customParamsRaw, d.Get("customer_managed_key_enabled").(bool), d.Get("infrastructure_encryption_enabled").(bool), backendPoolName, loadBalancerId)

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

	encryption, err := expandDatabricksWorkspaceEncryption(d)
	if err != nil {
		return fmt.Errorf("expanding workspace encryption: %+v", err)
	}

	workspace := workspaces.Workspace{
		Sku: &workspaces.Sku{
			Name: d.Get("sku").(string),
		},
		Location: location.Normalize(d.Get("location").(string)),
		Properties: workspaces.WorkspaceProperties{
			ComputeMode:                workspaces.ComputeModeHybrid,
			PublicNetworkAccess:        &publicNetworkAccess,
			ManagedResourceGroupId:     pointer.To(commonids.NewResourceGroupID(id.SubscriptionId, managedResourceGroupName).ID()),
			Parameters:                 customParams,
			Encryption:                 encryption,
			EnhancedSecurityCompliance: expandWorkspaceEnhancedSecurity(d.Get("enhanced_security_compliance").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if d.Get("default_storage_firewall_enabled").(bool) {
		accessConnectorId, err := accessconnector.ParseAccessConnectorID(d.Get("access_connector_id").(string))
		if err != nil {
			return err
		}

		accessConnector, err := acClient.Get(ctx, *accessConnectorId)
		if err != nil {
			return fmt.Errorf("retrieving Access Connector %s: %+v", accessConnectorId.AccessConnectorName, err)
		}

		accessConnectorProperties := workspaces.WorkspacePropertiesAccessConnector{}
		if model := accessConnector.Model; model != nil && model.Identity != nil {
			accIdentityId := ""
			for raw := range model.Identity.IdentityIds {
				id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
				if err != nil {
					return err
				}
				accIdentityId = id.ID()
				break
			}

			accessConnectorProperties.Id = pointer.From(model.Id)
			accessConnectorProperties.IdentityType = workspaces.IdentityType(model.Identity.Type)
			accessConnectorProperties.UserAssignedIdentityId = &accIdentityId
		}

		workspace.Properties.AccessConnector = &accessConnectorProperties
		workspace.Properties.DefaultStorageFirewall = pointer.To(workspaces.DefaultStorageFirewallEnabled)
	}

	if requireNsgRules := d.Get("network_security_group_rules_required").(string); requireNsgRules != "" {
		workspace.Properties.RequiredNsgRules = pointer.ToEnum[workspaces.RequiredNsgRules](requireNsgRules)
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, workspace, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// I have to set the custom_parameters so I can pass the public and private
	// subnet NSG association along with the backend Pool Id since they are not
	// returned in the read from Azure...
	custom, backendPoolReadId := flattenWorkspaceCustomParameters(customParams, pubSubAssoc, priSubAssoc)
	d.Set("load_balancer_backend_address_pool_id", backendPoolReadId)

	if err := d.Set("custom_parameters", custom); err != nil {
		return fmt.Errorf("setting `custom_parameters`: %+v", err)
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
		d.Set("location", location.Normalize(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku", sku.Name)
		}

		managedResourceGroupID, err := commonids.ParseResourceGroupIDInsensitively(pointer.From(model.Properties.ManagedResourceGroupId))
		if err != nil {
			return err
		}

		if !features.SixPointOh() {
			d.Set("managed_resource_group_id", model.Properties.ManagedResourceGroupId)
		} else {
			d.Set("managed_resource_group_id", managedResourceGroupID.ID())
		}

		d.Set("managed_resource_group_name", managedResourceGroupID.ResourceGroupName)

		if defaultStorageFirewall := model.Properties.DefaultStorageFirewall; defaultStorageFirewall != nil {
			d.Set("default_storage_firewall_enabled", *defaultStorageFirewall != workspaces.DefaultStorageFirewallDisabled)
			if model.Properties.AccessConnector != nil {
				d.Set("access_connector_id", model.Properties.AccessConnector.Id)
			}
		}

		if publicNetworkAccess := model.Properties.PublicNetworkAccess; publicNetworkAccess != nil {
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
			_, pubSubAssoc, priSubAssoc := expandWorkspaceCustomParameters(d.Get("custom_parameters").([]interface{}), cmkEnabled, infraEnabled, "", "")

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

		d.Set("workspace_url", pointer.From(model.Properties.WorkspaceURL))
		d.Set("workspace_id", pointer.From(model.Properties.WorkspaceId))

		// customer managed key for managed services
		var servicesKeyId string
		if encryption := model.Properties.Encryption; encryption != nil {
			if encryptionProps := encryption.Entities.ManagedServices; encryptionProps != nil {
				if encryptionProps.KeyVaultProperties.KeyVaultUri != "" {
					key, err := keyvault.NewNestedItemID(encryptionProps.KeyVaultProperties.KeyVaultUri, keyvault.NestedItemTypeKey, encryptionProps.KeyVaultProperties.KeyName, encryptionProps.KeyVaultProperties.KeyVersion)
					if err == nil {
						servicesKeyId = key.ID()
					}
				}
			}
		}

		// customer managed key for managed disk
		var diskKeyId string
		var encryptDiskRotationEnabled bool
		if encryption := model.Properties.Encryption; encryption != nil {
			if encryptionProps := encryption.Entities.ManagedDisk; encryptionProps != nil {
				if encryptionProps.KeyVaultProperties.KeyVaultUri != "" {
					key, err := keyvault.NewNestedItemID(encryptionProps.KeyVaultProperties.KeyVaultUri, keyvault.NestedItemTypeKey, encryptionProps.KeyVaultProperties.KeyName, encryptionProps.KeyVaultProperties.KeyVersion)
					if err == nil {
						diskKeyId = key.ID()
					}
				}

				encryptDiskRotationEnabled = *encryptionProps.RotationToLatestKeyVersionEnabled
			}
		}

		d.Set("enhanced_security_compliance", flattenWorkspaceEnhancedSecurity(model.Properties.EnhancedSecurityCompliance))
		d.Set("disk_encryption_set_id", pointer.From(model.Properties.DiskEncryptionSetId))

		// Always set these even if they are empty to keep the state file
		// consistent with the configuration file...
		d.Set("managed_services_cmk_key_vault_key_id", servicesKeyId)
		d.Set("managed_disk_cmk_key_vault_key_id", diskKeyId)
		d.Set("managed_disk_cmk_rotation_to_latest_version_enabled", encryptDiskRotationEnabled)

		if !features.SixPointOh() {
			d.Set("managed_services_cmk_key_vault_id", d.Get("managed_services_cmk_key_vault_id").(string))
			d.Set("managed_disk_cmk_key_vault_id", d.Get("managed_disk_cmk_key_vault_id").(string))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
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

	if d.HasChanges("default_storage_firewall_enabled", "access_connector_id") {
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
		props.RequiredNsgRules = pointer.ToEnum[workspaces.RequiredNsgRules](d.Get("network_security_group_rules_required").(string))
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

	if d.HasChanges("managed_services_cmk_key_vault_key_id", "managed_disk_cmk_key_vault_key_id", "managed_disk_cmk_rotation_to_latest_version_enabled") {
		encryption, err := expandDatabricksWorkspaceEncryption(d)
		if err != nil {
			return fmt.Errorf("expanding workspace encryption: %+v", err)
		}
		props.Encryption = encryption
	}

	if d.HasChange("enhanced_security_compliance") {
		props.EnhancedSecurityCompliance = expandWorkspaceEnhancedSecurity(d.Get("enhanced_security_compliance").([]interface{}))
	}

	model.Properties = props

	if err := client.CreateOrUpdateThenPoll(ctx, *id, model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceDatabricksWorkspaceRead(d, meta)
}

func expandDatabricksWorkspaceEncryption(d *pluginsdk.ResourceData) (*workspaces.WorkspacePropertiesEncryption, error) {
	diskCMK := d.Get("managed_disk_cmk_key_vault_key_id").(string)
	servicesCMK := d.Get("managed_services_cmk_key_vault_key_id").(string)
	if diskCMK == "" && servicesCMK == "" {
		return nil, nil
	}

	result := &workspaces.WorkspacePropertiesEncryption{
		Entities: workspaces.EncryptionEntitiesDefinition{},
	}

	if diskCMK != "" {
		key, err := keyvault.ParseNestedItemID(diskCMK, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return nil, err
		}

		rotateToLatest := (*bool)(nil)
		if d.Get("managed_disk_cmk_rotation_to_latest_version_enabled").(bool) {
			rotateToLatest = pointer.To(true)
		}

		result.Entities.ManagedDisk = &workspaces.ManagedDiskEncryption{
			KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
			KeyVaultProperties: workspaces.ManagedDiskEncryptionKeyVaultProperties{
				KeyName:     key.Name,
				KeyVaultUri: key.KeyVaultBaseURL,
				KeyVersion:  key.Version,
			},
			RotationToLatestKeyVersionEnabled: rotateToLatest,
		}
	}

	if servicesCMK != "" {
		key, err := keyvault.ParseNestedItemID(servicesCMK, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return nil, err
		}

		result.Entities.ManagedServices = &workspaces.EncryptionV2{
			KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
			KeyVaultProperties: &workspaces.EncryptionV2KeyVaultProperties{
				KeyName:     key.Name,
				KeyVaultUri: key.KeyVaultBaseURL,
				KeyVersion:  key.Version,
			},
		}
	}

	return result, nil
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
			if s == string(validate.ComplianceStandardNONE) {
				continue
			}
			standards.Add(s)
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

	complianceStandards := make([]string, 0)
	if standardSet, ok := config["compliance_security_profile_standards"].(*pluginsdk.Set); ok {
		for _, s := range standardSet.List() {
			complianceStandards = append(complianceStandards, s.(string))
		}
	}

	if complianceSecurityProfileEnabled == workspaces.ComplianceSecurityProfileValueEnabled && len(complianceStandards) == 0 {
		complianceStandards = append(complianceStandards, string(validate.ComplianceStandardNONE))
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
