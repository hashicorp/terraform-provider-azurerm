// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetAppVolumeBucketResource struct{}

var (
	_ sdk.Resource             = NetAppVolumeBucketResource{}
	_ sdk.ResourceWithIdentity = NetAppVolumeBucketResource{}
)

func (r NetAppVolumeBucketResource) Identity() resourceids.ResourceId {
	return &buckets.BucketId{}
}

func (r NetAppVolumeBucketResource) ModelObject() interface{} {
	return &netAppModels.NetAppVolumeBucketModel{}
}

func (r NetAppVolumeBucketResource) ResourceType() string {
	return "azurerm_netapp_volume_bucket"
}

func (r NetAppVolumeBucketResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return buckets.ValidateBucketID
}

func (r NetAppVolumeBucketResource) Arguments() map[string]*pluginsdk.Schema {
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
			ConflictsWith: []string{"server.0.certificate_pem"},
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

		"server": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"fqdn": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"certificate_pem": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						Sensitive:     true,
						ValidateFunc:  validation.StringIsBase64,
						ConflictsWith: []string{"key_vault"},
					},
					"on_certificate_conflict_action": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(buckets.OnCertificateConflictActionFail),
						ValidateFunc: validation.StringInSlice(buckets.PossibleValuesForOnCertificateConflictAction(), false),
					},
				},
			},
		},
	}
}

func (r NetAppVolumeBucketResource) Attributes() map[string]*pluginsdk.Schema {
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

func (r NetAppVolumeBucketResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model netAppModels.NetAppVolumeBucketModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			volumeID, err := volumes.ParseVolumeID(model.VolumeID)
			if err != nil {
				return fmt.Errorf("parsing volume id %s: %+v", model.VolumeID, err)
			}

			id := buckets.NewBucketID(subscriptionId, volumeID.ResourceGroupName, volumeID.NetAppAccountName, volumeID.CapacityPoolName, volumeID.VolumeName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			payload := buckets.Bucket{
				Properties: &buckets.BucketProperties{
					Path:        pointer.To(model.Path),
					Permissions: pointer.To(buckets.BucketPermissions(model.Permissions)),
					FileSystemUser: &buckets.FileSystemUser{
						NfsUser:  expandNetAppBucketNfsUser(model.FileSystemNfsUser),
						CifsUser: expandNetAppBucketCifsUser(model.FileSystemCifsUser),
					},
					AkvDetails: expandNetAppBucketAkvDetails(model.KeyVault),
					Server:     expandNetAppBucketServer(model.Server),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r NetAppVolumeBucketResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			id, err := buckets.ParseBucketID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r NetAppVolumeBucketResource) flatten(metadata sdk.ResourceMetaData, id *buckets.BucketId, bucket *buckets.Bucket) error {
	volumeID := volumes.NewVolumeID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName)

	model := netAppModels.NetAppVolumeBucketModel{
		Name:     id.BucketName,
		VolumeID: volumeID.ID(),
	}

	if bucket != nil && bucket.Properties != nil {
		props := bucket.Properties

		model.Path = pointer.From(props.Path)
		model.Permissions = string(pointer.From(props.Permissions))
		model.Status = string(pointer.From(props.Status))

		if props.FileSystemUser != nil {
			model.FileSystemNfsUser = flattenNetAppBucketNfsUser(props.FileSystemUser.NfsUser)
			model.FileSystemCifsUser = flattenNetAppBucketCifsUser(props.FileSystemUser.CifsUser)
		}
		model.KeyVault = flattenNetAppBucketAkvDetails(props.AkvDetails)
		model.Server = flattenNetAppBucketServer(props.Server)

		if props.Server != nil {
			model.ServerIPAddress = pointer.From(props.Server.IPAddress)
			model.ServerCertificateCommonName = pointer.From(props.Server.CertificateCommonName)
			model.ServerCertificateExpiryDate = pointer.From(props.Server.CertificateExpiryDate)
		}
	}

	// certificate_pem is never returned by the API; preserve from config/state.
	if v, ok := metadata.ResourceData.GetOk("server.0.certificate_pem"); ok {
		if len(model.Server) == 0 {
			model.Server = []netAppModels.NetAppVolumeBucketServer{{}}
		}
		model.Server[0].CertificatePem = v.(string)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&model)
}

func (r NetAppVolumeBucketResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			id, err := buckets.ParseBucketID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state netAppModels.NetAppVolumeBucketModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			patchProps := &buckets.BucketPatchProperties{}

			if metadata.ResourceData.HasChange("permissions") {
				patchProps.Permissions = pointer.To(buckets.BucketPatchPermissions(state.Permissions))
			}

			if metadata.ResourceData.HasChange("file_system_nfs_user") || metadata.ResourceData.HasChange("file_system_cifs_user") {
				patchProps.FileSystemUser = &buckets.FileSystemUser{
					NfsUser:  expandNetAppBucketNfsUser(state.FileSystemNfsUser),
					CifsUser: expandNetAppBucketCifsUser(state.FileSystemCifsUser),
				}
			}

			if metadata.ResourceData.HasChange("key_vault") {
				patchProps.AkvDetails = expandNetAppBucketAkvDetails(state.KeyVault)
			}

			if metadata.ResourceData.HasChange("server") {
				patchProps.Server = expandNetAppBucketServerPatch(state.Server, metadata.ResourceData)
			}

			payload := buckets.BucketPatch{
				Properties: patchProps,
			}

			if err := client.UpdateThenPoll(ctx, pointer.From(id), payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r NetAppVolumeBucketResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			id, err := buckets.ParseBucketID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, pointer.From(id)); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
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
