package clusters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListLanguageExtensionsOperationResponse struct {
	HttpResponse *http.Response
	Model        *LanguageExtensionsList
}

// ListLanguageExtensions ...
func (c ClustersClient) ListLanguageExtensions(ctx context.Context, id ClusterId) (result ListLanguageExtensionsOperationResponse, err error) {
	req, err := c.preparerForListLanguageExtensions(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListLanguageExtensions", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListLanguageExtensions", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForListLanguageExtensions(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "clusters.ClustersClient", "ListLanguageExtensions", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForListLanguageExtensions prepares the ListLanguageExtensions request.
func (c ClustersClient) preparerForListLanguageExtensions(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listLanguageExtensions", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListLanguageExtensions handles the response to the ListLanguageExtensions request. The method always
// closes the http.Response Body.
func (c ClustersClient) responderForListLanguageExtensions(resp *http.Response) (result ListLanguageExtensionsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
