package softwareupdateconfiguration

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationsCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *SoftwareUpdateConfiguration
}

type SoftwareUpdateConfigurationsCreateOperationOptions struct {
	ClientRequestId *string
}

func DefaultSoftwareUpdateConfigurationsCreateOperationOptions() SoftwareUpdateConfigurationsCreateOperationOptions {
	return SoftwareUpdateConfigurationsCreateOperationOptions{}
}

func (o SoftwareUpdateConfigurationsCreateOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ClientRequestId != nil {
		out["clientRequestId"] = *o.ClientRequestId
	}

	return out
}

func (o SoftwareUpdateConfigurationsCreateOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// SoftwareUpdateConfigurationsCreate ...
func (c SoftwareUpdateConfigurationClient) SoftwareUpdateConfigurationsCreate(ctx context.Context, id SoftwareUpdateConfigurationId, input SoftwareUpdateConfiguration, options SoftwareUpdateConfigurationsCreateOperationOptions) (result SoftwareUpdateConfigurationsCreateOperationResponse, err error) {
	req, err := c.preparerForSoftwareUpdateConfigurationsCreate(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSoftwareUpdateConfigurationsCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSoftwareUpdateConfigurationsCreate prepares the SoftwareUpdateConfigurationsCreate request.
func (c SoftwareUpdateConfigurationClient) preparerForSoftwareUpdateConfigurationsCreate(ctx context.Context, id SoftwareUpdateConfigurationId, input SoftwareUpdateConfiguration, options SoftwareUpdateConfigurationsCreateOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSoftwareUpdateConfigurationsCreate handles the response to the SoftwareUpdateConfigurationsCreate request. The method always
// closes the http.Response Body.
func (c SoftwareUpdateConfigurationClient) responderForSoftwareUpdateConfigurationsCreate(resp *http.Response) (result SoftwareUpdateConfigurationsCreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
