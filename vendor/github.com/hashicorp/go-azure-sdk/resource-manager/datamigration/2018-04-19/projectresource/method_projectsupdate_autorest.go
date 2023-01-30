package projectresource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectsUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Project
}

// ProjectsUpdate ...
func (c ProjectResourceClient) ProjectsUpdate(ctx context.Context, id ProjectId, input Project) (result ProjectsUpdateOperationResponse, err error) {
	req, err := c.preparerForProjectsUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForProjectsUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForProjectsUpdate prepares the ProjectsUpdate request.
func (c ProjectResourceClient) preparerForProjectsUpdate(ctx context.Context, id ProjectId, input Project) (*http.Request, error) {
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

// responderForProjectsUpdate handles the response to the ProjectsUpdate request. The method always
// closes the http.Response Body.
func (c ProjectResourceClient) responderForProjectsUpdate(resp *http.Response) (result ProjectsUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
