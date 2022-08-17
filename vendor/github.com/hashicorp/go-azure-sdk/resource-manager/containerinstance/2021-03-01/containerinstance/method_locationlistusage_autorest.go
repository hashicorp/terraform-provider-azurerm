package containerinstance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocationListUsageOperationResponse struct {
	HttpResponse *http.Response
	Model        *UsageListResult
}

// LocationListUsage ...
func (c ContainerInstanceClient) LocationListUsage(ctx context.Context, id LocationId) (result LocationListUsageOperationResponse, err error) {
	req, err := c.preparerForLocationListUsage(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListUsage", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListUsage", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLocationListUsage(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerinstance.ContainerInstanceClient", "LocationListUsage", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLocationListUsage prepares the LocationListUsage request.
func (c ContainerInstanceClient) preparerForLocationListUsage(ctx context.Context, id LocationId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/usages", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForLocationListUsage handles the response to the LocationListUsage request. The method always
// closes the http.Response Body.
func (c ContainerInstanceClient) responderForLocationListUsage(resp *http.Response) (result LocationListUsageOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
