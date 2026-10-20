package protectioncontainers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RefreshOperationResponse struct {
	HttpResponse *http.Response
}

type RefreshOperationOptions struct {
	Filter *string
}

func DefaultRefreshOperationOptions() RefreshOperationOptions {
	return RefreshOperationOptions{}
}

func (o RefreshOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o RefreshOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// Refresh ...
func (c ProtectionContainersClient) Refresh(ctx context.Context, id BackupFabricId, options RefreshOperationOptions) (result RefreshOperationResponse, err error) {
	req, err := c.preparerForRefresh(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Refresh", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Refresh", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRefresh(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protectioncontainers.ProtectionContainersClient", "Refresh", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRefresh prepares the Refresh request.
func (c ProtectionContainersClient) preparerForRefresh(ctx context.Context, id BackupFabricId, options RefreshOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/refreshContainers", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRefresh handles the response to the Refresh request. The method always
// closes the http.Response Body.
func (c ProtectionContainersClient) responderForRefresh(resp *http.Response) (result RefreshOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
