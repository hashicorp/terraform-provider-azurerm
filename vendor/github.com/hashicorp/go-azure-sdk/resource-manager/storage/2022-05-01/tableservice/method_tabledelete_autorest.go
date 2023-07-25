package tableservice

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TableDeleteOperationResponse struct {
	HttpResponse *http.Response
}

// TableDelete ...
func (c TableServiceClient) TableDelete(ctx context.Context, id TableId) (result TableDeleteOperationResponse, err error) {
	req, err := c.preparerForTableDelete(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableDelete", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableDelete", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForTableDelete(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "tableservice.TableServiceClient", "TableDelete", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForTableDelete prepares the TableDelete request.
func (c TableServiceClient) preparerForTableDelete(ctx context.Context, id TableId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForTableDelete handles the response to the TableDelete request. The method always
// closes the http.Response Body.
func (c TableServiceClient) responderForTableDelete(resp *http.Response) (result TableDeleteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusNoContent),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
