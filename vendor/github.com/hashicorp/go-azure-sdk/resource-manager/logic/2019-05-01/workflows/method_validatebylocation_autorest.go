package workflows

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidateByLocationOperationResponse struct {
	HttpResponse *http.Response
}

// ValidateByLocation ...
func (c WorkflowsClient) ValidateByLocation(ctx context.Context, id LocationWorkflowId, input Workflow) (result ValidateByLocationOperationResponse, err error) {
	req, err := c.preparerForValidateByLocation(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "ValidateByLocation", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "ValidateByLocation", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForValidateByLocation(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "ValidateByLocation", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForValidateByLocation prepares the ValidateByLocation request.
func (c WorkflowsClient) preparerForValidateByLocation(ctx context.Context, id LocationWorkflowId, input Workflow) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/validate", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForValidateByLocation handles the response to the ValidateByLocation request. The method always
// closes the http.Response Body.
func (c WorkflowsClient) responderForValidateByLocation(resp *http.Response) (result ValidateByLocationOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
