// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_flux_configuration

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2025-04-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: connectedclusters.ValidateConnectedClusterID,
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
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
					},

					"local_auth_reference": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.LocalAuthReference,
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
					},

					"sas_token": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
					},

					"service_principal": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						MaxItems:     1,
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
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
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(fluxconfiguration.PossibleValuesForScopeType(), false),
			Default:      string(fluxconfiguration.ScopeTypeNamespace),
		},

		"continuous_reconciliation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
