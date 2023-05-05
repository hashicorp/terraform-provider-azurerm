package arckubernetes

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	arckubernetes "github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2021-10-01/connectedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	FluxGitBranch       string = "branch"
	FluxGitCommit       string = "commit"
	FluxGitReferenceTag string = "tag"
	FluxGitSemverRange  string = "semver"
)

type ArcKubernetesFluxConfigurationModel struct {
	Name                            string                         `tfschema:"name"`
	ClusterID                       string                         `tfschema:"cluster_id"`
	AzureBlob                       []AzureBlobDefinitionModel     `tfschema:"azure_blob"`
	Bucket                          []BucketDefinitionModel        `tfschema:"bucket"`
	GitRepository                   []GitRepositoryDefinitionModel `tfschema:"git_repository"`
	Kustomizations                  []KustomizationDefinitionModel `tfschema:"kustomizations"`
	Namespace                       string                         `tfschema:"namespace"`
	Scope                           fluxconfiguration.ScopeType    `tfschema:"scope"`
	ContinuousReconciliationEnabled bool                           `tfschema:"continuous_reconciliation_enabled"`
}

type AzureBlobDefinitionModel struct {
	AccountKey            string                            `tfschema:"account_key"`
	ContainerName         string                            `tfschema:"container_name"`
	LocalAuthRef          string                            `tfschema:"local_auth_ref"`
	SasToken              string                            `tfschema:"sas_token"`
	ServicePrincipal      []ServicePrincipalDefinitionModel `tfschema:"service_principal"`
	SyncIntervalInSeconds int64                             `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64                             `tfschema:"timeout_in_seconds"`
	Url                   string                            `tfschema:"url"`
}

type ServicePrincipalDefinitionModel struct {
	ClientCertificate          string `tfschema:"client_certificate"`
	ClientCertificatePassword  string `tfschema:"client_certificate_password"`
	ClientCertificateSendChain bool   `tfschema:"client_certificate_send_chain"`
	ClientId                   string `tfschema:"client_id"`
	ClientSecret               string `tfschema:"client_secret"`
	TenantId                   string `tfschema:"tenant_id"`
}

type BucketDefinitionModel struct {
	AccessKey             string `tfschema:"access_key"`
	SecretKey             string `tfschema:"secret_key"`
	BucketName            string `tfschema:"bucket_name"`
	TlsEnabled            bool   `tfschema:"tls_enabled"`
	LocalAuthRef          string `tfschema:"local_auth_ref"`
	SyncIntervalInSeconds int64  `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64  `tfschema:"timeout_in_seconds"`
	Url                   string `tfschema:"url"`
}

