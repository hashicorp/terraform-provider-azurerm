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

type DeactivateRevisionOperationResponse struct {
	HttpResponse *http.Response
}

// DeactivateRevision ...
func (c ContainerAppsRevisionsClient) DeactivateRevision(ctx context.Context, id RevisionId) (result DeactivateRevisionOperationResponse, err error) {
	req, err := c.preparerForDeactivateRevision(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "DeactivateRevision", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "DeactivateRevision", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeactivateRevision(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerappsrevisions.ContainerAppsRevisionsClient", "DeactivateRevision", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeactivateRevision prepares the DeactivateRevision request.
func (c ContainerAppsRevisionsClient) preparerForDeactivateRevision(ctx context.Context, id RevisionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/deactivate", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDeactivateRevision handles the response to the DeactivateRevision request. The method always
// closes the http.Response Body.
func (c ContainerAppsRevisionsClient) responderForDeactivateRevision(resp *http.Response) (result DeactivateRevisionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
