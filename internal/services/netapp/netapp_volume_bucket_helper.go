// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/volumes"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// netAppBucketResourceCommonArguments returns the schema arguments shared by the
// `azurerm_netapp_volume_bucket` and `azurerm_netapp_volume_bucket_with_server`
// resources. The `server` block is intentionally not part of this set because it
// is only supported by the `_with_server` resource.
func netAppBucketResourceCommonArguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.BucketName,
		},

		"volume_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: volumes.ValidateVolumeID,
		},

		"file_system_nfs_user": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"group_id": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"user_id": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
			ExactlyOneOf: []string{"file_system_nfs_user", "file_system_cifs_user"},
		},

		"file_system_cifs_user": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ExactlyOneOf: []string{"file_system_nfs_user", "file_system_cifs_user"},
		},

		"key_vault": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_key_vault_uri": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},
					"certificate_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"credentials_key_vault_uri": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},
					"credentials_secret_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"path": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "/",
			ValidateFunc: netAppValidate.BucketPath,
		},

		"permissions": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(buckets.BucketPermissionsReadOnly),
			ValidateFunc: validation.StringInSlice(buckets.PossibleValuesForBucketPermissions(), false),
		},
	}
}

// netAppBucketResourceCommonAttributes returns the computed attributes shared by
// both bucket resources.
func netAppBucketResourceCommonAttributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"server_certificate_common_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"server_certificate_expiry_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"server_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

// netAppBucketDataSourceArguments returns the schema arguments shared by both
// bucket data sources.
func netAppBucketDataSourceArguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: netAppValidate.BucketName,
		},

		"netapp_volume_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: volumes.ValidateVolumeID,
		},
	}
}

// netAppBucketDataSourceCommonAttributes returns the computed attributes shared by
// both bucket data sources. The `server` block is intentionally not part of this
// set because it is only exported by the `_with_server` data source.
func netAppBucketDataSourceCommonAttributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"path": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"permissions": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"file_system_nfs_user": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"group_id": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
					"user_id": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"file_system_cifs_user": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"username": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"key_vault": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_key_vault_uri": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"certificate_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"credentials_key_vault_uri": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"credentials_secret_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"server_certificate_common_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"server_certificate_expiry_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"server_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func expandNetAppBucketNfsUser(input []netAppModels.NetAppVolumeBucketNfsUser) *buckets.NfsUser {
	if len(input) == 0 {
		return nil
	}

	return &buckets.NfsUser{
		GroupId: pointer.To(input[0].GroupID),
		UserId:  pointer.To(input[0].UserID),
	}
}

func expandNetAppBucketCifsUser(input []netAppModels.NetAppVolumeBucketCifsUser) *buckets.CifsUser {
	if len(input) == 0 {
		return nil
	}

	return &buckets.CifsUser{
		Username: pointer.To(input[0].Username),
	}
}

func flattenNetAppBucketNfsUser(input *buckets.NfsUser) []netAppModels.NetAppVolumeBucketNfsUser {
	if input == nil {
		return nil
	}

	return []netAppModels.NetAppVolumeBucketNfsUser{
		{
			GroupID: pointer.From(input.GroupId),
			UserID:  pointer.From(input.UserId),
		},
	}
}

func flattenNetAppBucketCifsUser(input *buckets.CifsUser) []netAppModels.NetAppVolumeBucketCifsUser {
	if input == nil {
		return nil
	}

	return []netAppModels.NetAppVolumeBucketCifsUser{
		{
			Username: pointer.From(input.Username),
		},
	}
}

