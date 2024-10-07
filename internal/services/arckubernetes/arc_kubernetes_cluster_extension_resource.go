// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arckubernetes

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	arckubernetes "github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ArcKubernetesClusterExtensionModel struct {
	Name                           string            `tfschema:"name"`
	ClusterID                      string            `tfschema:"cluster_id"`
	ConfigurationProtectedSettings map[string]string `tfschema:"configuration_protected_settings"`
	ConfigurationSettings          map[string]string `tfschema:"configuration_settings"`
	ExtensionType                  string            `tfschema:"extension_type"`
	ReleaseNamespace               string            `tfschema:"release_namespace"`
	ReleaseTrain                   string            `tfschema:"release_train"`
	TargetNamespace                string            `tfschema:"target_namespace"`
	Version                        string            `tfschema:"version"`
	CurrentVersion                 string            `tfschema:"current_version"`
}

type ArcKubernetesClusterExtensionResource struct{}

var _ sdk.ResourceWithUpdate = ArcKubernetesClusterExtensionResource{}

func (r ArcKubernetesClusterExtensionResource) ResourceType() string {
	return "azurerm_arc_kubernetes_cluster_extension"
}

func (r ArcKubernetesClusterExtensionResource) ModelObject() interface{} {
	return &ArcKubernetesClusterExtensionModel{}
}

func (r ArcKubernetesClusterExtensionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
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

		// validate the scope is a connected cluster id
		if _, err := arckubernetes.ParseConnectedClusterID(id.Scope); err != nil {
			errs = append(errs, fmt.Errorf("parsing %q as a Connected Cluster ID: %+v", idRaw, err))
			return
		}

		return
	}
}

func (r ArcKubernetesClusterExtensionResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: arckubernetes.ValidateConnectedClusterID,
		},

		"extension_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

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

		"release_train": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}

}

func (r ArcKubernetesClusterExtensionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"current_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ArcKubernetesClusterExtensionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ArcKubernetesClusterExtensionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ArcKubernetes.ExtensionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			clusterID, err := arckubernetes.ParseConnectedClusterID(model.ClusterID)
			if err != nil {
				return err
			}

			// defined as strings because they're not enums in the swagger https://github.com/Azure/azure-rest-api-specs/pull/23545
			connectedClusterId := arckubernetes.NewConnectedClusterID(subscriptionId, clusterID.ResourceGroupName, clusterID.ConnectedClusterName)
			id := extensions.NewScopedExtensionID(connectedClusterId.ID(), model.Name)
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
				Properties: &extensions.ExtensionProperties{
					AutoUpgradeMinorVersion:        &autoUpgradeMinorVersion,
					ConfigurationProtectedSettings: &model.ConfigurationProtectedSettings,
					ConfigurationSettings:          &model.ConfigurationSettings,
				},
			}

			identityValue, err := identity.ExpandSystemAssigned(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			properties.Identity = identityValue

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

			if err := client.CreateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ArcKubernetesClusterExtensionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ExtensionsClient

			id, err := extensions.ParseScopedExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ArcKubernetesClusterExtensionModel
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

func (r ArcKubernetesClusterExtensionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ExtensionsClient

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

			clusterId, err := arckubernetes.ParseConnectedClusterID(id.Scope)
			if err != nil {
				return fmt.Errorf("parsing %q as a Connected Cluster ID: %+v", id.Scope, err)
			}

			state := ArcKubernetesClusterExtensionModel{
				Name:      id.ExtensionName,
				ClusterID: clusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if err = metadata.ResourceData.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				if properties := model.Properties; properties != nil {
					var originalModel ArcKubernetesClusterExtensionModel
					if err := metadata.Decode(&originalModel); err != nil {
						return fmt.Errorf("decoding: %+v", err)
					}

					state.ConfigurationProtectedSettings = originalModel.ConfigurationProtectedSettings
					state.ConfigurationSettings = pointer.From(properties.ConfigurationSettings)
					state.CurrentVersion = pointer.From(properties.CurrentVersion)
					state.ExtensionType = pointer.From(properties.ExtensionType)
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

func (r ArcKubernetesClusterExtensionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ExtensionsClient

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
