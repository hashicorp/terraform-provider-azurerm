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

type MongoDBResourcesDeleteMongoUserDefinitionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesDeleteMongoUserDefinition ...
func (c MongorbacsClient) MongoDBResourcesDeleteMongoUserDefinition(ctx context.Context, id MongodbUserDefinitionId) (result MongoDBResourcesDeleteMongoUserDefinitionOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesDeleteMongoUserDefinition(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesDeleteMongoUserDefinition", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesDeleteMongoUserDefinition(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesDeleteMongoUserDefinition", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesDeleteMongoUserDefinitionThenPoll performs MongoDBResourcesDeleteMongoUserDefinition then polls until it's completed
func (c MongorbacsClient) MongoDBResourcesDeleteMongoUserDefinitionThenPoll(ctx context.Context, id MongodbUserDefinitionId) error {
	result, err := c.MongoDBResourcesDeleteMongoUserDefinition(ctx, id)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesDeleteMongoUserDefinition: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesDeleteMongoUserDefinition: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesDeleteMongoUserDefinition prepares the MongoDBResourcesDeleteMongoUserDefinition request.
func (c MongorbacsClient) preparerForMongoDBResourcesDeleteMongoUserDefinition(ctx context.Context, id MongodbUserDefinitionId) (*http.Request, error) {
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

// senderForMongoDBResourcesDeleteMongoUserDefinition sends the MongoDBResourcesDeleteMongoUserDefinition request. The method will close the
// http.Response Body if it receives an error.
func (c MongorbacsClient) senderForMongoDBResourcesDeleteMongoUserDefinition(ctx context.Context, req *http.Request) (future MongoDBResourcesDeleteMongoUserDefinitionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
