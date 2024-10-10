// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	components "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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
		Create: resourceMachineLearningWorkspaceCreate,
		Read:   resourceMachineLearningWorkspaceRead,
		Update: resourceMachineLearningWorkspaceUpdate,
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
				ValidateFunc: components.ValidateComponentID,
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

			"kind": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Default",
					"FeatureStore",
				}, false),
				Default: "Default",
			},

			"feature_store": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"computer_spark_runtime_version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"offline_connection_name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"online_connection_name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
					},
				},
			},

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
				Default:  true,
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
						"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
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

			"managed_network": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"isolation_mode": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(workspaces.PossibleValuesForIsolationMode(), false),
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
				ForceNew: true,
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

			"serverless_compute": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateSubnetID,
						},
						"public_ip_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
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
		resource.Schema["public_network_access_enabled"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"public_access_behind_virtual_network_enabled"},
		}
	}

	return resource
}

func resourceMachineLearningWorkspaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.Workspaces
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
		return tf.ImportAsExistsError("azurerm_machine_learning_workspace", id.ID())
	}

	expandedIdentity, err := expandMachineLearningWorkspaceIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	expandedEncryption := expandMachineLearningWorkspaceEncryption(d.Get("encryption").([]interface{}))

	networkAccessBehindVnetEnabled := workspaces.PublicNetworkAccessDisabled
	if !features.FourPointOhBeta() {
		// nolint: staticcheck
		if v, ok := d.GetOkExists("public_network_access_enabled"); ok && v.(bool) {
			networkAccessBehindVnetEnabled = workspaces.PublicNetworkAccessEnabled
		}
	} else {
		if v := d.Get("public_network_access_enabled").(bool); v {
			networkAccessBehindVnetEnabled = workspaces.PublicNetworkAccessEnabled
		}
	}

	workspace := workspaces.Workspace{
		Name:     pointer.To(id.WorkspaceName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku: &workspaces.Sku{
			Name: d.Get("sku_name").(string),
			Tier: pointer.To(workspaces.SkuTier(d.Get("sku_name").(string))),
		},
		Kind: pointer.To(d.Get("kind").(string)),

		Identity: expandedIdentity,
		Properties: &workspaces.WorkspaceProperties{
			ApplicationInsights: pointer.To(d.Get("application_insights_id").(string)),
			Encryption:          expandedEncryption,
			KeyVault:            pointer.To(d.Get("key_vault_id").(string)),
			ManagedNetwork:      expandMachineLearningWorkspaceManagedNetwork(d.Get("managed_network").([]interface{})),
			PublicNetworkAccess: pointer.To(networkAccessBehindVnetEnabled),
			StorageAccount:      pointer.To(d.Get("storage_account_id").(string)),
			V1LegacyMode:        pointer.To(d.Get("v1_legacy_mode_enabled").(bool)),
		},
	}

	serverlessCompute := expandMachineLearningWorkspaceServerlessCompute(d.Get("serverless_compute").([]interface{}))
	if serverlessCompute != nil {
		if *serverlessCompute.ServerlessComputeNoPublicIP && serverlessCompute.ServerlessComputeCustomSubnet == nil && networkAccessBehindVnetEnabled == workspaces.PublicNetworkAccessDisabled {
			return fmt.Errorf("`public_ip_enabled` must be set to  `true` if `subnet_id` is not set and `public_network_access_enabled` is `false`")
		}
	}

	workspace.Properties.ServerlessComputeSettings = serverlessCompute

	if v, ok := d.GetOk("description"); ok {
		workspace.Properties.Description = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("friendly_name"); ok {
		workspace.Properties.FriendlyName = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("container_registry_id"); ok {
		workspace.Properties.ContainerRegistry = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("high_business_impact"); ok {
		workspace.Properties.HbiWorkspace = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("image_build_compute_name"); ok {
		workspace.Properties.ImageBuildCompute = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("primary_user_assigned_identity"); ok {
		workspace.Properties.PrimaryUserAssignedIdentity = pointer.To(v.(string))
	}

	featureStore := expandMachineLearningWorkspaceFeatureStore(d.Get("feature_store").([]interface{}))
	if strings.EqualFold(*workspace.Kind, "Default") {
		if featureStore != nil {
			return fmt.Errorf("`feature_store` can only be set when `kind` is `FeatureStore`")
		}
	} else {
		if featureStore == nil {
			return fmt.Errorf("`feature_store` can not be empty when `kind` is `FeatureStore`")
		}
		workspace.Properties.FeatureStoreSettings = featureStore
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, workspace); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMachineLearningWorkspaceRead(d, meta)
}

func resourceMachineLearningWorkspaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.Workspaces
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
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("identity") {
		expandedIdentity, err := expandMachineLearningWorkspaceIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		payload.Identity = expandedIdentity
	}

	if d.HasChange("kind") {
		payload.Kind = pointer.To(d.Get("kind").(string))
	}

	if d.HasChange("feature_store") {
		featureStore := expandMachineLearningWorkspaceFeatureStore(d.Get("feature_store").([]interface{}))
		if strings.EqualFold(*payload.Kind, "Default") {
			if featureStore != nil {
				return fmt.Errorf("`feature_store` can only be set when `kind` is `FeatureStore`")
			}
		} else {
			if featureStore == nil {
				return fmt.Errorf("`feature_store` can not be empty when `kind` is `FeatureStore`")
			}
			payload.Properties.FeatureStoreSettings = featureStore
		}
	}

	if d.HasChange("primary_user_assigned_identity") {
		payload.Properties.PrimaryUserAssignedIdentity = pointer.To(d.Get("primary_user_assigned_identity").(string))
	}

	if d.HasChange("public_network_access_enabled") {
		if d.Get("public_network_access_enabled").(bool) {
			payload.Properties.PublicNetworkAccess = pointer.To(workspaces.PublicNetworkAccessEnabled)
		} else {
			payload.Properties.PublicNetworkAccess = pointer.To(workspaces.PublicNetworkAccessDisabled)
		}
	}

	if d.HasChange("image_build_compute_name") {
		payload.Properties.ImageBuildCompute = pointer.To(d.Get("image_build_compute_name").(string))
	}

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("friendly_name") {
		payload.Properties.FriendlyName = pointer.To(d.Get("friendly_name").(string))
	}

	if d.HasChange("managed_network") {
		payload.Properties.ManagedNetwork = expandMachineLearningWorkspaceManagedNetwork(d.Get("managed_network").([]interface{}))
	}

	if d.HasChange("sku_name") {
		payload.Sku = &workspaces.Sku{
			Name: d.Get("sku_name").(string),
			Tier: pointer.To(workspaces.SkuTier(d.Get("sku_name").(string))),
		}
	}

	if d.HasChange("v1_legacy_mode_enabled") {
		payload.Properties.V1LegacyMode = pointer.To(d.Get("v1_legacy_mode_enabled").(bool))
	}

	if d.HasChange("serverless_compute") {
		serverlessCompute := expandMachineLearningWorkspaceServerlessCompute(d.Get("serverless_compute").([]interface{}))
		if serverlessCompute != nil {
			networkAccessBehindVnetEnabled := false
			if v := payload.Properties.PublicNetworkAccess; v != nil && *v == workspaces.PublicNetworkAccessEnabled {
				networkAccessBehindVnetEnabled = true
			}
			if *serverlessCompute.ServerlessComputeNoPublicIP && serverlessCompute.ServerlessComputeCustomSubnet == nil && !networkAccessBehindVnetEnabled {
				return fmt.Errorf("`public_ip_enabled` must be set to  `true` if `subnet_id` is not set and `public_network_access_enabled` is `false`")
			}

			if serverlessCompute.ServerlessComputeCustomSubnet == nil {
				oldVal, newVal := d.GetChange("serverless_compute.0.public_ip_enabled")
				if oldVal.(bool) && !newVal.(bool) {
					return fmt.Errorf("`public_ip_enabled` cannot be updated from `true` to `false` when `subnet_id` is null or empty")
				}
			}
		}
		payload.Properties.ServerlessComputeSettings = serverlessCompute
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
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
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if sku := model.Sku; sku != nil {
			d.Set("sku_name", sku.Name)
		}
		d.Set("kind", model.Kind)

		flattenedIdentity, err := flattenMachineLearningWorkspaceIdentity(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			appInsightsId := ""
			if props.ApplicationInsights != nil {
				applicationInsightsId, err := components.ParseComponentIDInsensitively(*props.ApplicationInsights)
				if err != nil {
					return err
				}
				appInsightsId = applicationInsightsId.ID()
			}
			d.Set("application_insights_id", appInsightsId)
			d.Set("storage_account_id", props.StorageAccount)
			d.Set("container_registry_id", props.ContainerRegistry)
			d.Set("description", props.Description)
			d.Set("friendly_name", props.FriendlyName)
			d.Set("high_business_impact", props.HbiWorkspace)
			d.Set("image_build_compute_name", props.ImageBuildCompute)
			d.Set("discovery_url", props.DiscoveryURL)
			d.Set("primary_user_assigned_identity", props.PrimaryUserAssignedIdentity)
			d.Set("public_network_access_enabled", *props.PublicNetworkAccess == workspaces.PublicNetworkAccessEnabled)
			d.Set("v1_legacy_mode_enabled", props.V1LegacyMode)
			d.Set("workspace_id", props.WorkspaceId)
			d.Set("managed_network", flattenMachineLearningWorkspaceManagedNetwork(props.ManagedNetwork))
			d.Set("serverless_compute", flattenMachineLearningWorkspaceServerlessCompute(props.ServerlessComputeSettings))

			kvId, err := commonids.ParseKeyVaultIDInsensitively(*props.KeyVault)
			if err != nil {
				return err
			}
			d.Set("key_vault_id", kvId.ID())

			if !features.FourPointOhBeta() {
				d.Set("public_access_behind_virtual_network_enabled", props.AllowPublicAccessWhenBehindVnet)
			}

			featureStoreSettings := flattenMachineLearningWorkspaceFeatureStore(props.FeatureStoreSettings)
			if err := d.Set("feature_store", featureStoreSettings); err != nil {
				return fmt.Errorf("setting `feature_store`: %+v", err)
			}

			flattenedEncryption, err := flattenMachineLearningWorkspaceEncryption(props.Encryption)
			if err != nil {
				return fmt.Errorf("flattening `encryption`: %+v", err)
			}
			if err := d.Set("encryption", flattenedEncryption); err != nil {
				return fmt.Errorf("setting `encryption`: %+v", err)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceMachineLearningWorkspaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.Workspaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Machine Learning Workspace ID `%q`: %+v", d.Id(), err)
	}

	options := workspaces.DefaultDeleteOperationOptions()
	if meta.(*clients.Client).Features.MachineLearning.PurgeSoftDeletedWorkspaceOnDestroy {
		options = workspaces.DeleteOperationOptions{
			ForceToPurge: pointer.To(true),
		}
	}

	future, err := client.Delete(ctx, *id, options)
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
		out.Identity.UserAssignedIdentity = pointer.To(raw["user_assigned_identity_id"].(string))
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

func expandMachineLearningWorkspaceFeatureStore(input []interface{}) *workspaces.FeatureStoreSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	out := workspaces.FeatureStoreSettings{}

	if raw["computer_spark_runtime_version"].(string) != "" {
		out.ComputeRuntime = &workspaces.ComputeRuntimeDto{
			SparkRuntimeVersion: pointer.To(raw["computer_spark_runtime_version"].(string)),
		}
	}

	if raw["offline_connection_name"].(string) != "" {
		out.OfflineStoreConnectionName = pointer.To(raw["offline_connection_name"].(string))
	}

	if raw["online_connection_name"].(string) != "" {
		out.OnlineStoreConnectionName = pointer.To(raw["online_connection_name"].(string))
	}
	return &out
}

func flattenMachineLearningWorkspaceFeatureStore(input *workspaces.FeatureStoreSettings) *[]interface{} {
	if input == nil {
		return &[]interface{}{}
	}

	computerSparkRunTimeVersion := ""
	offlineConnectionName := ""
	onlineConnectionName := ""

	if input.ComputeRuntime != nil && input.ComputeRuntime.SparkRuntimeVersion != nil {
		computerSparkRunTimeVersion = *input.ComputeRuntime.SparkRuntimeVersion
	}
	if input.OfflineStoreConnectionName != nil {
		offlineConnectionName = *input.OfflineStoreConnectionName
	}

	if input.OnlineStoreConnectionName != nil {
		onlineConnectionName = *input.OnlineStoreConnectionName
	}

	return &[]interface{}{
		map[string]interface{}{
			"computer_spark_runtime_version": computerSparkRunTimeVersion,
			"offline_connection_name":        offlineConnectionName,
			"online_connection_name":         onlineConnectionName,
		},
	}
}

func expandMachineLearningWorkspaceManagedNetwork(i []interface{}) *workspaces.ManagedNetworkSettings {
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	v := i[0].(map[string]interface{})

	return &workspaces.ManagedNetworkSettings{
		IsolationMode: pointer.To(workspaces.IsolationMode(v["isolation_mode"].(string))),
	}
}

func flattenMachineLearningWorkspaceManagedNetwork(i *workspaces.ManagedNetworkSettings) *[]interface{} {
	if i == nil {
		return &[]interface{}{}
	}

	out := map[string]interface{}{}

	if i.IsolationMode != nil {
		out["isolation_mode"] = *i.IsolationMode
	}

	return &[]interface{}{out}
}

func expandMachineLearningWorkspaceServerlessCompute(i []interface{}) *workspaces.ServerlessComputeSettings {
	if len(i) == 0 || i[0] == nil {
		return nil
	}

	v := i[0].(map[string]interface{})

	serverlessCompute := workspaces.ServerlessComputeSettings{
		ServerlessComputeNoPublicIP: pointer.To(!v["public_ip_enabled"].(bool)),
	}

	if subnetId, ok := v["subnet_id"].(string); ok && subnetId != "" {
		serverlessCompute.ServerlessComputeCustomSubnet = pointer.To(subnetId)
	}

	return &serverlessCompute
}

func flattenMachineLearningWorkspaceServerlessCompute(i *workspaces.ServerlessComputeSettings) *[]interface{} {
	if i == nil {
		return &[]interface{}{}
	}

	out := map[string]interface{}{}

	if i.ServerlessComputeCustomSubnet != nil {
		out["subnet_id"] = *i.ServerlessComputeCustomSubnet
	}

	if i.ServerlessComputeNoPublicIP != nil {
		out["public_ip_enabled"] = !*i.ServerlessComputeNoPublicIP
	}

	return &[]interface{}{out}
}
