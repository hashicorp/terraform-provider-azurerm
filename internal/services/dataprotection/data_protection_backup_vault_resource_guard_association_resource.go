// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	resourceguardproxy "github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardproxybaseresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguardresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/resourceguards"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	dataProtectionBackupVaultResourceGuardAssociationDeleteRequestName = "default"
	dataProtectionBackupVaultResourceGuardAssociationProxyName         = "DppResourceGuardProxy"
	dataProtectionBackupVaultResourceGuardAssociationResourceType      = "azurerm_data_protection_backup_vault_resource_guard_association"
)

type DataProtectionBackupVaultResourceGuardAssociationModel struct {
	DataProtectionBackupVaultId   string `tfschema:"data_protection_backup_vault_id"`
	DataProtectionResourceGuardId string `tfschema:"data_protection_resource_guard_id"`
}

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_protection_backup_vault_resource_guard_association -service-package-name dataprotection -compare-values "subscription_id:data_protection_backup_vault_id,resource_group_name:data_protection_backup_vault_id,backup_vault_name:data_protection_backup_vault_id" -known-values "name:DppResourceGuardProxy"

type DataProtectionBackupVaultResourceGuardAssociationResource struct{}

var _ sdk.ResourceWithIdentity = DataProtectionBackupVaultResourceGuardAssociationResource{}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) Identity() resourceids.ResourceId {
	return &resourceguardproxy.BackupResourceGuardProxyId{}
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return resourceguardproxy.ValidateBackupResourceGuardProxyID
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) ModelObject() interface{} {
	return &DataProtectionBackupVaultResourceGuardAssociationModel{}
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) ResourceType() string {
	return dataProtectionBackupVaultResourceGuardAssociationResourceType
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"data_protection_backup_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&backupvaultresources.BackupVaultId{}),

		"data_protection_resource_guard_id": commonschema.ResourceIDReferenceRequiredForceNew(&resourceguardresources.ResourceGuardId{}),
	}
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DataProtectionBackupVaultResourceGuardAssociationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.ResourceGuardProxyClient

			vaultId, err := backupvaultresources.ParseBackupVaultID(model.DataProtectionBackupVaultId)
			if err != nil {
				return err
			}

			id := resourceguardproxy.NewBackupResourceGuardProxyID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, dataProtectionBackupVaultResourceGuardAssociationProxyName)

			existing, err := client.DppResourceGuardProxyGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			proxy := resourceguardproxy.ResourceGuardProxyBaseResource{
				Properties: &resourceguardproxy.ResourceGuardProxyBase{
					ResourceGuardResourceId: pointer.To(model.DataProtectionResourceGuardId),
				},
			}

			if _, err = client.DppResourceGuardProxyCreateOrUpdate(ctx, id, proxy); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.ResourceGuardProxyClient

			id, err := resourceguardproxy.ParseBackupResourceGuardProxyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.DppResourceGuardProxyGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			state := flattenDataProtectionBackupVaultResourceGuardAssociation(*id, resp.Model)

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupVaultResourceGuardAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DataProtectionBackupVaultResourceGuardAssociationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.ResourceGuardProxyClient

			id, err := resourceguardproxy.ParseBackupResourceGuardProxyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			guardId, err := resourceguardresources.ParseResourceGuardID(model.DataProtectionResourceGuardId)
			if err != nil {
				return err
			}

			requestId := resourceguards.NewDeleteResourceGuardProxyRequestID(guardId.SubscriptionId, guardId.ResourceGroupName, guardId.ResourceGuardName, dataProtectionBackupVaultResourceGuardAssociationDeleteRequestName)
			unlock := resourceguardproxy.UnlockDeleteRequest{
				ResourceGuardOperationRequests: pointer.To([]string{requestId.ID()}),
			}

			if _, err = client.DppResourceGuardProxyUnlockDelete(ctx, *id, unlock, resourceguardproxy.DefaultDppResourceGuardProxyUnlockDeleteOperationOptions()); err != nil {
				return fmt.Errorf("unlocking delete for %s: %+v", *id, err)
			}

			if _, err = client.DppResourceGuardProxyDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func flattenDataProtectionBackupVaultResourceGuardAssociation(id resourceguardproxy.BackupResourceGuardProxyId, input *resourceguardproxy.ResourceGuardProxyBaseResource) DataProtectionBackupVaultResourceGuardAssociationModel {
	vaultId := backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)
	state := DataProtectionBackupVaultResourceGuardAssociationModel{
		DataProtectionBackupVaultId: vaultId.ID(),
	}

	if input != nil && input.Properties != nil {
		state.DataProtectionResourceGuardId = pointer.From(input.Properties.ResourceGuardResourceId)
	}

	return state
}

func setDataProtectionBackupVaultResourceGuardAssociationResourceData(d *pluginsdk.ResourceData, id resourceguardproxy.BackupResourceGuardProxyId, input *resourceguardproxy.ResourceGuardProxyBaseResource) error {
	state := flattenDataProtectionBackupVaultResourceGuardAssociation(id, input)

	if err := d.Set("data_protection_backup_vault_id", state.DataProtectionBackupVaultId); err != nil {
		return fmt.Errorf("setting `data_protection_backup_vault_id`: %+v", err)
	}

	if err := d.Set("data_protection_resource_guard_id", state.DataProtectionResourceGuardId); err != nil {
		return fmt.Errorf("setting `data_protection_resource_guard_id`: %+v", err)
	}

	return pluginsdk.SetResourceIdentityData(d, &id)
}
