package mongorbacs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBResourcesListMongoUserDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *MongoUserDefinitionListResult
}

// MongoDBResourcesListMongoUserDefinitions ...
func (c MongorbacsClient) MongoDBResourcesListMongoUserDefinitions(ctx context.Context, id DatabaseAccountId) (result MongoDBResourcesListMongoUserDefinitionsOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesListMongoUserDefinitions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesListMongoUserDefinitions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesListMongoUserDefinitions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMongoDBResourcesListMongoUserDefinitions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesListMongoUserDefinitions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMongoDBResourcesListMongoUserDefinitions prepares the MongoDBResourcesListMongoUserDefinitions request.
func (c MongorbacsClient) preparerForMongoDBResourcesListMongoUserDefinitions(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/mongodbUserDefinitions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMongoDBResourcesListMongoUserDefinitions handles the response to the MongoDBResourcesListMongoUserDefinitions request. The method always
// closes the http.Response Body.
func (c MongorbacsClient) responderForMongoDBResourcesListMongoUserDefinitions(resp *http.Response) (result MongoDBResourcesListMongoUserDefinitionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
