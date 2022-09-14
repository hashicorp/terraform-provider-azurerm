package machinelearning

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/machinelearningservices/mgmt/2021-07-01/machinelearningservices"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	appInsightsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/validate"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkspaceSku string

const (
	Basic WorkspaceSku = "Basic"
)

func resourceMachineLearningWorkspace() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceMachineLearningWorkspaceCreate,
		Read:   resourceMachineLearningWorkspaceRead,
		Update: resourceMachineLearningWorkspaceUpdate,
		Delete: resourceMachineLearningWorkspaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspaceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"application_insights_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: appInsightsValidate.ComponentID,
				// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: keyVaultValidate.VaultID,
				// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
				// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

			"primary_user_assigned_identity": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			},

			"container_registry_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: containerValidate.RegistryID,
				// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"public_access_behind_virtual_network_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: !features.FourPointOhBeta(),
				ForceNew: true,
				ConflictsWith: func() []string {
					if !features.FourPointOhBeta() {
						return []string{"public_network_access_enabled"}
					}
					return []string{}
				}(),
			},

			"image_build_compute_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"encryption": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyVaultValidate.VaultID,
						},
						"key_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},
						"user_assigned_identity_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
							// TODO: remove this
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"friendly_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"high_business_impact": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"sku_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(Basic),
				ValidateFunc: validation.StringInSlice([]string{string(Basic)}, false),
			},

			"discovery_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}

	if !features.FourPointOhBeta() {
		// For the time being we should just deprecate and remove this property since it's broken in the API - it doesn't
		// actually set the property and also isn't returned by the API. Once https://github.com/Azure/azure-rest-api-specs/issues/18340
		// is fixed we can reassess how to deal with this field.
		resource.Schema["public_network_access_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			Deprecated:    "`public_network_access_enabled` will be removed in favour of the property `public_access_behind_virtual_network_enabled` in version 4.0 of the AzureRM Provider.",
			ConflictsWith: []string{"public_access_behind_virtual_network_enabled"},
		}
	}

	return resource
}

func resourceMachineLearningWorkspaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_machine_learning_workspace", id.ID())
	}

	expandedIdentity, err := expandMachineLearningWorkspaceIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	expandedEncryption := expandMachineLearningWorkspaceEncryption(d.Get("encryption").([]interface{}))

	networkAccessBehindVnetEnabled := false
	if !features.FourPointOhBeta() {
		if v, ok := d.GetOkExists("public_network_access_enabled"); ok {
			networkAccessBehindVnetEnabled = v.(bool)
		}
	}
	if v, ok := d.GetOkExists("public_access_behind_virtual_network_enabled"); ok {
		networkAccessBehindVnetEnabled = v.(bool)
	}

	workspace := machinelearningservices.Workspace{
		Name:     utils.String(id.Name),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku: &machinelearningservices.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
			Tier: utils.String(d.Get("sku_name").(string)),
		},
		Identity: expandedIdentity,
		WorkspaceProperties: &machinelearningservices.WorkspaceProperties{
			Encryption:                      expandedEncryption,
			StorageAccount:                  utils.String(d.Get("storage_account_id").(string)),
			ApplicationInsights:             utils.String(d.Get("application_insights_id").(string)),
			KeyVault:                        utils.String(d.Get("key_vault_id").(string)),
			AllowPublicAccessWhenBehindVnet: utils.Bool(networkAccessBehindVnetEnabled),
		},
	}

	if v, ok := d.GetOk("description"); ok {
		workspace.WorkspaceProperties.Description = utils.String(v.(string))
	}

	if v, ok := d.GetOk("friendly_name"); ok {
		workspace.WorkspaceProperties.FriendlyName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("container_registry_id"); ok {
		workspace.WorkspaceProperties.ContainerRegistry = utils.String(v.(string))
	}

	if v, ok := d.GetOk("high_business_impact"); ok {
		workspace.WorkspaceProperties.HbiWorkspace = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("image_build_compute_name"); ok {
		workspace.WorkspaceProperties.ImageBuildCompute = utils.String(v.(string))
	}

	if v, ok := d.GetOk("primary_user_assigned_identity"); ok {
		workspace.WorkspaceProperties.PrimaryUserAssignedIdentity = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, workspace)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMachineLearningWorkspaceRead(d, meta)
}

func resourceMachineLearningWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Workspace ID `%q`: %+v", d.Id(), err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if props := resp.WorkspaceProperties; props != nil {
		d.Set("application_insights_id", props.ApplicationInsights)
		d.Set("storage_account_id", props.StorageAccount)
		d.Set("key_vault_id", props.KeyVault)
		d.Set("container_registry_id", props.ContainerRegistry)
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("high_business_impact", props.HbiWorkspace)
		d.Set("image_build_compute_name", props.ImageBuildCompute)
		d.Set("discovery_url", props.DiscoveryURL)
		d.Set("primary_user_assigned_identity", props.PrimaryUserAssignedIdentity)
		d.Set("public_access_behind_virtual_network_enabled", props.AllowPublicAccessWhenBehindVnet)

		if !features.FourPointOhBeta() {
			d.Set("public_network_access_enabled", props.AllowPublicAccessWhenBehindVnet)
		}
	}

	flattenedIdentity, err := flattenMachineLearningWorkspaceIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	flattenedEncryption, err := flattenMachineLearningWorkspaceEncryption(resp.Encryption)
	if err != nil {
		return fmt.Errorf("flattening `encryption`: %+v", err)
	}
	if err := d.Set("encryption", flattenedEncryption); err != nil {
		return fmt.Errorf("flattening encryption on Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMachineLearningWorkspaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return err
	}

	update := machinelearningservices.WorkspaceUpdateParameters{
		WorkspacePropertiesUpdateParameters: &machinelearningservices.WorkspacePropertiesUpdateParameters{},
	}

	if d.HasChange("sku_name") {
		skuName := d.Get("sku_name").(string)
		update.Sku = &machinelearningservices.Sku{
			Name: &skuName,
			Tier: &skuName,
		}
	}

	if d.HasChange("description") {
		update.WorkspacePropertiesUpdateParameters.Description = utils.String(d.Get("description").(string))
	}

	if d.HasChange("friendly_name") {
		update.WorkspacePropertiesUpdateParameters.FriendlyName = utils.String(d.Get("friendly_name").(string))
	}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("identity") {
		identity, err := expandMachineLearningWorkspaceIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return err
		}
		update.Identity = identity
	}

	if d.HasChange("primary_user_assigned_identity") {
		update.WorkspacePropertiesUpdateParameters.PrimaryUserAssignedIdentity = utils.String(d.Get("primary_user_assigned_identity").(string))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, update); err != nil {
		return fmt.Errorf("updating Machine Learning Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceMachineLearningWorkspaceRead(d, meta)
}

func resourceMachineLearningWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Workspace ID `%q`: %+v", d.Id(), err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Machine Learning Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Machine Learning Workspace %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandMachineLearningWorkspaceIdentity(input []interface{}) (*machinelearningservices.Identity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := machinelearningservices.Identity{
		Type: machinelearningservices.ResourceIdentityType(string(expanded.Type)),
	}
	// api uses `SystemAssigned,UserAssigned` (no space), so convert it from the normalized value
	if expanded.Type == identity.TypeSystemAssignedUserAssigned {
		out.Type = machinelearningservices.ResourceIdentityTypeSystemAssignedUserAssigned
	}
	if len(expanded.IdentityIds) > 0 {
		out.UserAssignedIdentities = make(map[string]*machinelearningservices.UserAssignedIdentity)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &machinelearningservices.UserAssignedIdentity{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func flattenMachineLearningWorkspaceIdentity(input *machinelearningservices.Identity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}

		// api uses `SystemAssigned,UserAssigned` (no space), so normalize it back
		if input.Type == machinelearningservices.ResourceIdentityTypeSystemAssignedUserAssigned {
			transform.Type = identity.TypeSystemAssignedUserAssigned
		}

		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func expandMachineLearningWorkspaceEncryption(input []interface{}) *machinelearningservices.EncryptionProperty {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	out := machinelearningservices.EncryptionProperty{
		Identity: &machinelearningservices.IdentityForCmk{
			UserAssignedIdentity: nil,
		},
		KeyVaultProperties: &machinelearningservices.KeyVaultProperties{
			KeyVaultArmID: utils.String(raw["key_vault_id"].(string)),
			KeyIdentifier: utils.String(raw["key_id"].(string)),
		},
		Status: machinelearningservices.EncryptionStatusEnabled,
	}

	if raw["user_assigned_identity_id"].(string) != "" {
		out.Identity.UserAssignedIdentity = utils.String(raw["user_assigned_identity_id"].(string))
	}

	return &out
}

func flattenMachineLearningWorkspaceEncryption(input *machinelearningservices.EncryptionProperty) (*[]interface{}, error) {
	if input == nil || input.Status != machinelearningservices.EncryptionStatusEnabled {
		return &[]interface{}{}, nil
	}

	keyVaultId := ""
	keyVaultKeyId := ""
	if input.KeyVaultProperties != nil {
		if input.KeyVaultProperties.KeyIdentifier != nil {
			keyVaultKeyId = *input.KeyVaultProperties.KeyIdentifier
		}
		if input.KeyVaultProperties.KeyVaultArmID != nil {
			keyVaultId = *input.KeyVaultProperties.KeyVaultArmID
		}
	}

	userAssignedIdentityId := ""
	if input.Identity != nil && input.Identity.UserAssignedIdentity != nil {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.Identity.UserAssignedIdentity)
		if err != nil {
			return nil, fmt.Errorf("parsing userAssignedIdentityId %q: %+v", *input.Identity.UserAssignedIdentity, err)
		}

		userAssignedIdentityId = id.ID()
	}

	return &[]interface{}{
		map[string]interface{}{
			"user_assigned_identity_id": userAssignedIdentityId,
			"key_vault_id":              keyVaultId,
			"key_id":                    keyVaultKeyId,
		},
	}, nil
}
