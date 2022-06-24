package findrestorabletimeranges

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorableTimeRangesFindOperationResponse struct {
	HttpResponse *http.Response
	Model        *AzureBackupFindRestorableTimeRangesResponseResource
}

// RestorableTimeRangesFind ...
func (c FindRestorableTimeRangesClient) RestorableTimeRangesFind(ctx context.Context, id BackupInstanceId, input AzureBackupFindRestorableTimeRangesRequest) (result RestorableTimeRangesFindOperationResponse, err error) {
	req, err := c.preparerForRestorableTimeRangesFind(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "findrestorabletimeranges.FindRestorableTimeRangesClient", "RestorableTimeRangesFind", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "findrestorabletimeranges.FindRestorableTimeRangesClient", "RestorableTimeRangesFind", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForRestorableTimeRangesFind(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "findrestorabletimeranges.FindRestorableTimeRangesClient", "RestorableTimeRangesFind", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForRestorableTimeRangesFind prepares the RestorableTimeRangesFind request.
func (c FindRestorableTimeRangesClient) preparerForRestorableTimeRangesFind(ctx context.Context, id BackupInstanceId, input AzureBackupFindRestorableTimeRangesRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/findRestorableTimeRanges", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForRestorableTimeRangesFind handles the response to the RestorableTimeRangesFind request. The method always
// closes the http.Response Body.
func (c FindRestorableTimeRangesClient) responderForRestorableTimeRangesFind(resp *http.Response) (result RestorableTimeRangesFindOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
