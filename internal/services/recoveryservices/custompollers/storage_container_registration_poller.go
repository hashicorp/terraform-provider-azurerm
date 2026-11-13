// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectioncontainers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const maxRegistrationRetries = 3

// Register can fail asynchronously with CloudInternalError for newly created storage accounts.
// Resubmitting the same request after a delay can succeed.
type storageContainerRegistrationPoller struct {
	client          *protectioncontainers.ProtectionContainersClient
	operationClient *backup.OperationStatusesClient
	id              protectioncontainers.ProtectionContainerId
	parameters      protectioncontainers.ProtectionContainerResource
	operationID     string
	retries         int
	state           *pluginsdk.StateChangeConf

	// WaitForStateContext can return on cancellation while a refresh is finishing.
	lastErrorMu sync.Mutex
	lastError   error
}

func NewStorageContainerRegistrationPoller(client *protectioncontainers.ProtectionContainersClient, operationClient *backup.OperationStatusesClient, id protectioncontainers.ProtectionContainerId, parameters protectioncontainers.ProtectionContainerResource, timeout time.Duration) *storageContainerRegistrationPoller {
	return &storageContainerRegistrationPoller{
		client:          client,
		operationClient: operationClient,
		id:              id,
		parameters:      parameters,
		state: &pluginsdk.StateChangeConf{
			Pending:                   []string{"InProgress", "Retrying"},
			Target:                    []string{"Succeeded"},
			Timeout:                   timeout,
			PollInterval:              10 * time.Second,
			ContinuousTargetOccurence: 5, // Retain the buffer before dependent file share operations.
		},
	}
}

func (p *storageContainerRegistrationPoller) WaitForCompletion(ctx context.Context) error {
	p.state.Refresh = func() (any, string, error) {
		return p.refresh(ctx)
	}
	if _, err := p.state.WaitForStateContext(ctx); err != nil {
		p.lastErrorMu.Lock()
		defer p.lastErrorMu.Unlock()
		if p.lastError != nil {
			return fmt.Errorf("registering %s: %w (last operation error: %s)", p.id, err, p.lastError)
		}
		return fmt.Errorf("registering %s: %w", p.id, err)
	}
	return nil
}

func (p *storageContainerRegistrationPoller) refresh(ctx context.Context) (any, string, error) {
	if err := ctx.Err(); err != nil {
		return nil, "", err
	}

	if p.operationID == "" {
		resp, err := p.client.Register(ctx, p.id, p.parameters)
		if err != nil {
			return nil, "", err
		}
		if resp.HttpResponse == nil {
			return nil, "", fmt.Errorf("no response returned from registration")
		}
		locationURL, err := resp.HttpResponse.Location()
		if err != nil || locationURL == nil {
			return nil, "", fmt.Errorf("unable to determine registration operation URL: Location header missing or invalid")
		}
		// Work around https://github.com/Azure/azure-sdk-for-go/issues/2824.
		parsedLocation, err := azure.ParseAzureResourceID(strings.Replace(locationURL.Path, "/Subscriptions/", "/subscriptions/", 1))
		if err != nil {
			return nil, "", err
		}
		p.operationID = parsedLocation.Path["operationResults"]
		if p.operationID == "" {
			return nil, "", fmt.Errorf("registration operation URL is missing operationResults")
		}
		p.state.PollInterval = 10 * time.Second
		return resp, "InProgress", nil
	}

	resp, err := p.operationClient.Get(ctx, p.id.VaultName, p.id.ResourceGroupName, p.operationID)
	if err != nil {
		return nil, "", fmt.Errorf("reading registration operation %q: %w", p.operationID, err)
	}
	if resp.Error != nil || (resp.Status != backup.OperationStatusValuesInProgress && resp.Status != backup.OperationStatusValuesSucceeded) {
		err := fmt.Errorf("operation %q returned status %q", p.operationID, resp.Status)
		if resp.Error != nil {
			err = fmt.Errorf("%w: %s: %s", err, pointer.From(resp.Error.Code), pointer.From(resp.Error.Message))
		}
		if resp.Status == backup.OperationStatusValuesFailed && resp.Error != nil && pointer.From(resp.Error.Code) == "CloudInternalError" {
			p.lastErrorMu.Lock()
			p.lastError = err
			p.lastErrorMu.Unlock()
			if p.retries >= maxRegistrationRetries {
				return resp, "", fmt.Errorf("registration failed after %d retries: %w", p.retries, err)
			}
			p.retries++
			// A failed operation cannot recover by polling. Resubmit registration after a
			// backoff, within the same create timeout, and track the new operation ID.
			log.Printf("[DEBUG] Registration of %s failed: %s; retry %d/%d in one minute", p.id, err, p.retries, maxRegistrationRetries)
			p.operationID = ""
			p.state.PollInterval = time.Minute
			return resp, "Retrying", nil
		}
		return resp, "", err
	}
	return resp, string(resp.Status), nil
}
