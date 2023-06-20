package containerservices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListOrchestratorsOperationResponse struct {
	HttpResponse *http.Response
	Model        *OrchestratorVersionProfileListResult
}

type ListOrchestratorsOperationOptions struct {
	ResourceType *string
}

func DefaultListOrchestratorsOperationOptions() ListOrchestratorsOperationOptions {
	return ListOrchestratorsOperationOptions{}
}

func (o ListOrchestratorsOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListOrchestratorsOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.ResourceType != nil {
		out["resource-type"] = *o.ResourceType
	}

	return out
}

// ListOrchestrators ...
func (c ContainerServicesClient) ListOrchestrators(ctx context.Context, id LocationId, options ListOrchestratorsOperationOptions) (result ListOrchestratorsOperationResponse, err error) {
	req, err := c.preparerForListOrchestrators(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservices.ContainerServicesClient", "ListOrchestrators", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservices.ContainerServicesClient", "ListOrchestrators", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListOrchestrators(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerservices.ContainerServicesClient", "ListOrchestrators", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListOrchestrators prepares the ListOrchestrators request.
func (c ContainerServicesClient) preparerForListOrchestrators(ctx context.Context, id LocationId, options ListOrchestratorsOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/orchestrators", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListOrchestrators handles the response to the ListOrchestrators request. The method always
// closes the http.Response Body.
func (c ContainerServicesClient) responderForListOrchestrators(resp *http.Response) (result ListOrchestratorsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
