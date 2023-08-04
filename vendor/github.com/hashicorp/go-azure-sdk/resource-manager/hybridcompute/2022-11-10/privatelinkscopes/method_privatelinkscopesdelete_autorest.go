package privatelinkscopes

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

type PrivateLinkScopesDeleteOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// PrivateLinkScopesDelete ...
func (c PrivateLinkScopesClient) PrivateLinkScopesDelete(ctx context.Context, id ProviderPrivateLinkScopeId) (result PrivateLinkScopesDeleteOperationResponse, err error) {
	req, err := c.preparerForPrivateLinkScopesDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopes.PrivateLinkScopesClient", "PrivateLinkScopesDelete", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForPrivateLinkScopesDelete(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "privatelinkscopes.PrivateLinkScopesClient", "PrivateLinkScopesDelete", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// PrivateLinkScopesDeleteThenPoll performs PrivateLinkScopesDelete then polls until it's completed
func (c PrivateLinkScopesClient) PrivateLinkScopesDeleteThenPoll(ctx context.Context, id ProviderPrivateLinkScopeId) error {
	result, err := c.PrivateLinkScopesDelete(ctx, id)
	if err != nil {
		return fmt.Errorf("performing PrivateLinkScopesDelete: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after PrivateLinkScopesDelete: %+v", err)
	}

	return nil
}

// preparerForPrivateLinkScopesDelete prepares the PrivateLinkScopesDelete request.
func (c PrivateLinkScopesClient) preparerForPrivateLinkScopesDelete(ctx context.Context, id ProviderPrivateLinkScopeId) (*http.Request, error) {
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

// senderForPrivateLinkScopesDelete sends the PrivateLinkScopesDelete request. The method will close the
// http.Response Body if it receives an error.
func (c PrivateLinkScopesClient) senderForPrivateLinkScopesDelete(ctx context.Context, req *http.Request) (future PrivateLinkScopesDeleteOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
