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
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2023-05-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/containers"
)

const (
	FluxGitBranch       string = "branch"
	FluxGitCommit       string = "commit"
	FluxGitReferenceTag string = "tag"
	FluxGitSemverRange  string = "semver"
)

const (
	SubstituteFromKindConfigMap string = "ConfigMap"
	SubstituteFromKindSecret    string = "Secret"
)

type KubernetesFluxConfigurationModel struct {
	Name                            string                         `tfschema:"name"`
	ClusterID                       string                         `tfschema:"cluster_id"`
	BlobStorage                     []AzureBlobDefinitionModel     `tfschema:"blob_storage"`
	Bucket                          []BucketDefinitionModel        `tfschema:"bucket"`
	GitRepository                   []GitRepositoryDefinitionModel `tfschema:"git_repository"`
	Kustomizations                  []KustomizationDefinitionModel `tfschema:"kustomizations"`
	Namespace                       string                         `tfschema:"namespace"`
	Scope                           string                         `tfschema:"scope"`
	ContinuousReconciliationEnabled bool                           `tfschema:"continuous_reconciliation_enabled"`
}

type AzureBlobDefinitionModel struct {
	AccountKey            string                            `tfschema:"account_key"`
	ContainerID           string                            `tfschema:"container_id"`
	LocalAuthRef          string                            `tfschema:"local_auth_reference"`
	ManagedIdentity       []ManagedIdentityDefinitionModel  `tfschema:"managed_identity"`
	SasToken              string                            `tfschema:"sas_token"`
	ServicePrincipal      []ServicePrincipalDefinitionModel `tfschema:"service_principal"`
	SyncIntervalInSeconds int64                             `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64                             `tfschema:"timeout_in_seconds"`
}

type ServicePrincipalDefinitionModel struct {
	ClientCertificate          string `tfschema:"client_certificate_base64"`
	ClientCertificatePassword  string `tfschema:"client_certificate_password"`
	ClientCertificateSendChain bool   `tfschema:"client_certificate_send_chain"`
	ClientId                   string `tfschema:"client_id"`
	ClientSecret               string `tfschema:"client_secret"`
	TenantId                   string `tfschema:"tenant_id"`
}

type BucketDefinitionModel struct {
	AccessKey             string `tfschema:"access_key"`
	SecretKey             string `tfschema:"secret_key_base64"`
	BucketName            string `tfschema:"bucket_name"`
	TlsEnabled            bool   `tfschema:"tls_enabled"`
	LocalAuthRef          string `tfschema:"local_auth_reference"`
	SyncIntervalInSeconds int64  `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64  `tfschema:"timeout_in_seconds"`
	Url                   string `tfschema:"url"`
}

