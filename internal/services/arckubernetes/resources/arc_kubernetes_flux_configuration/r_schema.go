// Copyright IBM Corp. 2014, 2026
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
			Description: "The name which should be used for this Arc Kubernetes Flux Configuration.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,28}[a-z\d])?$`),
				"`name` must be between 1 and 30 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
			),
		},

		"cluster_id": {
			Description:  "The Cluster ID.",
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
						Description: "The name of the kustomization.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,28}[a-z\d])?$`),
							"`name` of `kustomizations` must be between 1 and 30 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						),
					},

					"path": {
						Description:  "The path in the source reference to reconcile on the cluster.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"timeout_in_seconds": {
						Description:  "The maximum time to attempt to reconcile the kustomization on the cluster.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"sync_interval_in_seconds": {
						Description:  "The interval at which to re-reconcile the kustomization on the cluster.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"retry_interval_in_seconds": {
						Description:  "The interval at which to re-reconcile the kustomization on the cluster in the event of failure on reconciliation.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"recreating_enabled": {
						Description: "Whether re-creating Kubernetes resources on the cluster is enabled when patching fails due to an immutable field change.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"garbage_collection_enabled": {
						Description: "Whether garbage collections of Kubernetes objects created by this kustomization is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"depends_on": {
						Description: "Specifies other kustomizations that this kustomization depends on. This kustomization will not reconcile until all dependencies have completed their reconciliation.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"namespace": {
			Description: "The namespace to which this configuration is installed to.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
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
						Description:  "The Azure Blob container ID.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: storageValidate.StorageContainerDataPlaneID,
					},

					"account_key": {
						Description:  "The account key (shared key) to access the Storage Account.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
					},

					"local_auth_reference": {
						Description:  "The name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.LocalAuthReference,
						ExactlyOneOf: []string{"blob_storage.0.account_key", "blob_storage.0.local_auth_reference", "blob_storage.0.sas_token", "blob_storage.0.service_principal"},
					},

					"sas_token": {
						Description:  "The shared access token to access the Storage Container.",
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
									Description:  "The client ID for authenticating a Service Principal.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"tenant_id": {
									Description:  "The tenant ID for authenticating a Service Principal.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"client_certificate_base64": {
									Description:  "Base64-encoded certificate used to authenticate a Service Principal.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									ExactlyOneOf: []string{"blob_storage.0.service_principal.0.client_certificate_base64", "blob_storage.0.service_principal.0.client_secret"},
								},

								"client_certificate_password": {
									Description:  "The password for the certificate used to authenticate a Service Principal.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
									RequiredWith: []string{"blob_storage.0.service_principal.0.client_certificate_base64"},
								},

								"client_certificate_send_chain": {
									Description: "Whether to include x5c header in client claims when acquiring a token to enable subject name / issuer based authentication for the client certificate.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},

								"client_secret": {
									Description:  "The client secret for authenticating a Service Principal.",
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
						Description:  "The interval at which to re-reconcile the cluster Azure Blob source with the remote.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"timeout_in_seconds": {
						Description:  "The maximum time to attempt to reconcile the cluster Azure Blob source with the remote.",
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
						Description: "The bucket name to sync from the URL endpoint for the flux configuration.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,61}[a-z\d])?$`),
							"`bucket_name` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						),
					},

					"url": {
						Description:  "The URL to sync for the flux configuration S3 bucket. It must start with `http://` or `https://`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"access_key": {
						Description:  "The plaintext access key used to securely access the S3 bucket.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						RequiredWith: []string{"bucket.0.secret_key_base64"},
						ExactlyOneOf: []string{"bucket.0.access_key", "bucket.0.local_auth_reference"},
					},

					"secret_key_base64": {
						Description:  "The Base64-encoded secret key used to authenticate with the bucket source.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsBase64,
						Sensitive:    true,
						RequiredWith: []string{"bucket.0.access_key"},
					},

					"local_auth_reference": {
						Description:  "The name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.LocalAuthReference,
						ExactlyOneOf: []string{"bucket.0.access_key", "bucket.0.local_auth_reference"},
					},

					"tls_enabled": {
						Description: "Whether to communicate with a bucket using TLS is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"sync_interval_in_seconds": {
						Description:  "The interval at which to re-reconcile the cluster git repository source with the remote.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"timeout_in_seconds": {
						Description:  "The maximum time to attempt to reconcile the cluster git repository source with the remote.",
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
						Description:  "The URL to sync for the flux configuration git repository. It must start with `http://`, `https://`, `git@` or `ssh://`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.KubernetesGitRepositoryUrl(),
					},

					"reference_type": {
						Description: "The source reference type for the GitRepository object.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							FluxGitBranch,
							FluxGitCommit,
							FluxGitSemverRange,
							FluxGitReferenceTag,
						}, false),
					},

					"reference_value": {
						Description:  "The source reference value for the GitRepository object.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"https_ca_cert_base64": {
						Description:  "The Base64-encoded HTTPS certificate authority contents used to access git private git repositories over HTTPS.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsBase64,
						Sensitive:    true,
						RequiredWith: []string{"git_repository.0.https_user"},
					},

					"https_user": {
						Description:   "The plaintext HTTPS username used to access private git repositories over HTTPS.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						RequiredWith:  []string{"git_repository.0.https_key_base64"},
						ConflictsWith: []string{"git_repository.0.local_auth_reference", "git_repository.0.ssh_private_key_base64", "git_repository.0.ssh_known_hosts_base64"},
					},

					"https_key_base64": {
						Description:  "The Base64-encoded HTTPS personal access token or password that will be used to access the repository.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsBase64,
						Sensitive:    true,
						RequiredWith: []string{"git_repository.0.https_user"},
					},

					"local_auth_reference": {
						Description:   "The name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validate.LocalAuthReference,
						ConflictsWith: []string{"git_repository.0.https_user", "git_repository.0.ssh_private_key_base64", "git_repository.0.ssh_known_hosts_base64"},
					},

					"ssh_private_key_base64": {
						Description:   "The Base64-encoded SSH private key in PEM format.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsBase64,
						Sensitive:     true,
						ConflictsWith: []string{"git_repository.0.https_user", "git_repository.0.local_auth_reference"},
					},

					"ssh_known_hosts_base64": {
						Description:   "The Base64-encoded known_hosts value containing public SSH keys required to access private git repositories over SSH.",
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsBase64,
						ConflictsWith: []string{"git_repository.0.https_user", "git_repository.0.local_auth_reference"},
					},

					"sync_interval_in_seconds": {
						Description: "The interval at which to re-reconcile the cluster git repository source with the remote.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Default:     600,
					},

					"timeout_in_seconds": {
						Description: "The maximum time to attempt to reconcile the cluster git repository source with the remote.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Default:     600,
					},
				},
			},
		},

		"scope": {
			Description:  "The scope at which the operator will be installed.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(fluxconfiguration.PossibleValuesForScopeType(), false),
			Default:      string(fluxconfiguration.ScopeTypeNamespace),
		},

		"continuous_reconciliation_enabled": {
			Description: "Whether the configuration will keep its reconciliation of its kustomizations and sources with the repository.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
