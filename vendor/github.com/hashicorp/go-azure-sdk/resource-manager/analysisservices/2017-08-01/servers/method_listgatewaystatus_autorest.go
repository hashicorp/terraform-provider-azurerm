package servers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListGatewayStatusOperationResponse struct {
	HttpResponse *http.Response
	Model        *GatewayListStatusLive
}

// ListGatewayStatus ...
func (c ServersClient) ListGatewayStatus(ctx context.Context, id ServerId) (result ListGatewayStatusOperationResponse, err error) {
	req, err := c.preparerForListGatewayStatus(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "ListGatewayStatus", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "ListGatewayStatus", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListGatewayStatus(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servers.ServersClient", "ListGatewayStatus", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListGatewayStatus prepares the ListGatewayStatus request.
func (c ServersClient) preparerForListGatewayStatus(ctx context.Context, id ServerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listGatewayStatus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListGatewayStatus handles the response to the ListGatewayStatus request. The method always
// closes the http.Response Body.
func (c ServersClient) responderForListGatewayStatus(resp *http.Response) (result ListGatewayStatusOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
