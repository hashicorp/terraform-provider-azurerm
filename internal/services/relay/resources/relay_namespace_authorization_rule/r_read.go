// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace_authorization_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	relayClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := relayClient.From(metadata.Client).NamespacesClient

			id, err := namespaces.ParseAuthorizationRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetAuthorizationRule(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			keysResp, err := client.ListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
			}

			state := RelayNamespaceAuthorizationRuleResourceModel{
				Name:              id.AuthorizationRuleName,
				NamespaceName:     id.NamespaceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Listen, state.Send, state.Manage = flattenAuthorizationRuleRights(props.Rights)
				}
			}

			if keysModel := keysResp.Model; keysModel != nil {
				state.PrimaryConnectionString = pointer.From(keysModel.PrimaryConnectionString)
				state.SecondaryConnectionString = pointer.From(keysModel.SecondaryConnectionString)
				state.PrimaryKey = pointer.From(keysModel.PrimaryKey)
				state.SecondaryKey = pointer.From(keysModel.SecondaryKey)
			}
			return metadata.Encode(&state)
		},
	}
}
