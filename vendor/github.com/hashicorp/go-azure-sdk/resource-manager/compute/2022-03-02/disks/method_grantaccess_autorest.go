package disks

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

type GrantAccessOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// GrantAccess ...
func (c DisksClient) GrantAccess(ctx context.Context, id DiskId, input GrantAccessData) (result GrantAccessOperationResponse, err error) {
	req, err := c.preparerForGrantAccess(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disks.DisksClient", "GrantAccess", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGrantAccess(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "disks.DisksClient", "GrantAccess", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GrantAccessThenPoll performs GrantAccess then polls until it's completed
func (c DisksClient) GrantAccessThenPoll(ctx context.Context, id DiskId, input GrantAccessData) error {
	result, err := c.GrantAccess(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing GrantAccess: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GrantAccess: %+v", err)
	}

	return nil
}

// preparerForGrantAccess prepares the GrantAccess request.
func (c DisksClient) preparerForGrantAccess(ctx context.Context, id DiskId, input GrantAccessData) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/beginGetAccess", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForGrantAccess sends the GrantAccess request. The method will close the
// http.Response Body if it receives an error.
func (c DisksClient) senderForGrantAccess(ctx context.Context, req *http.Request) (future GrantAccessOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
