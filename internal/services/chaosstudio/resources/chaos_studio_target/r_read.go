// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_target

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/targets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Targets
			schema := ChaosStudioTargetResourceSchema{}

			id, err := commonids.ParseChaosStudioTargetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.TargetResourceId = id.Scope
				schema.TargetType = id.TargetName
				if err := r.mapTargetToChaosStudioTargetResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r Resource) mapChaosStudioTargetResourceSchemaToTarget(input ChaosStudioTargetResourceSchema, output *targets.Target) error {
	output.Location = pointer.To(location.Normalize(input.Location))
	return nil
}

func (r Resource) mapTargetToChaosStudioTargetResourceSchema(input targets.Target, output *ChaosStudioTargetResourceSchema) error {
	output.Location = location.NormalizeNilable(input.Location)
	return nil
}
