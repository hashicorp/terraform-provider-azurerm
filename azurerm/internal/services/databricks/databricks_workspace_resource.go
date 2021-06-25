package databricks

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/databricks/mgmt/2018-04-01/databricks"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/databricks/validate"
	resourcesParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

			"custom_parameters": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"aml_workspace_id": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"customer_managed_key": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"key_source": {
										Type:             pluginsdk.TypeString,
										Optional:         true,
										Default:          "Default",
										DiffSuppressFunc: suppress.CaseDifference,
										ValidateFunc: validation.StringInSlice([]string{
											string(databricks.Default),
											string(databricks.MicrosoftKeyvault),
										}, true),
									},
									"key_name": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
									"key_version": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
									"key_vault_uri": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},
								},
							},
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"no_public_ip": {
							Type:         pluginsdk.TypeBool,
							ForceNew:     true,
							Optional:     true,
							Default:      false,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"public_subnet_name": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"private_subnet_name": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"virtual_network_id": {
							Type:         pluginsdk.TypeString,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"customer_managed_key_enabled": {
							Type:         pluginsdk.TypeBool,
							Optional:     true,
							Default:      false,
							AtLeastOneOf: workspaceCustomParametersString(),
						},

						"infrastructure_encryption_enabled": {
							Type:         pluginsdk.TypeBool,
							ForceNew:     true,
							Optional:     true,
							Default:      false,
							AtLeastOneOf: workspaceCustomParametersString(),
						},
					},
				},
			},

			"tags": tags.Schema(),

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
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			if d.HasChange("sku") {
				sku, changedSKU := d.GetChange("sku")

				if changedSKU == "trial" {
					log.Printf("[DEBUG] recreate databricks workspace, cannot be migrated to %s", changedSKU)
					d.ForceNew("sku")
				} else {
					log.Printf("[DEBUG] databricks workspace can be upgraded from %s to %s", sku, changedSKU)
				}
			}

			if d.HasChange("custom_parameters") {
				oldParamsRaw, paramsRaw := d.GetChange("custom_parameters")
				params := paramsRaw.([]interface{})
				oldParams := oldParamsRaw.([]interface{})

				if len(params) != 0 && params[0] != nil {
					_, changedSKU := d.GetChange("sku")
					config := params[0].(map[string]interface{})

					var oldConfig map[string]interface{}
					if len(oldParams) != 0 && oldParams[0] != nil {
						oldConfig = oldParams[0].(map[string]interface{})
					}

					cmkRaw := config["customer_managed_key"].([]interface{})
					if len(cmkRaw) != 0 && cmkRaw[0] != nil {
						cmk := cmkRaw[0].(map[string]interface{})

						cmkSource := cmk["key_source"].(string)
						cmkName := cmk["key_name"].(string)
						cmkVersion := cmk["key_version"].(string)
						cmkUri := cmk["key_vault_uri"].(string)

						cmkEncrypt := false
						oldCmkEncrypt := false
						infraEncypt := false

						if v, ok := config["enable_cmk_encryption"].(bool); ok {
							cmkEncrypt = v
						}

						if v, ok := oldConfig["enable_cmk_encryption"].(bool); ok {
							oldCmkEncrypt = v
						}

						if v, ok := config["enable_infrastructure_encryption"].(bool); ok {
							infraEncypt = v
						}

						if cmkEncrypt && infraEncypt {
							return fmt.Errorf("'enable_cmk_encryption' and 'enable_infrastructure_encryption' cannot both be 'true'")
						}

						if oldCmkEncrypt && !cmkEncrypt {
							return fmt.Errorf("'enable_cmk_encryption' cannot be set to 'false' once it has been set to 'true'")
						}

						// Key Source: Default
						if strings.EqualFold(cmkSource, "default") {
							if cmkEncrypt {
								return fmt.Errorf("'enable_cmk_encryption' is only valid if the 'customer_managed_key' 'key_source' is set to 'Microsoft.Keyvault', got %q", cmkSource)
							} else if !cmkEncrypt && (cmkName != "" || cmkVersion != "" || cmkUri != "") {
								return fmt.Errorf("'key_name', 'key_version' and 'key_vault_uri' must be empty if the 'customer_managed_key' 'key_source' is set to 'Default'")
							}
						}

						// Key Source: Microsoft.Keyvault
						if strings.EqualFold(cmkSource, "Microsoft.Keyvault") {
							if cmkEncrypt && (cmkName == "" || cmkVersion == "" || cmkUri == "") {
								return fmt.Errorf("'key_name', 'key_version' and 'key_vault_uri' must be set if the 'customer_managed_key' 'key_source' is set to 'Microsoft.Keyvault'")
							} else if !cmkEncrypt {
								return fmt.Errorf("'enable_cmk_encryption' cannot be 'false' if the 'customer_managed_key' 'key_source' is set to 'Microsoft.Keyvault'")
							}
						}

						if cmkEncrypt && !strings.EqualFold("premium", changedSKU.(string)) {
							return fmt.Errorf("'enable_cmk_encryption' is only available with a 'premium' workspace 'sku', got %q", changedSKU)
						}
					}
				}
			}

			return nil
		}),
	}
}

func resourceDatabricksWorkspaceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
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

	skuName := d.Get("sku").(string)
	managedResourceGroupName := d.Get("managed_resource_group_name").(string)

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	if managedResourceGroupName == "" {
		// no managed resource group name was provided, we use the default pattern
		log.Printf("[DEBUG][azurerm_databricks_workspace] no managed resource group id was provided, we use the default pattern.")
		managedResourceGroupName = fmt.Sprintf("databricks-rg-%s", id.ResourceGroup)
	}
	managedResourceGroupID := resourcesParse.NewResourceGroupID(subscriptionId, managedResourceGroupName).ID()

	customParamsRaw := d.Get("custom_parameters").([]interface{})
	customParams := expandWorkspaceCustomParameters(customParamsRaw)

	if len(customParamsRaw) != 0 && customParamsRaw[0] != nil {
		configCmk := false
		config := customParamsRaw[0].(map[string]interface{})

		if v, ok := config["customer_managed_key"].([]interface{}); ok {
			if len(v) != 0 {
				configCmk = true
			}

			if d.IsNewResource() && configCmk {
				return fmt.Errorf("%s: 'customer_managed_key' cannot be defined during workspace creation, you must define the 'customer_managed_key' once the workspace has been created and key vault access policies have been added.", id)
			}
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
			ManagedResourceGroupID: &managedResourceGroupID,
			Parameters:             customParams,
		},
		Tags: expandedTags,
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

		if err := d.Set("custom_parameters", flattenWorkspaceCustomParameters(props.Parameters)); err != nil {
			return fmt.Errorf("setting `custom_parameters`: %+v", err)
		}

		if err := d.Set("storage_account_identity", flattenWorkspaceStorageAccountIdentity(props.StorageAccountIdentity)); err != nil {
			return fmt.Errorf("setting `storage_account_identity`: %+v", err)
		}

		d.Set("workspace_url", props.WorkspaceURL)
		d.Set("workspace_id", props.WorkspaceID)
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

func flattenWorkspaceCustomParameters(input *databricks.WorkspaceCustomParameters) []interface{} {
	if input == nil {
		return nil
	}

	parameters := make(map[string]interface{})

	if v := input.AmlWorkspaceID; v != nil {
		if v.Value != nil {
			parameters["aml_workspace_id"] = *v.Value
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

	if v := input.Encryption; v != nil {
		e := make(map[string]interface{})

		if t := v.Value.KeySource; t != "" {
			e["key_source"] = t
		}
		if t := v.Value.KeyName; t != nil {
			e["key_name"] = *t
		}
		if t := v.Value.KeyVersion; t != nil {
			e["key_version"] = *t
		}
		if t := v.Value.KeyVaultURI; t != nil {
			e["key_vault_uri"] = *t
		}

		if len(e) != 0 {
			parameters["customer_managed_key"] = []interface{}{e}
		}
	}

	if v := input.PrepareEncryption; v != nil {
		if v.Value != nil {
			parameters["enable_cmk_encryption"] = *v.Value
		}
	}

	if v := input.RequireInfrastructureEncryption; v != nil {
		if v.Value != nil {
			parameters["enable_infrastructure_encryption"] = *v.Value
		}
	}

	return []interface{}{parameters}
}

func expandWorkspaceCustomParameters(input []interface{}) *databricks.WorkspaceCustomParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})
	parameters := databricks.WorkspaceCustomParameters{}

	if v, ok := config["aml_workspace_id"].(string); ok && v != "" {
		parameters.AmlWorkspaceID = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v, ok := config["customer_managed_key"]; ok {
		if v != nil {
			cmkRaw := v.([]interface{})
			if len(cmkRaw) > 0 {
				cmk := cmkRaw[0].(map[string]interface{})
				var keySource string
				var keyName string
				var keyVersion string
				var keyVaultURI string

				if t := cmk["key_source"].(string); t != "" {
					keySource = t
				}
				if t := cmk["key_name"].(string); t != "" {
					keyName = t
				}
				if t := cmk["key_version"].(string); t != "" {
					keyVersion = t
				}
				if t := cmk["key_vault_uri"].(string); t != "" {
					keyVaultURI = t
				}

				parameters.Encryption = &databricks.WorkspaceEncryptionParameter{
					Value: &databricks.Encryption{
						KeySource: databricks.KeySource(keySource),
					},
				}

				// Only set the values if they are not empty strings
				if keyName != "" {
					parameters.Encryption.Value.KeyName = &keyName
				}
				if keyVersion != "" {
					parameters.Encryption.Value.KeyVersion = &keyVersion
				}
				if keyVaultURI != "" {
					parameters.Encryption.Value.KeyVaultURI = &keyVaultURI
				}
			}
		}
	}

	if v, ok := config["no_public_ip"].(bool); ok {
		parameters.EnableNoPublicIP = &databricks.WorkspaceCustomBooleanParameter{
			Value: &v,
		}
	}

	if v := config["public_subnet_name"].(string); v != "" {
		parameters.CustomPublicSubnetName = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v, ok := config["enable_cmk_encryption"].(bool); ok {
		parameters.PrepareEncryption = &databricks.WorkspaceCustomBooleanParameter{
			Value: &v,
		}
	}

	if v, ok := config["enable_infrastructure_encryption"].(bool); ok {
		parameters.RequireInfrastructureEncryption = &databricks.WorkspaceCustomBooleanParameter{
			Value: &v,
		}
	}

	if v := config["private_subnet_name"].(string); v != "" {
		parameters.CustomPrivateSubnetName = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	if v := config["virtual_network_id"].(string); v != "" {
		parameters.CustomVirtualNetworkID = &databricks.WorkspaceCustomStringParameter{
			Value: &v,
		}
	}

	return &parameters
}

func workspaceCustomParametersString() []string {
	return []string{"custom_parameters.0.aml_workspace_id", "custom_parameters.0.customer_managed_key", "custom_parameters.0.no_public_ip", "custom_parameters.0.public_subnet_name",
		"custom_parameters.0.private_subnet_name", "custom_parameters.0.enable_cmk_encryption", "custom_parameters.0.enable_infrastructure_encryption", "custom_parameters.0.virtual_network_id",
	}
}
