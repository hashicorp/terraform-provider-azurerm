package softwareupdateconfiguration

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationsGetByNameOperationResponse struct {
	HttpResponse *http.Response
	Model        *SoftwareUpdateConfiguration
}

type SoftwareUpdateConfigurationsGetByNameOperationOptions struct {
	ClientRequestId *string
}

func DefaultSoftwareUpdateConfigurationsGetByNameOperationOptions() SoftwareUpdateConfigurationsGetByNameOperationOptions {
	return SoftwareUpdateConfigurationsGetByNameOperationOptions{}
}

func (o SoftwareUpdateConfigurationsGetByNameOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ClientRequestId != nil {
		out["clientRequestId"] = *o.ClientRequestId
	}

	return out
}

func (o SoftwareUpdateConfigurationsGetByNameOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// SoftwareUpdateConfigurationsGetByName ...
func (c SoftwareUpdateConfigurationClient) SoftwareUpdateConfigurationsGetByName(ctx context.Context, id SoftwareUpdateConfigurationId, options SoftwareUpdateConfigurationsGetByNameOperationOptions) (result SoftwareUpdateConfigurationsGetByNameOperationResponse, err error) {
	req, err := c.preparerForSoftwareUpdateConfigurationsGetByName(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsGetByName", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsGetByName", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSoftwareUpdateConfigurationsGetByName(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsGetByName", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSoftwareUpdateConfigurationsGetByName prepares the SoftwareUpdateConfigurationsGetByName request.
func (c SoftwareUpdateConfigurationClient) preparerForSoftwareUpdateConfigurationsGetByName(ctx context.Context, id SoftwareUpdateConfigurationId, options SoftwareUpdateConfigurationsGetByNameOperationOptions) (*http.Request, error) {
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

// responderForSoftwareUpdateConfigurationsGetByName handles the response to the SoftwareUpdateConfigurationsGetByName request. The method always
// closes the http.Response Body.
func (c SoftwareUpdateConfigurationClient) responderForSoftwareUpdateConfigurationsGetByName(resp *http.Response) (result SoftwareUpdateConfigurationsGetByNameOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
