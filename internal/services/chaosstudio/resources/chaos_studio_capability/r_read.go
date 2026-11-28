// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_capability

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Capabilities
			schema := ChaosStudioCapabilityResourceSchema{}

			id, err := commonids.ParseChaosStudioCapabilityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			targetId := commonids.NewChaosStudioTargetID(id.Scope, id.TargetName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema.ChaosStudioTargetId = targetId.ID()
			schema.CapabilityType = id.CapabilityName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					schema.Urn = pointer.From(props.Urn)
				}
			}

			return metadata.Encode(&schema)
		},
		Timeout: 5 * time.Minute,
	}
}

type ChaosStudioCapabilityResourceSchema struct {
	CapabilityType      string `tfschema:"capability_type"`
	ChaosStudioTargetId string `tfschema:"chaos_studio_target_id"`
	Urn                 string `tfschema:"urn"`
}
