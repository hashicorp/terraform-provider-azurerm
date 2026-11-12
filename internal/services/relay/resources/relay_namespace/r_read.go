// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package relay_namespace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/relay/2021-11-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Relay.NamespacesClient

			id, err := namespaces.ParseNamespaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := RelayNamespaceResourceModel{
				Name:              id.NamespaceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)

				if sku := model.Sku; sku != nil {
					state.SkuName = string(sku.Name)
				}

				if props := resp.Model.Properties; props != nil {
					state.MetricID = pointer.From(props.MetricId)
				}
			}

			authRuleId := namespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, "RootManageSharedAccessKey")

			keysResp, err := client.ListKeys(ctx, authRuleId)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
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
