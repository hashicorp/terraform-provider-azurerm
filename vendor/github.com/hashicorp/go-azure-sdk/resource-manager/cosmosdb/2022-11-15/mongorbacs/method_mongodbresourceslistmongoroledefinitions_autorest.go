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

type MongoDBResourcesListMongoRoleDefinitionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *MongoRoleDefinitionListResult
}

// MongoDBResourcesListMongoRoleDefinitions ...
func (c MongorbacsClient) MongoDBResourcesListMongoRoleDefinitions(ctx context.Context, id DatabaseAccountId) (result MongoDBResourcesListMongoRoleDefinitionsOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesListMongoRoleDefinitions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesListMongoRoleDefinitions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesListMongoRoleDefinitions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMongoDBResourcesListMongoRoleDefinitions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesListMongoRoleDefinitions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMongoDBResourcesListMongoRoleDefinitions prepares the MongoDBResourcesListMongoRoleDefinitions request.
func (c MongorbacsClient) preparerForMongoDBResourcesListMongoRoleDefinitions(ctx context.Context, id DatabaseAccountId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/mongodbRoleDefinitions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForMongoDBResourcesListMongoRoleDefinitions handles the response to the MongoDBResourcesListMongoRoleDefinitions request. The method always
// closes the http.Response Body.
func (c MongorbacsClient) responderForMongoDBResourcesListMongoRoleDefinitions(resp *http.Response) (result MongoDBResourcesListMongoRoleDefinitionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
