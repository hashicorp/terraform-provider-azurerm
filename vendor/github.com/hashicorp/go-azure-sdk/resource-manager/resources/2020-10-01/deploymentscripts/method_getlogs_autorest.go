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

type GetLogsOperationResponse struct {
	HttpResponse *http.Response
	Model        *ScriptLogsList
}

// GetLogs ...
func (c DeploymentScriptsClient) GetLogs(ctx context.Context, id DeploymentScriptId) (result GetLogsOperationResponse, err error) {
	req, err := c.preparerForGetLogs(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentscripts.DeploymentScriptsClient", "GetLogs", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentscripts.DeploymentScriptsClient", "GetLogs", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetLogs(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "deploymentscripts.DeploymentScriptsClient", "GetLogs", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetLogs prepares the GetLogs request.
func (c DeploymentScriptsClient) preparerForGetLogs(ctx context.Context, id DeploymentScriptId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/logs", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetLogs handles the response to the GetLogs request. The method always
// closes the http.Response Body.
func (c DeploymentScriptsClient) responderForGetLogs(resp *http.Response) (result GetLogsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
