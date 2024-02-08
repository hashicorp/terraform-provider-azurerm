package containerappsrevisions

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetRevisionOperationResponse struct {
	HttpResponse *http.Response
	Model        *Revision
}

// GetRevision ...
func (c ContainerAppsRevisionsClient) GetRevision(ctx context.Context, id RevisionId) (result GetRevisionOperationResponse, err error) {
	req, err := c.preparerForGetRevision(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "GetRevision", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "GetRevision", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetRevision(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "GetRevision", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetRevision prepares the GetRevision request.
func (c ContainerAppsRevisionsClient) preparerForGetRevision(ctx context.Context, id RevisionId) (*http.Request, error) {
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

// responderForGetRevision handles the response to the GetRevision request. The method always
// closes the http.Response Body.
func (c ContainerAppsRevisionsClient) responderForGetRevision(resp *http.Response) (result GetRevisionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
