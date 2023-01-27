package analysisservices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServersListSkusForNewOperationResponse struct {
	HttpResponse *http.Response
	Model        *SkuEnumerationForNewResourceResult
}

// ServersListSkusForNew ...
func (c AnalysisServicesClient) ServersListSkusForNew(ctx context.Context, id commonids.SubscriptionId) (result ServersListSkusForNewOperationResponse, err error) {
	req, err := c.preparerForServersListSkusForNew(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "analysisservices.AnalysisServicesClient", "ServersListSkusForNew", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "analysisservices.AnalysisServicesClient", "ServersListSkusForNew", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForServersListSkusForNew(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "analysisservices.AnalysisServicesClient", "ServersListSkusForNew", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForServersListSkusForNew prepares the ServersListSkusForNew request.
func (c AnalysisServicesClient) preparerForServersListSkusForNew(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.AnalysisServices/skus", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForServersListSkusForNew handles the response to the ServersListSkusForNew request. The method always
// closes the http.Response Body.
func (c AnalysisServicesClient) responderForServersListSkusForNew(resp *http.Response) (result ServersListSkusForNewOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
