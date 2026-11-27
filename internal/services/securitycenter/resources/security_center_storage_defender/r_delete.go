// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_storage_defender

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2025-06-01/defenderforstorage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SecurityCenter.DefenderForStorageClient

			id, err := commonids.ParseScopeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Get(ctx, *id); err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			input := defenderforstorage.DefenderForStorageSetting{
				Properties: &defenderforstorage.DefenderForStorageSettingProperties{
					IsEnabled: pointer.To(false),
				},
			}

			if _, err = client.Create(ctx, *id, input); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
