// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
)

var (
	_ pollers.PollerType = &frontDoorRouteCreatePoller{}
	_ pollers.PollerType = &frontDoorRouteUpdatePoller{}
	_ pollers.PollerType = &frontDoorRouteDeletePoller{}
)

// NOTE: Front Door can continue backend synchronization after the ARM write call
// returns. For route operations, `provisioningState` is the authoritative
// readiness signal. `deploymentStatus` is still useful for surfacing terminal
// failures, but it can legitimately remain `NotStarted` after the service has
// already applied the update.
const frontDoorRoutePollInterval = 30 * time.Second

type frontDoorRouteCreatePoller struct {
	endpointClient  *cdn.AFDEndpointsClient
	routeClient     *cdn.RoutesClient
	id              parse.FrontDoorRouteId
	input           cdn.Route
	operationIssued bool
}

type frontDoorRouteUpdatePoller struct {
	client            *cdn.RoutesClient
	workaroundsClient azuresdkhacks.CdnFrontDoorRoutesWorkaroundClient
	id                parse.FrontDoorRouteId
	input             azuresdkhacks.RouteUpdateParameters
	operationIssued   bool
}

type frontDoorRouteDeletePoller struct {
	client          *cdn.RoutesClient
	id              parse.FrontDoorRouteId
	operationIssued bool
}

func NewFrontDoorRouteCreatePoller(endpointClient *cdn.AFDEndpointsClient, routeClient *cdn.RoutesClient, id parse.FrontDoorRouteId, input cdn.Route) pollers.PollerType {
	return &frontDoorRouteCreatePoller{
		endpointClient: endpointClient,
		routeClient:    routeClient,
		id:             id,
		input:          input,
	}
}

func NewFrontDoorRouteUpdatePoller(client *cdn.RoutesClient, workaroundsClient azuresdkhacks.CdnFrontDoorRoutesWorkaroundClient, id parse.FrontDoorRouteId, input azuresdkhacks.RouteUpdateParameters) pollers.PollerType {
	return &frontDoorRouteUpdatePoller{
		client:            client,
		workaroundsClient: workaroundsClient,
		id:                id,
		input:             input,
	}
}

func NewFrontDoorRouteDeletePoller(client *cdn.RoutesClient, id parse.FrontDoorRouteId) pollers.PollerType {
	return &frontDoorRouteDeletePoller{
		client: client,
		id:     id,
	}
}

func (p *frontDoorRouteCreatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if !p.operationIssued {
		ready, err := afdEndpointSettledForRouteOperation(ctx, p.endpointClient, p.id)
		if err != nil {
			return nil, err
		}
		if !ready {
			return frontDoorRouteInProgress(), nil
		}

		future, err := p.routeClient.Create(ctx, p.id.ResourceGroup, p.id.ProfileName, p.id.AfdEndpointName, p.id.RouteName, p.input)
		if err != nil {
			if frontDoorRouteOperationInProgress(err) {
				return frontDoorRouteInProgress(), nil
			}

			return nil, fmt.Errorf("creating %s: %+v", p.id, err)
		}

		if err := future.WaitForCompletionRef(ctx, p.routeClient.Client); err != nil {
			if frontDoorRouteOperationInProgress(err) {
				return frontDoorRouteInProgress(), nil
			}

			return nil, fmt.Errorf("waiting for the creation of %s: %+v", p.id, err)
		}

		p.operationIssued = true
	}

	ready, err := frontDoorRouteSettledForCreate(ctx, p.routeClient, p.id)
	if err != nil {
		if routeNotFound(err) {
			return frontDoorRouteInProgress(), nil
		}

		return nil, err
	}
	if !ready {
		return frontDoorRouteInProgress(), nil
	}

	return frontDoorRouteSucceeded(), nil
}

func (p *frontDoorRouteUpdatePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if !p.operationIssued {
		ready, err := frontDoorRouteSettled(ctx, p.client, p.id)
		if err != nil {
			return nil, err
		}
		if !ready {
			return frontDoorRouteInProgress(), nil
		}

		future, err := p.workaroundsClient.Update(ctx, p.id.ResourceGroup, p.id.ProfileName, p.id.AfdEndpointName, p.id.RouteName, p.input)
		if err != nil {
			if frontDoorRouteOperationInProgress(err) {
				return frontDoorRouteInProgress(), nil
			}

			return nil, fmt.Errorf("updating %s: %+v", p.id, err)
		}

		if err := future.WaitForCompletionRef(ctx, p.client.Client); err != nil {
			if frontDoorRouteOperationInProgress(err) {
				return frontDoorRouteInProgress(), nil
			}

			return nil, fmt.Errorf("waiting for the update of %s: %+v", p.id, err)
		}

		p.operationIssued = true
	}

	ready, err := frontDoorRouteSettled(ctx, p.client, p.id)
	if err != nil {
		return nil, err
	}
	if !ready {
		return frontDoorRouteInProgress(), nil
	}

	return frontDoorRouteSucceeded(), nil
}