type GitRepositoryDefinitionModel struct {
	HttpsCACert           string `tfschema:"https_ca_cert"`
	HttpsUser             string `tfschema:"https_user"`
	HttpsKey              string `tfschema:"https_key"`
	LocalAuthRef          string `tfschema:"local_auth_ref"`
	ReferenceType         string `tfschema:"reference_type"`
	ReferenceValue        string `tfschema:"reference_value"`
	SshKnownHosts         string `tfschema:"ssh_known_hosts"`
	SshPrivateKey         string `tfschema:"ssh_private_key"`
	SyncIntervalInSeconds int64  `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64  `tfschema:"timeout_in_seconds"`
	Url                   string `tfschema:"url"`
}

type KustomizationDefinitionModel struct {
	Name                   string   `tfschema:"name"`
	Path                   string   `tfschema:"path"`
	TimeoutInSeconds       int64    `tfschema:"timeout_in_seconds"`
	SyncIntervalInSeconds  int64    `tfschema:"sync_interval_in_seconds"`
	RetryIntervalInSeconds int64    `tfschema:"retry_interval_in_seconds"`
	Force                  bool     `tfschema:"re_creating_enabled"`
	Prune                  bool     `tfschema:"garbage_collection_enabled"`
	DependsOn              []string `tfschema:"depends_on"`
}

type ArcKubernetesFluxConfigurationResource struct{}

var (
	_ sdk.ResourceWithUpdate = ArcKubernetesFluxConfigurationResource{}
)

func (r ArcKubernetesFluxConfigurationResource) ResourceType() string {
	return "azurerm_arc_kubernetes_flux_configuration"
}

func (r ArcKubernetesFluxConfigurationResource) ModelObject() interface{} {
	return &ArcKubernetesFluxConfigurationModel{}
}

func (r ArcKubernetesFluxConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fluxconfiguration.ValidateFluxConfigurationID
}

func (r ArcKubernetesFluxConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,28}[a-z\d])?$`),
				"`name` must be between 1 and 30 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
			),
		},

		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: arckubernetes.ValidateConnectedClusterID,
		},

		"kustomizations": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,28}[a-z\d])?$`),
							"`name` of `kustomizations` must be between 1 and 30 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						),
					},

					"path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"sync_interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"retry_interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"re_creating_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"garbage_collection_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"depends_on": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
			Set: func(v interface{}) int {
				var buf bytes.Buffer
				m := v.(map[string]interface{})
				buf.WriteString(m["name"].(string))
				return pluginsdk.HashString(buf.String())
			},
		},

		"azure_blob": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"azure_blob", "bucket", "git_repository"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"account_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"azure_blob.0.account_key", "azure_blob.0.local_auth_ref", "azure_blob.0.sas_token", "azure_blob.0.service_principal"},
					},

					"local_auth_ref": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"azure_blob.0.account_key", "azure_blob.0.local_auth_ref", "azure_blob.0.sas_token", "azure_blob.0.service_principal"},
					},

					"sas_token": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"azure_blob.0.account_key", "azure_blob.0.local_auth_ref", "azure_blob.0.sas_token", "azure_blob.0.service_principal"},
					},

					"service_principal": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"azure_blob.0.account_key", "azure_blob.0.local_auth_ref", "azure_blob.0.sas_token", "azure_blob.0.service_principal"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"tenant_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"client_certificate": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
									ExactlyOneOf: []string{"azure_blob.0.service_principal.0.client_certificate", "azure_blob.0.service_principal.0.client_secret"},
								},

								"client_certificate_password": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									RequiredWith: []string{"azure_blob.0.service_principal.0.client_certificate"},
								},

								"client_certificate_send_chain": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"client_secret": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									ExactlyOneOf: []string{"azure_blob.0.service_principal.0.client_certificate", "azure_blob.0.service_principal.0.client_secret"},
								},
							},
						},
					},

					"sync_interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},
				},
			},
		},

		"bucket": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"azure_blob", "bucket", "git_repository"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"bucket_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,61}[a-z\d])?$`),
							"`bucket_name` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						),
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"access_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						RequiredWith: []string{"bucket.0.secret_key"},
						ExactlyOneOf: []string{"bucket.0.access_key", "bucket.0.local_auth_ref"},
					},

					"secret_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsBase64,
						Sensitive:    true,
						RequiredWith: []string{"bucket.0.access_key"},
					},

					"local_auth_ref": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,61}[a-z\d])?$`),
							"`local_auth_ref` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						),
						ExactlyOneOf: []string{"bucket.0.access_key", "bucket.0.local_auth_ref"},
					},

					"tls_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"sync_interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},
				},
			},
		},

		"git_repository": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"azure_blob", "bucket", "git_repository"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.KubernetesGitRepositoryUrl(),
					},

					"reference_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							FluxGitBranch,
							FluxGitCommit,
							FluxGitSemverRange,
							FluxGitReferenceTag,
						}, false),
					},

					"reference_value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"https_ca_cert": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsBase64,
						Sensitive:    true,
						RequiredWith: []string{"git_repository.0.https_user"},
					},

					"https_user": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						RequiredWith:  []string{"git_repository.0.https_key"},
						ConflictsWith: []string{"git_repository.0.local_auth_ref", "git_repository.0.ssh_private_key", "git_repository.0.ssh_known_hosts"},
					},

					"https_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsBase64,
						Sensitive:    true,
						RequiredWith: []string{"git_repository.0.https_user"},
					},

					"local_auth_ref": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,61}[a-z\d])?$`),
							"`local_auth_ref` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						),
						ConflictsWith: []string{"git_repository.0.https_user", "git_repository.0.ssh_private_key", "git_repository.0.ssh_known_hosts"},
					},

					"ssh_private_key": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsBase64,
						Sensitive:     true,
						ConflictsWith: []string{"git_repository.0.https_user", "git_repository.0.local_auth_ref"},
					},

					"ssh_known_hosts": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsBase64,
						ConflictsWith: []string{"git_repository.0.https_user", "git_repository.0.local_auth_ref"},
					},

					"sync_interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  600,
					},

					"timeout_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  600,
					},
				},
			},
		},

		"namespace": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,61}[a-z\d])?$`),
				"`name` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
			),
			Default: "default",
		},

		"scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(fluxconfiguration.ScopeTypeNamespace),
				string(fluxconfiguration.ScopeTypeCluster),
			}, false),
			Default: string(fluxconfiguration.ScopeTypeNamespace),
		},

		"continuous_reconciliation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}

func (r ArcKubernetesFluxConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ArcKubernetesFluxConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ArcKubernetesFluxConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ArcKubernetes.FluxConfigurationClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			clusterID, err := arckubernetes.ParseConnectedClusterID(model.ClusterID)
			if err != nil {
				return err
			}

			// defined as strings because they're not enums in the swagger https://github.com/Azure/azure-rest-api-specs/pull/23545
			id := fluxconfiguration.NewFluxConfigurationID(subscriptionId, clusterID.ResourceGroupName, "Microsoft.Kubernetes", "connectedClusters", clusterID.ConnectedClusterName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &fluxconfiguration.FluxConfiguration{
				Properties: &fluxconfiguration.FluxConfigurationProperties{
					Kustomizations: expandKustomizationDefinitionModel(model.Kustomizations),
					Scope:          &model.Scope,
					Suspend:        utils.Bool(!model.ContinuousReconciliationEnabled),
				},
			}

			if _, exists := metadata.ResourceData.GetOk("git_repository"); exists {
				gitRepositoryValue, configurationProtectedSettings, err := expandGitRepositoryDefinitionModel(model.GitRepository)
				if err != nil {
					return err
				}

				properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeGitRepository)
				properties.Properties.GitRepository = gitRepositoryValue
				properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
			} else if _, exists = metadata.ResourceData.GetOk("bucket"); exists {
				properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeBucket)
				properties.Properties.Bucket, properties.Properties.ConfigurationProtectedSettings = expandBucketDefinitionModel(model.Bucket)
			} else if _, exists = metadata.ResourceData.GetOk("azure_blob"); exists {
				properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeAzureBlob)
				properties.Properties.AzureBlob = expandArcAzureBlobDefinitionModel(model.AzureBlob)
			}

			if model.Namespace != "" {
				properties.Properties.Namespace = &model.Namespace
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ArcKubernetesFluxConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.FluxConfigurationClient

			id, err := fluxconfiguration.ParseFluxConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ArcKubernetesFluxConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			properties.Properties.ConfigurationProtectedSettings = nil
			if metadata.ResourceData.HasChange("azure_blob") {
				properties.Properties.AzureBlob = expandArcAzureBlobDefinitionModel(model.AzureBlob)
				if properties.Properties.AzureBlob != nil {
					properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeAzureBlob)
				}
			}

			if metadata.ResourceData.HasChange("bucket") {
				bucketValue, configurationProtectedSettings := expandBucketDefinitionModel(model.Bucket)
				properties.Properties.Bucket = bucketValue
				if properties.Properties.Bucket != nil {
					properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeBucket)
					properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
				}
			}

			if metadata.ResourceData.HasChange("git_repository") {
				gitRepositoryValue, configurationProtectedSettings, err := expandGitRepositoryDefinitionModel(model.GitRepository)
				if err != nil {
					return err
				}

				properties.Properties.GitRepository = gitRepositoryValue
				if properties.Properties.GitRepository != nil {
					properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeGitRepository)
					properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
				}
			}

			if metadata.ResourceData.HasChange("kustomizations") {
				properties.Properties.Kustomizations = expandKustomizationDefinitionModel(model.Kustomizations)
			}

			if metadata.ResourceData.HasChange("continuous_reconciliation_enabled") {
				properties.Properties.Suspend = utils.Bool(!model.ContinuousReconciliationEnabled)
			}

			properties.SystemData = nil

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ArcKubernetesFluxConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.FluxConfigurationClient

			id, err := fluxconfiguration.ParseFluxConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var configModel ArcKubernetesFluxConfigurationModel
			if err := metadata.Decode(&configModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ArcKubernetesFluxConfigurationModel{
				Name:      id.FluxConfigurationName,
				ClusterID: arckubernetes.NewConnectedClusterID(metadata.Client.Account.SubscriptionId, id.ResourceGroupName, id.ClusterName).ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					state.AzureBlob = flattenArcAzureBlobDefinitionModel(properties.AzureBlob, configModel.AzureBlob)
					state.Bucket = flattenBucketDefinitionModel(properties.Bucket, configModel.Bucket)
					gitRepositoryValue, err := flattenGitRepositoryDefinitionModel(properties.GitRepository, configModel.GitRepository)
					if err != nil {
						return err
					}

					state.GitRepository = gitRepositoryValue
					state.Kustomizations = flattenKustomizationDefinitionModel(properties.Kustomizations)
					state.Namespace = pointer.From(properties.Namespace)
					state.Scope = pointer.From(properties.Scope)
					state.ContinuousReconciliationEnabled = !pointer.From(properties.Suspend)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ArcKubernetesFluxConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.FluxConfigurationClient

			id, err := fluxconfiguration.ParseFluxConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, fluxconfiguration.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandArcAzureBlobDefinitionModel(inputList []AzureBlobDefinitionModel) *fluxconfiguration.AzureBlobDefinition {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fluxconfiguration.AzureBlobDefinition{
		ServicePrincipal:      expandServicePrincipalDefinitionModel(input.ServicePrincipal),
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
	}
	if input.AccountKey != "" {
		output.AccountKey = &input.AccountKey
	}

	if input.ContainerName != "" {
		output.ContainerName = &input.ContainerName
	}

	if input.LocalAuthRef != "" {
		output.LocalAuthRef = &input.LocalAuthRef
	}

	if input.SasToken != "" {
		output.SasToken = &input.SasToken
	}

	if input.Url != "" {
		output.Url = &input.Url
	}

	return &output
}

func expandKustomizationDefinitionModel(inputList []KustomizationDefinitionModel) *map[string]fluxconfiguration.KustomizationDefinition {
	if len(inputList) == 0 {
		return nil
	}

	outputList := make(map[string]fluxconfiguration.KustomizationDefinition)
	for _, v := range inputList {
		input := v
		output := fluxconfiguration.KustomizationDefinition{
			DependsOn:              &input.DependsOn,
			Force:                  &input.Force,
			Name:                   &input.Name,
			Prune:                  &input.Prune,
			RetryIntervalInSeconds: &input.RetryIntervalInSeconds,
			SyncIntervalInSeconds:  &input.SyncIntervalInSeconds,
			TimeoutInSeconds:       &input.TimeoutInSeconds,
		}

		if input.Path != "" {
			output.Path = utils.String(input.Path)
		}

		outputList[input.Name] = output
	}

	return &outputList
}

func expandServicePrincipalDefinitionModel(inputList []ServicePrincipalDefinitionModel) *fluxconfiguration.ServicePrincipalDefinition {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fluxconfiguration.ServicePrincipalDefinition{
		ClientCertificateSendChain: &input.ClientCertificateSendChain,
	}

	if input.ClientCertificate != "" {
		output.ClientCertificate = &input.ClientCertificate
	}

	if input.ClientCertificatePassword != "" {
		output.ClientCertificatePassword = &input.ClientCertificatePassword
	}

	if input.ClientId != "" {
		output.ClientId = &input.ClientId
	}

	if input.ClientSecret != "" {
		output.ClientSecret = &input.ClientSecret
	}

	if input.TenantId != "" {
		output.TenantId = &input.TenantId
	}

	return &output
}

func expandBucketDefinitionModel(inputList []BucketDefinitionModel) (*fluxconfiguration.BucketDefinition, *map[string]string) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.BucketDefinition{
		Insecure:              utils.Bool(!input.TlsEnabled),
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
	}

	if input.AccessKey != "" {
		output.AccessKey = &input.AccessKey
	}

	if input.BucketName != "" {
		output.BucketName = &input.BucketName
	}

	if input.LocalAuthRef != "" {
		output.LocalAuthRef = &input.LocalAuthRef
	}

	if input.Url != "" {
		output.Url = &input.Url
	}

	var configSettings = make(map[string]string)
	if input.SecretKey != "" {
		configSettings["bucketSecretKey"] = input.SecretKey
	}

	var outputConfigSettings *map[string]string = nil
	if len(configSettings) > 0 {
		outputConfigSettings = &configSettings
	}

	return &output, outputConfigSettings
}

func expandGitRepositoryDefinitionModel(inputList []GitRepositoryDefinitionModel) (*fluxconfiguration.GitRepositoryDefinition, *map[string]string, error) {
	if len(inputList) == 0 {
		return nil, nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.GitRepositoryDefinition{
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
	}

	if input.HttpsCACert != "" {
		output.HTTPSCACert = &input.HttpsCACert
	}

	if input.HttpsUser != "" {
		output.HTTPSUser = &input.HttpsUser
	}

	if input.LocalAuthRef != "" {
		output.LocalAuthRef = &input.LocalAuthRef
	}

	repositoryRefValue, err := expandRepositoryRefDefinitionModel(input.ReferenceType, input.ReferenceValue)
	if err != nil {
		return nil, nil, err
	}

	output.RepositoryRef = repositoryRefValue

	if input.SshKnownHosts != "" {
		output.SshKnownHosts = &input.SshKnownHosts
	}

	if input.Url != "" {
		output.Url = &input.Url
	}

	var configSettings = make(map[string]string)
	if input.HttpsKey != "" {
		configSettings["httpsKey"] = input.HttpsKey
	}

	if input.SshPrivateKey != "" {
		configSettings["sshPrivateKey"] = input.SshPrivateKey
	}

	return &output, &configSettings, nil
}

func expandRepositoryRefDefinitionModel(referenceType string, referenceValue string) (*fluxconfiguration.RepositoryRefDefinition, error) {
	output := fluxconfiguration.RepositoryRefDefinition{}

	switch referenceType {
	case FluxGitBranch:
		output.Branch = &referenceValue
	case FluxGitCommit:
		output.Commit = &referenceValue
	case FluxGitSemverRange:
		output.Semver = &referenceValue
	case FluxGitReferenceTag:
		output.Tag = &referenceValue
	default:
		return &output, fmt.Errorf("reference type %s not defined", referenceType)
	}

	return &output, nil
}

func flattenArcAzureBlobDefinitionModel(input *fluxconfiguration.AzureBlobDefinition, azureBlob []AzureBlobDefinitionModel) []AzureBlobDefinitionModel {
	var outputList []AzureBlobDefinitionModel
	if input == nil {
		return outputList
	}

	output := AzureBlobDefinitionModel{
		ContainerName:         pointer.From(input.ContainerName),
		LocalAuthRef:          pointer.From(input.LocalAuthRef),
		SyncIntervalInSeconds: pointer.From(input.SyncIntervalInSeconds),
		TimeoutInSeconds:      pointer.From(input.TimeoutInSeconds),
		Url:                   pointer.From(input.Url),
	}

	var servicePrincipal []ServicePrincipalDefinitionModel
	if len(azureBlob) > 0 {
		output.AccountKey = azureBlob[0].AccountKey
		output.SasToken = azureBlob[0].SasToken
		servicePrincipal = azureBlob[0].ServicePrincipal
	}

	output.ServicePrincipal = flattenServicePrincipalDefinitionModel(input.ServicePrincipal, servicePrincipal)

	return append(outputList, output)
}

func flattenKustomizationDefinitionModel(inputList *map[string]fluxconfiguration.KustomizationDefinition) []KustomizationDefinitionModel {
	var outputList []KustomizationDefinitionModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := KustomizationDefinitionModel{
			DependsOn:              pointer.From(input.DependsOn),
			Force:                  pointer.From(input.Force),
			Name:                   pointer.From(input.Name),
			Path:                   pointer.From(input.Path),
			Prune:                  pointer.From(input.Prune),
			RetryIntervalInSeconds: pointer.From(input.RetryIntervalInSeconds),
			SyncIntervalInSeconds:  pointer.From(input.SyncIntervalInSeconds),
			TimeoutInSeconds:       pointer.From(input.TimeoutInSeconds),
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenServicePrincipalDefinitionModel(input *fluxconfiguration.ServicePrincipalDefinition, servicePrincipal []ServicePrincipalDefinitionModel) []ServicePrincipalDefinitionModel {
	var outputList []ServicePrincipalDefinitionModel
	if input == nil {
		return outputList
	}
	output := ServicePrincipalDefinitionModel{
		ClientCertificateSendChain: pointer.From(input.ClientCertificateSendChain),
		ClientId:                   pointer.From(input.ClientId),
		TenantId:                   pointer.From(input.TenantId),
	}

	if len(servicePrincipal) > 0 {
		output.ClientCertificate = servicePrincipal[0].ClientCertificate
		output.ClientCertificatePassword = servicePrincipal[0].ClientCertificatePassword
		output.ClientSecret = servicePrincipal[0].ClientSecret
	}

	return append(outputList, output)
}

func flattenBucketDefinitionModel(input *fluxconfiguration.BucketDefinition, bucket []BucketDefinitionModel) []BucketDefinitionModel {
	var outputList []BucketDefinitionModel
	if input == nil {
		return outputList
	}

	output := BucketDefinitionModel{
		AccessKey:             pointer.From(input.AccessKey),
		BucketName:            pointer.From(input.BucketName),
		TlsEnabled:            !pointer.From(input.Insecure),
		LocalAuthRef:          pointer.From(input.LocalAuthRef),
		SyncIntervalInSeconds: pointer.From(input.SyncIntervalInSeconds),
		TimeoutInSeconds:      pointer.From(input.TimeoutInSeconds),
		Url:                   pointer.From(input.Url),
	}

	if len(bucket) > 0 {
		output.SecretKey = bucket[0].SecretKey
	}

	return append(outputList, output)
}

func flattenGitRepositoryDefinitionModel(input *fluxconfiguration.GitRepositoryDefinition, gitRepository []GitRepositoryDefinitionModel) ([]GitRepositoryDefinitionModel, error) {
	var outputList []GitRepositoryDefinitionModel
	if input == nil {
		return outputList, nil
	}

	output := GitRepositoryDefinitionModel{
		HttpsCACert:           pointer.From(input.HTTPSCACert),
		HttpsUser:             pointer.From(input.HTTPSUser),
		LocalAuthRef:          pointer.From(input.LocalAuthRef),
		SshKnownHosts:         pointer.From(input.SshKnownHosts),
		SyncIntervalInSeconds: pointer.From(input.SyncIntervalInSeconds),
		TimeoutInSeconds:      pointer.From(input.TimeoutInSeconds),
		Url:                   pointer.From(input.Url),
	}

	referenceType, referenceValue, err := flattenRepositoryRefDefinitionModel(input.RepositoryRef)
	if err != nil {
		return nil, err
	}

	output.ReferenceType = referenceType
	output.ReferenceValue = referenceValue

	if len(gitRepository) > 0 {
		output.HttpsKey = gitRepository[0].HttpsKey
		output.SshPrivateKey = gitRepository[0].SshPrivateKey
	}

	return append(outputList, output), nil
}

func flattenRepositoryRefDefinitionModel(input *fluxconfiguration.RepositoryRefDefinition) (string, string, error) {
	if input == nil {
		return "", "", nil
	}

	var referenceType string
	var referenceValue string

	switch {
	case input.Branch != nil:
		referenceType = FluxGitBranch
		referenceValue = *input.Branch
	case input.Commit != nil:
		referenceType = FluxGitCommit
		referenceValue = *input.Commit
	case input.Semver != nil:
		referenceType = FluxGitSemverRange
		referenceValue = *input.Semver
	case input.Tag != nil:
		referenceType = FluxGitReferenceTag
		referenceValue = *input.Tag
	default:
		return "", "", fmt.Errorf("failed to retrieve git reference")
	}

	return referenceType, referenceValue, nil
}
