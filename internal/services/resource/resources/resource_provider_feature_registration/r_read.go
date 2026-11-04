// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_feature_registration

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.FeaturesClient

			featureId, err := features.ParseFeatureID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *featureId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(featureId)
				}

				return fmt.Errorf("retrieving %s: %+v", featureId, err)
			}

			registrationState := ""
			if model := resp.Model; model != nil && model.Properties != nil && model.Properties.State != nil {
				registrationState = *model.Properties.State
			}

			if !strings.EqualFold(registrationState, Registered) {
				return metadata.MarkAsGone(featureId)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, featureId); err != nil {
				return err
			}

			return metadata.Encode(&ResourceProviderFeatureRegistrationModel{
				Name:         featureId.FeatureName,
				ProviderName: featureId.ProviderName,
			})
		},
		Timeout: 5 * time.Minute,
	}
}
