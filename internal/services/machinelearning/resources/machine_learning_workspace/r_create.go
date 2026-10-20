// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMachineLearningWorkspaceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MachineLearning.Workspaces
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workspaces.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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

	expandedIdentity, err := expandMachineLearningWorkspaceIdentity(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	expandedEncryption := expandMachineLearningWorkspaceEncryption(d.Get("encryption").([]any))

	managedNetwork, provisionNetworkNow := expandMachineLearningWorkspaceManagedNetwork(d.Get("managed_network").([]any))

	networkAccessBehindVnetEnabled := workspaces.PublicNetworkAccessDisabled

	if v := d.Get("public_network_access_enabled").(bool); v {
		networkAccessBehindVnetEnabled = workspaces.PublicNetworkAccessEnabled
	}

	workspace := workspaces.Workspace{
		Name:     pointer.To(id.WorkspaceName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Sku: &workspaces.Sku{
			Name: d.Get("sku_name").(string),
			Tier: pointer.ToEnum[workspaces.SkuTier](d.Get("sku_name").(string)),
		},
		Kind:     pointer.To(d.Get("kind").(string)),
		Identity: expandedIdentity,
		Properties: &workspaces.WorkspaceProperties{
			ApplicationInsights:            pointer.To(d.Get("application_insights_id").(string)),
			Encryption:                     expandedEncryption,
			KeyVault:                       pointer.To(d.Get("key_vault_id").(string)),
			ManagedNetwork:                 managedNetwork,
			ProvisionNetworkNow:            pointer.To(provisionNetworkNow),
			PublicNetworkAccess:            pointer.To(networkAccessBehindVnetEnabled),
			EnableServiceSideCMKEncryption: pointer.To(d.Get("service_side_encryption_enabled").(bool)),
			StorageAccount:                 pointer.To(d.Get("storage_account_id").(string)),
			SystemDatastoresAuthMode:       pointer.ToEnum[workspaces.SystemDatastoresAuthMode](d.Get("storage_account_access_type").(string)),
			V1LegacyMode:                   pointer.To(d.Get("v1_legacy_mode_enabled").(bool)),
		},
	}

	serverlessCompute := expandMachineLearningWorkspaceServerlessCompute(d.Get("serverless_compute").([]any))
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
		workspace.Properties.HbiWorkspace = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("image_build_compute_name"); ok {
		workspace.Properties.ImageBuildCompute = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("primary_user_assigned_identity"); ok {
		workspace.Properties.PrimaryUserAssignedIdentity = pointer.To(v.(string))
	}

	featureStore := expandMachineLearningWorkspaceFeatureStore(d.Get("feature_store").([]any))
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

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, workspace, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMachineLearningWorkspaceRead(d, meta)
}
