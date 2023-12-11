package softwareupdateconfiguration

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetByNameOperationResponse struct {
	HttpResponse *http.Response
	Model        *SoftwareUpdateConfiguration
}

type GetByNameOperationOptions struct {
	ClientRequestId *string
}

func DefaultGetByNameOperationOptions() GetByNameOperationOptions {
	return GetByNameOperationOptions{}
}

func (o GetByNameOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ClientRequestId != nil {
		out["clientRequestId"] = *o.ClientRequestId
	}

	return out
}

func (o GetByNameOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// GetByName ...
func (c SoftwareUpdateConfigurationClient) GetByName(ctx context.Context, id SoftwareUpdateConfigurationId, options GetByNameOperationOptions) (result GetByNameOperationResponse, err error) {
	req, err := c.preparerForGetByName(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "GetByName", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "GetByName", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetByName(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "GetByName", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetByName prepares the GetByName request.
func (c SoftwareUpdateConfigurationClient) preparerForGetByName(ctx context.Context, id SoftwareUpdateConfigurationId, options GetByNameOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetByName handles the response to the GetByName request. The method always
// closes the http.Response Body.
func (c SoftwareUpdateConfigurationClient) responderForGetByName(resp *http.Response) (result GetByNameOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
