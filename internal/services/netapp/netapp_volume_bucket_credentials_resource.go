// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetAppVolumeBucketCredentialsResource struct{}

var _ sdk.Resource = NetAppVolumeBucketCredentialsResource{}

func (r NetAppVolumeBucketCredentialsResource) ModelObject() interface{} {
	return &netAppModels.NetAppVolumeBucketCredentialsModel{}
}

func (r NetAppVolumeBucketCredentialsResource) ResourceType() string {
	return "azurerm_netapp_volume_bucket_credentials"
}

func (r NetAppVolumeBucketCredentialsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return buckets.ValidateBucketID
}

func (r NetAppVolumeBucketCredentialsResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"bucket_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: buckets.ValidateBucketID,
		},

		"key_pair_expiry_days": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"store_in_key_vault": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
	}
}

func (r NetAppVolumeBucketCredentialsResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secret_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"key_pair_expiry": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r NetAppVolumeBucketCredentialsResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			var model netAppModels.NetAppVolumeBucketCredentialsModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := buckets.ParseBucketID(model.BucketID)
			if err != nil {
				return fmt.Errorf("parsing bucket id %s: %+v", model.BucketID, err)
			}

			input := buckets.BucketCredentialsExpiry{
				KeyPairExpiryDays: pointer.To(model.KeyPairExpiryDays),
			}

			if model.StoreInKeyVault {
				if err := client.GenerateAkvCredentialsThenPoll(ctx, pointer.From(id), input); err != nil {
					return fmt.Errorf("generating credentials in key vault for %s: %+v", id, err)
				}
			} else {
				resp, err := client.GenerateCredentials(ctx, pointer.From(id), input)
				if err != nil {
					return fmt.Errorf("generating credentials for %s: %+v", id, err)
				}
				if resp.Model != nil {
					model.AccessKey = pointer.From(resp.Model.AccessKey)
					model.SecretKey = pointer.From(resp.Model.SecretKey)
					model.KeyPairExpiry = pointer.From(resp.Model.KeyPairExpiry)
				}
			}

			metadata.SetID(pointer.From(id))

			// Refresh status from the bucket
			bucketResp, err := client.Get(ctx, pointer.From(id))
			if err == nil && bucketResp.Model != nil && bucketResp.Model.Properties != nil && bucketResp.Model.Properties.Status != nil {
				model.Status = string(pointer.From(bucketResp.Model.Properties.Status))
			}

			return metadata.Encode(&model)
		},
	}
}

func (r NetAppVolumeBucketCredentialsResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BucketsClient

			id, err := buckets.ParseBucketID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state netAppModels.NetAppVolumeBucketCredentialsModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.BucketID = id.ID()

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Status != nil {
				state.Status = string(pointer.From(resp.Model.Properties.Status))
			}

			// access_key, secret_key and key_pair_expiry are only returned at generation time;
			// preserve from existing state.
			return metadata.Encode(&state)
		},
	}
}

func (r NetAppVolumeBucketCredentialsResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			// There is no API to revoke generated credentials; expiry is the only control.
			// The credentials become invalid on the configured expiry date and are
			// implicitly revoked when new credentials are generated for the same bucket.
			return nil
		},
	}
}
