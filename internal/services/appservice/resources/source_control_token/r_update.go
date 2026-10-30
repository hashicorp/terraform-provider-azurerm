// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package source_control_token

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var sourceControlToken AppServiceSourceControlTokenModel

			id, err := resourceproviders.ParseSourceControlID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := metadata.Decode(&sourceControlToken); err != nil {
				return err
			}

			client := metadata.Client.AppService.ResourceProvidersClient

			sourceControlOAuth := resourceproviders.SourceControl{
				Properties: &resourceproviders.SourceControlProperties{
					Token:       pointer.To(sourceControlToken.Token),
					TokenSecret: pointer.To(sourceControlToken.TokenSecret),
				},
			}

			if _, err := client.UpdateSourceControl(ctx, *id, sourceControlOAuth); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