func (p *frontDoorRouteDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	if !p.operationIssued {
		ready, err := frontDoorRouteSettledForCreate(ctx, p.client, p.id)
		if err != nil {
			if routeNotFound(err) {
				return frontDoorRouteSucceeded(), nil
			}

			return nil, err
		}
		if !ready {
			return frontDoorRouteInProgress(), nil
		}

		future, err := p.client.Delete(ctx, p.id.ResourceGroup, p.id.ProfileName, p.id.AfdEndpointName, p.id.RouteName)
		if err != nil {
			if frontDoorRouteOperationInProgress(err) {
				return frontDoorRouteInProgress(), nil
			}
			if routeNotFound(err) {
				return frontDoorRouteSucceeded(), nil
			}

			return nil, fmt.Errorf("deleting %s: %+v", p.id, err)
		}

		if err := future.WaitForCompletionRef(ctx, p.client.Client); err != nil {
			if frontDoorRouteOperationInProgress(err) {
				return frontDoorRouteInProgress(), nil
			}

			return nil, fmt.Errorf("waiting for the deletion of %s: %+v", p.id, err)
		}

		p.operationIssued = true
	}

	resp, err := p.client.Get(ctx, p.id.ResourceGroup, p.id.ProfileName, p.id.AfdEndpointName, p.id.RouteName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return frontDoorRouteSucceeded(), nil
		}

		return nil, fmt.Errorf("checking deletion of %s: %+v", p.id, err)
	}

	if resp.RouteProperties != nil {
		if resp.ProvisioningState == cdn.AfdProvisioningStateFailed || resp.DeploymentStatus == cdn.DeploymentStatusFailed {
			return nil, fmt.Errorf("waiting for deletion of %s: route entered failed state with `deploymentStatus` `%s` and `provisioningState` `%s`", p.id, resp.DeploymentStatus, resp.ProvisioningState)
		}
	}

	return frontDoorRouteInProgress(), nil
}

func afdEndpointSettledForRouteOperation(ctx context.Context, client *cdn.AFDEndpointsClient, id parse.FrontDoorRouteId) (bool, error) {
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
	if err != nil {
		return false, fmt.Errorf("retrieving Front Door Endpoint %q for route operation on %s: %+v", id.AfdEndpointName, id, err)
	}

	props := resp.AFDEndpointProperties
	if props == nil {
		return false, nil
	}

	if props.ProvisioningState == cdn.AfdProvisioningStateFailed || props.DeploymentStatus == cdn.DeploymentStatusFailed {
		return false, fmt.Errorf("waiting for Front Door Endpoint %q before route operation on %s: endpoint entered failed state with `deploymentStatus` `%s` and `provisioningState` `%s`", id.AfdEndpointName, id, props.DeploymentStatus, props.ProvisioningState)
	}

	return props.ProvisioningState == cdn.AfdProvisioningStateSucceeded, nil
}

func frontDoorRouteSettled(ctx context.Context, client *cdn.RoutesClient, id parse.FrontDoorRouteId) (bool, error) {
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return false, fmt.Errorf("route %s was not found", id)
		}

		return false, fmt.Errorf("retrieving %s while waiting for route state to settle: %+v", id, err)
	}

	props := resp.RouteProperties
	settled, err := frontDoorRouteStatusesSettled(props)
	if err != nil {
		return false, fmt.Errorf("waiting for %s: %+v", id, err)
	}

	return settled, nil
}

func frontDoorRouteSettledForCreate(ctx context.Context, client *cdn.RoutesClient, id parse.FrontDoorRouteId) (bool, error) {
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return false, fmt.Errorf("route %s was not found", id)
		}

		return false, fmt.Errorf("retrieving %s while waiting for route state to settle after create: %+v", id, err)
	}

	props := resp.RouteProperties
	settled, err := frontDoorRouteStatusesSettled(props)
	if err != nil {
		return false, fmt.Errorf("waiting for %s: %+v", id, err)
	}

	return settled, nil
}

func frontDoorRouteStatusesSettled(props *cdn.RouteProperties) (bool, error) {
	if props == nil {
		return false, nil
	}

	if props.ProvisioningState == cdn.AfdProvisioningStateFailed || props.DeploymentStatus == cdn.DeploymentStatusFailed {
		return false, fmt.Errorf("route entered failed state with `deploymentStatus` `%s` and `provisioningState` `%s`", props.DeploymentStatus, props.ProvisioningState)
	}

	return props.ProvisioningState == cdn.AfdProvisioningStateSucceeded, nil
}

func frontDoorRouteInProgress() *pollers.PollResult {
	return &pollers.PollResult{
		PollInterval: frontDoorRoutePollInterval,
		Status:       pollers.PollingStatusInProgress,
	}
}

func frontDoorRouteSucceeded() *pollers.PollResult {
	return &pollers.PollResult{
		PollInterval: frontDoorRoutePollInterval,
		Status:       pollers.PollingStatusSucceeded,
	}
}

func routeNotFound(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(strings.ToLower(err.Error()), "was not found")
}

func frontDoorRouteOperationInProgress(err error) bool {
	if err == nil {
		return false
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "another operation is in progress") || strings.Contains(message, "operation is in progress")
}
