// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace_network_outbound_rule_service_tag

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/managednetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.ManagedNetwork
			id, err := managednetwork.ParseOutboundRuleID(metadata.ResourceData.Id())

			var model MachineLearningWorkspaceServiceTagOutboundRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if err != nil {
				return err
			}

			existing, err := client.SettingsRuleGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model
			serviceTagOutboundRule := payload.Properties.(managednetwork.ServiceTagOutboundRule)
			if metadata.ResourceData.HasChange("service_tag") {
				serviceTagOutboundRule.Destination.ServiceTag = pointer.To(model.ServiceTag)
			}

			if metadata.ResourceData.HasChange("protocol") {
				serviceTagOutboundRule.Destination.Protocol = pointer.To(model.Protocol)
			}

			if metadata.ResourceData.HasChange("port_ranges") {
				serviceTagOutboundRule.Destination.PortRanges = pointer.To(model.PortRanges)
			}

			payload.Properties = serviceTagOutboundRule
			if err := client.SettingsRuleCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}
