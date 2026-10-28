// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custompollers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
)

type managedEnvironmentReader interface {
	Get(context.Context, managedenvironments.ManagedEnvironmentId) (managedenvironments.GetOperationResponse, error)
}

type operationPoller interface {
	PollUntilDone(context.Context) error
}

// WaitForManagedEnvironmentDeletion waits for the delete operation to complete, returning any error the operation
// itself reports. Only when the SDK's poller cannot decode the operation's terminal response, which this API returns
// without a Content-Type, is the environment polled until it no longer exists instead.
func WaitForManagedEnvironmentDeletion(ctx context.Context, operation operationPoller, client managedEnvironmentReader, id managedenvironments.ManagedEnvironmentId) error {
	err := operation.PollUntilDone(ctx)
	if err == nil {
		return nil
	}

	if !isUndecodableOperationResponse(err) {
		return err
	}

	poller := NewManagedEnvironmentDeletePoller(client, id)
	if pollErr := poller.PollUntilDone(ctx); pollErr != nil {
		return fmt.Errorf("%+v (after the delete operation returned a response that could not be decoded: %+v)", pollErr, err)
	}
	return nil
}

// isUndecodableOperationResponse reports whether the SDK's long-running-operation poller rejected the terminal
// response because it has no supported Content-Type. The SDK returns this as an untyped error, so the message is
// matched rather than a type.
func isUndecodableOperationResponse(err error) bool {
	message := err.Error()
	return strings.Contains(message, "polling support for the Content-Type") && strings.Contains(message, "was not implemented")
}

type managedEnvironmentDeletePoller struct {
	client managedEnvironmentReader
	id     managedenvironments.ManagedEnvironmentId
}

func NewManagedEnvironmentDeletePoller(client managedEnvironmentReader, id managedenvironments.ManagedEnvironmentId) pollers.Poller {
	return pollers.NewPoller(&managedEnvironmentDeletePoller{client: client, id: id}, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
}

func (p *managedEnvironmentDeletePoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	result, err := p.client.Get(ctx, p.id)
	if response.WasNotFound(result.HttpResponse) {
		return &pollers.PollResult{Status: pollers.PollingStatusSucceeded}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("checking deletion of %s: %+v", p.id, err)
	}
	if result.HttpResponse == nil {
		return nil, fmt.Errorf("checking deletion of %s: response was nil", p.id)
	}
	// A failed delete leaves the environment in place, so surface that rather than waiting for the timeout.
	if result.Model != nil && result.Model.Properties != nil && result.Model.Properties.ProvisioningState != nil {
		state := string(*result.Model.Properties.ProvisioningState)
		if strings.EqualFold(state, string(managedenvironments.EnvironmentProvisioningStateFailed)) || strings.EqualFold(state, string(managedenvironments.EnvironmentProvisioningStateCanceled)) {
			return nil, pollers.PollingFailedError{
				HttpResponse: &client.Response{Response: result.HttpResponse},
				Message:      fmt.Sprintf("deleting %s: the environment was left in the provisioning state %q", p.id, state),
			}
		}
	}

	return &pollers.PollResult{
		Status:       pollers.PollingStatusInProgress,
		PollInterval: 10 * time.Second,
	}, nil
}
