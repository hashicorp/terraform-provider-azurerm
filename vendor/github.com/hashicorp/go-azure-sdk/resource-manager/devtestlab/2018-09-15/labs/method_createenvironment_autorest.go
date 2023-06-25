package labs

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

type CreateEnvironmentOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// CreateEnvironment ...
func (c LabsClient) CreateEnvironment(ctx context.Context, id LabId, input LabVirtualMachineCreationParameter) (result CreateEnvironmentOperationResponse, err error) {
	req, err := c.preparerForCreateEnvironment(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "CreateEnvironment", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCreateEnvironment(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "labs.LabsClient", "CreateEnvironment", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CreateEnvironmentThenPoll performs CreateEnvironment then polls until it's completed
func (c LabsClient) CreateEnvironmentThenPoll(ctx context.Context, id LabId, input LabVirtualMachineCreationParameter) error {
	result, err := c.CreateEnvironment(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CreateEnvironment: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CreateEnvironment: %+v", err)
	}

	return nil
}

// preparerForCreateEnvironment prepares the CreateEnvironment request.
func (c LabsClient) preparerForCreateEnvironment(ctx context.Context, id LabId, input LabVirtualMachineCreationParameter) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/createEnvironment", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCreateEnvironment sends the CreateEnvironment request. The method will close the
// http.Response Body if it receives an error.
func (c LabsClient) senderForCreateEnvironment(ctx context.Context, req *http.Request) (future CreateEnvironmentOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
