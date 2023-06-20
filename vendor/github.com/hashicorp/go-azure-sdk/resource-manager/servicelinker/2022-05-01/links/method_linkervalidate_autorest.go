package links

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

type LinkerValidateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// LinkerValidate ...
func (c LinksClient) LinkerValidate(ctx context.Context, id ScopedLinkerId) (result LinkerValidateOperationResponse, err error) {
	req, err := c.preparerForLinkerValidate(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.LinksClient", "LinkerValidate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForLinkerValidate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.LinksClient", "LinkerValidate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// LinkerValidateThenPoll performs LinkerValidate then polls until it's completed
func (c LinksClient) LinkerValidateThenPoll(ctx context.Context, id ScopedLinkerId) error {
	result, err := c.LinkerValidate(ctx, id)
	if err != nil {
		return fmt.Errorf("performing LinkerValidate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after LinkerValidate: %+v", err)
	}

	return nil
}

// preparerForLinkerValidate prepares the LinkerValidate request.
func (c LinksClient) preparerForLinkerValidate(ctx context.Context, id ScopedLinkerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/validateLinker", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForLinkerValidate sends the LinkerValidate request. The method will close the
// http.Response Body if it receives an error.
func (c LinksClient) senderForLinkerValidate(ctx context.Context, req *http.Request) (future LinkerValidateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
