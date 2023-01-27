package tableservice

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Table
}

// TableUpdate ...
func (c TableServiceClient) TableUpdate(ctx context.Context, id TableId, input Table) (result TableUpdateOperationResponse, err error) {
	req, err := c.preparerForTableUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTableUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTableUpdate prepares the TableUpdate request.
func (c TableServiceClient) preparerForTableUpdate(ctx context.Context, id TableId, input Table) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTableUpdate handles the response to the TableUpdate request. The method always
// closes the http.Response Body.
func (c TableServiceClient) responderForTableUpdate(resp *http.Response) (result TableUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
