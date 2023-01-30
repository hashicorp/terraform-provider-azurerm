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

type WorkflowVersionTriggersListCallbackUrlOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkflowTriggerCallbackUrl
}

// WorkflowVersionTriggersListCallbackUrl ...
func (c WorkflowTriggersClient) WorkflowVersionTriggersListCallbackUrl(ctx context.Context, id VersionTriggerId, input GetCallbackUrlParameters) (result WorkflowVersionTriggersListCallbackUrlOperationResponse, err error) {
	req, err := c.preparerForWorkflowVersionTriggersListCallbackUrl(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "WorkflowVersionTriggersListCallbackUrl", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "WorkflowVersionTriggersListCallbackUrl", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForWorkflowVersionTriggersListCallbackUrl(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "WorkflowVersionTriggersListCallbackUrl", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForWorkflowVersionTriggersListCallbackUrl prepares the WorkflowVersionTriggersListCallbackUrl request.
func (c WorkflowTriggersClient) preparerForWorkflowVersionTriggersListCallbackUrl(ctx context.Context, id VersionTriggerId, input GetCallbackUrlParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listCallbackUrl", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForWorkflowVersionTriggersListCallbackUrl handles the response to the WorkflowVersionTriggersListCallbackUrl request. The method always
// closes the http.Response Body.
func (c WorkflowTriggersClient) responderForWorkflowVersionTriggersListCallbackUrl(resp *http.Response) (result WorkflowVersionTriggersListCallbackUrlOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
