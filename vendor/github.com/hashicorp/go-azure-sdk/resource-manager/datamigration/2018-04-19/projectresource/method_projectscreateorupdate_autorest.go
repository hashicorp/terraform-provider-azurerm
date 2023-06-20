package projectresource

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectsCreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *Project
}

// ProjectsCreateOrUpdate ...
func (c ProjectResourceClient) ProjectsCreateOrUpdate(ctx context.Context, id ProjectId, input Project) (result ProjectsCreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForProjectsCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsCreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsCreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForProjectsCreateOrUpdate(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "projectresource.ProjectResourceClient", "ProjectsCreateOrUpdate", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForProjectsCreateOrUpdate prepares the ProjectsCreateOrUpdate request.
func (c ProjectResourceClient) preparerForProjectsCreateOrUpdate(ctx context.Context, id ProjectId, input Project) (*http.Request, error) {
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

// responderForProjectsCreateOrUpdate handles the response to the ProjectsCreateOrUpdate request. The method always
// closes the http.Response Body.
func (c ProjectResourceClient) responderForProjectsCreateOrUpdate(resp *http.Response) (result ProjectsCreateOrUpdateOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusCreated, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
