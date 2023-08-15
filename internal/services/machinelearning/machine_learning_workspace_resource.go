// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2023-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	appInsightsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
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
		Create: resourceMachineLearningWorkspaceCreateOrUpdate,
		Read:   resourceMachineLearningWorkspaceRead,
		Update: resourceMachineLearningWorkspaceCreateOrUpdate,
		Delete: resourceMachineLearningWorkspaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := workspaces.ParseWorkspaceID(id)
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

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

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
				ValidateFunc: commonids.ValidateKeyVaultID,
				// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
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
				ValidateFunc: registries.ValidateRegistryID,
				// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
				ConflictsWith: func() []string {
					if !features.FourPointOhBeta() {
						return []string{"public_access_behind_virtual_network_enabled"}
					}
					return []string{}
				}(),
			},

			"image_build_compute_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
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
						"key_vault_id": commonschema.ResourceIDReferenceRequired(commonids.KeyVaultId{}),
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

			"v1_legacy_mode_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"discovery_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"workspace_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}

	if !features.FourPointOhBeta() {
		// For the time being we should just deprecate and remove this property since it's broken in the API - it doesn't
		// actually set the property and also isn't returned by the API. Once https://github.com/Azure/azure-rest-api-specs/issues/18340
		// is fixed we can reassess how to deal with this field.
		resource.Schema["public_access_behind_virtual_network_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			ForceNew:      true,
			Deprecated:    "`public_access_behind_virtual_network_enabled` will be removed in favour of the property `public_network_access_enabled` in version 4.0 of the AzureRM Provider.",
			ConflictsWith: []string{"public_network_access_enabled"},
		}
	}

	return resource
}

func resourceMachineLearningWorkspaceCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.Workspaces
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
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
			return tf.ImportAsExistsError("azurerm_machine_learning_workspace", id.ID())
		}
	}

	expandedIdentity, err := expandMachineLearningWorkspaceIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	expandedEncryption := expandMachineLearningWorkspaceEncryption(d.Get("encryption").([]interface{}))

	networkAccessBehindVnetEnabled := false

	// nolint: staticcheck
	if v, ok := d.GetOkExists("public_network_access_enabled"); ok {
		networkAccessBehindVnetEnabled = v.(bool)
	}

	workspace := workspaces.Workspace{
		Name:     utils.String(id.WorkspaceName),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku: &workspaces.Sku{
			Name: d.Get("sku_name").(string),
			Tier: utils.ToPtr(workspaces.SkuTier(d.Get("sku_name").(string))),
		},

		Identity: expandedIdentity,
		Properties: &workspaces.WorkspaceProperties{
			V1LegacyMode:        utils.ToPtr(d.Get("v1_legacy_mode_enabled").(bool)),
			Encryption:          expandedEncryption,
			StorageAccount:      utils.String(d.Get("storage_account_id").(string)),
			ApplicationInsights: utils.String(d.Get("application_insights_id").(string)),
			KeyVault:            utils.String(d.Get("key_vault_id").(string)),
			PublicNetworkAccess: utils.ToPtr(workspaces.PublicNetworkAccessDisabled),
		},
	}

	if networkAccessBehindVnetEnabled {
		workspace.Properties.PublicNetworkAccess = utils.ToPtr(workspaces.PublicNetworkAccessEnabled)
	}

	if v, ok := d.GetOk("description"); ok {
		workspace.Properties.Description = utils.String(v.(string))
	}

	if v, ok := d.GetOk("friendly_name"); ok {
		workspace.Properties.FriendlyName = utils.String(v.(string))
	}

	if v, ok := d.GetOk("container_registry_id"); ok {
		workspace.Properties.ContainerRegistry = utils.String(v.(string))
	}

	if v, ok := d.GetOk("high_business_impact"); ok {
		workspace.Properties.HbiWorkspace = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("image_build_compute_name"); ok {
		workspace.Properties.ImageBuildCompute = utils.String(v.(string))
	}

	if v, ok := d.GetOk("primary_user_assigned_identity"); ok {
		workspace.Properties.PrimaryUserAssignedIdentity = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id, workspace)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMachineLearningWorkspaceRead(d, meta)
}

func resourceMachineLearningWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.Workspaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Workspace ID `%q`: %+v", d.Id(), err)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
	}

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if location := resp.Model.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Model.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if props := resp.Model.Properties; props != nil {
		d.Set("application_insights_id", props.ApplicationInsights)
		d.Set("storage_account_id", props.StorageAccount)
		d.Set("container_registry_id", props.ContainerRegistry)
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("high_business_impact", props.HbiWorkspace)
		d.Set("image_build_compute_name", props.ImageBuildCompute)
		d.Set("discovery_url", props.DiscoveryUrl)
		d.Set("primary_user_assigned_identity", props.PrimaryUserAssignedIdentity)
		d.Set("public_network_access_enabled", *props.PublicNetworkAccess == workspaces.PublicNetworkAccessEnabled)
		d.Set("v1_legacy_mode_enabled", props.V1LegacyMode)
		d.Set("workspace_id", props.WorkspaceId)

		kvId, err := commonids.ParseKeyVaultIDInsensitively(*props.KeyVault)
		if err != nil {
			return err
		}
		d.Set("key_vault_id", kvId.ID())

		if !features.FourPointOhBeta() {
			d.Set("public_access_behind_virtual_network_enabled", props.AllowPublicAccessWhenBehindVnet)
		}
	}

	flattenedIdentity, err := flattenMachineLearningWorkspaceIdentity(resp.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}

	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	flattenedEncryption, err := flattenMachineLearningWorkspaceEncryption(resp.Model.Properties.Encryption)
	if err != nil {
		return fmt.Errorf("flattening `encryption`: %+v", err)
	}
	if err := d.Set("encryption", flattenedEncryption); err != nil {
		return fmt.Errorf("flattening encryption on Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
	}

	return tags.FlattenAndSet(d, resp.Model.Tags)
}

func resourceMachineLearningWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.Workspaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Workspace ID `%q`: %+v", d.Id(), err)
	}

	future, err := client.Delete(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting Machine Learning Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
	}

	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of Machine Learning Workspace %q (Resource Group %q): %+v", id.WorkspaceName, id.ResourceGroupName, err)
	}

	return nil
}

func expandMachineLearningWorkspaceIdentity(input []interface{}) (*identity.LegacySystemAndUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := identity.LegacySystemAndUserAssignedMap{
		Type: expanded.Type,
	}
	if len(expanded.IdentityIds) > 0 {
		out.IdentityIds = map[string]identity.UserAssignedIdentityDetails{}
		for k := range expanded.IdentityIds {
			out.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func flattenMachineLearningWorkspaceIdentity(input *identity.LegacySystemAndUserAssignedMap) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        input.Type,
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}

		// api uses `SystemAssigned,UserAssigned` (no space), so normalize it back
		if input.Type == "SystemAssigned,UserAssigned" {
			transform.Type = identity.TypeSystemAssignedUserAssigned
		}

		if input.PrincipalId != "" {
			transform.PrincipalId = input.PrincipalId
		}
		if input.TenantId != "" {
			transform.TenantId = input.TenantId
		}

		if input != nil && input.IdentityIds != nil {
			for k, v := range input.IdentityIds {
				transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
					ClientId:    v.ClientId,
					PrincipalId: v.PrincipalId,
				}
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}

func expandMachineLearningWorkspaceEncryption(input []interface{}) *workspaces.EncryptionProperty {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	out := workspaces.EncryptionProperty{
		Identity: &workspaces.IdentityForCmk{
			UserAssignedIdentity: nil,
		},
		KeyVaultProperties: workspaces.EncryptionKeyVaultProperties{
			KeyVaultArmId: raw["key_vault_id"].(string),
			KeyIdentifier: raw["key_id"].(string),
		},
		Status: workspaces.EncryptionStatusEnabled,
	}

	if raw["user_assigned_identity_id"].(string) != "" {
		out.Identity.UserAssignedIdentity = utils.String(raw["user_assigned_identity_id"].(string))
	}

	return &out
}

func flattenMachineLearningWorkspaceEncryption(input *workspaces.EncryptionProperty) (*[]interface{}, error) {
	if input == nil || input.Status != workspaces.EncryptionStatusEnabled {
		return &[]interface{}{}, nil
	}

	keyVaultId := ""
	keyVaultKeyId := ""

	if input.KeyVaultProperties.KeyIdentifier != "" {
		keyVaultKeyId = input.KeyVaultProperties.KeyIdentifier
	}
	if input.KeyVaultProperties.KeyVaultArmId != "" {
		keyVaultId = input.KeyVaultProperties.KeyVaultArmId
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
