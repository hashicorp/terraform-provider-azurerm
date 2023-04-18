package iotconnectors

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

type IotConnectorFhirDestinationCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// IotConnectorFhirDestinationCreateOrUpdate ...
func (c IotConnectorsClient) IotConnectorFhirDestinationCreateOrUpdate(ctx context.Context, id FhirDestinationId, input IotFhirDestination) (result IotConnectorFhirDestinationCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForIotConnectorFhirDestinationCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "IotConnectorFhirDestinationCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForIotConnectorFhirDestinationCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "IotConnectorFhirDestinationCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// IotConnectorFhirDestinationCreateOrUpdateThenPoll performs IotConnectorFhirDestinationCreateOrUpdate then polls until it's completed
func (c IotConnectorsClient) IotConnectorFhirDestinationCreateOrUpdateThenPoll(ctx context.Context, id FhirDestinationId, input IotFhirDestination) error {
	result, err := c.IotConnectorFhirDestinationCreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing IotConnectorFhirDestinationCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after IotConnectorFhirDestinationCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForIotConnectorFhirDestinationCreateOrUpdate prepares the IotConnectorFhirDestinationCreateOrUpdate request.
func (c IotConnectorsClient) preparerForIotConnectorFhirDestinationCreateOrUpdate(ctx context.Context, id FhirDestinationId, input IotFhirDestination) (*http.Request, error) {
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

// senderForIotConnectorFhirDestinationCreateOrUpdate sends the IotConnectorFhirDestinationCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c IotConnectorsClient) senderForIotConnectorFhirDestinationCreateOrUpdate(ctx context.Context, req *http.Request) (future IotConnectorFhirDestinationCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
