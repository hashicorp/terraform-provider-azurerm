// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package function_app_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceConnector.ServiceLinkerClient
			id, err := servicelinker.ParseScopedLinkerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.LinkerGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			pwd := metadata.ResourceData.Get("authentication.0.secret").(string)

			if model := resp.Model; model != nil {
				props := model.Properties
				if props.AuthInfo == nil || props.TargetService == nil {
					return nil
				}

				state := FunctionAppConnectorResourceModel{
					Name:             id.LinkerName,
					FunctionAppId:    id.ResourceUri,
					TargetResourceId: helpers.FlattenTargetService(props.TargetService),
					AuthInfo:         helpers.FlattenServiceConnectorAuthInfo(props.AuthInfo, pwd),
				}

				if props.ClientType != nil {
					state.ClientType = string(*props.ClientType)
				}

				if props.VNetSolution != nil && props.VNetSolution.Type != nil {
					state.VnetSolution = string(*props.VNetSolution.Type)
				}

				if props.SecretStore != nil {
					state.SecretStore = helpers.FlattenSecretStore(*props.SecretStore)
				}

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}
