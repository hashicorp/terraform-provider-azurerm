package kubernetesconfiguration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ArcKubernetesClusterExtensionModel struct {
	Name                           string            `tfschema:"name"`
	ResourceGroupName              string            `tfschema:"resource_group_name"`
	ClusterName                    string            `tfschema:"cluster_name"`
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
	return extensions.ValidateExtensionID
}

func (r ArcKubernetesClusterExtensionResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := commonArguments()
	arguments["identity"] = commonschema.SystemAssignedIdentityRequiredForceNew()
	return arguments

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

			client := metadata.Client.KubernetesConfiguration.ExtensionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := extensions.NewExtensionID(subscriptionId, model.ResourceGroupName, "Microsoft.Kubernetes", "connectedClusters", model.ClusterName, model.Name)
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
			client := metadata.Client.KubernetesConfiguration.ExtensionsClient

			id, err := extensions.ParseExtensionID(metadata.ResourceData.Id())
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
			client := metadata.Client.KubernetesConfiguration.ExtensionsClient

			id, err := extensions.ParseExtensionID(metadata.ResourceData.Id())
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

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if model.Identity != nil {
				if err = metadata.ResourceData.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}
			}

			state := ArcKubernetesClusterExtensionModel{
				Name:              id.ExtensionName,
				ResourceGroupName: id.ResourceGroupName,
				ClusterName:       id.ClusterName,
			}

			if properties := model.Properties; properties != nil {
				var originalModel ArcKubernetesClusterExtensionModel
				if err := metadata.Decode(&originalModel); err != nil {
					return fmt.Errorf("decoding: %+v", err)
				}

				state.ConfigurationProtectedSettings = originalModel.ConfigurationProtectedSettings

				if properties.ConfigurationSettings != nil {
					state.ConfigurationSettings = *properties.ConfigurationSettings
				}

				if properties.CurrentVersion != nil {
					state.CurrentVersion = *properties.CurrentVersion
				}

				if properties.ExtensionType != nil {
					state.ExtensionType = *properties.ExtensionType
				}

				if properties.ReleaseTrain != nil {
					state.ReleaseTrain = *properties.ReleaseTrain
				}

				if properties.Scope != nil {
					if properties.Scope.Cluster != nil && properties.Scope.Cluster.ReleaseNamespace != nil {
						state.ReleaseNamespace = *properties.Scope.Cluster.ReleaseNamespace
					}

					if properties.Scope.Namespace != nil && properties.Scope.Namespace.TargetNamespace != nil {
						state.TargetNamespace = *properties.Scope.Namespace.TargetNamespace
					}
				}

				if properties.Version != nil {
					state.Version = *properties.Version
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ArcKubernetesClusterExtensionResource) Delete() sdk.ResourceFunc {
	return deleteExtension()
}
