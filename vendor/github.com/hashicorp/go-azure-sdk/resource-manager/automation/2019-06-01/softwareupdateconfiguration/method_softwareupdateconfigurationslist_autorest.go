package softwareupdateconfiguration

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationsListOperationResponse struct {
	HttpResponse *http.Response
	Model        *SoftwareUpdateConfigurationListResult
}

type SoftwareUpdateConfigurationsListOperationOptions struct {
	ClientRequestId *string
	Filter          *string
}

func DefaultSoftwareUpdateConfigurationsListOperationOptions() SoftwareUpdateConfigurationsListOperationOptions {
	return SoftwareUpdateConfigurationsListOperationOptions{}
}

func (o SoftwareUpdateConfigurationsListOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ClientRequestId != nil {
		out["clientRequestId"] = *o.ClientRequestId
	}

	return out
}

func (o SoftwareUpdateConfigurationsListOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// SoftwareUpdateConfigurationsList ...
func (c SoftwareUpdateConfigurationClient) SoftwareUpdateConfigurationsList(ctx context.Context, id AutomationAccountId, options SoftwareUpdateConfigurationsListOperationOptions) (result SoftwareUpdateConfigurationsListOperationResponse, err error) {
	req, err := c.preparerForSoftwareUpdateConfigurationsList(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsList", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsList", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForSoftwareUpdateConfigurationsList(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "softwareupdateconfiguration.SoftwareUpdateConfigurationClient", "SoftwareUpdateConfigurationsList", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForSoftwareUpdateConfigurationsList prepares the SoftwareUpdateConfigurationsList request.
func (c SoftwareUpdateConfigurationClient) preparerForSoftwareUpdateConfigurationsList(ctx context.Context, id AutomationAccountId, options SoftwareUpdateConfigurationsListOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/softwareUpdateConfigurations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForSoftwareUpdateConfigurationsList handles the response to the SoftwareUpdateConfigurationsList request. The method always
// closes the http.Response Body.
func (c SoftwareUpdateConfigurationClient) responderForSoftwareUpdateConfigurationsList(resp *http.Response) (result SoftwareUpdateConfigurationsListOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
