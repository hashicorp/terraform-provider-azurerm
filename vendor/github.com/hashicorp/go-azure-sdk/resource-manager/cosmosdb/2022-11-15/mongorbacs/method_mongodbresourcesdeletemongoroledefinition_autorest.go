package mongorbacs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDBResourcesDeleteMongoRoleDefinitionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesDeleteMongoRoleDefinition ...
func (c MongorbacsClient) MongoDBResourcesDeleteMongoRoleDefinition(ctx context.Context, id MongodbRoleDefinitionId) (result MongoDBResourcesDeleteMongoRoleDefinitionOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesDeleteMongoRoleDefinition(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesDeleteMongoRoleDefinition", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesDeleteMongoRoleDefinition(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesDeleteMongoRoleDefinition", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesDeleteMongoRoleDefinitionThenPoll performs MongoDBResourcesDeleteMongoRoleDefinition then polls until it's completed
func (c MongorbacsClient) MongoDBResourcesDeleteMongoRoleDefinitionThenPoll(ctx context.Context, id MongodbRoleDefinitionId) error {
	result, err := c.MongoDBResourcesDeleteMongoRoleDefinition(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesDeleteMongoRoleDefinition: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesDeleteMongoRoleDefinition: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesDeleteMongoRoleDefinition prepares the MongoDBResourcesDeleteMongoRoleDefinition request.
func (c MongorbacsClient) preparerForMongoDBResourcesDeleteMongoRoleDefinition(ctx context.Context, id MongodbRoleDefinitionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMongoDBResourcesDeleteMongoRoleDefinition sends the MongoDBResourcesDeleteMongoRoleDefinition request. The method will close the
// http.Response Body if it receives an error.
func (c MongorbacsClient) senderForMongoDBResourcesDeleteMongoRoleDefinition(ctx context.Context, req *http.Request) (future MongoDBResourcesDeleteMongoRoleDefinitionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