type GitRepositoryDefinitionModel struct {
	HttpsCACert           string `tfschema:"https_ca_cert_base64"`
	HttpsUser             string `tfschema:"https_user"`
	HttpsKey              string `tfschema:"https_key_base64"`
	LocalAuthRef          string `tfschema:"local_auth_reference"`
	ReferenceType         string `tfschema:"reference_type"`
	ReferenceValue        string `tfschema:"reference_value"`
	SshKnownHosts         string `tfschema:"ssh_known_hosts_base64"`
	SshPrivateKey         string `tfschema:"ssh_private_key_base64"`
	SyncIntervalInSeconds int64  `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64  `tfschema:"timeout_in_seconds"`
	Url                   string `tfschema:"url"`
}

type KustomizationDefinitionModel struct {
	Name                   string                     `tfschema:"name"`
	Path                   string                     `tfschema:"path"`
	TimeoutInSeconds       int64                      `tfschema:"timeout_in_seconds"`
	SyncIntervalInSeconds  int64                      `tfschema:"sync_interval_in_seconds"`
	RetryIntervalInSeconds int64                      `tfschema:"retry_interval_in_seconds"`
	Force                  bool                       `tfschema:"recreating_enabled"`
	Prune                  bool                       `tfschema:"garbage_collection_enabled"`
	DependsOn              []string                   `tfschema:"depends_on"`
	PostBuild              []PostBuildDefinitionModel `tfschema:"post_build"`
	Wait                   bool                       `tfschema:"wait"`
}

type PostBuildDefinitionModel struct {
	Substitute     map[string]string               `tfschema:"substitute"`
	SubstituteFrom []SubstituteFromDefinitionModel `tfschema:"substitute_from"`
}

type SubstituteFromDefinitionModel struct {
	Kind     string `tfschema:"kind"`
	Name     string `tfschema:"name"`
	Optional bool   `tfschema:"optional"`
}

type ManagedIdentityDefinitionModel struct {
	ClientId string `tfschema:"client_id"`
}

type KubernetesFluxConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = KubernetesFluxConfigurationResource{}

func (r KubernetesFluxConfigurationResource) ResourceType() string {
	return "azurerm_kubernetes_flux_configuration"
}

func (r KubernetesFluxConfigurationResource) ModelObject() interface{} {
	return &KubernetesFluxConfigurationModel{}
}

func (r KubernetesFluxConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(val interface{}, key string) (warns []string, errs []error) {
		idRaw, ok := val.(string)
		if !ok {
			errs = append(errs, fmt.Errorf("expected `id` to be a string but got %+v", val))
			return
		}

		id, err := fluxconfiguration.ParseScopedFluxConfigurationID(idRaw)
		if err != nil {
			errs = append(errs, fmt.Errorf("parsing %q: %+v", idRaw, err))
			return
		}

		// validate the scope is a connected cluster id
		if _, err := commonids.ParseKubernetesClusterID(id.Scope); err != nil {
			errs = append(errs, fmt.Errorf("parsing %q as a Kubernetes Cluster ID: %+v", idRaw, err))
			return
		}

		return
	}
}

func (r KubernetesFluxConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: commonids.ValidateKubernetesClusterID,
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

					"recreating_enabled": {
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

					"post_build": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"substitute": {
									Type:     pluginsdk.TypeMap,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"substitute_from": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"kind": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													SubstituteFromKindConfigMap,
													SubstituteFromKindSecret,
												}, false),
											},
											"name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"optional": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												Default:  false,
											},
										},
									},
								},
							},
						},
					},
					"wait": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"namespace": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,61}[a-z\d])?$`),
				"`name` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
			),
		},

		"blob_storage": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"blob_storage", "bucket", "git_repository"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"container_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: storageValidate.StorageContainerDataPlaneID,
					},

					"account_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.managed_identity", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
					},

					"local_auth_reference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.LocalAuthReference,
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.managed_identity", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
					},

					"managed_identity": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.managed_identity", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"sas_token": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.managed_identity", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
					},

					"service_principal": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.managed_identity", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
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

								"client_certificate_base64": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									ExactlyOneOf: []string{"blob_storage.0.service_principal.0.client_certificate_base64", "blob_storage.0.service_principal.0.client_secret"},
								},

								"client_certificate_password": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									RequiredWith: []string{"blob_storage.0.service_principal.0.client_certificate_base64"},
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
									ExactlyOneOf: []string{"blob_storage.0.service_principal.0.client_certificate_base64", "blob_storage.0.service_principal.0.client_secret"},
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
			ExactlyOneOf: []string{"blob_storage", "bucket", "git_repository"},
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
						RequiredWith: []string{"bucket.0.secret_key_base64"},
						ExactlyOneOf: []string{"bucket.0.access_key", "bucket.0.local_auth_reference"},
					},

					"secret_key_base64": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsBase64,
						Sensitive:    true,
						RequiredWith: []string{"bucket.0.access_key"},
					},

					"local_auth_reference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.LocalAuthReference,
						ExactlyOneOf: []string{"bucket.0.access_key", "bucket.0.local_auth_reference"},
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
			ExactlyOneOf: []string{"blob_storage", "bucket", "git_repository"},
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

					"https_ca_cert_base64": {
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
						RequiredWith:  []string{"git_repository.0.https_key_base64"},
						ConflictsWith: []string{"git_repository.0.local_auth_reference", "git_repository.0.ssh_private_key_base64", "git_repository.0.ssh_known_hosts_base64"},
					},

					"https_key_base64": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsBase64,
						Sensitive:    true,
						RequiredWith: []string{"git_repository.0.https_user"},
					},

					"local_auth_reference": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validate.LocalAuthReference,
						ConflictsWith: []string{"git_repository.0.https_user", "git_repository.0.ssh_private_key_base64", "git_repository.0.ssh_known_hosts_base64"},
					},

					"ssh_private_key_base64": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsBase64,
						Sensitive:     true,
						ConflictsWith: []string{"git_repository.0.https_user", "git_repository.0.local_auth_reference"},
					},

					"ssh_known_hosts_base64": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsBase64,
						ConflictsWith: []string{"git_repository.0.https_user", "git_repository.0.local_auth_reference"},
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

