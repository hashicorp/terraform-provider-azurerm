package digitaltwinsinstance

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DigitalTwinsGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *DigitalTwinsDescription
}

// DigitalTwinsGet ...
func (c DigitalTwinsInstanceClient) DigitalTwinsGet(ctx context.Context, id DigitalTwinsInstanceId) (result DigitalTwinsGetOperationResponse, err error) {
	req, err := c.preparerForDigitalTwinsGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDigitalTwinsGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "digitaltwinsinstance.DigitalTwinsInstanceClient", "DigitalTwinsGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDigitalTwinsGet prepares the DigitalTwinsGet request.
func (c DigitalTwinsInstanceClient) preparerForDigitalTwinsGet(ctx context.Context, id DigitalTwinsInstanceId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDigitalTwinsGet handles the response to the DigitalTwinsGet request. The method always
// closes the http.Response Body.
func (c DigitalTwinsInstanceClient) responderForDigitalTwinsGet(resp *http.Response) (result DigitalTwinsGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
