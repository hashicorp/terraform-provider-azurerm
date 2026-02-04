// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-07-01/deploymentsafeguards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = KubernetesClusterDeploymentSafeguardResource{}
	_ sdk.ResourceWithUpdate = KubernetesClusterDeploymentSafeguardResource{}
)

type KubernetesClusterDeploymentSafeguardResource struct{}

type KubernetesClusterDeploymentSafeguardResourceModel struct {
	KubernetesClusterId       string   `tfschema:"kubernetes_cluster_id"`
	Level                     string   `tfschema:"level"`
	ExcludedNamespaces        []string `tfschema:"excluded_namespaces"`
	PodSecurityStandardsLevel string   `tfschema:"pod_security_standards_level"`
}

func (r KubernetesClusterDeploymentSafeguardResource) ModelObject() interface{} {
	return &KubernetesClusterDeploymentSafeguardResourceModel{}
}

func (r KubernetesClusterDeploymentSafeguardResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateKubernetesClusterID
}

func (r KubernetesClusterDeploymentSafeguardResource) ResourceType() string {
	return "azurerm_kubernetes_cluster_deployment_safeguard"
}

func (r KubernetesClusterDeploymentSafeguardResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"kubernetes_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKubernetesClusterID,
		},

		"level": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(deploymentsafeguards.PossibleValuesForDeploymentSafeguardsLevel(), false),
		},

		"excluded_namespaces": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"pod_security_standards_level": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      deploymentsafeguards.PodSecurityStandardsLevelPrivileged,
			ValidateFunc: validation.StringInSlice(deploymentsafeguards.PossibleValuesForPodSecurityStandardsLevel(), false),
		},
	}
}

func (r KubernetesClusterDeploymentSafeguardResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r KubernetesClusterDeploymentSafeguardResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.DeploymentSafeguardsClient

			var model KubernetesClusterDeploymentSafeguardResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(model.KubernetesClusterId)
			if err != nil {
				return err
			}

			scopeId := commonids.NewScopeID(kubernetesClusterId.ID())

			existing, err := client.Get(ctx, scopeId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", kubernetesClusterId, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), kubernetesClusterId)
			}

			payload := deploymentsafeguards.DeploymentSafeguard{
				Properties: &deploymentsafeguards.DeploymentSafeguardsProperties{
					Level: deploymentsafeguards.DeploymentSafeguardsLevel(model.Level),
				},
			}

			if len(model.ExcludedNamespaces) > 0 {
				payload.Properties.ExcludedNamespaces = pointer.To(model.ExcludedNamespaces)
			}

			if model.PodSecurityStandardsLevel != "" {
				payload.Properties.PodSecurityStandardsLevel = pointer.ToEnum[deploymentsafeguards.PodSecurityStandardsLevel](model.PodSecurityStandardsLevel)
			}

			if err := client.CreateThenPoll(ctx, scopeId, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", kubernetesClusterId, err)
			}

			metadata.SetID(kubernetesClusterId)
			return nil
		},
	}
}

func (r KubernetesClusterDeploymentSafeguardResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.DeploymentSafeguardsClient

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			scopeId := commonids.NewScopeID(kubernetesClusterId.ID())

			resp, err := client.Get(ctx, scopeId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(kubernetesClusterId)
				}
				return fmt.Errorf("retrieving %s: %+v", kubernetesClusterId, err)
			}

			state := KubernetesClusterDeploymentSafeguardResourceModel{
				KubernetesClusterId: kubernetesClusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Level = string(props.Level)
					state.ExcludedNamespaces = pointer.From(props.ExcludedNamespaces)
					state.PodSecurityStandardsLevel = pointer.FromEnum(props.PodSecurityStandardsLevel)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r KubernetesClusterDeploymentSafeguardResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.DeploymentSafeguardsClient

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KubernetesClusterDeploymentSafeguardResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			scopeId := commonids.NewScopeID(kubernetesClusterId.ID())

			existing, err := client.Get(ctx, scopeId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", kubernetesClusterId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", kubernetesClusterId)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", kubernetesClusterId)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("level") {
				payload.Properties.Level = deploymentsafeguards.DeploymentSafeguardsLevel(model.Level)
			}

			if metadata.ResourceData.HasChange("excluded_namespaces") {
				payload.Properties.ExcludedNamespaces = pointer.To(model.ExcludedNamespaces)
			}

			if metadata.ResourceData.HasChange("pod_security_standards_level") {
				payload.Properties.PodSecurityStandardsLevel = pointer.ToEnum[deploymentsafeguards.PodSecurityStandardsLevel](model.PodSecurityStandardsLevel)
			}

			if err := client.CreateThenPoll(ctx, scopeId, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", kubernetesClusterId, err)
			}

			return nil
		},
	}
}

func (r KubernetesClusterDeploymentSafeguardResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.DeploymentSafeguardsClient

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			scopeId := commonids.NewScopeID(kubernetesClusterId.ID())

			if err := client.DeleteThenPoll(ctx, scopeId); err != nil {
				return fmt.Errorf("deleting %s: %+v", kubernetesClusterId, err)
			}

			return nil
		},
	}
}
