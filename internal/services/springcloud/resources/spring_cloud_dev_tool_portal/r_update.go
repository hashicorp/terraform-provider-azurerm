// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_dev_tool_portal

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func (s Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := springcloudClient.From(metadata.Client).DevToolPortalClient

			id, err := parse.SpringCloudDevToolPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SpringCloudDevToolPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if _, err = client.Get(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName); err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			DevToolPortalResource := appplatform.DevToolPortalResource{
				Properties: &appplatform.DevToolPortalProperties{
					Public:        pointer.To(model.PublicNetworkAccessEnabled),
					SsoProperties: expandSpringCloudDevToolPortalSsoProperties(model.Sso),
					Features:      expandSpringCloudDevToolPortalFeatures(model.ApplicationAcceleratorEnabled, model.ApplicationLiveViewEnabled),
				},
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName, DevToolPortalResource)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}

			return nil
		},
	}
}
