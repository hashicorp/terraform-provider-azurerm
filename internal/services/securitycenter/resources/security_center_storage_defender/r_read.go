// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_storage_defender

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := securitycenterClient.From(metadata.Client).DefenderForStorageClient

			id, err := commonids.ParseScopeID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			storageAccountId, err := commonids.ParseStorageAccountID(id.Scope)
			if err != nil {
				return err
			}
			state := StorageDefenderModel{
				StorageAccountId: storageAccountId.ID(),
			}

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					if !pointer.From(prop.IsEnabled) {
						return metadata.MarkAsGone(id)
					}

					state.OverrideSubscriptionSettings = pointer.From(prop.OverrideSubscriptionLevelSettings)

					if ms := prop.MalwareScanning; ms != nil {
						if onUpload := ms.OnUpload; onUpload != nil {
							state.MalwareScanningOnUploadEnabled = pointer.From(onUpload.IsEnabled)
							state.MalwareScanningOnUploadCapPerMon = pointer.From(onUpload.CapGBPerMonth)
						}
						if ms.ScanResultsEventGridTopicResourceId != nil {
							topicId, err := topics.ParseTopicID(*ms.ScanResultsEventGridTopicResourceId)
							if err != nil {
								return err
							}
							state.ScanResultsEventGridTopicId = topicId.ID()
						}
					}

					if sdd := prop.SensitiveDataDiscovery; sdd != nil {
						state.SensitiveDataDiscoveryEnabled = pointer.From(sdd.IsEnabled)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
