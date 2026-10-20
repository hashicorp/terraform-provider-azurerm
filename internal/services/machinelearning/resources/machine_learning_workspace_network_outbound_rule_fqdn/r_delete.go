// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace_network_outbound_rule_fqdn

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/managednetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.ManagedNetwork

			id, err := managednetwork.ParseOutboundRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.SettingsRuleDelete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting Machine Learning Workspace FQDN Network Outbound Rule %q (Resource Group %q, Workspace %q): %+v", id.OutboundRuleName, id.ResourceGroupName, id.WorkspaceName, err)
			}

			if err = future.Poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of Machine Learning Workspace FQDN Network Outbound Rule %q (Resource Group %q, Workspace %q): %+v", id.OutboundRuleName, id.ResourceGroupName, id.WorkspaceName, err)
			}

			return nil
		},
	}
}
