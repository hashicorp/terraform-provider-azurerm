// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package advisor_suppression

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/suppressions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Advisor.SuppressionsClient

			state := AdvisorSuppressionResourceModel{}

			id, err := suppressions.ParseScopedSuppressionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state.Name = id.SuppressionName
			state.ResourceID = id.ResourceUri
			state.RecommendationID = id.RecommendationId

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.TTL = pointer.From(props.Ttl)
					state.SuppressionID = pointer.From(props.SuppressionId)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
