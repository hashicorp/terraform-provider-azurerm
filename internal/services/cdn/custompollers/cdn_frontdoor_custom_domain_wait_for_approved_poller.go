// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-04-15/afdcustomdomains"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &frontDoorCustomDomainWaitForApprovedPoller{}

type frontDoorCustomDomainWaitForApprovedPoller struct {
	client *afdcustomdomains.AFDCustomDomainsClient
	id     afdcustomdomains.CustomDomainId
}

var (
	frontDoorCustomDomainWaitForApprovedSuccess = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusSucceeded,
	}
	frontDoorCustomDomainWaitForApprovedInProgress = pollers.PollResult{
		PollInterval: 30 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}
)

func NewFrontDoorCustomDomainWaitForApprovedPoller(client *afdcustomdomains.AFDCustomDomainsClient, id afdcustomdomains.CustomDomainId) pollers.PollerType {
	return &frontDoorCustomDomainWaitForApprovedPoller{
		client: client,
		id:     id,
	}
}

func (p frontDoorCustomDomainWaitForApprovedPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.Get(ctx, p.id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s while waiting for domain validation approval: %+v", p.id, err)
	}

	model := resp.Model
	if model == nil || model.Properties == nil {
		log.Printf("[DEBUG] AFD Custom Domain %s waiting for approval; properties are nil", p.id)
		return &frontDoorCustomDomainWaitForApprovedInProgress, nil
	}

	deploymentStatus := ""
	if model.Properties.DeploymentStatus != nil {
		deploymentStatus = string(*model.Properties.DeploymentStatus)
	}

	provisioningState := ""
	if model.Properties.ProvisioningState != nil {
		provisioningState = string(*model.Properties.ProvisioningState)
	}

	if model.Properties.DomainValidationState == nil {
		log.Printf("[DEBUG] AFD Custom Domain %s waiting for approval; domainValidationState is nil (deploymentStatus=%q provisioningState=%q)", p.id, deploymentStatus, provisioningState)
		return &frontDoorCustomDomainWaitForApprovedInProgress, nil
	}

	state := *model.Properties.DomainValidationState
	switch state {
	case afdcustomdomains.DomainValidationStateApproved:
		log.Printf("[DEBUG] AFD Custom Domain %s approved (deploymentStatus=%q provisioningState=%q)", p.id, deploymentStatus, provisioningState)
		return &frontDoorCustomDomainWaitForApprovedSuccess, nil
	case afdcustomdomains.DomainValidationStateRejected, afdcustomdomains.DomainValidationStateTimedOut, afdcustomdomains.DomainValidationStateInternalError:
		log.Printf("[DEBUG] AFD Custom Domain %s domain validation terminal state=%q (deploymentStatus=%q provisioningState=%q)", p.id, state, deploymentStatus, provisioningState)
		return nil, fmt.Errorf("domain validation for %s failed with `domainValidationState` `%s`", p.id, state)
	default:
		log.Printf("[DEBUG] AFD Custom Domain %s waiting for approval; domainValidationState=%q (deploymentStatus=%q provisioningState=%q)", p.id, state, deploymentStatus, provisioningState)
		return &frontDoorCustomDomainWaitForApprovedInProgress, nil
	}
}
