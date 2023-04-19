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

type IotConnectorFhirDestinationDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// IotConnectorFhirDestinationDelete ...
func (c IotConnectorsClient) IotConnectorFhirDestinationDelete(ctx context.Context, id FhirDestinationId) (result IotConnectorFhirDestinationDeleteOperationResponse, err error) {
	req, err := c.preparerForIotConnectorFhirDestinationDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "IotConnectorFhirDestinationDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForIotConnectorFhirDestinationDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "iotconnectors.IotConnectorsClient", "IotConnectorFhirDestinationDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// IotConnectorFhirDestinationDeleteThenPoll performs IotConnectorFhirDestinationDelete then polls until it's completed
func (c IotConnectorsClient) IotConnectorFhirDestinationDeleteThenPoll(ctx context.Context, id FhirDestinationId) error {
	result, err := c.IotConnectorFhirDestinationDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing IotConnectorFhirDestinationDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after IotConnectorFhirDestinationDelete: %+v", err)
	}

	return nil
}

// preparerForIotConnectorFhirDestinationDelete prepares the IotConnectorFhirDestinationDelete request.
func (c IotConnectorsClient) preparerForIotConnectorFhirDestinationDelete(ctx context.Context, id FhirDestinationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForIotConnectorFhirDestinationDelete sends the IotConnectorFhirDestinationDelete request. The method will close the
// http.Response Body if it receives an error.
func (c IotConnectorsClient) senderForIotConnectorFhirDestinationDelete(ctx context.Context, req *http.Request) (future IotConnectorFhirDestinationDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
