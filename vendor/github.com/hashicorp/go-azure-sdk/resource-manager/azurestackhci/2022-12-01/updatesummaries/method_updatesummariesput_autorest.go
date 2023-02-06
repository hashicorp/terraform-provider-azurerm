package updatesummaries

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateSummariesPutOperationResponse struct {
	HttpResponse *http.Response
	Model        *UpdateSummaries
}

// UpdateSummariesPut ...
func (c UpdateSummariesClient) UpdateSummariesPut(ctx context.Context, id ClusterId, input UpdateSummaries) (result UpdateSummariesPutOperationResponse, err error) {
	req, err := c.preparerForUpdateSummariesPut(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesPut", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesPut", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUpdateSummariesPut(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesPut", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUpdateSummariesPut prepares the UpdateSummariesPut request.
func (c UpdateSummariesClient) preparerForUpdateSummariesPut(ctx context.Context, id ClusterId, input UpdateSummaries) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/updateSummaries/default", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUpdateSummariesPut handles the response to the UpdateSummariesPut request. The method always
// closes the http.Response Body.
func (c UpdateSummariesClient) responderForUpdateSummariesPut(resp *http.Response) (result UpdateSummariesPutOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
