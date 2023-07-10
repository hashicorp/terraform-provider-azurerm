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

type UpdateSummariesGetOperationResponse struct {
	HttpResponse *http.Response
	Model        *UpdateSummaries
}

// UpdateSummariesGet ...
func (c UpdateSummariesClient) UpdateSummariesGet(ctx context.Context, id ClusterId) (result UpdateSummariesGetOperationResponse, err error) {
	req, err := c.preparerForUpdateSummariesGet(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesGet", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesGet", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForUpdateSummariesGet(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "updatesummaries.UpdateSummariesClient", "UpdateSummariesGet", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForUpdateSummariesGet prepares the UpdateSummariesGet request.
func (c UpdateSummariesClient) preparerForUpdateSummariesGet(ctx context.Context, id ClusterId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/updateSummaries/default", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForUpdateSummariesGet handles the response to the UpdateSummariesGet request. The method always
// closes the http.Response Body.
func (c UpdateSummariesClient) responderForUpdateSummariesGet(resp *http.Response) (result UpdateSummariesGetOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
