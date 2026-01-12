// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageMoverSmbMountEndpointModel struct {
	Name           string `tfschema:"name"`
	StorageMoverId string `tfschema:"storage_mover_id"`
	Host           string `tfschema:"host"`
	ShareName      string `tfschema:"share_name"`
	UsernameUri    string `tfschema:"username_uri"`
	PasswordUri    string `tfschema:"password_uri"`
	Description    string `tfschema:"description"`
}

type StorageMoverSmbMountEndpointResource struct{}

var _ sdk.ResourceWithUpdate = StorageMoverSmbMountEndpointResource{}

func (r StorageMoverSmbMountEndpointResource) ResourceType() string {
	return "azurerm_storage_mover_smb_mount_endpoint"
}

func (r StorageMoverSmbMountEndpointResource) ModelObject() interface{} {
	return &StorageMoverSmbMountEndpointModel{}
}

func (r StorageMoverSmbMountEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return endpoints.ValidateEndpointID
}

func (r StorageMoverSmbMountEndpointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9a-zA-Z][-_0-9a-zA-Z]{0,63}$`),
				`The name must be between 1 and 64 characters in length, begin with a letter or number, and may contain letters, numbers, dashes and underscore.`,
			),
		},

		"storage_mover_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: storagemovers.ValidateStorageMoverID,
		},

		"host": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"share_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"username_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"password_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r StorageMoverSmbMountEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageMoverSmbMountEndpointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StorageMoverSmbMountEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.StorageMover.EndpointsClient
			storageMoverId, err := storagemovers.ParseStorageMoverID(model.StorageMoverId)
			if err != nil {
				return err
			}

			id := endpoints.NewEndpointID(storageMoverId.SubscriptionId, storageMoverId.ResourceGroupName, storageMoverId.StorageMoverName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			endpointProperties := endpoints.SmbMountEndpointProperties{
				Host:      model.Host,
				ShareName: model.ShareName,
			}

			if model.UsernameUri != "" || model.PasswordUri != "" {
				endpointProperties.Credentials = &endpoints.AzureKeyVaultSmbCredentials{
					Type:        endpoints.CredentialTypeAzureKeyVaultSmb,
					UsernameUri: utils.String(model.UsernameUri),
					PasswordUri: utils.String(model.PasswordUri),
				}
			}

			if model.Description != "" {
				endpointProperties.Description = utils.String(model.Description)
			}

			properties := endpoints.Endpoint{
				Properties: endpointProperties,
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r StorageMoverSmbMountEndpointResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.EndpointsClient

			id, err := endpoints.ParseEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageMoverSmbMountEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if v, ok := properties.Properties.(endpoints.SmbMountEndpointProperties); ok {
				if metadata.ResourceData.HasChange("description") {
					v.Description = utils.String(model.Description)
				}

				if metadata.ResourceData.HasChange("username_uri") || metadata.ResourceData.HasChange("password_uri") {
					if model.UsernameUri != "" || model.PasswordUri != "" {
						v.Credentials = &endpoints.AzureKeyVaultSmbCredentials{
							Type:        endpoints.CredentialTypeAzureKeyVaultSmb,
							UsernameUri: utils.String(model.UsernameUri),
							PasswordUri: utils.String(model.PasswordUri),
						}
					} else {
						v.Credentials = nil
					}
				}

				properties.Properties = v
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r StorageMoverSmbMountEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.EndpointsClient

			id, err := endpoints.ParseEndpointID(metadata.ResourceData.Id())
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

			state := StorageMoverSmbMountEndpointModel{
				Name:           id.EndpointName,
				StorageMoverId: storagemovers.NewStorageMoverID(id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName).ID(),
			}

			if model := resp.Model; model != nil {
				if v, ok := model.Properties.(endpoints.SmbMountEndpointProperties); ok {
					state.Host = v.Host
					state.ShareName = v.ShareName

					if v.Credentials != nil {
						if v.Credentials.UsernameUri != nil {
							state.UsernameUri = *v.Credentials.UsernameUri
						}
						if v.Credentials.PasswordUri != nil {
							state.PasswordUri = *v.Credentials.PasswordUri
						}
					}

					des := ""
					if v.Description != nil {
						des = *v.Description
					}
					state.Description = des
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r StorageMoverSmbMountEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageMover.EndpointsClient

			id, err := endpoints.ParseEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
