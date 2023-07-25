package containerappsrevisions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivateRevisionOperationResponse struct {
	HttpResponse *http.Response
}

// ActivateRevision ...
func (c ContainerAppsRevisionsClient) ActivateRevision(ctx context.Context, id RevisionId) (result ActivateRevisionOperationResponse, err error) {
	req, err := c.preparerForActivateRevision(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "ActivateRevision", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "ActivateRevision", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForActivateRevision(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "ActivateRevision", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForActivateRevision prepares the ActivateRevision request.
func (c ContainerAppsRevisionsClient) preparerForActivateRevision(ctx context.Context, id RevisionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/activate", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForActivateRevision handles the response to the ActivateRevision request. The method always
// closes the http.Response Body.
func (c ContainerAppsRevisionsClient) responderForActivateRevision(resp *http.Response) (result ActivateRevisionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
