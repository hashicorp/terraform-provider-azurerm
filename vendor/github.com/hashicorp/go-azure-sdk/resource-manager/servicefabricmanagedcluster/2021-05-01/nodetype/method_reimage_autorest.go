package nodetype

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

type ReimageOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// Reimage ...
func (c NodeTypeClient) Reimage(ctx context.Context, id NodeTypeId, input NodeTypeActionParameters) (result ReimageOperationResponse, err error) {
	req, err := c.preparerForReimage(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nodetype.NodeTypeClient", "Reimage", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForReimage(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nodetype.NodeTypeClient", "Reimage", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// ReimageThenPoll performs Reimage then polls until it's completed
func (c NodeTypeClient) ReimageThenPoll(ctx context.Context, id NodeTypeId, input NodeTypeActionParameters) error {
	result, err := c.Reimage(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing Reimage: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after Reimage: %+v", err)
	}

	return nil
}

// preparerForReimage prepares the Reimage request.
func (c NodeTypeClient) preparerForReimage(ctx context.Context, id NodeTypeId, input NodeTypeActionParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/reimage", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForReimage sends the Reimage request. The method will close the
// http.Response Body if it receives an error.
func (c NodeTypeClient) senderForReimage(ctx context.Context, req *http.Request) (future ReimageOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
