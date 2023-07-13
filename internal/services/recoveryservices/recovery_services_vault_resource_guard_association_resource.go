// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2022-04-01/resourceguards"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2022-10-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const VaultGuardResourceType = "Microsoft.RecoveryServices/vaults/backupResourceGuardProxies"
const VaultGuardProxyDeleteRequestName = "default" // this name does not matter, this value comes from Portal.

type VaultGuardProxyResource struct{}

var _ sdk.Resource = VaultGuardProxyResource{}

type VaultGuardProxyModel struct {
	Name            string `tfschema:"name"`
	VaultId         string `tfschema:"vault_id"`
	ResourceGuardId string `tfschema:"resource_guard_id"`
}

func (r VaultGuardProxyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return resourceguardproxy.ValidateBackupResourceGuardProxyID
}

func (r VaultGuardProxyResource) ModelObject() interface{} {
	return &VaultGuardProxyModel{}
}

func (r VaultGuardProxyResource) ResourceType() string {
	return "azurerm_recovery_services_vault_resource_guard_association"
}

func (r VaultGuardProxyResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(vaults.VaultId{}),

		"resource_guard_id": commonschema.ResourceIDReferenceRequiredForceNew(resourceguards.ResourceGuardId{}),
	}
}

func (r VaultGuardProxyResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
func (r VaultGuardProxyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan VaultGuardProxyModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %w", err)
			}
			client := metadata.Client.RecoveryServices.ResourceGuardProxyClient

			vaultId, err := vaults.ParseVaultID(plan.VaultId)
			if err != nil {
				return fmt.Errorf("parsing vault id %w", err)
			}

			id := resourceguardproxy.NewBackupResourceGuardProxyID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.VaultName, plan.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking presence of %s:%+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_recovery_services_vault_resource_guard_association", id.ID())
			}

			proxy := resourceguardproxy.ResourceGuardProxyBaseResource{
				Id:   pointer.To(id.ID()),
				Type: pointer.To(VaultGuardResourceType),
				Properties: pointer.To(resourceguardproxy.ResourceGuardProxyBase{
					ResourceGuardResourceId: pointer.To(plan.ResourceGuardId),
				}),
			}

			if _, err = client.Put(ctx, id, proxy); err != nil {
				return fmt.Errorf("creating %s:%w", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r VaultGuardProxyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ResourceGuardProxyClient

			id, err := resourceguardproxy.ParseBackupResourceGuardProxyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %q:%+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("creating %s:%+v", id, err)
			}

			vaultId := vaults.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)
			state := VaultGuardProxyModel{
				VaultId: vaultId.ID(),
				Name:    id.BackupResourceGuardProxyName,
			}

			if resp.Model != nil && resp.Model.Properties != nil {
				state.ResourceGuardId = pointer.From(resp.Model.Properties.ResourceGuardResourceId)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r VaultGuardProxyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan VaultGuardProxyModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %w", err)
			}
			client := metadata.Client.RecoveryServices.ResourceGuardProxyClient

			id, err := resourceguardproxy.ParseBackupResourceGuardProxyID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %q:%+v", metadata.ResourceData.Id(), err)
			}

			guardId, err := resourceguards.ParseResourceGuardID(plan.ResourceGuardId)
			if err != nil {
				return fmt.Errorf("parsing %q:%+v", plan.ResourceGuardId, err)
			}

			requestId := resourceguards.NewDeleteResourceGuardProxyRequestID(guardId.SubscriptionId, guardId.ResourceGroupName, guardId.ResourceGuardName, VaultGuardProxyDeleteRequestName)

			unlock := resourceguardproxy.UnlockDeleteRequest{
				ResourceGuardOperationRequests: pointer.To([]string{requestId.ID()}),
			}

			if _, err = client.UnlockDelete(ctx, *id, unlock); err != nil {
				return fmt.Errorf("unlocking delete %s:%+v", id, err)
			}

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s:%+v", id, err)
			}

			return nil
		},
	}
}
