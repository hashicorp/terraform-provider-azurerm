// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storagemover

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/storagemovers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StorageMoverSmbMountEndpointModel struct {
	Name                     string `tfschema:"name"`
	StorageMoverId           string `tfschema:"storage_mover_id"`
	Host                     string `tfschema:"host"`
	ShareName                string `tfschema:"share_name"`
	UsernameKeyVaultSecretId string `tfschema:"username_key_vault_secret_id"`
	PasswordKeyVaultSecretId string `tfschema:"password_key_vault_secret_id"`
	Description              string `tfschema:"description"`
}

type StorageMoverSmbMountEndpointResource struct{}

var (
	_ sdk.ResourceWithUpdate   = StorageMoverSmbMountEndpointResource{}
	_ sdk.ResourceWithIdentity = StorageMoverSmbMountEndpointResource{}
)

func (r StorageMoverSmbMountEndpointResource) ResourceType() string {
	return "azurerm_storage_mover_smb_mount_endpoint"
}

func (r StorageMoverSmbMountEndpointResource) ModelObject() interface{} {
	return &StorageMoverSmbMountEndpointModel{}
}

func (r StorageMoverSmbMountEndpointResource) Identity() resourceids.ResourceId {
	return &endpoints.EndpointId{}
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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				validation.IsIPv4Address,
				validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9.-]*[a-zA-Z0-9])?$`),
					"Host must be a valid IPv4 address or hostname/FQDN (letters, numbers, dots, hyphens only).",
				),
			),
		},

		"share_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 80),
				validation.StringMatch(
					regexp.MustCompile(`^[^\\/\[\]:<> +=;,*?\x00-\x1f\x7f]+$`),
					"Share name must be 1-80 characters and cannot contain: \\ / [ ] : < > + = ; , * ? or control characters.",
				),
			),
		},

		"username_key_vault_secret_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			RequiredWith: []string{"password_key_vault_secret_id"},
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
		},

		"password_key_vault_secret_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			RequiredWith: []string{"username_key_vault_secret_id"},
			ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 1024),
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

			if model.UsernameKeyVaultSecretId != "" || model.PasswordKeyVaultSecretId != "" {
				if model.UsernameKeyVaultSecretId == "" || model.PasswordKeyVaultSecretId == "" {
					return fmt.Errorf("both `username_key_vault_secret_id` and `password_key_vault_secret_id` must be specified together when configuring SMB mount endpoint credentials")
				}
				endpointProperties.Credentials = &endpoints.AzureKeyVaultSmbCredentials{
					Type:        endpoints.CredentialTypeAzureKeyVaultSmb,
					UsernameUri: pointer.To(model.UsernameKeyVaultSecretId),
					PasswordUri: pointer.To(model.PasswordKeyVaultSecretId),
				}
			}

			if model.Description != "" {
				endpointProperties.Description = pointer.To(model.Description)
			}

			properties := endpoints.Endpoint{
				Properties: endpointProperties,
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return metadata.Encode(&model)
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
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if v, ok := properties.Properties.(endpoints.SmbMountEndpointProperties); ok {
				if metadata.ResourceData.HasChange("description") {
					v.Description = pointer.To(model.Description)
				}

				if metadata.ResourceData.HasChange("username_key_vault_secret_id") || metadata.ResourceData.HasChange("password_key_vault_secret_id") {
					bothSet := model.UsernameKeyVaultSecretId != "" && model.PasswordKeyVaultSecretId != ""
					bothEmpty := model.UsernameKeyVaultSecretId == "" && model.PasswordKeyVaultSecretId == ""
					switch {
					case bothSet:
						v.Credentials = &endpoints.AzureKeyVaultSmbCredentials{
							Type:        endpoints.CredentialTypeAzureKeyVaultSmb,
							UsernameUri: pointer.To(model.UsernameKeyVaultSecretId),
							PasswordUri: pointer.To(model.PasswordKeyVaultSecretId),
						}
					case bothEmpty:
						v.Credentials = nil
					default:
						return fmt.Errorf("both `username_key_vault_secret_id` and `password_key_vault_secret_id` must be specified together")
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
						state.UsernameKeyVaultSecretId = pointer.From(v.Credentials.UsernameUri)
						state.PasswordKeyVaultSecretId = pointer.From(v.Credentials.PasswordUri)
					}
					state.Description = pointer.From(v.Description)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
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
