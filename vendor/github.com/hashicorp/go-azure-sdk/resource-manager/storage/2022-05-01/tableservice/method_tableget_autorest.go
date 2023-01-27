package tableservice

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *Table
}

// TableGet ...
func (c TableServiceClient) TableGet(ctx context.Context, id TableId) (result TableGetOperationResponse, err error) {
	req, err := c.preparerForTableGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTableGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTableGet prepares the TableGet request.
func (c TableServiceClient) preparerForTableGet(ctx context.Context, id TableId) (*http.Request, error) {
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

// responderForTableGet handles the response to the TableGet request. The method always
// closes the http.Response Body.
func (c TableServiceClient) responderForTableGet(resp *http.Response) (result TableGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
