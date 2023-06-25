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

type GetSchemaJsonOperationResponse struct {
	HttpResponse *http.Response
	Model        *JsonSchema
}

// GetSchemaJson ...
func (c WorkflowTriggersClient) GetSchemaJson(ctx context.Context, id TriggerId) (result GetSchemaJsonOperationResponse, err error) {
	req, err := c.preparerForGetSchemaJson(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "GetSchemaJson", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "GetSchemaJson", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetSchemaJson(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflowtriggers.WorkflowTriggersClient", "GetSchemaJson", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetSchemaJson prepares the GetSchemaJson request.
func (c WorkflowTriggersClient) preparerForGetSchemaJson(ctx context.Context, id TriggerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/schemas/json", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetSchemaJson handles the response to the GetSchemaJson request. The method always
// closes the http.Response Body.
func (c WorkflowTriggersClient) responderForGetSchemaJson(resp *http.Response) (result GetSchemaJsonOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
