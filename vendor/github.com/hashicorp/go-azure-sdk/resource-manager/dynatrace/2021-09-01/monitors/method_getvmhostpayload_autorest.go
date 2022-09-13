package monitors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetVMHostPayloadOperationResponse struct {
	HttpResponse *http.Response
	Model        *VMExtensionPayload
}

// GetVMHostPayload ...
func (c MonitorsClient) GetVMHostPayload(ctx context.Context, id MonitorId) (result GetVMHostPayloadOperationResponse, err error) {
	req, err := c.preparerForGetVMHostPayload(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "GetVMHostPayload", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "GetVMHostPayload", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetVMHostPayload(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "monitors.MonitorsClient", "GetVMHostPayload", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetVMHostPayload prepares the GetVMHostPayload request.
func (c MonitorsClient) preparerForGetVMHostPayload(ctx context.Context, id MonitorId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/getVMHostPayload", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetVMHostPayload handles the response to the GetVMHostPayload request. The method always
// closes the http.Response Body.
func (c MonitorsClient) responderForGetVMHostPayload(resp *http.Response) (result GetVMHostPayloadOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
