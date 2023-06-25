package virtualmachines

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetrieveBootDiagnosticsDataOperationResponse struct {
	HttpResponse *http.Response
	Model        *RetrieveBootDiagnosticsDataResult
}

type RetrieveBootDiagnosticsDataOperationOptions struct {
	SasUriExpirationTimeInMinutes *int64
}

func DefaultRetrieveBootDiagnosticsDataOperationOptions() RetrieveBootDiagnosticsDataOperationOptions {
	return RetrieveBootDiagnosticsDataOperationOptions{}
}

func (o RetrieveBootDiagnosticsDataOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o RetrieveBootDiagnosticsDataOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.SasUriExpirationTimeInMinutes != nil {
		out["sasUriExpirationTimeInMinutes"] = *o.SasUriExpirationTimeInMinutes
	}

	return out
}

// RetrieveBootDiagnosticsData ...
func (c VirtualMachinesClient) RetrieveBootDiagnosticsData(ctx context.Context, id VirtualMachineId, options RetrieveBootDiagnosticsDataOperationOptions) (result RetrieveBootDiagnosticsDataOperationResponse, err error) {
	req, err := c.preparerForRetrieveBootDiagnosticsData(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "RetrieveBootDiagnosticsData", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "RetrieveBootDiagnosticsData", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRetrieveBootDiagnosticsData(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "virtualmachines.VirtualMachinesClient", "RetrieveBootDiagnosticsData", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRetrieveBootDiagnosticsData prepares the RetrieveBootDiagnosticsData request.
func (c VirtualMachinesClient) preparerForRetrieveBootDiagnosticsData(ctx context.Context, id VirtualMachineId, options RetrieveBootDiagnosticsDataOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/retrieveBootDiagnosticsData", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRetrieveBootDiagnosticsData handles the response to the RetrieveBootDiagnosticsData request. The method always
// closes the http.Response Body.
func (c VirtualMachinesClient) responderForRetrieveBootDiagnosticsData(resp *http.Response) (result RetrieveBootDiagnosticsDataOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
