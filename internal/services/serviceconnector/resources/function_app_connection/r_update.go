// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package function_app_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/links"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceConnector.LinksClient
			id, err := links.ParseScopedLinkerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state FunctionAppConnectorResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding state: %+v", err)
			}

			linkerProps := links.LinkerProperties{}
			d := metadata.ResourceData

			if d.HasChange("client_type") {
				linkerProps.ClientType = pointer.ToEnum[links.ClientType](state.ClientType)
			}

			if d.HasChange("vnet_solution") {
				vnetSolution := links.VNetSolution{
					Type: pointer.ToEnum[links.VNetSolutionType](state.VnetSolution),
				}
				linkerProps.VNetSolution = &vnetSolution
			}

			if d.HasChange("secret_store") {
				linkerProps.SecretStore = pointer.To(links.SecretStore{KeyVaultId: helpers.ExpandSecretStore(state.SecretStore).KeyVaultId})
			}

			if d.HasChange("authentication") {
				authInfo, err := helpers.ExpandServiceConnectorAuthInfoForUpdate(state.AuthInfo)
				if err != nil {
					return fmt.Errorf("expanding `authentication`: %+v", err)
				}

				linkerProps.AuthInfo = authInfo
			}

			props := links.LinkerPatch{
				Properties: &linkerProps,
			}

			if err := client.LinkerUpdateThenPoll(ctx, *id, props); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
