package softwareupdateconfiguration

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationsDeleteOperationResponse struct {
	HttpResponse *http.Response
}

type SoftwareUpdateConfigurationsDeleteOperationOptions struct {
	ClientRequestId *string
}

func DefaultSoftwareUpdateConfigurationsDeleteOperationOptions() SoftwareUpdateConfigurationsDeleteOperationOptions {
	return SoftwareUpdateConfigurationsDeleteOperationOptions{}
}

func (o SoftwareUpdateConfigurationsDeleteOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ClientRequestId != nil {
		out["clientRequestId"] = *o.ClientRequestId
	}

	return out
}

func (o SoftwareUpdateConfigurationsDeleteOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// SoftwareUpdateConfigurationsDelete ...
func (c SoftwareUpdateConfigurationClient) SoftwareUpdateConfigurationsDelete(ctx context.Context, id SoftwareUpdateConfigurationId, options SoftwareUpdateConfigurationsDeleteOperationOptions) (result SoftwareUpdateConfigurationsDeleteOperationResponse, err error) {
	req, err := c.preparerForSoftwareUpdateConfigurationsDelete(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSoftwareUpdateConfigurationsDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSoftwareUpdateConfigurationsDelete prepares the SoftwareUpdateConfigurationsDelete request.
func (c SoftwareUpdateConfigurationClient) preparerForSoftwareUpdateConfigurationsDelete(ctx context.Context, id SoftwareUpdateConfigurationId, options SoftwareUpdateConfigurationsDeleteOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSoftwareUpdateConfigurationsDelete handles the response to the SoftwareUpdateConfigurationsDelete request. The method always
// closes the http.Response Body.
func (c SoftwareUpdateConfigurationClient) responderForSoftwareUpdateConfigurationsDelete(resp *http.Response) (result SoftwareUpdateConfigurationsDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent, http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
