// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMachineLearningWorkspaceUpdate(d *pluginsdk.ResourceData, meta any) error {
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
		expandedIdentity, err := expandMachineLearningWorkspaceIdentity(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		payload.Identity = expandedIdentity
	}

	if d.HasChange("kind") {
		payload.Kind = pointer.To(d.Get("kind").(string))
	}

	if d.HasChange("feature_store") {
		featureStore := expandMachineLearningWorkspaceFeatureStore(d.Get("feature_store").([]any))
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
		payload.Properties.ManagedNetwork, _ = expandMachineLearningWorkspaceManagedNetwork(d.Get("managed_network").([]any))
	}

	if d.HasChange("sku_name") {
		payload.Sku = &workspaces.Sku{
			Name: d.Get("sku_name").(string),
			Tier: pointer.ToEnum[workspaces.SkuTier](d.Get("sku_name").(string)),
		}
	}

	if d.HasChange("v1_legacy_mode_enabled") {
		payload.Properties.V1LegacyMode = pointer.To(d.Get("v1_legacy_mode_enabled").(bool))
	}

	if d.HasChange("storage_account_access_type") {
		payload.Properties.SystemDatastoresAuthMode = pointer.ToEnum[workspaces.SystemDatastoresAuthMode](d.Get("storage_account_access_type").(string))
	}

	if d.HasChange("serverless_compute") {
		serverlessCompute := expandMachineLearningWorkspaceServerlessCompute(d.Get("serverless_compute").([]any))
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
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMachineLearningWorkspaceRead(d, meta)
}
