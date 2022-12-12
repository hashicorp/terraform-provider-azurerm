package nginxdeployment

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *NginxDeployment
}

// DeploymentsGet ...
func (c NginxDeploymentClient) DeploymentsGet(ctx context.Context, id NginxDeploymentId) (result DeploymentsGetOperationResponse, err error) {
	req, err := c.preparerForDeploymentsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeploymentsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "nginxdeployment.NginxDeploymentClient", "DeploymentsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeploymentsGet prepares the DeploymentsGet request.
func (c NginxDeploymentClient) preparerForDeploymentsGet(ctx context.Context, id NginxDeploymentId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDeploymentsGet handles the response to the DeploymentsGet request. The method always
// closes the http.Response Body.
func (c NginxDeploymentClient) responderForDeploymentsGet(resp *http.Response) (result DeploymentsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
