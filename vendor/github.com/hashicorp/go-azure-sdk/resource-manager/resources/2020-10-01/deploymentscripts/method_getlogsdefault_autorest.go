package deploymentscripts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetLogsDefaultOperationResponse struct {
	HttpResponse *http.Response
	Model        *ScriptLog
}

type GetLogsDefaultOperationOptions struct {
	Tail *int64
}

func DefaultGetLogsDefaultOperationOptions() GetLogsDefaultOperationOptions {
	return GetLogsDefaultOperationOptions{}
}

func (o GetLogsDefaultOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o GetLogsDefaultOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Tail != nil {
		out["tail"] = *o.Tail
	}

	return out
}

// GetLogsDefault ...
func (c DeploymentScriptsClient) GetLogsDefault(ctx context.Context, id DeploymentScriptId, options GetLogsDefaultOperationOptions) (result GetLogsDefaultOperationResponse, err error) {
	req, err := c.preparerForGetLogsDefault(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentscripts.DeploymentScriptsClient", "GetLogsDefault", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentscripts.DeploymentScriptsClient", "GetLogsDefault", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetLogsDefault(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentscripts.DeploymentScriptsClient", "GetLogsDefault", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetLogsDefault prepares the GetLogsDefault request.
func (c DeploymentScriptsClient) preparerForGetLogsDefault(ctx context.Context, id DeploymentScriptId, options GetLogsDefaultOperationOptions) (*http.Request, error) {
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
		autorest.WithPath(fmt.Sprintf("%s/logs/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetLogsDefault handles the response to the GetLogsDefault request. The method always
// closes the http.Response Body.
func (c DeploymentScriptsClient) responderForGetLogsDefault(resp *http.Response) (result GetLogsDefaultOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
