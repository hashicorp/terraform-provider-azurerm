// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

var _ pollers.PollerType = &frontDoorCustomDomainWaitForApprovedPoller{}

type frontDoorCustomDomainWaitForApprovedPoller struct {
	client *afddomains.AFDDomainsClient
	id     afddomains.CustomDomainId
}

func NewFrontDoorCustomDomainWaitForApprovedPoller(client *afddomains.AFDDomainsClient, id afddomains.CustomDomainId) pollers.PollerType {
	return &frontDoorCustomDomainWaitForApprovedPoller{
		client: client,
		id:     id,
	}
}

func (p frontDoorCustomDomainWaitForApprovedPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	resp, err := p.client.AFDCustomDomainsGet(ctx, p.id)
	pollInterval := 30 * time.Second
	if resp.HttpResponse != nil {
		if retryAfter := resp.HttpResponse.Header.Get("Retry-After"); retryAfter != "" {
			if parsedSeconds, parseErr := strconv.Atoi(retryAfter); parseErr == nil {
				pollInterval = time.Duration(parsedSeconds) * time.Second
			}
		}
	}

	if err != nil {
		if resp.HttpResponse != nil && resp.HttpResponse.StatusCode == http.StatusTooManyRequests {
			log.Printf("[DEBUG] 429 Too Many Requests retrieving %s. Retrying after %s", p.id, pollInterval)
			return &pollers.PollResult{
				PollInterval: pollInterval,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}
		return nil, fmt.Errorf("retrieving %s while waiting for domain validation approval: %+v", p.id, err)
	}

	model := resp.Model
	if model == nil || model.Properties == nil {
		log.Printf("[DEBUG] AFD Custom Domain %s waiting for approval; properties are nil", p.id)
		return &pollers.PollResult{
			PollInterval: pollInterval,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	deploymentStatus := pointer.From(model.Properties.DeploymentStatus)
	provisioningState := pointer.From(model.Properties.ProvisioningState)

	if provisioningState == afddomains.AfdProvisioningStateFailed {
		log.Printf("[DEBUG] AFD Custom Domain %s provisioning failed (deploymentStatus=%q provisioningState=%q)", p.id, string(deploymentStatus), string(provisioningState))
		return nil, fmt.Errorf("provisioning for %s failed with `provisioningState` `%s`", p.id, provisioningState)
	}

	if model.Properties.DomainValidationState == nil {
		log.Printf("[DEBUG] AFD Custom Domain %s waiting for approval; domainValidationState is nil (deploymentStatus=%q provisioningState=%q)", p.id, deploymentStatus, provisioningState)
		return &pollers.PollResult{
			PollInterval: pollInterval,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}

	state := *model.Properties.DomainValidationState
	switch state {
	case afddomains.DomainValidationStateApproved:
		if deploymentStatus != afddomains.DeploymentStatusSucceeded {
			log.Printf("[DEBUG] AFD Custom Domain %s validation approved but deployment not succeeded yet (deploymentStatus=%q provisioningState=%q)", p.id, string(deploymentStatus), string(provisioningState))
			return &pollers.PollResult{
				PollInterval: pollInterval,
				Status:       pollers.PollingStatusInProgress,
			}, nil
		}

		log.Printf("[DEBUG] AFD Custom Domain %s approved and deployed (deploymentStatus=%q provisioningState=%q)", p.id, string(deploymentStatus), string(provisioningState))
		return &pollers.PollResult{
			PollInterval: pollInterval,
			Status:       pollers.PollingStatusSucceeded,
		}, nil
	case afddomains.DomainValidationStateRejected, afddomains.DomainValidationStateTimedOut, afddomains.DomainValidationStateInternalError:
		log.Printf("[DEBUG] AFD Custom Domain %s domain validation terminal state=%q (deploymentStatus=%q provisioningState=%q)", p.id, state, string(deploymentStatus), string(provisioningState))
		return nil, fmt.Errorf("domain validation for %s failed with `domainValidationState` `%s`", p.id, state)
	default:
		log.Printf("[DEBUG] AFD Custom Domain %s waiting for approval; domainValidationState=%q (deploymentStatus=%q provisioningState=%q)", p.id, state, string(deploymentStatus), string(provisioningState))
		return &pollers.PollResult{
			PollInterval: pollInterval,
			Status:       pollers.PollingStatusInProgress,
		}, nil
	}
}
