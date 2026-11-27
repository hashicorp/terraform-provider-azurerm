// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/buckets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// netAppBucketResourceCommonArguments returns the schema arguments shared by the
// `azurerm_netapp_volume_bucket` and `azurerm_netapp_volume_bucket_with_server`
// resources. The `server` block is intentionally not part of this set because it
// is only supported by the `_with_server` resource.
func NetAppBucketResourceCommonArguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The S3-compatible name of the bucket.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.BucketName,
		},

		"volume_id": {
			Description:  "The ARM ID of the parent NetApp Volume the bucket attaches to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: volumes.ValidateVolumeID,
		},

		"file_system_nfs_user": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"file_system_nfs_user", "file_system_cifs_username"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"group_id": {
						Description:  "The POSIX group ID used by the bucket when accessing volume data over NFS.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"user_id": {
						Description:  "The POSIX user ID used by the bucket when accessing volume data over NFS.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"file_system_cifs_username": {
			Description:  "The CIFS username used by the bucket when accessing volume data over SMB.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"file_system_nfs_user", "file_system_cifs_username"},
		},

		"key_vault": {
			Description: "Used to store the generated bucket credentials in Azure Key Vault.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_key_vault_uri": {
						Description:  "The URI of the Azure Key Vault that stores the bucket server certificate.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},
					"certificate_name": {
						Description:  "The name of the certificate object inside the Key Vault.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"credentials_key_vault_uri": {
						Description:  "The URI of the Azure Key Vault used to store the generated bucket access and secret keys. May be the same vault as `certificate_key_vault_uri` but the documentation recommends using two separate vaults.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},
					"credentials_secret_name": {
						Description:  "The name of the secret in `credentials_key_vault_uri` that stores the generated bucket credentials. The Key Vault Secret value is a JSON document with `access_key_id` and `secret_access_key` properties.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"path": {
			Description:  "The volume sub-path mounted inside the bucket.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "/",
			ValidateFunc: validate.BucketPath,
		},

		"permissions": {
			Description:  "The bucket permission level.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(buckets.BucketPermissionsReadOnly),
			ValidateFunc: validation.StringInSlice(buckets.PossibleValuesForBucketPermissions(), false),
		},
	}
}

// netAppBucketResourceCommonAttributes returns the computed attributes shared by
// both bucket resources.
func NetAppBucketResourceCommonAttributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"server_certificate_common_name": {
			Description: "The Common Name (CN) of the bucket server certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"server_certificate_expiry_date": {
			Description: "The expiry date of the bucket server certificate, in RFC3339 format.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"server_ip_address": {
			Description: "The IP address that backs the bucket endpoint.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"status": {
			Description: "The credentials status of the bucket.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

// netAppBucketDataSourceArguments returns the schema arguments shared by both
// bucket data sources.
func NetAppBucketDataSourceArguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the NetApp Volume Bucket.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.BucketName,
		},

		"netapp_volume_id": {
			Description:  "The ARM ID of the parent NetApp Volume.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: volumes.ValidateVolumeID,
		},
	}
}

// netAppBucketDataSourceCommonAttributes returns the computed attributes shared by
// both bucket data sources. The `server` block is intentionally not part of this
// set because it is only exported by the `_with_server` data source.
func NetAppBucketDataSourceCommonAttributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"path": {
			Description: "The volume sub-path mounted inside the bucket.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"permissions": {
			Description: "The bucket permission level (`ReadOnly` or `ReadWrite`).",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"file_system_nfs_user": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"group_id": {
						Description: "The POSIX group ID used by the bucket.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},
					"user_id": {
						Description: "The POSIX user ID used by the bucket.",
						Type:        pluginsdk.TypeInt,
						Computed:    true,
					},
				},
			},
		},

		"file_system_cifs_username": {
			Description: "The CIFS username used by the bucket (only set when the bucket is configured for CIFS).",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"key_vault": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_key_vault_uri": {
						Description: "The URI of the Azure Key Vault that stores the bucket server certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"certificate_name": {
						Description: "The name of the certificate object inside the Key Vault.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"credentials_key_vault_uri": {
						Description: "The URI of the Azure Key Vault used to store the generated bucket access and secret keys.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"credentials_secret_name": {
						Description: "The name of the secret in `credentials_key_vault_uri` that stores the generated bucket credentials.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"server_certificate_common_name": {
			Description: "The Common Name (CN) of the bucket server certificate.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"server_certificate_expiry_date": {
			Description: "The expiry date of the bucket server certificate, in RFC3339 format.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"server_ip_address": {
			Description: "The IP address that backs the bucket endpoint.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"status": {
			Description: "The credentials status of the bucket.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func ExpandNetAppBucketNfsUser(input []models.NetAppVolumeBucketNfsUser) *buckets.NfsUser {
	if len(input) == 0 {
		return nil
	}

	return &buckets.NfsUser{
		GroupId: pointer.To(input[0].GroupID),
		UserId:  pointer.To(input[0].UserID),
	}
}

func ExpandNetAppBucketCifsUser(input string) *buckets.CifsUser {
	if input == "" {
		return nil
	}

	return &buckets.CifsUser{
		Username: pointer.To(input),
	}
}

func FlattenNetAppBucketNfsUser(input *buckets.NfsUser) []models.NetAppVolumeBucketNfsUser {
	if input == nil {
		return []models.NetAppVolumeBucketNfsUser{}
	}

	return []models.NetAppVolumeBucketNfsUser{
		{
			GroupID: pointer.From(input.GroupId),
			UserID:  pointer.From(input.UserId),
		},
	}
}

func FlattenNetAppBucketCifsUser(input *buckets.CifsUser) string {
	if input == nil {
		return ""
	}

	return pointer.From(input.Username)
}

func ExpandNetAppBucketAkvDetails(input []models.NetAppVolumeBucketKeyVault) *buckets.AzureKeyVaultDetails {
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

func FlattenNetAppBucketAkvDetails(input *buckets.AzureKeyVaultDetails) []models.NetAppVolumeBucketKeyVault {
	if input == nil {
		return []models.NetAppVolumeBucketKeyVault{}
	}

	out := models.NetAppVolumeBucketKeyVault{}

	if input.CertificateAkvDetails != nil {
		out.CertificateKeyVaultUri = pointer.From(input.CertificateAkvDetails.CertificateKeyVaultUri)
		out.CertificateName = pointer.From(input.CertificateAkvDetails.CertificateName)
	}

	if input.CredentialsAkvDetails != nil {
		out.CredentialsKeyVaultUri = pointer.From(input.CredentialsAkvDetails.CredentialsKeyVaultUri)
		out.CredentialsSecretName = pointer.From(input.CredentialsAkvDetails.SecretName)
	}

	if out.CertificateKeyVaultUri == "" && out.CertificateName == "" && out.CredentialsKeyVaultUri == "" && out.CredentialsSecretName == "" {
		return []models.NetAppVolumeBucketKeyVault{}
	}

	return []models.NetAppVolumeBucketKeyVault{out}
}
