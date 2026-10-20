// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace_network_outbound_rule_service_tag

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/managednetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.ManagedNetwork

			id, err := managednetwork.ParseOutboundRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.SettingsRuleGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := MachineLearningWorkspaceServiceTagOutboundRuleModel{
				Name:        id.OutboundRuleName,
				WorkspaceId: managednetwork.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if prop, ok := props.(managednetwork.ServiceTagOutboundRule); ok && prop.Destination != nil {
						state.ServiceTag = pointer.From(prop.Destination.ServiceTag)
						state.Protocol = pointer.From(prop.Destination.Protocol)
						state.PortRanges = pointer.From(prop.Destination.PortRanges)
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}
