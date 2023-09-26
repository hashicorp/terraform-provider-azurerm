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
	mlworkspace "github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	loadBalancerParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	loadBalancerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	resourcesParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDatabricksWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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

			"managed_services_cmk_key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.KeyVaultChildID,
			},

			"managed_disk_cmk_key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: keyVaultValidate.KeyVaultChildID,
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

			"managed_disk_cmk_rotation_to_latest_version_enabled": {
				Type:         pluginsdk.TypeBool,
				Optional:     true,
				RequiredWith: []string{"managed_disk_cmk_key_vault_key_id"},
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
				ValidateFunc: loadBalancerValidate.LoadBalancerBackendAddressPoolID,
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
		backendPoolId, err := loadBalancerParse.LoadBalancerBackendAddressPoolID(backendPool)
		if err != nil {
			return err
		}

		// Generate the load balancer ID from the Backend Address Pool Id...
		lbId := loadBalancerParse.NewLoadBalancerID(backendPoolId.SubscriptionId, backendPoolId.ResourceGroup, backendPoolId.LoadBalancerName)

		backendPoolName = backendPoolId.BackendAddressPoolName
		loadBalancerId = lbId.ID()

		locks.ByID(backendPoolId.ID())
		defer locks.UnlockByID(backendPoolId.ID())

		locks.ByID(lbId.ID())
		defer locks.UnlockByID(lbId.ID())

		// check to make sure the load balancer exists as referred to by the Backend Address Pool...
		lb, err := lbClient.Get(ctx, lbId.ResourceGroup, lbId.Name, "")
		if err != nil {
			if utils.ResponseWasNotFound(lb.Response) {
				return fmt.Errorf("load balancer %q for Backend Address Pool %q was not found", lbId, backendPoolId)
			}
			return fmt.Errorf("failed to retrieve Load Balancer %q for Backend Address Pool %q: %+v", lbId, backendPoolId, err)
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
	servicesKeyIdRaw := d.Get("managed_services_cmk_key_vault_key_id").(string)
	if servicesKeyIdRaw != "" {
		setEncrypt = true
		key, err := keyVaultParse.ParseNestedItemID(servicesKeyIdRaw)
		if err != nil {
			return err
		}

		// make sure the key vault exists
		keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, meta.(*clients.Client).Resource, key.KeyVaultBaseUrl)
		if err != nil || keyVaultIdRaw == nil {
			return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed services Key Vault at URL %q: %+v", key.KeyVaultBaseUrl, err)
		}

		encrypt.Entities.ManagedServices = &workspaces.EncryptionV2{
			// There is only one valid source for this field at this point in time so I have hardcoded the value
			KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
			KeyVaultProperties: &workspaces.EncryptionV2KeyVaultProperties{
				KeyName:     key.Name,
				KeyVersion:  key.Version,
				KeyVaultUri: key.KeyVaultBaseUrl,
			},
		}
	}

	diskKeyIdRaw := d.Get("managed_disk_cmk_key_vault_key_id").(string)
	if diskKeyIdRaw != "" {
		setEncrypt = true
		key, err := keyVaultParse.ParseNestedItemID(diskKeyIdRaw)
		if err != nil {
			return err
		}

		// make sure the key vault exists
		keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, meta.(*clients.Client).Resource, key.KeyVaultBaseUrl)
		if err != nil || keyVaultIdRaw == nil {
			return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed disk Key Vault at URL %q: %+v", key.KeyVaultBaseUrl, err)
		}

		encrypt.Entities.ManagedDisk = &workspaces.ManagedDiskEncryption{
			// There is only one valid source for this field at this point in time so I have hardcoded the value
			KeySource: workspaces.EncryptionKeySourceMicrosoftPointKeyvault,
			KeyVaultProperties: workspaces.ManagedDiskEncryptionKeyVaultProperties{
				KeyName:     key.Name,
				KeyVersion:  key.Version,
				KeyVaultUri: key.KeyVaultBaseUrl,
			},
		}

		rotationEnabled := d.Get("managed_disk_cmk_rotation_to_latest_version_enabled").(bool)
		if rotationEnabled {
			encrypt.Entities.ManagedDisk.RotationToLatestKeyVersionEnabled = utils.Bool(rotationEnabled)
		}
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

	if encrypt != nil && servicesKeyIdRaw != "" {
		d.Set("managed_services_cmk_key_vault_key_id", servicesKeyIdRaw)
	}

	if encrypt != nil && diskKeyIdRaw != "" {
		d.Set("managed_disk_cmk_key_vault_key_id", diskKeyIdRaw)
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

		if model.Properties.WorkspaceUrl != nil {
			d.Set("workspace_url", model.Properties.WorkspaceUrl)
		}

		if model.Properties.WorkspaceId != nil {
			d.Set("workspace_id", model.Properties.WorkspaceId)
		}

		// customer managed key for managed services
		encryptKeyName := ""
		encryptKeyVersion := ""
		encryptKeyVaultURI := ""

		if encryption := model.Properties.Encryption; encryption != nil {
			if encryptionProps := encryption.Entities.ManagedServices; encryptionProps != nil {
				encryptKeyName = encryptionProps.KeyVaultProperties.KeyName
				encryptKeyVersion = encryptionProps.KeyVaultProperties.KeyVersion
				encryptKeyVaultURI = encryptionProps.KeyVaultProperties.KeyVaultUri
			}
		}

		if encryptKeyVaultURI != "" {
			key, err := keyVaultParse.NewNestedItemID(encryptKeyVaultURI, keyVaultParse.NestedItemTypeKey, encryptKeyName, encryptKeyVersion)
			if err == nil {
				d.Set("managed_services_cmk_key_vault_key_id", key.ID())
			}
		}
		// customer managed key for managed disk
		encryptDiskKeyName := ""
		encryptDiskKeyVersion := ""
		encryptDiskKeyVaultURI := ""
		encryptDiskRotationEnabled := false

		if encryption := model.Properties.Encryption; encryption != nil {
			if encryptionProps := encryption.Entities.ManagedDisk; encryptionProps != nil {
				encryptDiskKeyName = encryptionProps.KeyVaultProperties.KeyName
				encryptDiskKeyVersion = encryptionProps.KeyVaultProperties.KeyVersion
				encryptDiskKeyVaultURI = encryptionProps.KeyVaultProperties.KeyVaultUri
				encryptDiskRotationEnabled = *encryptionProps.RotationToLatestKeyVersionEnabled
			}

		}

		if encryptDiskKeyVaultURI != "" {
			key, err := keyVaultParse.NewNestedItemID(encryptDiskKeyVaultURI, keyVaultParse.NestedItemTypeKey, encryptDiskKeyName, encryptDiskKeyVersion)
			if err == nil {
				d.Set("managed_disk_cmk_key_vault_key_id", key.ID())
			}
			d.Set("managed_disk_cmk_rotation_to_latest_version_enabled", encryptDiskRotationEnabled)
			d.Set("disk_encryption_set_id", model.Properties.DiskEncryptionSetId)
		}

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

	if v := input; v != nil {
		if t := v.PrincipalId; t != nil {
			if t != nil {
				e["principal_id"] = *t
			}
		}

		if t := v.TenantId; t != nil {
			e["tenant_id"] = *t
		}

		if t := v.Type; t != nil {
			if t != nil {
				e["type"] = *t
			}
		}

		if len(e) != 0 {
			return []interface{}{e}
		}
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

	lbId, err := loadBalancerParse.LoadBalancerIDInsensitively(loadBalancerId)

	if err == nil {
		backendId := loadBalancerParse.NewLoadBalancerBackendAddressPoolID(lbId.SubscriptionId, lbId.ResourceGroup, lbId.Name, backendName)
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
