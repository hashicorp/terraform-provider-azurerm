package configurations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateOnCoordinatorOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// UpdateOnCoordinator ...
func (c ConfigurationsClient) UpdateOnCoordinator(ctx context.Context, id CoordinatorConfigurationId, input ServerConfiguration) (result UpdateOnCoordinatorOperationResponse, err error) {
	req, err := c.preparerForUpdateOnCoordinator(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "UpdateOnCoordinator", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForUpdateOnCoordinator(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "configurations.ConfigurationsClient", "UpdateOnCoordinator", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// UpdateOnCoordinatorThenPoll performs UpdateOnCoordinator then polls until it's completed
func (c ConfigurationsClient) UpdateOnCoordinatorThenPoll(ctx context.Context, id CoordinatorConfigurationId, input ServerConfiguration) error {
	result, err := c.UpdateOnCoordinator(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing UpdateOnCoordinator: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after UpdateOnCoordinator: %+v", err)
	}

	return nil
}

// preparerForUpdateOnCoordinator prepares the UpdateOnCoordinator request.
func (c ConfigurationsClient) preparerForUpdateOnCoordinator(ctx context.Context, id CoordinatorConfigurationId, input ServerConfiguration) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForUpdateOnCoordinator sends the UpdateOnCoordinator request. The method will close the
// http.Response Body if it receives an error.
func (c ConfigurationsClient) senderForUpdateOnCoordinator(ctx context.Context, req *http.Request) (future UpdateOnCoordinatorOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
