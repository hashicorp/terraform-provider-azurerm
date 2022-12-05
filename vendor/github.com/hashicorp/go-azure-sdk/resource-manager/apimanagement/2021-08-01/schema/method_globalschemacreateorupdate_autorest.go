package schema

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

type GlobalSchemaCreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
}

type GlobalSchemaCreateOrUpdateOperationOptions struct {
	IfMatch *string
}

func DefaultGlobalSchemaCreateOrUpdateOperationOptions() GlobalSchemaCreateOrUpdateOperationOptions {
	return GlobalSchemaCreateOrUpdateOperationOptions{}
}

func (o GlobalSchemaCreateOrUpdateOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	if o.IfMatch != nil {
		out["If-Match"] = *o.IfMatch
	}

	return out
}

func (o GlobalSchemaCreateOrUpdateOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

// GlobalSchemaCreateOrUpdate ...
func (c SchemaClient) GlobalSchemaCreateOrUpdate(ctx context.Context, id SchemaId, input GlobalSchemaContract, options GlobalSchemaCreateOrUpdateOperationOptions) (result GlobalSchemaCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForGlobalSchemaCreateOrUpdate(ctx, id, input, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForGlobalSchemaCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "schema.SchemaClient", "GlobalSchemaCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// GlobalSchemaCreateOrUpdateThenPoll performs GlobalSchemaCreateOrUpdate then polls until it's completed
func (c SchemaClient) GlobalSchemaCreateOrUpdateThenPoll(ctx context.Context, id SchemaId, input GlobalSchemaContract, options GlobalSchemaCreateOrUpdateOperationOptions) error {
	result, err := c.GlobalSchemaCreateOrUpdate(ctx, id, input, options)
	if err != nil {
		return fmt.Errorf("performing GlobalSchemaCreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after GlobalSchemaCreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForGlobalSchemaCreateOrUpdate prepares the GlobalSchemaCreateOrUpdate request.
func (c SchemaClient) preparerForGlobalSchemaCreateOrUpdate(ctx context.Context, id SchemaId, input GlobalSchemaContract, options GlobalSchemaCreateOrUpdateOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForGlobalSchemaCreateOrUpdate sends the GlobalSchemaCreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c SchemaClient) senderForGlobalSchemaCreateOrUpdate(ctx context.Context, req *http.Request) (future GlobalSchemaCreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
