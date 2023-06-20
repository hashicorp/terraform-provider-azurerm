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

type GenerateUpgradedDefinitionOperationResponse struct {
	HttpResponse *http.Response
	Model        *interface{}
}

// GenerateUpgradedDefinition ...
func (c WorkflowsClient) GenerateUpgradedDefinition(ctx context.Context, id WorkflowId, input GenerateUpgradedDefinitionParameters) (result GenerateUpgradedDefinitionOperationResponse, err error) {
	req, err := c.preparerForGenerateUpgradedDefinition(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "GenerateUpgradedDefinition", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "GenerateUpgradedDefinition", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGenerateUpgradedDefinition(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "workflows.WorkflowsClient", "GenerateUpgradedDefinition", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGenerateUpgradedDefinition prepares the GenerateUpgradedDefinition request.
func (c WorkflowsClient) preparerForGenerateUpgradedDefinition(ctx context.Context, id WorkflowId, input GenerateUpgradedDefinitionParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/generateUpgradedDefinition", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGenerateUpgradedDefinition handles the response to the GenerateUpgradedDefinition request. The method always
// closes the http.Response Body.
func (c WorkflowsClient) responderForGenerateUpgradedDefinition(resp *http.Response) (result GenerateUpgradedDefinitionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
