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

type MongoDBResourcesCreateUpdateMongoUserDefinitionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesCreateUpdateMongoUserDefinition ...
func (c MongorbacsClient) MongoDBResourcesCreateUpdateMongoUserDefinition(ctx context.Context, id MongodbUserDefinitionId, input MongoUserDefinitionCreateUpdateParameters) (result MongoDBResourcesCreateUpdateMongoUserDefinitionOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesCreateUpdateMongoUserDefinition(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesCreateUpdateMongoUserDefinition", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesCreateUpdateMongoUserDefinition(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesCreateUpdateMongoUserDefinition", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesCreateUpdateMongoUserDefinitionThenPoll performs MongoDBResourcesCreateUpdateMongoUserDefinition then polls until it's completed
func (c MongorbacsClient) MongoDBResourcesCreateUpdateMongoUserDefinitionThenPoll(ctx context.Context, id MongodbUserDefinitionId, input MongoUserDefinitionCreateUpdateParameters) error {
	result, err := c.MongoDBResourcesCreateUpdateMongoUserDefinition(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesCreateUpdateMongoUserDefinition: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesCreateUpdateMongoUserDefinition: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesCreateUpdateMongoUserDefinition prepares the MongoDBResourcesCreateUpdateMongoUserDefinition request.
func (c MongorbacsClient) preparerForMongoDBResourcesCreateUpdateMongoUserDefinition(ctx context.Context, id MongodbUserDefinitionId, input MongoUserDefinitionCreateUpdateParameters) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForMongoDBResourcesCreateUpdateMongoUserDefinition sends the MongoDBResourcesCreateUpdateMongoUserDefinition request. The method will close the
// http.Response Body if it receives an error.
func (c MongorbacsClient) senderForMongoDBResourcesCreateUpdateMongoUserDefinition(ctx context.Context, req *http.Request) (future MongoDBResourcesCreateUpdateMongoUserDefinitionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
