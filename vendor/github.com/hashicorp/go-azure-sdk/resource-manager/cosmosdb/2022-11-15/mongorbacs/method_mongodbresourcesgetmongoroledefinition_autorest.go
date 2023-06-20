package mongorbacs

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBResourcesGetMongoRoleDefinitionOperationResponse struct {
	HttpResponse *http.Response
	Model        *MongoRoleDefinitionGetResults
}

// MongoDBResourcesGetMongoRoleDefinition ...
func (c MongorbacsClient) MongoDBResourcesGetMongoRoleDefinition(ctx context.Context, id MongodbRoleDefinitionId) (result MongoDBResourcesGetMongoRoleDefinitionOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesGetMongoRoleDefinition(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesGetMongoRoleDefinition", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesGetMongoRoleDefinition", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForMongoDBResourcesGetMongoRoleDefinition(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesGetMongoRoleDefinition", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForMongoDBResourcesGetMongoRoleDefinition prepares the MongoDBResourcesGetMongoRoleDefinition request.
func (c MongorbacsClient) preparerForMongoDBResourcesGetMongoRoleDefinition(ctx context.Context, id MongodbRoleDefinitionId) (*http.Request, error) {
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

// responderForMongoDBResourcesGetMongoRoleDefinition handles the response to the MongoDBResourcesGetMongoRoleDefinition request. The method always
// closes the http.Response Body.
func (c MongorbacsClient) responderForMongoDBResourcesGetMongoRoleDefinition(resp *http.Response) (result MongoDBResourcesGetMongoRoleDefinitionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
