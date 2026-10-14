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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.ResourceProvidersClient

			id, err := resourceproviders.ParseSourceControlID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetSourceControl(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := AppServiceSourceControlTokenModel{}

			state.Type = id.SourceControlName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Token = pointer.From(props.Token)
					state.TokenSecret = pointer.From(props.TokenSecret)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
