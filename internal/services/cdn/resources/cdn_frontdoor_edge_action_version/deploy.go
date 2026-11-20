// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_version

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/custompollers"
)

func deployVersionCode(ctx context.Context, client *edgeactionversions.EdgeActionVersionsClient, id edgeactionversions.VersionId, config EdgeActionVersionModel) error {
	input := edgeactionversions.VersionCode{
		Name:    config.CodeFileName,
		Content: config.CodeContentBase64,
	}

	if err := client.DeployVersionCodeThenPoll(ctx, id, input); err != nil {
		return fmt.Errorf("deploying code for %s: %+v", id, err)
	}

	// the version continues provisioning after the `deployVersionCode` operation completes,
	// so wait for it to finish before the version can be used (e.g. by an execution filter)
	pollerType := custompollers.NewFrontDoorEdgeActionVersionProvisioningPoller(client, id)
	poller := pollers.NewPoller(pollerType, 30*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for code deployment of %s: %+v", id, err)
	}

	return nil
}
