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

type MongoDBResourcesCreateUpdateMongoRoleDefinitionOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

// MongoDBResourcesCreateUpdateMongoRoleDefinition ...
func (c MongorbacsClient) MongoDBResourcesCreateUpdateMongoRoleDefinition(ctx context.Context, id MongodbRoleDefinitionId, input MongoRoleDefinitionCreateUpdateParameters) (result MongoDBResourcesCreateUpdateMongoRoleDefinitionOperationResponse, err error) {
	req, err := c.preparerForMongoDBResourcesCreateUpdateMongoRoleDefinition(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesCreateUpdateMongoRoleDefinition", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForMongoDBResourcesCreateUpdateMongoRoleDefinition(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mongorbacs.MongorbacsClient", "MongoDBResourcesCreateUpdateMongoRoleDefinition", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// MongoDBResourcesCreateUpdateMongoRoleDefinitionThenPoll performs MongoDBResourcesCreateUpdateMongoRoleDefinition then polls until it's completed
func (c MongorbacsClient) MongoDBResourcesCreateUpdateMongoRoleDefinitionThenPoll(ctx context.Context, id MongodbRoleDefinitionId, input MongoRoleDefinitionCreateUpdateParameters) error {
	result, err := c.MongoDBResourcesCreateUpdateMongoRoleDefinition(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing MongoDBResourcesCreateUpdateMongoRoleDefinition: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after MongoDBResourcesCreateUpdateMongoRoleDefinition: %+v", err)
	}

	return nil
}

// preparerForMongoDBResourcesCreateUpdateMongoRoleDefinition prepares the MongoDBResourcesCreateUpdateMongoRoleDefinition request.
func (c MongorbacsClient) preparerForMongoDBResourcesCreateUpdateMongoRoleDefinition(ctx context.Context, id MongodbRoleDefinitionId, input MongoRoleDefinitionCreateUpdateParameters) (*http.Request, error) {
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

// senderForMongoDBResourcesCreateUpdateMongoRoleDefinition sends the MongoDBResourcesCreateUpdateMongoRoleDefinition request. The method will close the
// http.Response Body if it receives an error.
func (c MongorbacsClient) senderForMongoDBResourcesCreateUpdateMongoRoleDefinition(ctx context.Context, req *http.Request) (future MongoDBResourcesCreateUpdateMongoRoleDefinitionOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
