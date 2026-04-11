// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = KubernetesClusterAutomaticResource{}

type KubernetesClusterAutomaticResource struct{}

func (r KubernetesClusterAutomaticResource) ModelObject() interface{} {
	return &KubernetesClusterAutomaticModel{}
}

type KubernetesClusterAutomaticModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	KubernetesVersion string            `tfschema:"kubernetes_version"`
	NodeResourceGroup string            `tfschema:"node_resource_group"`
	Tags              map[string]string `tfschema:"tags"`

	FQDN              string `tfschema:"fqdn"`
	PortalFQDN        string `tfschema:"portal_fqdn"`
	CurrentKubeConfig string `tfschema:"kube_config_raw"`
}

func (r KubernetesClusterAutomaticResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateKubernetesClusterID
}

func (r KubernetesClusterAutomaticResource) ResourceType() string {
	return "azurerm_kubernetes_cluster_automatic"
}

func (r KubernetesClusterAutomaticResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9]$|^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,61}[a-zA-Z0-9]$`),
				"AKS Cluster names must be between 1 and 63 characters in length, must begin and end with an alphanumeric character, and may contain only alphanumeric characters, underscores, and hyphens",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemOrUserAssignedIdentityRequired(),

		"kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C Azure will assign the latest recommended version if not specified
			Computed: true,
		},

		"node_resource_group": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C Azure creates a managed resource group with a generated name if not specified
			Computed: true,
			ForceNew: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r KubernetesClusterAutomaticResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kube_config_raw": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"portal_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r KubernetesClusterAutomaticResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model KubernetesClusterAutomaticModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := commonids.NewKubernetesClusterID(subscriptionId, model.ResourceGroupName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			expandedIdentity, err := expandKubernetesClusterAutomaticIdentity(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			agentPoolProfiles := []managedclusters.ManagedClusterAgentPoolProfile{
				{
					Name:  "systempool",
					Mode:  pointer.To(managedclusters.AgentPoolModeSystem),
					Count: pointer.To(int64(3)),
				},
			}

			parameters := managedclusters.ManagedCluster{
				Location: location.Normalize(model.Location),
				Sku: &managedclusters.ManagedClusterSKU{
					Name: pointer.To(managedclusters.ManagedClusterSKUNameAutomatic),
					Tier: pointer.To(managedclusters.ManagedClusterSKUTierStandard),
				},
				Identity: expandedIdentity,
				Properties: &managedclusters.ManagedClusterProperties{
					AgentPoolProfiles: &agentPoolProfiles,
					KubernetesVersion: pointer.To(model.KubernetesVersion),
					NodeResourceGroup: pointer.To(model.NodeResourceGroup),
				},
				Tags: pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, managedclusters.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KubernetesClusterAutomaticResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesClustersClient

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			credentialsResp, err := client.ListClusterUserCredentials(ctx, *id, managedclusters.DefaultListClusterUserCredentialsOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving user credentials for %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			state := KubernetesClusterAutomaticModel{
				Name:              id.ManagedClusterName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
				Tags:              pointer.From(model.Tags),
			}

			if err := metadata.ResourceData.Set("identity", flattenKubernetesClusterAutomaticIdentity(model.Identity)); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if props := model.Properties; props != nil {
				state.KubernetesVersion = pointer.From(props.KubernetesVersion)
				state.NodeResourceGroup = pointer.From(props.NodeResourceGroup)
				state.FQDN = pointer.From(props.Fqdn)
				state.PortalFQDN = pointer.From(props.AzurePortalFQDN)
			}

			if credentialsModel := credentialsResp.Model; credentialsModel != nil {
				if kubeconfigs := credentialsModel.Kubeconfigs; kubeconfigs != nil && len(*kubeconfigs) > 0 {
					adminKubeConfigRaw := (*kubeconfigs)[0].Value
					if adminKubeConfigRaw != nil {
						rawConfig := *adminKubeConfigRaw
						if base64IsEncoded(*adminKubeConfigRaw) {
							rawConfig = base64Decode(*adminKubeConfigRaw)
						}
						state.CurrentKubeConfig = rawConfig
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r KubernetesClusterAutomaticResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesClustersClient

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, managedclusters.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesClusterAutomaticResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesClustersClient

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KubernetesClusterAutomaticModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: `model` was nil", *id)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("kubernetes_version") {
				if payload.Properties == nil {
					payload.Properties = &managedclusters.ManagedClusterProperties{}
				}
				payload.Properties.KubernetesVersion = pointer.To(model.KubernetesVersion)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := expandKubernetesClusterAutomaticIdentity(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload, managedclusters.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandKubernetesClusterAutomaticIdentity(input []interface{}) (*identity.SystemOrUserAssignedMap, error) {
	return identity.ExpandSystemOrUserAssignedMap(input)
}

func flattenKubernetesClusterAutomaticIdentity(input *identity.SystemOrUserAssignedMap) *[]interface{} {
	if input == nil {
		return &[]interface{}{}
	}

	result, err := identity.FlattenSystemOrUserAssignedMap(input)
	if err != nil {
		return &[]interface{}{}
	}

	return result
}
