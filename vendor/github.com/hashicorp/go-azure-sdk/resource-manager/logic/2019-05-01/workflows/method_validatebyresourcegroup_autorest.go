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

type ValidateByResourceGroupOperationResponse struct {
	HttpResponse *http.Response
}

// ValidateByResourceGroup ...
func (c WorkflowsClient) ValidateByResourceGroup(ctx context.Context, id WorkflowId, input Workflow) (result ValidateByResourceGroupOperationResponse, err error) {
	req, err := c.preparerForValidateByResourceGroup(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "ValidateByResourceGroup", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "ValidateByResourceGroup", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForValidateByResourceGroup(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "ValidateByResourceGroup", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForValidateByResourceGroup prepares the ValidateByResourceGroup request.
func (c WorkflowsClient) preparerForValidateByResourceGroup(ctx context.Context, id WorkflowId, input Workflow) (*http.Request, error) {
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

// responderForValidateByResourceGroup handles the response to the ValidateByResourceGroup request. The method always
// closes the http.Response Body.
func (c WorkflowsClient) responderForValidateByResourceGroup(resp *http.Response) (result ValidateByResourceGroupOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
