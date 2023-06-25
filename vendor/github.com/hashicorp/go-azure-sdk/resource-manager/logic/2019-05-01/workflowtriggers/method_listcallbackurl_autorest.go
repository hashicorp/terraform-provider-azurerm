package workflowtriggers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListCallbackUrlOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkflowTriggerCallbackUrl
}

// ListCallbackUrl ...
func (c WorkflowTriggersClient) ListCallbackUrl(ctx context.Context, id TriggerId) (result ListCallbackUrlOperationResponse, err error) {
	req, err := c.preparerForListCallbackUrl(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "ListCallbackUrl", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "ListCallbackUrl", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListCallbackUrl(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "ListCallbackUrl", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListCallbackUrl prepares the ListCallbackUrl request.
func (c WorkflowTriggersClient) preparerForListCallbackUrl(ctx context.Context, id TriggerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listCallbackUrl", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListCallbackUrl handles the response to the ListCallbackUrl request. The method always
// closes the http.Response Body.
func (c WorkflowTriggersClient) responderForListCallbackUrl(resp *http.Response) (result ListCallbackUrlOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
