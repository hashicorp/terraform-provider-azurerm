// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMachineLearningWorkspaceRead(d *pluginsdk.ResourceData, meta any) error {
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
				applicationInsightsId, err := componentsapis.ParseComponentIDInsensitively(*props.ApplicationInsights)
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
			d.Set("service_side_encryption_enabled", props.EnableServiceSideCMKEncryption)
			d.Set("v1_legacy_mode_enabled", props.V1LegacyMode)
			d.Set("storage_account_access_type", pointer.FromEnum(props.SystemDatastoresAuthMode))
			d.Set("workspace_id", props.WorkspaceId)
			d.Set("managed_network", flattenMachineLearningWorkspaceManagedNetwork(props.ManagedNetwork, props.ProvisionNetworkNow))
			d.Set("serverless_compute", flattenMachineLearningWorkspaceServerlessCompute(props.ServerlessComputeSettings))

			kvId, err := commonids.ParseKeyVaultIDInsensitively(*props.KeyVault)
			if err != nil {
				return err
			}
			d.Set("key_vault_id", kvId.ID())

			if err := d.Set("feature_store", flattenMachineLearningWorkspaceFeatureStore(props.FeatureStoreSettings)); err != nil {
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
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

type WorkspaceSku string

const (
	Basic WorkspaceSku = "Basic"
)
