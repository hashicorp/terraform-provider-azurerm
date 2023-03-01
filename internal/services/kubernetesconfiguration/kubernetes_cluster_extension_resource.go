package kubernetesconfiguration

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ClusterResourceType string

const (
	AksResource       ClusterResourceType = "managedClusters"
	ArcResource       ClusterResourceType = "connectedClusters"
	AksHybridResource ClusterResourceType = "provisionedClusters"
)

type ClusterProviderType string

const (
	AksProvider       ClusterProviderType = "Microsoft.ContainerService"
	ArcProvider       ClusterProviderType = "Microsoft.Kubernetes"
	AksHybridProvider ClusterProviderType = "Microsoft.HybridContainerService"
)

type KubernetesClusterExtensionModel struct {
	Name                           string                                        `tfschema:"name"`
	ResourceGroupName              string                                        `tfschema:"resource_group_name"`
	ClusterName                    string                                        `tfschema:"cluster_name"`
	ClusterResourceName            string                                        `tfschema:"cluster_resource_name"`
	AksAssignedIdentity            []ExtensionPropertiesAksAssignedIdentityModel `tfschema:"aks_assigned_identity"`
	ConfigurationProtectedSettings map[string]string                             `tfschema:"configuration_protected_settings"`
	ConfigurationSettings          map[string]string                             `tfschema:"configuration_settings"`
	ExtensionType                  string                                        `tfschema:"extension_type"`
	Plan                           []PlanModel                                   `tfschema:"plan"`
	ReleaseNamespace               string                                        `tfschema:"release_namespace"`
	ReleaseTrain                   string                                        `tfschema:"release_train"`
	TargetNamespace                string                                        `tfschema:"target_namespace"`
	Version                        string                                        `tfschema:"version"`
	CurrentVersion                 string                                        `tfschema:"current_version"`
}

type ExtensionPropertiesAksAssignedIdentityModel struct {
	PrincipalId string                     `tfschema:"principal_id"`
	TenantId    string                     `tfschema:"tenant_id"`
	Type        extensions.AKSIdentityType `tfschema:"type"`
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
	return extensions.ValidateExtensionID
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

		"resource_group_name": commonschema.ResourceGroupName(),

		"cluster_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cluster_resource_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(AksResource),
				string(ArcResource),
				string(AksHybridResource),
			}, false),
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

		"identity": commonschema.SystemAssignedIdentityOptionalForceNew(),

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

			providerName, err := getProviderName(model.ClusterResourceName)
			if err != nil {
				return err
			}

			client := metadata.Client.KubernetesConfiguration.ExtensionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := extensions.NewExtensionID(subscriptionId, model.ResourceGroupName, providerName, model.ClusterResourceName, model.ClusterName, model.Name)
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

			if _, ok := metadata.ResourceData.GetOk("identity"); ok {
				if model.ClusterResourceName != string(ArcResource) {
					return fmt.Errorf("`identity` should not be set when `cluster_resource_name` is `%s`", model.ClusterResourceName)
				}

				identityValue, err := identity.ExpandSystemAssigned(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				properties.Identity = identityValue
			} else {
				if model.ClusterResourceName == string(ArcResource) {
					return fmt.Errorf("`identity` must be set when `cluster_resource_name` is `%s`", model.ClusterResourceName)
				}
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

			if err := client.CreateThenPoll(ctx, id, *properties); err != nil {
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
			client := metadata.Client.KubernetesConfiguration.ExtensionsClient

			id, err := extensions.ParseExtensionID(metadata.ResourceData.Id())
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
				if err := metadata.ResourceData.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}
			}

			state := KubernetesClusterExtensionModel{
				Name:                id.ExtensionName,
				ResourceGroupName:   id.ResourceGroupName,
				ClusterName:         id.ClusterName,
				ClusterResourceName: id.ClusterResourceName,
				Plan:                flattenPlanModel(model.Plan),
			}

			if properties := model.Properties; properties != nil {
				var originalModel KubernetesClusterExtensionModel
				if err := metadata.Decode(&originalModel); err != nil {
					return fmt.Errorf("decoding: %+v", err)
				}

				state.AksAssignedIdentity = flattenExtensionPropertiesAksAssignedIdentityModel(properties.AksAssignedIdentity)
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

func (r KubernetesClusterExtensionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.KubernetesConfiguration.ExtensionsClient

			id, err := extensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, extensions.DeleteOperationOptions{}); err != nil {
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

func flattenExtensionPropertiesAksAssignedIdentityModel(input *extensions.ExtensionPropertiesAksAssignedIdentity) []ExtensionPropertiesAksAssignedIdentityModel {
	var outputList []ExtensionPropertiesAksAssignedIdentityModel
	if input == nil {
		return outputList
	}
	output := ExtensionPropertiesAksAssignedIdentityModel{}
	if input.PrincipalId != nil {
		output.PrincipalId = *input.PrincipalId
	}

	if input.TenantId != nil {
		output.TenantId = *input.TenantId
	}

	if input.Type != nil {
		output.Type = *input.Type
	}

	return append(outputList, output)
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
	if input.PromotionCode != nil {
		output.PromotionCode = *input.PromotionCode
	}

	if input.Version != nil {
		output.Version = *input.Version
	}

	return append(outputList, output)
}

func getProviderName(clusterResourceName string) (string, error) {
	if clusterResourceName == string(AksResource) {
		return string(AksProvider), nil
	}

	if clusterResourceName == string(ArcResource) {
		return string(ArcProvider), nil
	}

	if clusterResourceName == string(AksHybridResource) {
		return string(AksHybridProvider), nil
	}

	return "", fmt.Errorf("provider name not found for cluster resource: `%s`", clusterResourceName)
}
