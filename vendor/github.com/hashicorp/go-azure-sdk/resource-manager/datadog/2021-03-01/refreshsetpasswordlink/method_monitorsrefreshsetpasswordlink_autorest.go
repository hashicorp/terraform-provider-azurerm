package refreshsetpasswordlink

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsRefreshSetPasswordLinkOperationResponse struct {
	HttpResponse *http.Response
	Model        *DatadogSetPasswordLink
}

// MonitorsRefreshSetPasswordLink ...
func (c RefreshSetPasswordLinkClient) MonitorsRefreshSetPasswordLink(ctx context.Context, id MonitorId) (result MonitorsRefreshSetPasswordLinkOperationResponse, err error) {
	req, err := c.preparerForMonitorsRefreshSetPasswordLink(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "refreshsetpasswordlink.RefreshSetPasswordLinkClient", "MonitorsRefreshSetPasswordLink", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "refreshsetpasswordlink.RefreshSetPasswordLinkClient", "MonitorsRefreshSetPasswordLink", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMonitorsRefreshSetPasswordLink(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "refreshsetpasswordlink.RefreshSetPasswordLinkClient", "MonitorsRefreshSetPasswordLink", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMonitorsRefreshSetPasswordLink prepares the MonitorsRefreshSetPasswordLink request.
func (c RefreshSetPasswordLinkClient) preparerForMonitorsRefreshSetPasswordLink(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/refreshSetPasswordLink", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMonitorsRefreshSetPasswordLink handles the response to the MonitorsRefreshSetPasswordLink request. The method always
// closes the http.Response Body.
func (c RefreshSetPasswordLinkClient) responderForMonitorsRefreshSetPasswordLink(resp *http.Response) (result MonitorsRefreshSetPasswordLinkOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
