package tableservice

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableCreateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Table
}

// TableCreate ...
func (c TableServiceClient) TableCreate(ctx context.Context, id TableId, input Table) (result TableCreateOperationResponse, err error) {
	req, err := c.preparerForTableCreate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableCreate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableCreate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTableCreate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableCreate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTableCreate prepares the TableCreate request.
func (c TableServiceClient) preparerForTableCreate(ctx context.Context, id TableId, input Table) (*http.Request, error) {
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

// responderForTableCreate handles the response to the TableCreate request. The method always
// closes the http.Response Body.
func (c TableServiceClient) responderForTableCreate(resp *http.Response) (result TableCreateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
