package databricks

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/databricks/mgmt/2021-04-01-preview/databricks"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	loadBalancerParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	resourcesParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
			_, err := parse.WorkspaceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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
				ForceNew: true,
				Optional: true,
				Default:  false,
			},

			"managed_services_cmk_key_vault_key_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.KeyVaultChildID,
			},

			"infrastructure_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  true,
			},

			"network_security_group_rules_required": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(databricks.RequiredNsgRulesAllRules),
					string(databricks.RequiredNsgRulesNoAzureDatabricksRules),
					string(databricks.RequiredNsgRulesNoAzureServiceRules),
				}, false),
			},

			"load_balancer_backend_address_pool_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
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
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
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
							ForceNew:     true,
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
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"storage_account_name": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							Computed:     true,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

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

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			_, customerEncryptionEnabled := d.GetChange("customer_managed_key_enabled")
			_, infrastructureEncryptionEnabled := d.GetChange("infrastructure_encryption_enabled")
			_, publicNetworkAccess := d.GetChange("public_network_access_enabled")
			_, requireNsgRules := d.GetChange("network_security_group_rules_required")
			_, backendPool := d.GetChange("load_balancer_backend_address_pool_id")
			_, managedServicesCMK := d.GetChange("managed_services_cmk_key_vault_key_id")

			oldSku, newSku := d.GetChange("sku")

			// Disabling Public Network Access means that this is a Private Endpoint Workspace
			// Having a Load Balancer Backend Address Pool means the this is a Secure Cluster Connectivity Workspace
			// You cannot have a Private Enpoint Workspace and a Secure Cluster Connectivity Workspace definitions in
			// the same workspace configuration...
			if !publicNetworkAccess.(bool) {
				if requireNsgRules.(string) == string(databricks.RequiredNsgRulesAllRules) {
					return fmt.Errorf("having 'network_security_group_rules_required' set to %q and 'public_network_access_enabled' set to 'false' is an invalid configuration", string(databricks.RequiredNsgRulesAllRules))
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

			if (customerEncryptionEnabled.(bool) || infrastructureEncryptionEnabled.(bool) || managedServicesCMK.(string) != "") && !strings.EqualFold("premium", newSku.(string)) {
				return fmt.Errorf("'customer_managed_key_enabled', 'infrastructure_encryption_enabled' and 'managed_services_cmk_key_vault_key_id' are only available with a 'premium' workspace 'sku', got %q", newSku)
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

	id := parse.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_databricks_workspace", id.ID())
		}
	}

	var backendPoolName, loadBalancerId string
	skuName := d.Get("sku").(string)
	managedResourceGroupName := d.Get("managed_resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	backendPool := d.Get("load_balancer_backend_address_pool_id").(string)
	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

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
				return fmt.Errorf("Load Balancer %q for Backend Address Pool %q was not found", lbId, backendPoolId)
			}
			return fmt.Errorf("failed to retrieve Load Balancer %q for Backend Address Pool %q: %+v", lbId, backendPoolId, err)
		}
	}

	if managedResourceGroupName == "" {
		// no managed resource group name was provided, we use the default pattern
		log.Printf("[DEBUG][azurerm_databricks_workspace] no managed resource group id was provided, we use the default pattern.")
		managedResourceGroupName = fmt.Sprintf("databricks-rg-%s", id.ResourceGroup)
	}

	managedResourceGroupID := resourcesParse.NewResourceGroupID(subscriptionId, managedResourceGroupName).ID()
	customerEncryptionEnabled := d.Get("customer_managed_key_enabled").(bool)
	infrastructureEncryptionEnabled := d.Get("infrastructure_encryption_enabled").(bool)
	publicNetowrkAccessRaw := d.Get("public_network_access_enabled").(bool)
	publicNetworkAccess := databricks.PublicNetworkAccessDisabled
	if publicNetowrkAccessRaw {
		publicNetworkAccess = databricks.PublicNetworkAccessEnabled
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
	encrypt := &databricks.WorkspacePropertiesEncryption{}
	keyIdRaw := d.Get("managed_services_cmk_key_vault_key_id").(string)
	if keyIdRaw != "" {
		key, err := keyVaultParse.ParseNestedItemID(keyIdRaw)
		if err != nil {
			return err
		}

		encrypt.Entities = &databricks.EncryptionEntitiesDefinition{
			ManagedServices: &databricks.EncryptionV2{
				// There is only one valid source for this field at this point in time so I have hardcoded the value
				KeySource: utils.String(string(databricks.KeySourceMicrosoftKeyvault)),
				KeyVaultProperties: &databricks.EncryptionV2KeyVaultProperties{
					KeyName:     utils.String(key.Name),
					KeyVersion:  utils.String(key.Version),
					KeyVaultURI: utils.String(key.KeyVaultBaseUrl),
				},
			},
		}

		// make sure the key vault exists
		keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, meta.(*clients.Client).Resource, key.KeyVaultBaseUrl)
		if err != nil || keyVaultIdRaw == nil {
			return fmt.Errorf("retrieving the Resource ID for the customer-managed keys for managed services Key Vault at URL %q: %+v", key.KeyVaultBaseUrl, err)
		}
	}

	// Including the Tags in the workspace parameters will update the tags on
	// the workspace only
	workspace := databricks.Workspace{
		Sku: &databricks.Sku{
			Name: utils.String(skuName),
		},
		Location: utils.String(location),
		WorkspaceProperties: &databricks.WorkspaceProperties{
			PublicNetworkAccess:    publicNetworkAccess,
			ManagedResourceGroupID: &managedResourceGroupID,
			Parameters:             customParams,
		},
		Tags: expandedTags,
	}

	if requireNsgRules != "" {
		workspace.WorkspaceProperties.RequiredNsgRules = databricks.RequiredNsgRules(requireNsgRules)
	}

	if encrypt.Entities != nil {
		workspace.WorkspaceProperties.Encryption = encrypt
	}

	future, err := client.CreateOrUpdate(ctx, workspace, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", id, err)
	}

	// Only call Update(e.g. PATCH) if it is not a new resource and the Tags have changed
	// this will cause the updated tags to be propagated to all of the connected
	// workspace resources.
	// TODO: can be removed once https://github.com/Azure/azure-sdk-for-go/issues/14571 is fixed
	if !d.IsNewResource() && d.HasChange("tags") {
		workspaceUpdate := databricks.WorkspaceUpdate{
			Tags: expandedTags,
		}

		future, err := client.Update(ctx, workspaceUpdate, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("updating %s Tags: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for %s Tags to be updated: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	// I have to set the custom_parameters so I can pass the public and private
	// subnet NSG association along with the backend Pool Id since they are not
	// returned in the read from Azure...
	custom, backendPoolReadId := flattenWorkspaceCustomParameters(customParams, pubSubAssoc, priSubAssoc)

	if backendPoolReadId != nil {
		d.Set("load_balancer_backend_address_pool_id", backendPoolReadId)
	}

	if err := d.Set("custom_parameters", custom); err != nil {
		return fmt.Errorf("setting `custom_parameters`: %+v", err)
	}

	if encrypt != nil && keyIdRaw != "" {
		d.Set("managed_services_cmk_key_vault_key_id", keyIdRaw)
	}

	return resourceDatabricksWorkspaceRead(d, meta)
}

func resourceDatabricksWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", sku.Name)
	}

	if props := resp.WorkspaceProperties; props != nil {
		managedResourceGroupID, err := resourcesParse.ResourceGroupID(*props.ManagedResourceGroupID)
		if err != nil {
			return err
		}
		d.Set("managed_resource_group_id", props.ManagedResourceGroupID)
		d.Set("managed_resource_group_name", managedResourceGroupID.ResourceGroup)
		d.Set("public_network_access_enabled", (props.PublicNetworkAccess == databricks.PublicNetworkAccessEnabled))

		if props.PublicNetworkAccess == databricks.PublicNetworkAccessDisabled {
			d.Set("network_security_group_rules_required", string(props.RequiredNsgRules))
		}
		var cmkEnabled, infraEnabled *bool

		if props.Parameters != nil {
			if props.Parameters.PrepareEncryption != nil {
				cmkEnabled = props.Parameters.PrepareEncryption.Value
				d.Set("customer_managed_key_enabled", &props.Parameters.PrepareEncryption.Value)
			}

			if props.Parameters.RequireInfrastructureEncryption != nil {
				infraEnabled = props.Parameters.RequireInfrastructureEncryption.Value
				d.Set("infrastructure_encryption_enabled", &props.Parameters.RequireInfrastructureEncryption.Value)
			}

			// The subnet associations only exist in the statefile, so we need to do a Get before we Set
			// with what has come back from the Azure response...
			customParamsRaw := d.Get("custom_parameters").([]interface{})
			_, pubSubAssoc, priSubAssoc := expandWorkspaceCustomParameters(customParamsRaw, *cmkEnabled, *infraEnabled, "", "")

			custom, backendPoolReadId := flattenWorkspaceCustomParameters(props.Parameters, pubSubAssoc, priSubAssoc)
			if err := d.Set("custom_parameters", custom); err != nil {
				return fmt.Errorf("setting `custom_parameters`: %+v", err)
			}

			if backendPoolReadId != nil {
				d.Set("load_balancer_backend_address_pool_id", backendPoolReadId)
			}
		}

		if err := d.Set("storage_account_identity", flattenWorkspaceStorageAccountIdentity(props.StorageAccountIdentity)); err != nil {
			return fmt.Errorf("setting `storage_account_identity`: %+v", err)
		}

		if props.WorkspaceURL != nil {
			d.Set("workspace_url", props.WorkspaceURL)
		}

		if props.WorkspaceID != nil {
			d.Set("workspace_id", props.WorkspaceID)
		}
	}

	// customer managed key for managed services
	encryptKeyName := ""
	encryptKeyVersion := ""
	encryptKeyVaultURI := ""

	if resp.WorkspaceProperties.Encryption != nil {
		if props := resp.WorkspaceProperties.Encryption.Entities; props != nil {
			if props.ManagedServices.KeyVaultProperties.KeyName != nil {
				encryptKeyName = *props.ManagedServices.KeyVaultProperties.KeyName
			}
			if props.ManagedServices.KeyVaultProperties.KeyVersion != nil {
				encryptKeyVersion = *props.ManagedServices.KeyVaultProperties.KeyVersion
			}
			if props.ManagedServices.KeyVaultProperties.KeyVaultURI != nil {
				encryptKeyVaultURI = *props.ManagedServices.KeyVaultProperties.KeyVaultURI
			}
		}
	}

	if encryptKeyVaultURI != "" {
		key, err := keyVaultParse.NewNestedItemID(encryptKeyVaultURI, "keys", encryptKeyName, encryptKeyVersion)
		if err == nil {
			d.Set("managed_services_cmk_key_vault_key_id", key.ID())
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDatabricksWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
		}
	}

	return nil
}

