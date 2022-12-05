package functions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RetrieveDefaultDefinitionOperationResponse struct {
	HttpResponse *http.Response
	Model        *Function
}

// RetrieveDefaultDefinition ...
func (c FunctionsClient) RetrieveDefaultDefinition(ctx context.Context, id FunctionId, input FunctionRetrieveDefaultDefinitionParameters) (result RetrieveDefaultDefinitionOperationResponse, err error) {
	req, err := c.preparerForRetrieveDefaultDefinition(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "functions.FunctionsClient", "RetrieveDefaultDefinition", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "functions.FunctionsClient", "RetrieveDefaultDefinition", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRetrieveDefaultDefinition(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "functions.FunctionsClient", "RetrieveDefaultDefinition", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRetrieveDefaultDefinition prepares the RetrieveDefaultDefinition request.
func (c FunctionsClient) preparerForRetrieveDefaultDefinition(ctx context.Context, id FunctionId, input FunctionRetrieveDefaultDefinitionParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/retrieveDefaultDefinition", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRetrieveDefaultDefinition handles the response to the RetrieveDefaultDefinition request. The method always
// closes the http.Response Body.
func (c FunctionsClient) responderForRetrieveDefaultDefinition(resp *http.Response) (result RetrieveDefaultDefinitionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
