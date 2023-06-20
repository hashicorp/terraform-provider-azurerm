package replicationfabrics

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

type ReassociateGatewayOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// ReassociateGateway ...
func (c ReplicationFabricsClient) ReassociateGateway(ctx context.Context, id ReplicationFabricId, input FailoverProcessServerRequest) (result ReassociateGatewayOperationResponse, err error) {
	req, err := c.preparerForReassociateGateway(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationfabrics.ReplicationFabricsClient", "ReassociateGateway", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForReassociateGateway(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "replicationfabrics.ReplicationFabricsClient", "ReassociateGateway", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ReassociateGatewayThenPoll performs ReassociateGateway then polls until it's completed
func (c ReplicationFabricsClient) ReassociateGatewayThenPoll(ctx context.Context, id ReplicationFabricId, input FailoverProcessServerRequest) error {
	result, err := c.ReassociateGateway(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing ReassociateGateway: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after ReassociateGateway: %+v", err)
	}

	return nil
}

// preparerForReassociateGateway prepares the ReassociateGateway request.
func (c ReplicationFabricsClient) preparerForReassociateGateway(ctx context.Context, id ReplicationFabricId, input FailoverProcessServerRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/reassociateGateway", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForReassociateGateway sends the ReassociateGateway request. The method will close the
// http.Response Body if it receives an error.
func (c ReplicationFabricsClient) senderForReassociateGateway(ctx context.Context, req *http.Request) (future ReassociateGatewayOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