func (r KubernetesFluxConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r KubernetesFluxConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model KubernetesFluxConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err := validateKubernetesFluxConfigurationModel(&model); err != nil {
				return err
			}

			client := metadata.Client.Containers.KubernetesFluxConfigurationClient
			clusterID, err := commonids.ParseKubernetesClusterID(model.ClusterID)
			if err != nil {
				return err
			}

			// defined as strings because they're not enums in the swagger https://github.com/Azure/azure-rest-api-specs/pull/23545
			id := fluxconfiguration.NewScopedFluxConfigurationID(clusterID.ID(), model.Name)
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
					Scope:          pointer.To(fluxconfiguration.ScopeType(model.Scope)),
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
			} else if _, exists = metadata.ResourceData.GetOk("blob_storage"); exists {
				properties.Properties.SourceKind = pointer.To(fluxconfiguration.SourceKindTypeAzureBlob)
				azureBlob, err := expandAzureBlobDefinitionModel(model.BlobStorage, metadata.Client.Storage.StorageDomainSuffix)
				if err != nil {
					return fmt.Errorf("expanding `blob_storage`: %+v", err)
				}

				properties.Properties.AzureBlob = azureBlob
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

func (r KubernetesFluxConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesFluxConfigurationClient

			id, err := fluxconfiguration.ParseScopedFluxConfigurationIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KubernetesFluxConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err := validateKubernetesFluxConfigurationModel(&model); err != nil {
				return err
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
			if metadata.ResourceData.HasChange("blob_storage") {
				azureBlob, err := expandAzureBlobDefinitionModel(model.BlobStorage, metadata.Client.Storage.StorageDomainSuffix)
				if err != nil {
					return fmt.Errorf("expanding `blob_storage`: %+v", err)
				}

				properties.Properties.AzureBlob = azureBlob
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

			if properties.Properties.ConfigurationProtectedSettings == nil {
				if err := setConfigurationProtectedSettings(metadata, model, properties); err != nil {
					return err
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesFluxConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesFluxConfigurationClient

			id, err := fluxconfiguration.ParseScopedFluxConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var configModel KubernetesFluxConfigurationModel
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

			clusterId, err := commonids.ParseKubernetesClusterID(id.Scope)
			if err != nil {
				return err
			}
			state := KubernetesFluxConfigurationModel{
				Name:      id.FluxConfigurationName,
				ClusterID: clusterId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					blobStorage, err := flattenAzureBlobDefinitionModel(properties.AzureBlob, configModel.BlobStorage, metadata.Client.Storage.StorageDomainSuffix)
					if err != nil {
						return fmt.Errorf("flattening `blob_storage`: %+v", err)
					}

					state.BlobStorage = blobStorage
					state.Bucket = flattenBucketDefinitionModel(properties.Bucket, configModel.Bucket)
					gitRepositoryValue, err := flattenGitRepositoryDefinitionModel(properties.GitRepository, configModel.GitRepository)
					if err != nil {
						return err
					}

					state.GitRepository = gitRepositoryValue
					state.Kustomizations = flattenKustomizationDefinitionModel(properties.Kustomizations)
					state.Namespace = pointer.From(properties.Namespace)
					state.Scope = string(pointer.From(properties.Scope))
					state.ContinuousReconciliationEnabled = !pointer.From(properties.Suspend)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r KubernetesFluxConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesFluxConfigurationClient

			id, err := fluxconfiguration.ParseScopedFluxConfigurationID(metadata.ResourceData.Id())
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

func expandAzureBlobDefinitionModel(inputList []AzureBlobDefinitionModel, storageDomainSuffix string) (*fluxconfiguration.AzureBlobDefinition, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.AzureBlobDefinition{
		ManagedIdentity:       expandManagedIdentityDefinitionModel(input.ManagedIdentity),
		ServicePrincipal:      expandServicePrincipalDefinitionModel(input.ServicePrincipal),
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
	}
	if input.AccountKey != "" {
		output.AccountKey = &input.AccountKey
	}

	if input.ContainerID != "" {
		id, err := containers.ParseContainerID(input.ContainerID, storageDomainSuffix)
		if err != nil {
			return nil, err
		}

		output.ContainerName = &id.ContainerName
		output.Url = pointer.To(id.AccountId.ID())
	}

	if input.LocalAuthRef != "" {
		output.LocalAuthRef = &input.LocalAuthRef
	}

	if input.SasToken != "" {
		output.SasToken = &input.SasToken
	}

	return &output, nil
}

func expandManagedIdentityDefinitionModel(inputList []ManagedIdentityDefinitionModel) *fluxconfiguration.ManagedIdentityDefinition {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	output := fluxconfiguration.ManagedIdentityDefinition{}
	if input.ClientId != "" {
		output.ClientId = &input.ClientId
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
			PostBuild:              expandPostBuildDefinitionModel(input.PostBuild),
			Prune:                  &input.Prune,
			RetryIntervalInSeconds: &input.RetryIntervalInSeconds,
			SyncIntervalInSeconds:  &input.SyncIntervalInSeconds,
			TimeoutInSeconds:       &input.TimeoutInSeconds,
			Wait:                   &input.Wait,
		}

		if input.Path != "" {
			output.Path = utils.String(input.Path)
		}

		outputList[input.Name] = output
	}

	return &outputList
}

func expandPostBuildDefinitionModel(inputList []PostBuildDefinitionModel) *fluxconfiguration.PostBuildDefinition {
	if len(inputList) == 0 {
		return nil
	}

	input := inputList[0]

	output := fluxconfiguration.PostBuildDefinition{}

	if len(input.Substitute) > 0 {
		output.Substitute = &input.Substitute
	}

	if len(input.SubstituteFrom) > 0 {
		output.SubstituteFrom = expandSubstituteFromDefinitionModel(input.SubstituteFrom)
	}

	return &output
}

func expandSubstituteFromDefinitionModel(inputList []SubstituteFromDefinitionModel) *[]fluxconfiguration.SubstituteFromDefinition {
	if len(inputList) == 0 {
		return nil
	}

	input := inputList
	output := make([]fluxconfiguration.SubstituteFromDefinition, 0)

	for _, v := range input {
		output = append(output, fluxconfiguration.SubstituteFromDefinition{
			Kind:     &v.Kind,
			Name:     &v.Name,
			Optional: &v.Optional,
		})
	}

	return &output
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

	configSettings := make(map[string]string)
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

	configSettings := make(map[string]string)
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

func flattenAzureBlobDefinitionModel(input *fluxconfiguration.AzureBlobDefinition, azureBlob []AzureBlobDefinitionModel, storageDomainSuffix string) ([]AzureBlobDefinitionModel, error) {
	outputList := make([]AzureBlobDefinitionModel, 0)
	if input == nil {
		return outputList, nil
	}

	accountId, err := accounts.ParseAccountID(pointer.From(input.Url), storageDomainSuffix)
	if err != nil {
		return nil, fmt.Errorf("parsing account %q: %+v", pointer.From(input.Url), err)
	}

	id := containers.NewContainerID(*accountId, pointer.From(input.ContainerName))

	output := AzureBlobDefinitionModel{
		ContainerID:           id.ID(),
		LocalAuthRef:          pointer.From(input.LocalAuthRef),
		ManagedIdentity:       flattenManagedIdentityDefinitionModel(input.ManagedIdentity),
		SyncIntervalInSeconds: pointer.From(input.SyncIntervalInSeconds),
		TimeoutInSeconds:      pointer.From(input.TimeoutInSeconds),
	}

	var servicePrincipal []ServicePrincipalDefinitionModel
	if len(azureBlob) > 0 {
		output.AccountKey = azureBlob[0].AccountKey
		output.SasToken = azureBlob[0].SasToken
		servicePrincipal = azureBlob[0].ServicePrincipal
	}

	output.ServicePrincipal = flattenServicePrincipalDefinitionModel(input.ServicePrincipal, servicePrincipal)

	return append(outputList, output), nil
}

func flattenManagedIdentityDefinitionModel(input *fluxconfiguration.ManagedIdentityDefinition) []ManagedIdentityDefinitionModel {
	outputList := make([]ManagedIdentityDefinitionModel, 0)
	if input == nil {
		return outputList
	}
	output := ManagedIdentityDefinitionModel{
		ClientId: pointer.From(input.ClientId),
	}

	return append(outputList, output)
}

func flattenKustomizationDefinitionModel(inputList *map[string]fluxconfiguration.KustomizationDefinition) []KustomizationDefinitionModel {
	outputList := make([]KustomizationDefinitionModel, 0)
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := KustomizationDefinitionModel{
			DependsOn:              pointer.From(input.DependsOn),
			Force:                  pointer.From(input.Force),
			Name:                   pointer.From(input.Name),
			Path:                   pointer.From(input.Path),
			PostBuild:              flattenPostBuildDefinitionModel(input.PostBuild),
			Prune:                  pointer.From(input.Prune),
			RetryIntervalInSeconds: pointer.From(input.RetryIntervalInSeconds),
			SyncIntervalInSeconds:  pointer.From(input.SyncIntervalInSeconds),
			TimeoutInSeconds:       pointer.From(input.TimeoutInSeconds),
			Wait:                   pointer.From(input.Wait),
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenPostBuildDefinitionModel(input *fluxconfiguration.PostBuildDefinition) []PostBuildDefinitionModel {
	outputList := make([]PostBuildDefinitionModel, 0)

	if input == nil {
		return outputList
	}

	output := PostBuildDefinitionModel{
		Substitute:     pointer.From(input.Substitute),
		SubstituteFrom: flattenSubstituteFromDefinitionModel(input.SubstituteFrom),
	}

	return append(outputList, output)
}

func flattenSubstituteFromDefinitionModel(input *[]fluxconfiguration.SubstituteFromDefinition) []SubstituteFromDefinitionModel {
	outputList := make([]SubstituteFromDefinitionModel, 0)
	if input == nil {
		return outputList
	}

	for _, v := range *input {
		output := SubstituteFromDefinitionModel{
			Kind:     pointer.From(v.Kind),
			Name:     pointer.From(v.Name),
			Optional: pointer.From(v.Optional),
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenServicePrincipalDefinitionModel(input *fluxconfiguration.ServicePrincipalDefinition, servicePrincipal []ServicePrincipalDefinitionModel) []ServicePrincipalDefinitionModel {
	outputList := make([]ServicePrincipalDefinitionModel, 0)
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
	outputList := make([]BucketDefinitionModel, 0)
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
	outputList := make([]GitRepositoryDefinitionModel, 0)
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

func validateKubernetesFluxConfigurationModel(model *KubernetesFluxConfigurationModel) error {
	allKeys := make(map[string]bool)
	for _, k := range model.Kustomizations {
		if _, exists := allKeys[k.Name]; exists {
			return fmt.Errorf("kustomization name `%s` is not unique", k.Name)
		}

		allKeys[k.Name] = true
	}

	return nil
}

func setConfigurationProtectedSettings(metadata sdk.ResourceMetaData, model KubernetesFluxConfigurationModel, properties *fluxconfiguration.FluxConfiguration) error {
	if _, exists := metadata.ResourceData.GetOk("git_repository"); exists {
		_, configurationProtectedSettings, err := expandGitRepositoryDefinitionModel(model.GitRepository)
		if err != nil {
			return err
		}
		properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
	} else if _, exists = metadata.ResourceData.GetOk("bucket"); exists {
		_, properties.Properties.ConfigurationProtectedSettings = expandBucketDefinitionModel(model.Bucket)
	}
	return nil
}
