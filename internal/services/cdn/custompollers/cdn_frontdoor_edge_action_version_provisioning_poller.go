// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &frontDoorEdgeActionVersionProvisioningPoller{}

// frontDoorEdgeActionVersionProvisioningPoller waits for an Edge Action Version to finish provisioning, since
// deploying the version code continues in the background after the `deployVersionCode` operation has completed
type frontDoorEdgeActionVersionProvisioningPoller struct {
	client *edgeactionversions.EdgeActionVersionsClient
	id     edgeactionversions.VersionId
}

func NewFrontDoorEdgeActionVersionProvisioningPoller(client *edgeactionversions.EdgeActionVersionsClient, id edgeactionversions.VersionId) pollers.PollerType {
	return &frontDoorEdgeActionVersionProvisioningPoller{
		client: client,
		id:     id,
	}
}

func (p frontDoorEdgeActionVersionProvisioningPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", p.id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return nil, fmt.Errorf("retrieving %s: `properties` was nil", p.id)
	}

	props := resp.Model.Properties
	if props.ValidationStatus != nil && *props.ValidationStatus == edgeactionversions.EdgeActionVersionValidationStatusFailed {
		return nil, fmt.Errorf("the code for %s failed validation", p.id)
	}

	switch state := pointer.From(props.ProvisioningState); state {
	case edgeactionversions.ProvisioningStateSucceeded:
		return &pollers.PollResult{
			PollInterval: 30 * time.Second,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	case edgeactionversions.ProvisioningStateFailed, edgeactionversions.ProvisioningStateCanceled:
		return nil, fmt.Errorf("%s finished provisioning with state %q", p.id, string(state))
	}

	return &pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}, nil
}
