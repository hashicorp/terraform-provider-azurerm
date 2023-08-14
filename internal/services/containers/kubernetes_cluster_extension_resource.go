// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type KubernetesClusterExtensionModel struct {
	Name                           string            `tfschema:"name"`
	ClusterID                      string            `tfschema:"cluster_id"`
	ConfigurationProtectedSettings map[string]string `tfschema:"configuration_protected_settings"`
	ConfigurationSettings          map[string]string `tfschema:"configuration_settings"`
	ExtensionType                  string            `tfschema:"extension_type"`
	Plan                           []PlanModel       `tfschema:"plan"`
	ReleaseNamespace               string            `tfschema:"release_namespace"`
	ReleaseTrain                   string            `tfschema:"release_train"`
	TargetNamespace                string            `tfschema:"target_namespace"`
	Version                        string            `tfschema:"version"`
	CurrentVersion                 string            `tfschema:"current_version"`
}

type PlanModel struct {
	Name          string `tfschema:"name"`
	Product       string `tfschema:"product"`
	PromotionCode string `tfschema:"promotion_code"`
	Publisher     string `tfschema:"publisher"`
	Version       string `tfschema:"version"`
}

type KubernetesClusterExtensionResource struct{}

var _ sdk.ResourceWithUpdate = KubernetesClusterExtensionResource{}

func (r KubernetesClusterExtensionResource) ResourceType() string {
	return "azurerm_kubernetes_cluster_extension"
}

func (r KubernetesClusterExtensionResource) ModelObject() interface{} {
	return &KubernetesClusterExtensionModel{}
}

func (r KubernetesClusterExtensionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(val interface{}, key string) (warns []string, errs []error) {
		idRaw, ok := val.(string)
		if !ok {
			errs = append(errs, fmt.Errorf("expected `id` to be a string but got %+v", val))
			return
		}

		id, err := extensions.ParseScopedExtensionID(idRaw)
		if err != nil {
			errs = append(errs, fmt.Errorf("parsing %q: %+v", idRaw, err))
			return
		}

		// validate the scope is a kubernetes cluster id
		if _, err := commonids.ParseKubernetesClusterID(id.Scope); err != nil {
			errs = append(errs, fmt.Errorf("parsing %q as a Kubernetes Cluster ID: %+v", idRaw, err))
			return
		}

		return
	}
}

func (r KubernetesClusterExtensionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-.]{0,252}$"),
				"name must be between 1 and 253 characters in length and may contain only letters, numbers, periods (.), hyphens (-), and must begin with a letter or number.",
			),
		},

		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKubernetesClusterID,
		},

		"extension_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"configuration_protected_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"configuration_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"product": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"publisher": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"promotion_code": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"release_train": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"version"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"release_namespace": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"target_namespace"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"target_namespace": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"release_namespace"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"version": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"release_train"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},
	}
}

