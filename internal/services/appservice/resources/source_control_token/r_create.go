// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package source_control_token

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var sourceControlToken AppServiceSourceControlTokenModel

			if err := metadata.Decode(&sourceControlToken); err != nil {
				return err
			}

			id := resourceproviders.NewSourceControlID(sourceControlToken.Type)

			client := metadata.Client.AppService.ResourceProvidersClient

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetSourceControl(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("%s not found", id)
					}
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
				if existing.Model.Properties != nil && pointer.From(existing.Model.Properties.Token) != "" {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			sourceControlOAuth := resourceproviders.SourceControl{
				Properties: &resourceproviders.SourceControlProperties{
					Token:       pointer.To(sourceControlToken.Token),
					TokenSecret: pointer.To(sourceControlToken.TokenSecret),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, id, sourceControlOAuth); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
