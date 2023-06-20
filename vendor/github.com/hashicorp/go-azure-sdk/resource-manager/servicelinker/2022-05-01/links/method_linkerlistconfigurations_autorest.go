package links

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkerListConfigurationsOperationResponse struct {
	HttpResponse *http.Response
	Model        *SourceConfigurationResult
}

// LinkerListConfigurations ...
func (c LinksClient) LinkerListConfigurations(ctx context.Context, id ScopedLinkerId) (result LinkerListConfigurationsOperationResponse, err error) {
	req, err := c.preparerForLinkerListConfigurations(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.LinksClient", "LinkerListConfigurations", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.LinksClient", "LinkerListConfigurations", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForLinkerListConfigurations(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "links.LinksClient", "LinkerListConfigurations", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForLinkerListConfigurations prepares the LinkerListConfigurations request.
func (c LinksClient) preparerForLinkerListConfigurations(ctx context.Context, id ScopedLinkerId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/listConfigurations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForLinkerListConfigurations handles the response to the LinkerListConfigurations request. The method always
// closes the http.Response Body.
func (c LinksClient) responderForLinkerListConfigurations(resp *http.Response) (result LinkerListConfigurationsOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