func expandNetAppBucketAkvDetails(input []netAppModels.NetAppVolumeBucketKeyVault) *buckets.AzureKeyVaultDetails {
	if len(input) == 0 {
		return nil
	}

	kv := input[0]
	return &buckets.AzureKeyVaultDetails{
		CertificateAkvDetails: &buckets.CertificateAkvDetails{
			CertificateKeyVaultUri: pointer.To(kv.CertificateKeyVaultUri),
			CertificateName:        pointer.To(kv.CertificateName),
		},
		CredentialsAkvDetails: &buckets.CredentialsAkvDetails{
			CredentialsKeyVaultUri: pointer.To(kv.CredentialsKeyVaultUri),
			SecretName:             pointer.To(kv.CredentialsSecretName),
		},
	}
}

func flattenNetAppBucketAkvDetails(input *buckets.AzureKeyVaultDetails) []netAppModels.NetAppVolumeBucketKeyVault {
	if input == nil {
		return []netAppModels.NetAppVolumeBucketKeyVault{}
	}

	out := netAppModels.NetAppVolumeBucketKeyVault{}

	if input.CertificateAkvDetails != nil {
		out.CertificateKeyVaultUri = pointer.From(input.CertificateAkvDetails.CertificateKeyVaultUri)
		out.CertificateName = pointer.From(input.CertificateAkvDetails.CertificateName)
	}

	if input.CredentialsAkvDetails != nil {
		out.CredentialsKeyVaultUri = pointer.From(input.CredentialsAkvDetails.CredentialsKeyVaultUri)
		out.CredentialsSecretName = pointer.From(input.CredentialsAkvDetails.SecretName)
	}

	if out.CertificateKeyVaultUri == "" && out.CertificateName == "" && out.CredentialsKeyVaultUri == "" && out.CredentialsSecretName == "" {
		return []netAppModels.NetAppVolumeBucketKeyVault{}
	}

	return []netAppModels.NetAppVolumeBucketKeyVault{out}
}

func expandNetAppBucketServer(input []netAppModels.NetAppVolumeBucketServer) *buckets.BucketServerProperties {
	if len(input) == 0 {
		return nil
	}

	srv := input[0]
	out := &buckets.BucketServerProperties{}

	if srv.Fqdn != "" {
		out.Fqdn = pointer.To(srv.Fqdn)
	}
	if srv.CertificatePem != "" {
		out.CertificateObject = pointer.To(srv.CertificatePem)
	}
	if srv.OnCertificateConflictAction != "" {
		out.OnCertificateConflictAction = pointer.To(buckets.OnCertificateConflictAction(srv.OnCertificateConflictAction))
	}

	return out
}

func expandNetAppBucketServerPatch(input []netAppModels.NetAppVolumeBucketServer, rd *pluginsdk.ResourceData) *buckets.BucketServerPatchProperties {
	if len(input) == 0 {
		return nil
	}

	srv := input[0]
	out := &buckets.BucketServerPatchProperties{}

	if rd.HasChange("server.0.fqdn") {
		out.Fqdn = pointer.To(srv.Fqdn)
	}
	if rd.HasChange("server.0.certificate_pem") && srv.CertificatePem != "" {
		out.CertificateObject = pointer.To(srv.CertificatePem)
	}
	if rd.HasChange("server.0.on_certificate_conflict_action") && srv.OnCertificateConflictAction != "" {
		out.OnCertificateConflictAction = pointer.To(buckets.OnCertificateConflictAction(srv.OnCertificateConflictAction))
	}

	return out
}

func flattenNetAppBucketServer(input *buckets.BucketServerProperties) []netAppModels.NetAppVolumeBucketServer {
	if input == nil {
		return []netAppModels.NetAppVolumeBucketServer{}
	}

	out := netAppModels.NetAppVolumeBucketServer{
		Fqdn: pointer.From(input.Fqdn),
	}
	if input.OnCertificateConflictAction != nil {
		out.OnCertificateConflictAction = string(pointer.From(input.OnCertificateConflictAction))
	}

	if out.Fqdn == "" && out.OnCertificateConflictAction == "" {
		return []netAppModels.NetAppVolumeBucketServer{}
	}

	return []netAppModels.NetAppVolumeBucketServer{out}
}
