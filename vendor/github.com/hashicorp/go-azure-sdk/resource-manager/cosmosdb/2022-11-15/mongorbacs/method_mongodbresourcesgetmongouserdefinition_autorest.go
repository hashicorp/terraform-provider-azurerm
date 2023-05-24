package mongorbacs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBResourcesGetMongoUserDefinitionOperationResponse struct {
	HttpResponse *http.Response
	Model        *MongoUserDefinitionGetResults
}

// MongoDBResourcesGetMongoUserDefinition ...
func (c MongorbacsClient) MongoDBResourcesGetMongoUserDefinition(ctx context.Context, id MongodbUserDefinitionId) (result MongoDBResourcesGetMongoUserDefinitionOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesGetMongoUserDefinition(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesGetMongoUserDefinition", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesGetMongoUserDefinition", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMongoDBResourcesGetMongoUserDefinition(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesGetMongoUserDefinition", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMongoDBResourcesGetMongoUserDefinition prepares the MongoDBResourcesGetMongoUserDefinition request.
func (c MongorbacsClient) preparerForMongoDBResourcesGetMongoUserDefinition(ctx context.Context, id MongodbUserDefinitionId) (*http.Request, error) {
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

// responderForMongoDBResourcesGetMongoUserDefinition handles the response to the MongoDBResourcesGetMongoUserDefinition request. The method always
// closes the http.Response Body.
func (c MongorbacsClient) responderForMongoDBResourcesGetMongoUserDefinition(resp *http.Response) (result MongoDBResourcesGetMongoUserDefinitionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