func flattenWorkspaceStorageAccountIdentity(input *databricks.ManagedIdentityConfiguration) []interface{} {
	if input == nil {
		return nil
	}

	e := make(map[string]interface{})

	if v := input; v != nil {
		if t := v.PrincipalID; t != nil {
			if t != nil {
				e["principal_id"] = t.String()
			}
		}

		if t := v.TenantID; t != nil {
			if t != nil {
				e["tenant_id"] = t.String()
			}
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

func flattenWorkspaceCustomParameters(input *databricks.WorkspaceCustomParameters, publicSubnetAssociation, privateSubnetAssociation *string) ([]interface{}, *string) {
	if input == nil {
		return nil, nil
	}

	var backendAddressPoolId, backendName, loadBalancerId *string
	parameters := make(map[string]interface{})

	if publicSubnetAssociation != nil && *publicSubnetAssociation != "" {
		parameters["public_subnet_network_security_group_association_id"] = *publicSubnetAssociation
	}

	if privateSubnetAssociation != nil && *privateSubnetAssociation != "" {
		parameters["private_subnet_network_security_group_association_id"] = *privateSubnetAssociation
	}

	if v := input.LoadBalancerBackendPoolName; v != nil {
		if v.Value != nil {
			backendName = v.Value
		}
	}

	if v := input.LoadBalancerID; v != nil {
		if v.Value != nil {
			loadBalancerId = v.Value
		}
	}

	if v := input.NatGatewayName; v != nil {
		if v.Value != nil {
			parameters["nat_gateway_name"] = *v.Value
		}
	}

	if v := input.PublicIPName; v != nil {
		if v.Value != nil {
			parameters["public_ip_name"] = *v.Value
		}
	}

	if v := input.StorageAccountName; v != nil {
		if v.Value != nil {
			parameters["storage_account_name"] = *v.Value
		}
	}

	if v := input.StorageAccountSkuName; v != nil {
		if v.Value != nil {
			parameters["storage_account_sku_name"] = *v.Value
		}
	}

	if v := input.VnetAddressPrefix; v != nil {
		if v.Value != nil {
			parameters["vnet_address_prefix"] = *v.Value
		}
	}

	if v := input.AmlWorkspaceID; v != nil {
		if v.Value != nil {
			parameters["machine_learning_workspace_id"] = *v.Value
		}
	}

	if v := input.EnableNoPublicIP; v != nil {
		if v.Value != nil {
			parameters["no_public_ip"] = *v.Value
		}
	}

	if v := input.CustomPrivateSubnetName; v != nil {
		if v.Value != nil {
			parameters["private_subnet_name"] = *v.Value
		}
	}

	if v := input.CustomPublicSubnetName; v != nil {
		if v.Value != nil {
			parameters["public_subnet_name"] = *v.Value
		}
	}

	if v := input.CustomVirtualNetworkID; v != nil {
		if v.Value != nil {
			parameters["virtual_network_id"] = *v.Value
		}
	}

	if backendName != nil && loadBalancerId != nil {
		lbId, err := loadBalancerParse.LoadBalancerID(*loadBalancerId)

		if err == nil {
			backendId := loadBalancerParse.NewLoadBalancerBackendAddressPoolID(lbId.SubscriptionId, lbId.ResourceGroup, lbId.Name, *backendName)
			backendAddressPoolId = utils.String(backendId.ID())
		}
	}

	return []interface{}{parameters}, backendAddressPoolId
}

func expandWorkspaceCustomParameters(input []interface{}, customerManagedKeyEnabled, infrastructureEncryptionEnabled bool, backendAddressPoolName, loadBalancerId string) (workspaceCustomParameters *databricks.WorkspaceCustomParameters, publicSubnetAssociation, privateSubnetAssociation *string) {
	if len(input) == 0 || input[0] == nil {
		// This will be hit when there are no custom params set but we still
		// need to pass the customerManagedKeyEnabled and infrastructureEncryptionEnabled
		// flags anyway...
		parameters := databricks.WorkspaceCustomParameters{}

		parameters.PrepareEncryption = &databricks.WorkspaceCustomBooleanParameter{
			Value: &customerManagedKeyEnabled,
		}

		parameters.RequireInfrastructureEncryption = &databricks.WorkspaceCustomBooleanParameter{
			Value: &infrastructureEncryptionEnabled,
		}

		return &parameters, nil, nil
	}

	config := input[0].(map[string]interface{})
	var pubSubnetAssoc, priSubnetAssoc *string
	parameters := databricks.WorkspaceCustomParameters{}

	if v, ok := config["public_subnet_network_security_group_association_id"].(string); ok && v != "" {
		pubSubnetAssoc = &v
	}

	if v, ok := config["private_subnet_network_security_group_association_id"].(string); ok && v != "" {
		priSubnetAssoc = &v
	}

	if backendAddressPoolName != "" {
		parameters.LoadBalancerBackendPoolName = &databricks.WorkspaceCustomStringParameter{
			Value: utils.String(backendAddressPoolName),
		}
	}

	if loadBalancerId != "" {
		parameters.LoadBalancerID = &databricks.WorkspaceCustomStringParameter{
			Value: utils.String(loadBalancerId),
		}
	}

	if v, ok := config["nat_gateway_name"].(string); ok && v != "" {
		parameters.NatGatewayName = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v, ok := config["public_ip_name"].(string); ok && v != "" {
		parameters.PublicIPName = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v, ok := config["storage_account_name"].(string); ok && v != "" {
		parameters.StorageAccountName = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v, ok := config["storage_account_sku_name"].(string); ok && v != "" {
		parameters.StorageAccountSkuName = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v, ok := config["vnet_address_prefix"].(string); ok && v != "" {
		parameters.VnetAddressPrefix = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v, ok := config["machine_learning_workspace_id"].(string); ok && v != "" {
		parameters.AmlWorkspaceID = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v, ok := config["no_public_ip"].(bool); ok {
		parameters.EnableNoPublicIP = &databricks.WorkspaceCustomBooleanParameter{
			Value: &v,
		}
	}

	if v, ok := config["public_subnet_name"].(string); ok && v != "" {
		parameters.CustomPublicSubnetName = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	parameters.PrepareEncryption = &databricks.WorkspaceCustomBooleanParameter{
		Value: &customerManagedKeyEnabled,
	}

	parameters.RequireInfrastructureEncryption = &databricks.WorkspaceCustomBooleanParameter{
		Value: &infrastructureEncryptionEnabled,
	}

	if v, ok := config["private_subnet_name"].(string); ok && v != "" {
		parameters.CustomPrivateSubnetName = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v, ok := config["virtual_network_id"].(string); ok && v != "" {
		parameters.CustomVirtualNetworkID = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	return &parameters, pubSubnetAssoc, priSubnetAssoc
}

func workspaceCustomParametersString() []string {
	return []string{"custom_parameters.0.machine_learning_workspace_id", "custom_parameters.0.no_public_ip",
		"custom_parameters.0.public_subnet_name", "custom_parameters.0.private_subnet_name", "custom_parameters.0.virtual_network_id",
		"custom_parameters.0.public_subnet_network_security_group_association_id", "custom_parameters.0.private_subnet_network_security_group_association_id",
		"custom_parameters.0.nat_gateway_name", "custom_parameters.0.public_ip_name", "custom_parameters.0.storage_account_name", "custom_parameters.0.storage_account_sku_name",
		"custom_parameters.0.vnet_address_prefix",
	}
}
