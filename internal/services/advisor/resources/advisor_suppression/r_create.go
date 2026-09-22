// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package advisor_suppression

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/suppressions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Advisor.SuppressionsClient

			var model AdvisorSuppressionResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := suppressions.NewScopedSuppressionID(model.ResourceID, model.RecommendationID, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			param := suppressions.SuppressionContract{
				Name: new(model.Name),
				Properties: &suppressions.SuppressionProperties{
					SuppressionId: new(model.SuppressionID),
				},
			}

			if model.TTL != "" {
				param.Properties.Ttl = new(model.TTL)
			}

			if _, err := client.Create(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