func (r KubernetesClusterExtensionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"aks_assigned_identity": commonschema.SystemAssignedIdentityComputed(),

		"current_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r KubernetesClusterExtensionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model KubernetesClusterExtensionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Containers.KubernetesExtensionsClient
			clusterID, err := commonids.ParseKubernetesClusterID(model.ClusterID)
			if err != nil {
				return err
			}

			// defined as strings because they're not enums in the swagger https://github.com/Azure/azure-rest-api-specs/pull/23545
			id := extensions.NewScopedExtensionID(clusterID.ID(), model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			autoUpgradeMinorVersion := false
			if model.Version == "" {
				autoUpgradeMinorVersion = true
			}

			properties := &extensions.Extension{
				Plan: expandPlanModel(model.Plan),
				Properties: &extensions.ExtensionProperties{
					AutoUpgradeMinorVersion:        &autoUpgradeMinorVersion,
					ConfigurationProtectedSettings: &model.ConfigurationProtectedSettings,
					ConfigurationSettings:          &model.ConfigurationSettings,
				},
			}

			if model.ExtensionType != "" {
				properties.Properties.ExtensionType = &model.ExtensionType
			}

			if model.ReleaseNamespace != "" {
				properties.Properties.Scope = &extensions.Scope{
					Cluster: &extensions.ScopeCluster{
						ReleaseNamespace: &model.ReleaseNamespace,
					},
				}
			}

			if model.ReleaseTrain != "" {
				properties.Properties.ReleaseTrain = &model.ReleaseTrain
			}

			if model.TargetNamespace != "" {
				properties.Properties.Scope = &extensions.Scope{
					Namespace: &extensions.ScopeNamespace{
						TargetNamespace: &model.TargetNamespace,
					},
				}
			}

			if model.Version != "" {
				properties.Properties.Version = &model.Version
			}

			if err = client.CreateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KubernetesClusterExtensionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesExtensionsClient

			id, err := extensions.ParseScopedExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KubernetesClusterExtensionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := &extensions.PatchExtension{
				Properties: &extensions.PatchExtensionProperties{},
			}

			if metadata.ResourceData.HasChange("configuration_protected_settings") {
				properties.Properties.ConfigurationProtectedSettings = &model.ConfigurationProtectedSettings
			}

			if metadata.ResourceData.HasChange("configuration_settings") {
				properties.Properties.ConfigurationSettings = &model.ConfigurationSettings
			}

			if err := client.UpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesClusterExtensionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesExtensionsClient

			id, err := extensions.ParseScopedExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(id.Scope)
			if err != nil {
				return fmt.Errorf("parsing %q as a Kubernetes Cluster ID: %+v", id.Scope, err)
			}
			state := KubernetesClusterExtensionModel{
				Name:      id.ExtensionName,
				ClusterID: kubernetesClusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					var originalModel KubernetesClusterExtensionModel
					if err := metadata.Decode(&originalModel); err != nil {
						return fmt.Errorf("decoding: %+v", err)
					}

					if err = metadata.ResourceData.Set("aks_assigned_identity", flattenAksAssignedIdentity(properties.AksAssignedIdentity)); err != nil {
						return fmt.Errorf("setting `aks_assigned_identity`: %+v", err)
					}

					state.ConfigurationProtectedSettings = originalModel.ConfigurationProtectedSettings
					state.ConfigurationSettings = pointer.From(properties.ConfigurationSettings)
					state.CurrentVersion = pointer.From(properties.CurrentVersion)
					state.ExtensionType = pointer.From(properties.ExtensionType)
					state.Plan = flattenPlanModel(model.Plan)
					state.ReleaseTrain = pointer.From(properties.ReleaseTrain)

					if properties.Scope != nil {
						if properties.Scope.Cluster != nil {
							state.ReleaseNamespace = pointer.From(properties.Scope.Cluster.ReleaseNamespace)
						}

						if properties.Scope.Namespace != nil {
							state.TargetNamespace = pointer.From(properties.Scope.Namespace.TargetNamespace)
						}
					}

					state.Version = pointer.From(properties.Version)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r KubernetesClusterExtensionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesExtensionsClient

			id, err := extensions.ParseScopedExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, extensions.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandPlanModel(inputList []PlanModel) *extensions.Plan {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := extensions.Plan{
		Name:      input.Name,
		Product:   input.Product,
		Publisher: input.Publisher,
	}
	if input.PromotionCode != "" {
		output.PromotionCode = &input.PromotionCode
	}

	if input.Version != "" {
		output.Version = &input.Version
	}

	return &output
}

func flattenPlanModel(input *extensions.Plan) []PlanModel {
	var outputList []PlanModel
	if input == nil {
		return outputList
	}
	output := PlanModel{
		Name:      input.Name,
		Product:   input.Product,
		Publisher: input.Publisher,
	}

	output.PromotionCode = pointer.From(input.PromotionCode)
	output.Version = pointer.From(input.Version)

	return append(outputList, output)
}

func flattenAksAssignedIdentity(input *extensions.ExtensionPropertiesAksAssignedIdentity) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := identity.SystemAssigned{
		Type:        identity.TypeSystemAssigned,
		PrincipalId: pointer.From(input.PrincipalId),
		TenantId:    pointer.From(input.TenantId),
	}

	if input.Type != nil && *input.Type == extensions.AKSIdentityTypeUserAssigned {
		output.Type = identity.TypeUserAssigned
	}

	return identity.FlattenSystemAssigned(&output)
}
