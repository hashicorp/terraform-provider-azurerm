package containerapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListCustomHostNameAnalysisOperationResponse struct {
	HttpResponse *http.Response
	Model        *CustomHostnameAnalysisResult
}

type ListCustomHostNameAnalysisOperationOptions struct {
	CustomHostname *string
}

func DefaultListCustomHostNameAnalysisOperationOptions() ListCustomHostNameAnalysisOperationOptions {
	return ListCustomHostNameAnalysisOperationOptions{}
}

func (o ListCustomHostNameAnalysisOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListCustomHostNameAnalysisOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.CustomHostname != nil {
		out["customHostname"] = *o.CustomHostname
	}

	return out
}

// ListCustomHostNameAnalysis ...
func (c ContainerAppsClient) ListCustomHostNameAnalysis(ctx context.Context, id ContainerAppId, options ListCustomHostNameAnalysisOperationOptions) (result ListCustomHostNameAnalysisOperationResponse, err error) {
	req, err := c.preparerForListCustomHostNameAnalysis(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "ListCustomHostNameAnalysis", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "ListCustomHostNameAnalysis", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListCustomHostNameAnalysis(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerapps.ContainerAppsClient", "ListCustomHostNameAnalysis", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListCustomHostNameAnalysis prepares the ListCustomHostNameAnalysis request.
func (c ContainerAppsClient) preparerForListCustomHostNameAnalysis(ctx context.Context, id ContainerAppId, options ListCustomHostNameAnalysisOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/listCustomHostNameAnalysis", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListCustomHostNameAnalysis handles the response to the ListCustomHostNameAnalysis request. The method always
// closes the http.Response Body.
func (c ContainerAppsClient) responderForListCustomHostNameAnalysis(resp *http.Response) (result ListCustomHostNameAnalysisOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
